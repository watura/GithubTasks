package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v25/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// loadEnv load .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

// GithubToken get token
func GithubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

// GithubClient Personal Access token をセットした github.Client を返す
func GithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GithubToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func main() {
	ctx := context.Background()

	loadEnv()
	issues, _, err := GithubClient().Issues.ListByRepo(ctx, "watura", "GithubTasks", nil)
	if err != nil {
		panic(err)
	}

	// labels, _, err := GithubClient().Issues.ListLabels(ctx, "watura", "GithubTasks", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, l := range labels {
	// 	fmt.Println(l.GetName())
	// }

	labels := [7]string{
		"DUE TODAY DO IT NOW",
		"@Day: 1",
		"@Day: 2",
		"@Day: 3",
		"@Day: 4",
		"@Day: 5",
		"@Day: 6",
	}

	r := regexp.MustCompile("^Due: ")
	loc, err1 := time.LoadLocation("Asia/Tokyo")
	if err1 != nil {
		panic(err1)
	}
	df := "Due: 2006/01/02"
	today := time.Now()

	for _, issue := range issues {
		bodies := strings.Split(issue.GetBody(), "\n")
		due := ""
		for _, text := range bodies {
			if r.MatchString(text) {
				due = text
				break
			}
		}
		if due == "" {
			continue
		}

		for _, l := range issue.Labels {
			for _, label := range labels {
				if l.GetName() == label {
					_, err = GithubClient().Issues.RemoveLabelForIssue(ctx, "watura", "GithubTasks", *issue.Number, label)
					if err != nil {
						panic(err)
					}
				}
			}
		}
		date, err := time.ParseInLocation(df, due, loc)
		if err != nil {
			continue
		}
		d := int(date.Sub(today).Hours() / 24)
		fmt.Println(date, labels[d])
		_, _, err = GithubClient().Issues.AddLabelsToIssue(ctx, "watura", "GithubTasks", *issue.Number, []string{labels[d]})
		if err != nil {
			panic(err)
		}
	}
}
