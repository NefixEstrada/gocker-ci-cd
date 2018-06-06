// Package dockerhub provides the ability to interact with
package dockerhub

import (
	"net/http"
	"fmt"
	"errors"
	"encoding/json"
	"time"
	"reflect"
)

const Protocol = "https"
const DomainName = "hub.docker.com"
const ApiVersion = "v2"
const BaseUrl = "repositories"

// Repo is the struct for the repositories in DockerHub
type Repo struct {
	User string
	Name string
	Namespace string
	LastUpdated time.Time
	Tags RepoTags
}

// repoJSONResponse is the struct of the response of the Docker Hub API when querying a repo
type repoJSONResponse struct {
	User string `json:"User"`
	Name string `json:"Name"`
	Namespace string `json:"Namespace"`
	RepositoryType string `json:"repository_type"`
	Status int `json:"status"`
	Description string `json:"description"`
	IsPrivate bool `json:"is_private"`
	IsAutomated bool `json:"is_automated"`
	CanEdit bool `json:"can_edit"`
	StarCount int `json:"star_count"`
	PullCount int `json:"pull_count"`
	LastUpdated time.Time `json:"last_updated"`
	HasStarred bool `json:"has_starred"`
	FullDescription string `json:"full_description"`
	Affiliation interface{} `json:"affiliation"`
	Permissions struct {
		Read bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
}

// GetRepo queries the Docker Hub API searching for the repository specified in the parameters
func GetRepo(username, name string) (r Repo, err error) {
	if username == "_" {
		username = "library"
	}

	url := fmt.Sprintf("%s://%s/%s/%s/%s/%s", Protocol, DomainName, ApiVersion, BaseUrl, username, name)
	resp, err := http.Get(url)

	if err != nil {
		return
	}

	statusCode := resp.StatusCode
	if statusCode == 404 {
		err = errors.New("not found: the repository wasn't found")
		return
	}

	var repoJSON repoJSONResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&repoJSON)

	if err != nil {
		return
	}

	tags, err := GetRepoTags(repoJSON.User, repoJSON.Name)

	if err != nil {
		return
	}

	r = Repo{repoJSON.User, repoJSON.Name, repoJSON.Namespace, repoJSON.LastUpdated, tags}

	return
}

// RepoTags is an slice of Tags. It's what is going to be returned in the functions
type RepoTags []Tag

// Tag is the struct for the image tags in the Docker Registry
type Tag struct {
	Name string
	ID int
	LastUpdated time.Time
	Images RepoTagImages
}

// RepoTagImages is an slice of TagImages
type RepoTagImages []TagImage

// TagImage is the struct for the images inside the RepoTags
type TagImage struct {
	Os string
	Architecture string
	Size int
}

// repoTagsJSONResponse is the struct of the response of the Docker Hub API when querying the tags of a repo
type repoTagsJSONResponse struct {
	Count int `json:"count"`
	Next interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results []struct {
		Name string `json:"Name"`
		FullSize int `json:"full_size"`
		Images []struct {
			Size int `json:"size"`
			Architecture string `json:"architecture"`
			Variant interface{} `json:"variant"`
			Features interface{} `json:"features"`
			Os string `json:"os"`
			OsVersion interface{} `json:"os_version"`
			OsFeatures interface{} `json:"os_features"`
		} `json:"images"`
		ID int `json:"id"`
		Repository int `json:"repository"`
		Creator int `json:"creator"`
		LastUpdater int `json:"last_updater"`
		LastUpdated time.Time `json:"last_updated"`
		ImageID interface{} `json:"image_id"`
		V2 bool `json:"v2"`
	} `json:"results"`
}

// GetRepoTags queries the Docker Registry searching for the tags of the image
func GetRepoTags(username, name string) (rT RepoTags, err error) {
	url := fmt.Sprintf("%s://%s/%s/%s/%s/%s/tags?page_size=100", Protocol, DomainName, ApiVersion, BaseUrl, username, name)
	resp, err := http.Get(url)

	if err != nil {
		return
	}

	statusCode := resp.StatusCode
	if statusCode == 404 {
		err = errors.New("not found: the repo wasn't found")
		return
	}

	var repoTagsJSON repoTagsJSONResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&repoTagsJSON)

	if err != nil {
		return
	}

	for _, tag := range repoTagsJSON.Results {
		var tagImages RepoTagImages
		for _, image := range tag.Images {
			tagImages = append(tagImages, TagImage{image.Os, image.Architecture, image.Size})
		}

		rT = append(rT, Tag{tag.Name, tag.ID, tag.LastUpdated, tagImages})
	}

	return
}

// GetRegistryTag queries the Docker Registry searching for the tag in the image
func (r *Repo) GetTag(repo Repo, name string) (t Tag, err error) {
	for _, tag := range r.Tags {
		if tag.Name == name {
			t = tag
		}
	}

	if reflect.DeepEqual(t, Tag{}) {
		err = errors.New("not found: the tag wasn't found in the repo")
		return
	}

	return
}
