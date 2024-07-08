package api

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

type TknStruct struct {
	tkn         string `cookie:"tkn"`
	githubState string `cookie:"github_state"`
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
		"http://localhost:3000/callback",
		state,
	)
	return c.Redirect(link, http.StatusTemporaryRedirect)
}

func GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	fmt.Println(c.Query("code"))
	fmt.Println(c.Query("state"))
	if c.Query("state") != c.Cookies("github_state") {
		return c.Status(400).SendString("error - bad state")
	}

	t := TknStruct{}
	if err := c.CookieParser(&t); err != nil {
		return err
	}
	fmt.Println(t)

	tkn := getGithubAccessToken(code)
	if tkn == "" {
		return c.Status(400).SendString("error")
	}
	c.Cookie(&fiber.Cookie{
		Name:  "tkn",
		Value: tkn,
	})
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

func getGithubData(accessToken string) string {
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

	// Convert byte slice to string and return
	return string(respbody)
}

/*
Get user repos

curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/users/USERNAME/repos


Get repo files

Post repo changes
*/
