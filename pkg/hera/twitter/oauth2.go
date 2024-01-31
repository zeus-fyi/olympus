package hera_twitter

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Replace these constants with your own data
const (
	RedirectURI      = "http://localhost:9000/oauth2/callback"
	AuthorizationURL = "https://twitter.com/i/oauth2/authorize"
	TokenURL1        = "https://api.twitter.com/2/oauth2/token"
)

var (
	scopes = []string{"bookmark.write", "bookmark.read", "tweet.read", "users.read", "offline.access", "follows.read"}
	state  = "your-random-state"
)

func addToken(consumerKey, consumerSecret string) string {
	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(consumerKey, consumerSecret).
		SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=utf-8").
		SetBody("grant_type=client_credentials").
		Post("https://api.twitter.com/oauth2/token")

	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		panic(fmt.Sprintf("Failed to get bearer token, status code %d", resp.StatusCode()))
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		panic(err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		panic(err)
	}
	return accessToken
}

func createCodeVerifierAndChallenge() (string, string) {
	verifierBytes := make([]byte, 30)
	_, err := rand.Read(verifierBytes)
	if err != nil {
		panic(err)
	}
	codeVerifier := base64.URLEncoding.EncodeToString(verifierBytes)
	codeVerifier = regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(codeVerifier, "")

	s256 := sha256.New()
	s256.Write([]byte(codeVerifier))
	codeChallenge := base64.URLEncoding.EncodeToString(s256.Sum(nil))
	codeChallenge = strings.TrimRight(codeChallenge, "=")

	return codeVerifier, codeChallenge
}

func getBookmarks(client *resty.Client, userID, accessToken string) error {
	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/bookmarks", userID)
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request returned an error: %d %s", resp.StatusCode(), resp.String())
	}

	fmt.Printf("Response code: %d\n", resp.StatusCode())
	fmt.Println(resp.String())

	return nil
}

func fetchUserID(client *resty.Client) (string, error) {
	resp, err := client.R().
		Get("https://api.twitter.com/2/users/me")

	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	userData, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data format")
	}

	userID, ok := userData["id"].(string)
	if !ok {
		return "", fmt.Errorf("user ID is not a string")
	}

	return userID, nil
}
