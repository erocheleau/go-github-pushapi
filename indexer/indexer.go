package indexer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/coveo/go-coveo/pushapi"
	"github.com/fatih/structs"
	"github.com/google/go-github/github"
)

type ReposIndexer interface {
	IndexRepositoriesByOrgAsync(orgname string)
	pushRepositories([]*github.Repository)
}

type ReposIndexerOptions struct {
	Debug    bool
	SourceID string
}

func NewReposIndexer(pushConf pushapi.Config, httpclient *http.Client, opt ReposIndexerOptions) (ReposIndexer, error) {

	if opt.Debug {
		InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	ghclient := github.NewClient(httpclient)

	pushClient, err := pushapi.NewClient(pushConf)
	if err != nil {
		return nil, err
	}

	return &reposIndexer{
		ghclient:   ghclient,
		HTTPClient: httpclient,
		debug:      opt.Debug,
		pushClient: pushClient,
		sourceID:   opt.SourceID,
	}, nil
}

type reposIndexer struct {
	ghclient   *github.Client
	HTTPClient *http.Client
	debug      bool
	pushClient pushapi.Client
	sourceID   string
}

func (r *reposIndexer) IndexRepositoriesByOrgAsync(orgname string) {

	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := r.ghclient.Repositories.ListByOrg(orgname, opt)

	if r.debug {
		Info.Println("Indexing", len(repos), "repositories of", orgname)
	}

	if err != nil {
		fmt.Println(err)
	}

	r.pushRepositories(repos)
}

func (r *reposIndexer) pushRepositories(repos []*github.Repository) {
	if r.debug {
		Info.Println("Pushing", len(repos), "repositories")
	}
	for _, repo := range repos {
		if r.debug {
			Info.Println("Pushing repository", *repo.Name)
		}

		toPush, err := BuildRepositoryObject(repo)
		if err != nil {
			Error.Fatal(err)
		}

		doc := pushapi.Document{
			DocumentID: toPush.URL,
			Fields:     structs.Map(toPush),
		}
		// Pushing the repositories with the pushapi
		if _, err = r.pushClient.PushDocument(doc, r.sourceID); err != nil {
			Error.Fatal(err)
		}
	}
	if r.debug {
		Info.Println("Done indexing")
	}
}
