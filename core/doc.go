/*
Package core allows to use pgrebase as a library in your own code.

To use it, your first need to initialize it:

  err := core.Init(databaseUrl, sqlDir)

databaseUrl is the postgres connection string, sqlDir is the path to your sql sources,
and watch is a flag you can set to true to keep watching sqlDir for changes.

Once pgrebase is initialized, call Process() to load your source files into
database.

If you want to keep watching FS for changes after than, you can call Watch().
*/
package core
