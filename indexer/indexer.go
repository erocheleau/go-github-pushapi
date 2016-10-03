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

// ReposIndexer interface to index github objects. @see NewReposIndexer for config
type ReposIndexer interface {
	IndexOrganizationByName(orgname string)
	IndexRepositoryByNameAsync(owner, name string)
	IndexRepositoriesByOrgAsync(orgname string)
	pushRepositories([]*github.Repository)
}

// ReposIndexerOptions options to start the indexer
type ReposIndexerOptions struct {
	Debug    bool
	SourceID string
}

// NewReposIndexer will create an indexer process to index different github objects
// it requires a pushConf which is the configuration of a pushapi client @see github.com/coveo/go-coveo/pushapi
// it can receive an httpClient to handle github authentication, and some opt which are the `SourceID` and a `Debug` flag
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

///////////////////////////////////////////////
/////////////// ORGANIZATIONS /////////////////
///////////////////////////////////////////////

func (r *reposIndexer) IndexOrganizationByName(orgname string) {

	if r.debug {
		Info.Println("Indexing", orgname)
	}

	organization, _, err := r.ghclient.Organizations.Get(orgname)
	if err != nil {
		fmt.Println(err)
	}

	r.pushOrganization(organization)
}

func (r *reposIndexer) pushOrganization(org *github.Organization) {
	toPush, err := BuildOrganizationObject(org)
	if err != nil {
		Error.Fatal(err)
	}

	doc := pushapi.Document{
		DocumentID: toPush.URL,
		Fields:     structs.Map(toPush),
	}

	if _, err = r.pushClient.PushDocument(doc, r.sourceID); err != nil {
		Error.Fatal(err)
	}

	if r.debug {
		Info.Println("Done pushing organization", *org.Name)
	}
}

///////////////////////////////////////////////
/////////////// REPOSITORIES //////////////////
///////////////////////////////////////////////

// IndexRepositoryByNameAsync will index a single repository from a owner.
func (r *reposIndexer) IndexRepositoryByNameAsync(owner, name string) {
	repos := make([]*github.Repository, 1)
	repo, _, err := r.ghclient.Repositories.Get(owner, name)
	if err != nil {
		fmt.Println(err)
	}

	if r.debug {
		Info.Println("Indexing repository ", owner, "/", name)
	}

	repos[0] = repo

	r.pushRepositories(repos)
}

// IndexRepositoriesByOrgAsync will index all the repositories from a specific
// orgname and will index them in a push source
// @param orgname The name of the organization to index the repositories of.
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
