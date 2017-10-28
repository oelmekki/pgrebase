/*
PgRebase is a tool that allows you to easily handle your postgres codebase for
functions, triggers and views.

PgRebase allows you to manage your functions/triggers/views as plain files in
filesystem. You put them in a `sql/` directory, one file per
function/trigger/type/view.

  $ tree sql/
  sql/
  ├── functions/
  │   └── assign_user_to_team.sql
  ├── triggers/
  │   └── user_updated_at.sql
  └── views/
      └── user_json.sql

No need to add drop statement in those files, PgRebase will take care of it.

In watch mode (useful for development), just save your file, pgrebase will
update your database. In normal mode (useful for deployment), pgrebase will
recreate all functions/triggers/views found in your filesystem directory.

You can now work with postgres codebase live reload, then call pgrebase just
after your migration task in your deployment pipeline.

Note: this is the documentation for pgrebase as a command tool. If you
want to use pgrebase as a library, see https://github.com/oelmekki/pgrebase/core
(or core subpackage in godoc).


Usage

Here is an example session with pgrebase:

	$ export DATABASE_URL=postgres://user:pass@host/db

	$ ./pgrebase sql/
	Loaded 10 functions
	Loaded 25 views
	Loaded 5 triggers - 1 trigger with error
		error while loading sql/triggers/user_updated_at.sql
		column users.updated_at does not exist


	$ ./pgrebase -w sql/
	Loaded 10 functions
	Loaded 25 views
	Loaded 6 triggers
	Watching filesystem for changes...
	FS changed. Building.

When working in development environment, you'll probably want to use watch mode
(`-w`) to have your changes automatically loaded.

For deployment, add `pgrebase` to your repos and call it after your usual
migrations step:

	DATABASE_URL=your_config ./pgrebase ./sql


Handling dependencies

You can specify dependencies for files using require statement, provided those
files are of the same kind. That is, function files can specify dependencies on
other function files, type files can define dependencies on other type files,
etc.

Here is an example about how to do it. Let's say your `sql/functions/foo.sql`
files depends on `sql/functions/whatever/bar.sql`:

	$ cat sql/functions/foo.sql
	-- require "whatever/bar.sql"
	CREATE FUNCTION foo()
	[...]

Filenames are always relative to your target directory (`sql/` in that
example), and within in, to the code kind (`functions/` here).

Do not try to do funky things like adding `./` or `../`, this is no path
resolution, it just tries to match filenames.

You can add multiple require lines:

	-- require "common.sql"
	-- require "hello/world.sql"
	-- require "whatever/bar.sql"
	CREATE FUNCTION foo()
	[...]

There is no advanced debugging for circular dependencies for now, so be sure
not to get too wild, here (or else, you will have a "maybe there's circular
dependencies?" message and you will have to figure it out for yourself).


Caveats

pgrebase doesn't keep any state about your codebase and does not delete what
is in your database and is not in your codebase. This means that if you want
to remove a trigger/type/view/function, deleting its file is not enough. You
have to use your usual way to migrate db and remove it.

trigger files should contain both trigger creation and the function it uses.
This is to avoid dropping function still used by trigger (if processing
functions first) or create trigger before its function (if triggers are
processed first).

files should only contain the definition of the view/function/type/trigger
they're named after (with the exception of trigger files declaring the
function they use). Hazardous results will ensue if it's not the case: only
the first definition will be dropped, but the whole file will be loaded in
pg.

pgrebase single top concern is to not delete any data. This means that no
DROP will CASCADE. This means that if your database structure depends on
anything defined in pgrebase codebase, it will fail to reload it when it
implies dropping it (that is, most of the time). In other terms, do not
use pgrebase managed functions as default value for your tables' fields.
*/
package main
