package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func getTopThreeRepos(client *github.Client,
	ctx context.Context,
	user string,
	repos []*github.Repository) [3]*github.Repository {
	repoViews := make([]int, len(repos))

	for index, repo := range repos {
		views, _, _ := client.Repositories.ListTrafficViews(ctx, user, repo.GetName(), nil)

		repoViews[index] = views.GetCount()
	}

	var topThreeRepos [3]*github.Repository

	for i := 0; i < 3; i++ {
		var topIndex = 0
		for j := 0; j < len(repoViews); j++ {
			if repoViews[j] > repoViews[topIndex] {
				topIndex = j
			}
		}
		topThreeRepos[i] = repos[topIndex]
		repoViews[topIndex] = 0
	}

	return topThreeRepos
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("No .env file found.")
	}

	user := os.Getenv("GITHUB_USER")

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	auth := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(auth)

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, _, err := client.Repositories.List(ctx, user, opt)
	if err != nil {
		return
	}
	repoViews := make([]int, len(repos))
	for index, repo := range repos {
		views, _, err := client.Repositories.ListTrafficViews(ctx, user, repo.GetName(), nil)
		if err != nil {
			return
		}
		repoViews[index] = views.GetCount()
	}

	topThree := getTopThreeRepos(client, ctx, user, repos)

	for _, repo := range topThree {
		views, _, _ := client.Repositories.ListTrafficViews(ctx, user, repo.GetName(), nil)
		var tempBuff bytes.Buffer
		fmt.Fprintf(&tempBuff, "%s with %d total views and %d unqiue views", repo.GetName(), views.GetCount(), views.GetUniques())
		fmt.Println(&tempBuff)

	}
}
