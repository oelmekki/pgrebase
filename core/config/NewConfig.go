package config

// NewConfig creates a config data structure with provided values.
//
// databaseUrl is a postgres connection url.
//
// sqlDir is the path to the directory where your sql sources are.
//
// watch is a flag you may set to true to keep watching for change
// in sqlDir after processing it a first time.
func NewConfig(databaseUrl, sqlDir string) (config Config) {
	config = Config{DatabaseUrl: databaseUrl, SqlDirPath: sanitizeSqlPath(sqlDir)}
	return
}
