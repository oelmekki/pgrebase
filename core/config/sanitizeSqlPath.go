package config

// parseSqlDir retrieves sql source directory.
func sanitizeSqlPath(path string) (newPath string) {
	if string(path[len(path)-1]) != "/" {
		path += "/"
	}

	newPath = path

	return
}
