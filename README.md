# PgRebase

PgRebase is a tool that allows you to easily handle your postgres codebase for
functions, triggers and views.


## Why

If you started outsourcing data manipulation to your database through
postgresql cool features, you probably realized this is painful. This is not
your usual codebase, it lives in postgres, and you often have to drop your
function/trigger/view if you want to edit it. Migrating servers is not easy
either.

The classic tool for this is to have migration files. This is great for
handling tables, not so great to make frequent changes to your functions. Can
we do better?


## What

PgRebase allows you to manage your functions/triggers/views as plain files in
filesystem. You put them in a `sql/` directory, one file per
function/trigger/view.

```
$ tree sql/
sql/
├── functions/
│   └── assign_user_to_team.sql
├── triggers/
│   └── user_updated_at.sql
└── views/
    └── user_json.sql
```

No need to add drop statement in those files, PgRebase will take care of it.


## Build

PgRebase use [gb](https://getgb.io/) as a building tool.

To build it:

```
go get github.com/constabulary/gb/...  # if you don't have gb yet
gb build all
```

Binary will be in `bin/pgrebase`.

If you don't want to build the project, a binary for linux/amd64 is already
there.


## Usage

```
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
```

When working in development environment, you'll probably want to use watch mode
(`-w`) to have your changes automatically loaded.

For deployment, add `pgrebase` to your repos and call it after your usual
migrations step:

```
DATABASE_URL=your_config ./pgrebase ./sql
```


## Caveats

* pgrebase doesn't keep any state about your codebase and does not delete what
  is in your database and is not in your codebase. This means that if you want
  to remove a trigger/view/function, deleting its file is not enough. You have
  to use your usual way to migrate db and remove it.

* trigger files should contain both trigger creation and the function it uses.
  This is to avoid dropping function still used by trigger (if processing
  functions first) or create trigger before its function (if triggers are
  processed first)

* files should only contain the definition of the view/function/trigger they're
  named after (with the exception of trigger files declaring the function they
  use). Hazardous results will ensue if it's not the case: only the first
  definition will be dropped, but the whole file will be loaded in pg.


## Credits

PgRebase was born after discussing with Derek Sivers about moving business logic
to the database. Make sure to read [his research](https://sivers.org/pg), it's
awesome!
