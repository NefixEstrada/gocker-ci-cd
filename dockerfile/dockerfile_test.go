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

	t.Run("Opens the file", func(t *testing.T) {
		path := "/home/nefix/Docker/docker-compose/webs/nginx-server"
		dockerfile, err := NewDockerfile(path)

		assert.Nil(t, err)
		assert.NotNil(t, dockerfile.file, "There should be a file")
	})
}

func TestGetFromTag(t *testing.T) {
	type dockerfileTests []struct{
		path string
		expected []string
	}

	t.Run("Returns the actual FROM(s) tag", func(t *testing.T) {
		tests := dockerfileTests {
			{"/home/nefix/Docker/docker-compose/webs/docsify-server/Dockerfile", []string{"alpine"}},
			{"/home/nefix/Docker/docker-compose/webs/gitea-server/Dockerfile", []string{"arm32v6/golang:1.10-alpine3.7", "alpine:3.7"}},
		}

		for _, tt := range tests {
			dockerfile, _ := NewDockerfile(tt.path)

			assert.Equal(t, tt.expected, dockerfile.GetFromTag())
		}
	})
}