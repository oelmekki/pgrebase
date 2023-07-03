package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type sourceFile struct {
	path         string
	dependencies []string
}

// parseDependencies reads dependencies from source file.
func (source *sourceFile) parseDependencies(base string) (err error) {
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

// resolved checks if all dependencies of current file are resolved
func (source *sourceFile) resolved(readyFiles []string) bool {
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

// dependencyResolver holds info about dependencies, like their resolve order and the
// current state of resolving.
type dependencyResolver struct {
	base         string       // the path to resolving root
	initialFiles []string     // list of found file, unordered
	sortedFiles  []string     // list of found file, sorted by resolving order
	pendingFiles []sourceFile // list of found files we're not sure yet of resolving order
}

// resolve is the actual resolve looping.
func (resolver *dependencyResolver) resolve() (sortedFiles []string, err error) {
	for _, file := range resolver.initialFiles {
		source := sourceFile{path: file}
		err = source.parseDependencies(resolver.base)
		if err != nil {
			return
		}

		if source.resolved(resolver.sortedFiles) {
			resolver.sortedFiles = append(resolver.sortedFiles, source.path)
			resolver.removePending(source)
			resolver.processPendings()
		} else {
			resolver.pendingFiles = append(resolver.pendingFiles, source)
		}
	}

	if len(resolver.pendingFiles) > 0 {
		for i := 0; i < len(resolver.pendingFiles); i++ {
			resolver.processPendings()
			if len(resolver.pendingFiles) == 0 {
				break
			}
		}
	}

	if len(resolver.pendingFiles) > 0 {
		err = fmt.Errorf("Can't resolve dependencies in %s. Circular dependencies?", resolver.base)
	} else {
		sortedFiles = resolver.sortedFiles
	}

	return
}

// processPendings checks if previously unresolved dependencies now are.
func (resolver *dependencyResolver) processPendings() {
	for _, source := range resolver.pendingFiles {
		if source.resolved(resolver.sortedFiles) {
			resolver.sortedFiles = append(resolver.sortedFiles, source.path)
			resolver.removePending(source)
		}
	}
}

// removePending removes a resolved source file from pending files.
func (resolver *dependencyResolver) removePending(source sourceFile) {
	newPendings := make([]sourceFile, 0)

	for _, pending := range resolver.pendingFiles {
		if pending.path != source.path {
			newPendings = append(newPendings, pending)
		}
	}

	resolver.pendingFiles = newPendings
}

// resolveDependencies parses files to find their dependencies requirements, and return them.
// sorted accordingly.
func resolveDependencies(files []string, base string) (sortedFiles []string, err error) {
	resolver := dependencyResolver{initialFiles: files, base: base}
	return resolver.resolve()
}
