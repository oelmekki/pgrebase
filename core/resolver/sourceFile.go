package resolver

import (
	"io/ioutil"
	"regexp"
)

type sourceFile struct {
	path         string
	dependencies []string
}

// ParseDependencies reads dependencies from source file.
func (source *sourceFile) ParseDependencies(base string) (err error) {
	source.dependencies = make([]string, 0)

	file, err := ioutil.ReadFile(source.path)
	dependencyFinder := regexp.MustCompile(`--\s+require\s+['"](.*)['"]`)
	dependencies := dependencyFinder.FindAllStringSubmatch(string(file), -1)

	for _, submatches := range dependencies {
		if len(submatches) > 1 {
			dependency := base + "/" + submatches[1]
			alreadyExists := false

			for _, existing := range source.dependencies {
				if existing == dependency {
					alreadyExists = true
				}
			}

			if !alreadyExists {
				source.dependencies = append(source.dependencies, dependency)
			}
		}
	}

	return
}

// Resolved checks if all dependencies of current file are resolved
func (source *sourceFile) Resolved(readyFiles []string) bool {
	for _, file := range source.dependencies {
		resolved := false

		for _, readyFile := range readyFiles {
			if readyFile == file {
				resolved = true
			}
		}

		if !resolved {
			return false
		}
	}

	return true
}
