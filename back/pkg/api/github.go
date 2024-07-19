package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v62/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Github struct {
	GithubAppID    string
	GithubClientID string
	GithubSecret   string
}

type UploadStruct struct {
	Repo     string `json:"repo"`
	Content  string `json:"content"`
	File string `json:"file"`
}

type GitHubMeStruct struct {
	Login string `json:"login"`
	Id    uint   `json:"id"`
	Email string `json:"email"`
}

func GetGithubEnv() Github {
	return Github{
		os.Getenv("GITHUB_APP_ID"),
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_SECRET"),
	}
}

func GithubRedirect(c *fiber.Ctx) error {
	// create a CSRF token and store it locally
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := hex.EncodeToString(b)
	c.Cookie(&fiber.Cookie{
		Name:  "github_state",
		Value: state,
	})
	scope := "repo:status+repo+public_repo+admin:repo_hook+admin:org+admin:public_key+admin:org_hook+user+user:follow+read:gpg_key+delete:packages+read:discussion+workflow+admin:gpg_key+write:packages+delete_repo+read:user+gist+write:public_key+write:org+write:repo_hook+repo:invite+write:gpg_key+read:packages+write:discussion+user:email+notifications+read:public_key+read:org+read:repo_hook+security_events+repo_deployment"

	d := GetGithubEnv()
	link := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%s",
		d.GithubClientID,
		scope,
		fmt.Sprintf("%s/callback", os.Getenv("DOMAIN_NAME")),
		state,
	)
	return c.Redirect(link, http.StatusTemporaryRedirect)
}

func GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if c.Query("state") != c.Cookies("github_state") {
		return c.Status(400).SendString("error - bad state")
	}

	tkn := getGithubAccessToken(code)
	if tkn == "" {
		return c.Status(400).SendString("error")
	}
	c.Cookie(&fiber.Cookie{
		Name:  "tkn",
		Value: tkn,
	})
	return c.Redirect("/repos", http.StatusTemporaryRedirect)
}

func GithubSetToken(c *fiber.Ctx) error {
	tkn := c.Query("tkn")
	data := getGithubData(tkn)
	if data.Id == 0 || data.Login == "" {
		c.Status(400).JSON(fiber.Map{"error": "bad token"})
		return nil
	}
	c.Cookie(&fiber.Cookie{
		Name:  "tkn",
		Value: tkn,
	})
	c.Status(200).JSON(fiber.Map{"status": "success"})
	return nil
}

func GithubMyRepos(c *fiber.Ctx) error {
	tkn := c.Cookies("tkn")
	data := getGithubData(tkn)
	client := github.NewClient(nil).WithAuthToken(tkn)
	opts := github.RepositoryListByUserOptions{Sort: "created", Type: "owner"}
	list, _, err := client.Repositories.ListByUser(context.Background(), data.Login, &opts)
	if err != nil {
		return c.Status(200).SendString(fmt.Sprintf("error, %e", err))
	}

	r := []string{}
	for _, l := range list {
		r = append(r, *l.Name)
	}
	return c.Status(200).JSON(r)
}

func GithubRepoFiles(c *fiber.Ctx) error {
	tkn := c.Cookies("tkn")
	repoName := c.Query("repo")

	res, err := getRepoFiles(tkn, repoName, "/")
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("error, %e", err))
	}
	return c.Status(200).JSON(res)
}

func GithubRepoFile(c *fiber.Ctx) error {
	tkn := c.Cookies("tkn")
	repoName := c.Query("repo")
	path := c.Query("path")

	res, err := getFile(tkn, repoName, path)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("error, %e", err))
	}
	return c.Status(200).SendString(res)
}

func GithubSendFile(c *fiber.Ctx) error {
	tkn := c.Cookies("tkn")
	payload := UploadStruct{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	err := sendFile(tkn, payload.Repo, payload.Content, payload.File, fmt.Sprintf("Update %s", payload.File))
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("error, %e", err))
	}
	return c.Status(200).SendString("success")
}

func getRepoFiles(tkn, repo, path string) (res []string, err error) {
	data := getGithubData(tkn)
	client := github.NewClient(nil).WithAuthToken(tkn)
	ctx := context.Background()

	_, d, resp, err := client.Repositories.GetContents(ctx, data.Login, repo, path, nil)
	if resp.StatusCode != 200 || err != nil {
		return res, errors.New("File does not exist")
	}
	for _, fd := range d {
		if fd.GetType() == "file" {
			res = append(res, fd.GetPath())
		}
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

func getFile(tkn, repo, fileName string) (string, error) {
	data := getGithubData(tkn)
	client := github.NewClient(nil).WithAuthToken(tkn)
	ctx := context.Background()

	f, _, resp, err := client.Repositories.GetContents(ctx, data.Login, repo, fileName, nil)
	if resp.StatusCode != 200 || err != nil || f.Content == nil {
		return "", errors.New("File does not exist")
	}
	res, err := f.GetContent()
	if err != nil {
		return "", err
	}
	return res, nil
}

func sendFile(tkn, repo, content, fileName, commitMessage string) error {
	data := getGithubData(tkn)
	client := github.NewClient(nil).WithAuthToken(tkn)
	ctx := context.Background()

	f, _, resp, err := client.Repositories.GetContents(ctx, data.Login, repo, fileName, nil)
	if resp.StatusCode != 200 || err != nil {
		opts := &github.RepositoryContentFileOptions{
			Message:   github.String(commitMessage),
			Content:   []byte(content),
			Committer: &github.CommitAuthor{Name: github.String(data.Login), Email: github.String(data.Email)},
		}
		commit, resp, err := client.Repositories.CreateFile(ctx, data.Login, repo, fileName, opts)
		if resp.StatusCode != 200 {
			return errors.New(fmt.Sprintf("Error creating file: %s", resp.Status))
		} else if err != nil {
			return err
		}
		fmt.Println(commit)
		return nil
	} else {
		opts := &github.RepositoryContentFileOptions{
			Message:   github.String(commitMessage),
			SHA:       f.SHA,
			Content:   []byte(content),
			Committer: &github.CommitAuthor{Name: github.String(data.Login), Email: github.String(data.Email)},
		}
		commit, resp, err := client.Repositories.CreateFile(ctx, data.Login, repo, fileName, opts)
		if resp.StatusCode != 200 {
			return errors.New(fmt.Sprintf("Error creating file: %s", resp.Status))
		} else if err != nil {
			return err
		}
		fmt.Println(commit)
		return nil
	}
}

func getGithubAccessToken(code string) string {
	d := GetGithubEnv()
	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     d.GithubClientID,
		"client_secret": d.GithubSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Response body converted to stringified JSON
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)
	fmt.Println("ghresp.Scope: ", ghresp.Scope)

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) GitHubMeStruct {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)
	res := GitHubMeStruct{}
	json.Unmarshal(respbody, &res)

	//return string(respbody)
	return res
}
