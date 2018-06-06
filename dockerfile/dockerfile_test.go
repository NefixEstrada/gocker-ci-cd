package dockerfile

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestNewDockerfile(t *testing.T) {
	t.Run("OS Error", func(t *testing.T) {
		_, err := NewDockerfile("/unexisting/path")

		assert.Error(t, err, "The function needs to return an error if there's one")
	})

	t.Run("Opens the Path", func(t *testing.T) {
		path := "/home/nefix/Docker/docker-compose/webs/nginx-server"
		dockerfile, err := NewDockerfile(path)

		assert.Nil(t, err)
		assert.NotNil(t, dockerfile.Path, "There should be a Path")
	})
}

func TestGetFromTag(t *testing.T) {
	type dockerfileTests []struct{
		path string
		expected []string
	}
	
	t.Run("No FROM tag error", func(t *testing.T) {
		path := "/home/nefix/Docker/docker-compose/testing/EmptyDockerfile"
		dockerfile, err := NewDockerfile(path)
		assert.NoError(t, err)

		err = dockerfile.GetFromTag()
		assert.Error(t, err)
	})

	t.Run("Returns the actual FROM(s) tag", func(t *testing.T) {
		tests := dockerfileTests {
			{"/home/nefix/Docker/docker-compose/webs/docsify-server/Dockerfile", []string{"alpine"}},
			{"/home/nefix/Docker/docker-compose/webs/gitea-server/Dockerfile", []string{"arm32v6/golang:1.10-alpine3.7", "alpine:3.7"}},
		}

		for _, tt := range tests {
			dockerfile, err := NewDockerfile(tt.path)
			assert.NoError(t, err)

			err = dockerfile.GetFromTag()
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, dockerfile.FromTags)
		}
	})
}