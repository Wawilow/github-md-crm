package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
	File     string `json:"file"`
	FileName string `json:"file_name"`
}

type GitHubMeStrict struct {
	Login string `json:"login"`
	Id    uint   `json:"id"`
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

	d := GetGithubEnv()
	link := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&response_type=code&scope=repo&redirect_uri=%s&state=%s",
		d.GithubClientID,
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

func GithubSendFile(c *fiber.Ctx) error {
	tkn := c.Cookies("tkn")
	payload := UploadStruct{}
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println(string(c.Body()), err)
		return err
	}
	data := getGithubData(tkn)
	client := github.NewClient(nil).WithAuthToken(tkn)

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("This is my commit message"),
		Content:   []byte(payload.File),
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	_, _, err := client.Repositories.CreateFile(context.Background(), data.Login, payload.Repo, payload.FileName, opts)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("error, %e", err))
	}
	return c.Status(200).SendString("success")
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

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) GitHubMeStrict {
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
	res := GitHubMeStrict{}
	json.Unmarshal(respbody, &res)

	//return string(respbody)
	return res
}
