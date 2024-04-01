package zeus_v1_ai

import (
	"fmt"
	"net/url"
)

func isValidURL(inputURL string) (*url.URL, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("URL must be http or https, got: %s", u.Scheme)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("URL must have a host")
	}
	return u, nil
}
func convertToHTTPS(inputURL string) (string, error) {
	u, err := isValidURL(inputURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "http" {
		u.Scheme = "https"
	}
	return u.String(), nil
}
