package resolver

// DependencyResolver parses files to find their dependencies requirements, and return them.
// sorted accordingly.
func ResolveDependencies(files []string, base string) (sortedFiles []string, err error) {
	resolver := dependencyResolver{initialFiles: files, Base: base}
	return resolver.Resolve()
}
