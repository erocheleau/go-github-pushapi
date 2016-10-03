package indexer

import (
	"encoding/json"

	"github.com/google/go-github/github"
)

// BuildOrganizationObject builds an organization object ready to be pushed to Coveo
func BuildOrganizationObject(toBuild *github.Organization) (*Organization, error) {
	marsh, err := json.Marshal(toBuild)
	if err != nil {
		return nil, err
	}

	toPush := &Organization{}
	if err := json.Unmarshal(marsh, toPush); err != nil {
		return nil, err
	}

	toPush.Data = toPush.Description
	toPush.FileType = toPush.Type

	return toPush, nil
}

// BuildRepositoryObject builds a Repository object ready to be pushed to Coveo
func BuildRepositoryObject(toBuild *github.Repository) (*Repository, error) {
	marsh, err := json.Marshal(toBuild)
	if err != nil {
		return nil, err
	}

	toPush := &Repository{}
	if err := json.Unmarshal(marsh, toPush); err != nil {
		return nil, err
	}

	toPush.OwnerID = toPush.Owner.ID
	toPush.OwnerName = toPush.Owner.Login
	toPush.OwnerType = toPush.Owner.Type
	toPush.OwnerAvatarURL = toPush.Owner.AvatarURL
	toPush.Data = toPush.Description
	toPush.FileType = "Repository"

	return toPush, nil
}

// BuildRepositoryObject builds a Repository object ready to be pushed to Coveo
func BuildPRObject(toBuild *github.PullRequest) (*PullRequest, error) {
	marsh, err := json.Marshal(toBuild)
	if err != nil {
		return nil, err
	}

	toPush := &PullRequest{}
	if err := json.Unmarshal(marsh, toPush); err != nil {
		return nil, err
	}

	toPush.Data = toPush.Title
	toPush.FileType = "Pull Request"

	return toPush, nil
}
