// Package dockerfile provides functions to read and interact with Dockerfiles.
package dockerfile

import (
	"bufio"
	"os"
	"strings"
	"errors"
)

// Dockerfile is the main struct for the Dockerfile package.
type Dockerfile struct {
	Path *os.File
	FromTags []string
}

// NewDockerfile creates a new Dockerfile instance using the path provided as an argument
func NewDockerfile (path string) (Dockerfile, error){
	var (
		d Dockerfile
		err error
	)

	d.Path, err = os.Open(path)

	return d, err
}

// GetFromTag returns an array of all the FROM tags content. Usually, it's just going to return a signle value
func (d *Dockerfile) GetFromTag() error {
	scanner := bufio.NewScanner(d.Path)

	for scanner.Scan() {
		splittedLine := strings.Split(scanner.Text(), " ")
		if splittedLine[0] == "FROM" {
			d.FromTags = append(d.FromTags, splittedLine[1])
		}
	}

	if len(d.FromTags) == 0 {
		return errors.New("not found: there is no FROM tag inside the Dockerfile you provided")
	}

	return nil
}
