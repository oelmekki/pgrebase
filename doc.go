/*
PgRebase is a postgres codebase injection tool.

Usage:

  DATABASE_URL=url pgrebase [-w] sql_dir/

Expected target structure:

	<sql_dir>/
		functions/
		triggers/
		types/
		views/

At least one of functions/triggers/types/views/ should exist.
*/
package main
