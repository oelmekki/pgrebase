package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// loadViews loads or reloads all views found in FS.
func loadViews() (err error) {
	successfulCount := len(conf.viewFiles)
	errors := make([]string, 0)
	bypass := make(map[string]bool)

	files, err := resolveDependencies(conf.viewFiles, conf.sqlDirPath+"views")
	if err != nil {
		return err
	}

	views := make([]*view, 0)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		v := view{}
		v.path = file
		views = append(views, &v)

		err = downPass(&v, v.path)
		if err != nil {
			successfulCount--
			errors = append(errors, fmt.Sprintf("%v\n", err))
			bypass[v.path] = true
		}
	}

	for i := len(views) - 1; i >= 0; i-- {
		v := views[i]
		if _, ignore := bypass[v.path]; !ignore {
			err = upPass(v, v.path)
			if err != nil {
				successfulCount--
				errors = append(errors, fmt.Sprintf("%v\n", err))
			}
		}
	}

	report("views", successfulCount, len(conf.viewFiles), errors)

	return
}

// view is the code unit for views.
type view struct {
	codeUnit
}

// load loads view definition from file.
func (v *view) load() (err error) {
	definition, err := ioutil.ReadFile(v.path)
	if err != nil {
		return err
	}
	v.definition = string(definition)

	return
}

// parse parses view for name.
func (v *view) parse() (err error) {
	nameFinder := regexp.MustCompile(`(?is)CREATE(?:\s+OR\s+REPLACE)?\s+VIEW\s+(\S+)`)
	subMatches := nameFinder.FindStringSubmatch(v.definition)

	if len(subMatches) < 2 {
		return fmt.Errorf("Can't find a view in %s", v.path)
	}

	v.name = subMatches[1]

	return
}

// drop removes existing view from pg.
func (v *view) drop() (err error) {
	return v.codeUnit.drop(`DROP VIEW IF EXISTS ` + v.name)
}

// create adds the view in pg.
func (v *view) create() (err error) {
	return v.codeUnit.create(v.definition)
}
