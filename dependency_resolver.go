package main

import (
	"io/ioutil"
	"regexp"
	"fmt"
)

/*
 * Parse files to find their dependencies requirements, and return them
 * sorted accordingly.
 */
func ResolveDependencies( files []string, base string ) ( sortedFiles []string, err error ) {
	resolver := DependencyResolver{ initialFiles: files, Base: base }
	return resolver.Resolve()
}

type DependencyResolver struct {
	Base         string
	initialFiles []string
	sortedFiles  []string
	pendingFiles []SourceFile
}

/*
 * Actual resolve looping
 */
func ( resolver *DependencyResolver ) Resolve() ( sortedFiles []string, err error ) {
	for _, file := range resolver.initialFiles {
		source := SourceFile{ path: file }
		err = source.ParseDependencies( resolver.Base )
		if err != nil { return }

		if source.Resolved( resolver.sortedFiles ) {
			resolver.sortedFiles = append( resolver.sortedFiles, source.path )
			resolver.RemovePending( source )
			resolver.ProcessPendings()
		} else {
			resolver.pendingFiles = append( resolver.pendingFiles, source )
		}
	}

	if len( resolver.pendingFiles ) > 0 {
		for i := 0 ; i < len( resolver.pendingFiles ) ; i++ {
			resolver.ProcessPendings()
			if len( resolver.pendingFiles ) == 0 {
				break
			}
		}
	}

	if len( resolver.pendingFiles ) > 0 {
		err = fmt.Errorf( "Can't resolve dependencies in %s. Circular dependencies?", resolver.Base )
	} else {
		sortedFiles = resolver.sortedFiles
	}

	return
}

/*
 * Check if previously unresolved dependencies now are
 */
func ( resolver *DependencyResolver ) ProcessPendings() {
	for _, source := range resolver.pendingFiles {
		if source.Resolved( resolver.sortedFiles ) {
			resolver.sortedFiles = append( resolver.sortedFiles, source.path )
			resolver.RemovePending( source )
		}
	}
}

/*
 * Remove a resolved source file from pending files
 */
func ( resolver *DependencyResolver ) RemovePending( source SourceFile ) {
	newPendings := make( []SourceFile, 0 )

	for _, pending := range resolver.pendingFiles {
		if pending.path != source.path {
			newPendings = append( newPendings, pending )
		}
	}

	resolver.pendingFiles = newPendings
}

type SourceFile struct {
	path         string
	dependencies []string
}

/*
 * Read dependencies from source file
 */
func ( source *SourceFile ) ParseDependencies( base string ) ( err error ) {
	source.dependencies = make( []string, 0 )

	file, err := ioutil.ReadFile( source.path )
	dependencyFinder := regexp.MustCompile( `--\s+require\s+['"](.*)['"]` )
	dependencies := dependencyFinder.FindAllStringSubmatch( string( file ), -1 )

	for _, submatches := range dependencies {
		if len( submatches ) > 1 {
			dependency :=  base + "/" + submatches[1]
			alreadyExists := false

			for _, existing := range source.dependencies {
				if existing == dependency {
					alreadyExists = true
				}
			}

			if ! alreadyExists {
				source.dependencies = append( source.dependencies, dependency )
			}
		}
	}

	return
}

/*
 * Check if all dependencies of current file are resolved
 */
func ( source *SourceFile ) Resolved( readyFiles []string ) bool {
	for _, file := range source.dependencies {
		resolved := false

		for _, readyFile := range readyFiles {
			if readyFile == file {
				resolved = true
			}
		}

		if ! resolved {
			return false
		}
	}

	return true
}
