package main

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// GithubAuthTokens are for making requests to GiHub's API.
var GithubAuthTokens = []string{"GITHUB_AUTH_TOKEN", "GITHUB_TOKEN", "GH_TOKEN", "GH_AUTH_TOKEN"}

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("requires two arguments owner and repository. Example ossf scorecard")
	}
	owner := args[1]
	repo := args[2]
	ctx := context.Background()
	token, result := readGitHubTokens()
	if !result {
		log.Fatal("GITHUB_TOKEN ENV variable has to be set with PAT")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	wf, _, e := client.Actions.ListWorkflows(ctx, owner, repo, nil)
	if e != nil {
		log.Fatal(e)
	}
	for _, w := range wf.Workflows {
		u, _, e := client.Actions.ListWorkflowRunsByID(ctx, owner, repo, w.GetID(), nil)
		if e != nil {
			log.Fatal(e)
		}
		if u.GetTotalCount() > 0 {
			filename := fmt.Sprintf("run-log-%d-%d.zip", w.GetID(), w.CreatedAt.Unix())
			filepath := filepath.Join(os.TempDir(), "cache-foo", filename)
			fmt.Println(w.GetName(), w.GetPath())
			r, e := getLog(tc, u.WorkflowRuns[0].GetLogsURL())
			if e != nil {
				log.Println("Unable to fetch logs", e)
				continue
			}
			e = Create(filepath, r)
			if e != nil {
				log.Println(e)
			}
			z, e := Open(filepath)
			if e != nil {
				log.Println(e)
			}
			for _, f := range z.File {
				if strings.Contains(f.Name, "1_Set") {
					logfile, err := f.Open()
					if err != nil {
						log.Println(err)
					}
					scanner := bufio.NewScanner(logfile)
					shouldPrint := false
					for scanner.Scan() {
						if strings.Contains(scanner.Text(), "##[group]GITHUB_TOKEN") {
							shouldPrint = true
							fmt.Println("permissions:")
							continue
						} else if shouldPrint && strings.Contains(scanner.Text(), "##[endgroup]") {
							break
						}
						if shouldPrint {
							data := strings.Split(scanner.Text(), " ")
							fmt.Println("  ", strings.ToLower(data[1]), data[2])
						}
					}
					break
				}
			}
		}
	}
}

func readGitHubTokens() (string, bool) {
	for _, name := range GithubAuthTokens {
		if token, exists := os.LookupEnv(name); exists && token != "" {
			return token, exists
		}
	}
	return "", false
}

func Create(path string, content io.ReadCloser) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("could not create cache: %w", err)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, content)
	return err
}

func getLog(httpClient *http.Client, logURL string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", logURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("log not found")
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("not a 200 status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func Open(path string) (*zip.ReadCloser, error) {
	return zip.OpenReader(path)
}
