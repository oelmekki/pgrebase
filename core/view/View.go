package view

import (
	"fmt"
	"github.com/oelmekki/pgrebase/core/codeunit"
	"github.com/oelmekki/pgrebase/core/resolver"
	"github.com/oelmekki/pgrebase/core/utils"
	"io/ioutil"
	"regexp"
)

// LoadViews loads or reloads all views found in FS.
func LoadViews() (err error) {
	successfulCount := len(conf.ViewFiles)
	errors := make([]string, 0)
	bypass := make(map[string]bool)

	files, err := resolver.ResolveDependencies(conf.ViewFiles, conf.SqlDirPath+"views")
	if err != nil {
		return err
	}

	views := make([]*View, 0)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		view := View{}
		view.Path = file
		views = append(views, &view)

		err = codeunit.DownPass(&view, view.Path)
		if err != nil {
			successfulCount--
			errors = append(errors, fmt.Sprintf("%v\n", err))
			bypass[view.Path] = true
		}
	}

	for i := len(views) - 1; i >= 0; i-- {
		view := views[i]
		if _, ignore := bypass[view.Path]; !ignore {
			err = codeunit.UpPass(view, view.Path)
			if err != nil {
				successfulCount--
				errors = append(errors, fmt.Sprintf("%v\n", err))
			}
		}
	}

	utils.Report("views", successfulCount, len(conf.ViewFiles), errors)

	return
}

// View is the code unit for views.
type View struct {
	codeunit.CodeUnit
}

// Load loads view definition from file.
func (view *View) Load() (err error) {
	definition, err := ioutil.ReadFile(view.Path)
	if err != nil {
		return err
	}
	view.Definition = string(definition)

	return
}

// Parse parses view for name.
func (view *View) Parse() (err error) {
	nameFinder := regexp.MustCompile(`(?is)CREATE(?:\s+OR\s+REPLACE)?\s+VIEW\s+(\S+)`)
	subMatches := nameFinder.FindStringSubmatch(view.Definition)

	if len(subMatches) < 2 {
		return fmt.Errorf("Can't find a view in %s", view.Path)
	}

	view.Name = subMatches[1]

	return
}

// Drop removes existing view from pg.
func (view *View) Drop() (err error) {
	return view.CodeUnit.Drop(`DROP VIEW IF EXISTS ` + view.Name)
}

// Create adds the view in pg.
func (view *View) Create() (err error) {
	return view.CodeUnit.Create(view.Definition)
}
