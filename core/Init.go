package core

import (
	"gitlab.com/oelmekki/pgrebase/core/config"
	"gitlab.com/oelmekki/pgrebase/core/connection"
	"gitlab.com/oelmekki/pgrebase/core/function"
	"gitlab.com/oelmekki/pgrebase/core/trigger"
	"gitlab.com/oelmekki/pgrebase/core/view"
)

// conf is the global level configuration data structure.
var conf config.Config

// Init stores the global config object.
//
// databaseUrl should be a connection string to the database (eg: postgres://postgres:@localhost/database).
//
// sqlDir is the path to the directory where sql source files live.
//
// watch should be true if you want to keep watching for changes in source files rather
// than just loading them once.
func Init(databaseUrl, sqlDir string) (err error) {
	conf = config.NewConfig(databaseUrl, sqlDir)
	connection.Init(&conf)
	function.Init(&conf)
	trigger.Init(&conf)
	view.Init(&conf)

	checker := sanity{}
	err = checker.Check()

	return
}
