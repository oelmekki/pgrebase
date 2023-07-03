Package core allows to use pgrebase as a library in your own code.

To use it, your first need to initialize it:

    err := core.Init(databaseUrl, sqlDir)

`databaseUrl` is the postgres connection string as used by
`database/sql.Open`, `sqlDir` is the path to your sql sources, and watch is
a flag you can set to true to keep watching sqlDir for changes.

Once pgrebase is initialized, call `Process()` to load your source files
into database.

If you want to keep watching FS for changes after than, you can call
`Watch()`. Note that this function won't return, so unless that's what you
want, you should run it in a goroutine.

