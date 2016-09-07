package indexer

import (
	"encoding/json"
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
		Info.Printf("Indexing %v repositories of %s", len(repos), orgname)
	}

	if err != nil {
		fmt.Println(err)
	}

	r.pushRepositories(repos)
}

func (r *reposIndexer) pushRepositories(repos []*github.Repository) {
	if r.debug {
		Info.Printf("Pushing %v repositories", len(repos))
	}
	for _, repo := range repos {
		if r.debug {
			Info.Printf("Indexing repository %s", *repo.Name)
		}

		toPush, err := buildRepositoryObject(repo)
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

		if r.debug {
			Info.Printf("Pushed repository %v", *repo.Name)
		}
	}
	if r.debug {
		Info.Printf("Done indexing")
	}
}

func buildRepositoryObject(ghrepo *github.Repository) (*Repository, error) {
	marshalledRepo, err := json.Marshal(ghrepo)
	if err != nil {
		return nil, err
	}

	toPush := &Repository{}
	if err := json.Unmarshal(marshalledRepo, toPush); err != nil {
		return nil, err
	}

	toPush.OwnerID = toPush.Owner.ID
	toPush.OwnerName = toPush.Owner.Login
	toPush.OwnerType = toPush.Owner.Type
	toPush.OwnerAvatarURL = toPush.Owner.AvatarURL
	toPush.Data = toPush.Description

	return toPush, nil
}
