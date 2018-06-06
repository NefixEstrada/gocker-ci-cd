package dockerhub

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestGetRepo(t *testing.T) {
	t.Run("Not found", func(t *testing.T) {
		_, err := GetRepo("unexisting-User", "and-repo")

		assert.Error(t, err, "The function needs to return an error if the repo wasn't found")

	})

	t.Run("Returns the repo from the Docker Hub", func(t *testing.T) {

		lastUpdated, _ := time.Parse(time.RFC3339, "2018-04-09T06:42:36.910134Z")

		tags, _ := GetRepoTags("codefresh", "kubectl")

		expected := Repo{
			User:        "codefresh",
			Name:        "kubectl",
			Namespace:   "codefresh",
			LastUpdated: lastUpdated,
			Tags: tags,
		}
		repo, _ := GetRepo("codefresh", "kubectl")

		assert.Equal(t, expected, repo)
	})
}

func TestGetRepoTags(t *testing.T) {
	t.Run("Not found", func(t *testing.T) {
		_, err := GetRepoTags("unexisting-User", "and-repo")

		assert.Error(t, err, "The function needs to return an error if the tag wasn't found")
	})

	t.Run("Returns the tags of a repo", func(t *testing.T) {

		beforeLastUpdated, _ := time.Parse(time.RFC3339, "2017-01-11T01:55:56.729502Z")
		beforeTagImages := RepoTagImages{{"linux", "amd64", 28956873}}

		afterLastUpdated, _ := time.Parse(time.RFC3339, "2017-01-11T01:52:18.752138Z")
		afterTagImages := RepoTagImages{{"linux", "amd64", 28956875}}

		latestLastUpdated, _ := time.Parse(time.RFC3339, "2017-01-04T23:46:28.186635Z")
		latestTagImages := RepoTagImages{{"linux", "amd64", 28929249}}


		expected := RepoTags{
			{"before", 7327046, beforeLastUpdated, beforeTagImages},
			{"after", 7326984, afterLastUpdated, afterTagImages},
			{"latest", 7201078, latestLastUpdated, latestTagImages},
		}


		repoTags, _ := GetRepoTags("dockersamples", "examplevotingapp_vote")

		assert.Equal(t, expected, repoTags)
	})
}

func TestGetTag(t *testing.T) {
	t.Run("Not found", func(t *testing.T) {
		repo, _ := GetRepo("codefresh", "kubectl")
		_, err := repo.GetTag(repo, "invalid.tag")

		assert.Error(t, err, "The function needs to return an error if the tag wasn't found")
	})

	t.Run("Returns a specific tag from a given repo", func(t *testing.T) {
		repo, _ := GetRepo("codefresh", "kubectl")
		tag, _ := repo.GetTag(repo, "1.10.0")

		beforeLastUpdated, _ := time.Parse(time.RFC3339, "2018-04-08T10:39:18.901693Z")
		expected := Tag{
			"1.10.0",
			25242207,
			beforeLastUpdated,
			RepoTagImages{
				{"linux", "amd64", 20398206},
			},
		}

		assert.Equal(t, expected, tag)
	})
}
