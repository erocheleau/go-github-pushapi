package indexer

import (
	"encoding/json"

	"github.com/google/go-github/github"
)

func BuildRepositoryObject(ghrepo *github.Repository) (*Repository, error) {
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
	toPush.FileType = "Repository"

	return toPush, nil
}
