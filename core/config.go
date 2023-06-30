package core

// config is the global configuration for execution.
type config struct {
	databaseUrl   string   // connection info for the database
	sqlDirPath    string   // place where to find the code units
	functionFiles []string // paths of all function files
	triggerFiles  []string // paths of all trigger files
	viewFiles     []string // paths of all view files
}

// parseSqlDir retrieves sql source directory.
func sanitizeSqlPath(path string) (newPath string) {
	if string(path[len(path)-1]) != "/" {
		path += "/"
	}

	newPath = path

	return
}

// newConfig creates a config data structure with provided values.
//
// databaseUrl is a postgres connection url.
//
// sqlDir is the path to the directory where your sql sources are.
//
// watch is a flag you may set to true to keep watching for change
// in sqlDir after processing it a first time.
func newConfig(databaseUrl, sqlDir string) (conf config) {
	conf = config{databaseUrl: databaseUrl, sqlDirPath: sanitizeSqlPath(sqlDir)}
	return
}
