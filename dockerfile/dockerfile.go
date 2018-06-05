// Package dockerfile provides functions to read and interact with Dockerfiles.
package dockerfile

import (
	"bufio"
	"os"
	"strings"
)

// Dockerfile is the main struct for the Dockerfile package.
type Dockerfile struct {
	file *os.File
}

// NewDockerfile creates a new Dockerfile instance using the path provided as an argument
func NewDockerfile (path string) (d Dockerfile, err error){
	d.file, err = os.Open(path)

	return
}

// GetFromTag returns an array of all the FROM tags content. Usually, it's just going to return a signle value
func (d Dockerfile) GetFromTag() (fromTags []string) {
	scanner := bufio.NewScanner(d.file)

	for scanner.Scan() {
		splittedLine := strings.Split(scanner.Text(), " ")
		if splittedLine[0] == "FROM" {
			fromTags = append(fromTags, splittedLine[1])
		}
	}

	return
}
