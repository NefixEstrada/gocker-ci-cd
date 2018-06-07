package settings

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestReadSettings(t *testing.T) {
	t.Run("If there's no config, write the default one and parse it", func(t *testing.T) {
		err := os.RemoveAll("/etc/gocker-ci-cd/")
		assert.NoError(t, err)

		var settings Settings
		settings, err = ReadSettings()
		assert.NoError(t, err)

		assert.Equal(t, defaultSettings, settings)
	})
}

func TestWriteSettings(t *testing.T) {
	t.Run("Create file and directory if they don't exist", func(t *testing.T) {
		err := os.RemoveAll("/etc/gocker-ci-cd/")
		assert.NoError(t, err)

		err = WriteSettings(defaultSettings)
		assert.NoError(t, err)

		assert.DirExists(t, "/etc/gocker-ci-cd")
		assert.FileExists(t, "/etc/gocker-ci-cd/config.toml")
	})

	t.Run("Update settings", func(t *testing.T) {
		settings := defaultSettings
		settings.User = "nefix"
		err := WriteSettings(settings)

		assert.NoError(t, err)

		var writtenSettings Settings
		writtenSettings, err = ReadSettings()

		assert.NoError(t, err)
		assert.Equal(t, settings, writtenSettings)
	})
}
