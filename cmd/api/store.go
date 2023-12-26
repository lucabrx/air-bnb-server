package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func (app *application) githubConfig() *oauth2.Config {
	githubOauthConfig := &oauth2.Config{
		ClientID:     app.config.GithubClientId,
		ClientSecret: app.config.GithubClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/v1/auth/github/callback",
		Scopes:       []string{"user:email"},
	}

	return githubOauthConfig
}
