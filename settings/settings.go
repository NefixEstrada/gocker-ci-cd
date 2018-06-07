package settings

import (
	"github.com/NefixEstrada/gocker-ci-cd/dockerfile"
	"os"
	"github.com/BurntSushi/toml"
	"bytes"
	"io/ioutil"
)

// Settings contains all the settings of the project.
type Settings struct {

	// Updateinterval contains the interval of checking if there are new updates (in minutes)
	UpdateInterval int

	// User is the username that the program is going to use to set the local tags
	User string

	// Email contains all the settings related with the email delivery
	Email EmailSettings

	// Applications contains all the images / projects / containers / call it what you want that are going to be tracked
	Applications []App
}

// EmailSettings contains all the information related with the email sending
type EmailSettings struct {
	SenderAddress string
	SenderPassword string
	SenderServer string
	SenderPort int
	RecieverAddress string
	EmailHeaderPrefix string
}

// App is the struct of all the settings related with the Apps to track
type App struct {

	// Name os the name that the container is going to have.
	Name string

	// Params is a string that contains all the params that go before the image name,
	// such as volumes, ports, environment variables...
	// Fields that can't be provided:
	//   - name
	//   - detached mode
	Params string

	// Cmd is the command that is going to be executed in the container. It's what goes after the image name. Normally,
	// you aren't going to use this since you normally use the entrypoint for that propose
	Cmd string

	// Dockerfile is the actual dockerfile of the App
	Dockerfile dockerfile.Dockerfile

	// Tag is the tag name that the program is going to put on the new build
	Tag string

	// Notify sets if the program is going to send an email when the App gets updated
	// or if there's a new update available
	Notify bool

	// AutoBuild sets if the program is going to automatically build and tag the new image
	AutoBulid bool

	// AutoRestart sets if the program is going to restart the container automatically
	AutoRestart bool

	// Compose contains all the settings related with the Docker Compose (if the App is part of)
	Compose ComposeSettings

	// Git contains all the settings related with the Git repostiory (if the Dockerfile is part of)
	Git GitSettings

}

// ComposeSettings are all those settings related with a possible Docker Compose configuration
type ComposeSettings struct {
	// If the App forms part of a Docker Compose, the string is going to contain the path of the directory containing the docker-compose.yml file
	// If not, the string is going to be empty
	Path string

	// AutoRestart sets if the program is going to restart automatically the whole Docker Compose
	AutoRestart bool

	// AutoUpdateFile sets if the program is going to edit the Docker Compose file in order to update the image tag
	AutoUpdateFile bool
}

// GitSettings are those settings related with a possible Git repository
type GitSettings struct {
	// If the Dockerfile is part of a Git repository, the program is going to make a pull before rebuilding the image
	// If there are issues, it's going to send a notification
	// The string is going to contain the path of the repository
	// If the Dockerfile doesn't form part of a Git repository, the string is going to be empty
	Path string

	// Notify sets if the program is going to send an email when the repo gets updated
	// or if there's a new pull avaliable
	Notify bool

	// AutoPull sets if the program is going to automatically pull the changes from the origin
	AutoPull bool

	// TODO: Add more fields:
	// 	- Stick with a version
	//  - Branch
	// 	- ...

}

// defaultSettings contains the default settings for the program
var defaultSettings = Settings {

	// Check once per day
	UpdateInterval: 1440,

	User: "username",

	Email: EmailSettings{
		SenderAddress: "username@mailprovider.com",
		SenderPassword: "p4$$w0rd",
		SenderServer: "smtp.mailprovider.com",
		SenderPort: 587,
		RecieverAddress:  "reciever@mailprovider.com",
		EmailHeaderPrefix: "[GOCKER CI/CD]",
	},


	Applications: []App{},

}


// WriteSettings writes to the configuration file the settings provided as a parameter
func WriteSettings(s Settings) error {

	// Check if the configuration directory exists and if it doesn't it creates it
	_, err := os.Stat("/etc/gocker-ci-cd")

	if os.IsNotExist(err) {
		err = os.MkdirAll("/etc/gocker-ci-cd", os.ModePerm)

		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(s)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/etc/gocker-ci-cd/config.toml", buf.Bytes(), 0644)

	return err
}

// ReadSettings reads the configuration file and returns them
func ReadSettings() (Settings, error) {
	
	var s Settings

	// Check if the configuration file exists and if it doesn't, write the default settings
	_, err := os.Stat("/etc/gocker-ci-cd/config.toml")

	if os.IsNotExist(err) {
		err = WriteSettings(defaultSettings)

		if err != nil {
			return s, err
		}
	}

	if err != nil {
		return s, err
	}

	_, err = toml.DecodeFile("/etc/gocker-ci-cd/config.toml", &s)

	return s, err
}
