package main

import (
	"fmt"
	"os"

	"github.com/coveo/go-coveo/pushapi"
	"github.com/erocheleau/go-github-pushapi/indexer"
)

func main() {

	sourceID := os.Getenv("SOURCEID")
	organizationID := os.Getenv("ORGANIZATIONID")
	APIKey := os.Getenv("APIKEY")

	if len(sourceID) == 0 || len(organizationID) == 0 || len(APIKey) == 0 {
		fmt.Println("SOURCEID, ORGANIZATIONID, APIKEY must be defined")
	}

	indexerOptions := indexer.ReposIndexerOptions{
		Debug:    true,
		SourceID: sourceID,
	}
	pushConf := pushapi.Config{
		OrganizationID: organizationID,
		APIKey:         APIKey,
	}
	indexer, err := indexer.NewReposIndexer(pushConf, nil, indexerOptions)
	if err != nil {
		panic(err)
	}
	indexer.IndexRepositoriesByOrgAsync("Coveo")
}

// func readGithub(org string) error {
// 	c := make(chan *github.Repository)
// 	quit := make(chan int)

// 	client := github.NewClient(nil)

// 	opt := &github.RepositoryListByOrgOptions{Type: "public"}
// 	fmt.Println("Sending request to github ListRepositories by Org")
// 	repos, _, err := client.Repositories.ListByOrg(org, opt)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Received response from github ListRepositories by Org: ", len(repos))
// 	go func() {
// 		for _, repo := range repos {
// 			fmt.Println("Looping through the repos async ", *repo.Name)
// 			c <- repo
// 		}
// 		quit <- 0
// 	}()

// 	printGithub(c, quit)
// 	return nil
// }

// func printGithub(c <-chan *github.Repository, quit chan int) {
// 	for {
// 		select {
// 		case repo := <-c:
// 			fmt.Println("Repo: ", *repo.Name)
// 		case <-quit:
// 			fmt.Println("Done")
// 			return
// 		}
// 	}
// }
