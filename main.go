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

	indexer.IndexOrganizationByName("Coveo")
	indexer.IndexRepositoriesByOrgAsync("Coveo")
}
