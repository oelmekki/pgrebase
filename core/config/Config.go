package config

// Config is the global configuration for execution.
type Config struct {
	DatabaseUrl   string   // connection info for the database
	SqlDirPath    string   // place where to find the code units
	FunctionFiles []string // paths of all function files
	TriggerFiles  []string // paths of all trigger files
	TypeFiles     []string // paths of all type files
	ViewFiles     []string // paths of all view files
}
