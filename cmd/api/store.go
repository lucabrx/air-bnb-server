package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func (app *application) githubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     app.config.GithubClientId,
		ClientSecret: app.config.GithubClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/v1/auth/github/callback",
		Scopes:       []string{"user:email"},
	}
}

func (app *application) googleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     app.config.GoogleClientId,
		ClientSecret: app.config.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/v1/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}
}
