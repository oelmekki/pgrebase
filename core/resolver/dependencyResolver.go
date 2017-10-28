package resolver

import (
	"fmt"
)

// dependencyResolver holds info about dependencies, like their resolve order and the
// current state of resolving.
type dependencyResolver struct {
	Base         string       // the path to resolving root
	initialFiles []string     // list of found file, unordered
	sortedFiles  []string     // list of found file, sorted by resolving order
	pendingFiles []sourceFile // list of found files we're not sure yet of resolving order
}

// Resolve is the actual resolve looping.
func (resolver *dependencyResolver) Resolve() (sortedFiles []string, err error) {
	for _, file := range resolver.initialFiles {
		source := sourceFile{path: file}
		err = source.ParseDependencies(resolver.Base)
		if err != nil {
			return
		}

		if source.Resolved(resolver.sortedFiles) {
			resolver.sortedFiles = append(resolver.sortedFiles, source.path)
			resolver.RemovePending(source)
			resolver.ProcessPendings()
		} else {
			resolver.pendingFiles = append(resolver.pendingFiles, source)
		}
	}

	if len(resolver.pendingFiles) > 0 {
		for i := 0; i < len(resolver.pendingFiles); i++ {
			resolver.ProcessPendings()
			if len(resolver.pendingFiles) == 0 {
				break
			}
		}
	}

	if len(resolver.pendingFiles) > 0 {
		err = fmt.Errorf("Can't resolve dependencies in %s. Circular dependencies?", resolver.Base)
	} else {
		sortedFiles = resolver.sortedFiles
	}

	return
}

// ProcessPending checks if previously unresolved dependencies now are.
func (resolver *dependencyResolver) ProcessPendings() {
	for _, source := range resolver.pendingFiles {
		if source.Resolved(resolver.sortedFiles) {
			resolver.sortedFiles = append(resolver.sortedFiles, source.path)
			resolver.RemovePending(source)
		}
	}
}

// RemovePending removes a resolved source file from pending files.
func (resolver *dependencyResolver) RemovePending(source sourceFile) {
	newPendings := make([]sourceFile, 0)

	for _, pending := range resolver.pendingFiles {
		if pending.path != source.path {
			newPendings = append(newPendings, pending)
		}
	}

	resolver.pendingFiles = newPendings
}
