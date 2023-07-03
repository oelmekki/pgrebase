package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// loadTriggers loads or reloads all triggers found in FS.
func loadTriggers() (err error) {
	successfulCount := len(conf.triggerFiles)
	errors := make([]string, 0)
	bypass := make(map[string]bool)

	files, err := resolveDependencies(conf.triggerFiles, conf.sqlDirPath+"triggers")
	if err != nil {
		return err
	}

	triggers := make([]*trigger, 0)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		t := trigger{}
		t.path = file
		triggers = append(triggers, &t)

		err = downPass(&t, t.path)
		if err != nil {
			successfulCount--
			errors = append(errors, fmt.Sprintf("%v\n", err))
			bypass[t.path] = true
		}
	}

	for i := len(triggers) - 1; i >= 0; i-- {
		t := triggers[i]
		if _, ignore := bypass[t.path]; !ignore {
			err = upPass(t, t.path)
			if err != nil {
				successfulCount--
				errors = append(errors, fmt.Sprintf("%v\n", err))
			}
		}
	}

	report("triggers", successfulCount, len(conf.triggerFiles), errors)

	return
}

// trigger is the code unit for triggers.
type trigger struct {
	codeUnit
	table    string   // name of the table for the trigger
	function function // related function for trigger
}

// load loads trigger definition from file.
func (t *trigger) load() (err error) {
	definition, err := ioutil.ReadFile(t.path)
	if err != nil {
		return err
	}
	t.definition = string(definition)

	return
}

// parse parses trigger for name and signature.
func (t *trigger) parse() (err error) {
	triggerFinder := regexp.MustCompile(`(?is)CREATE(?:\s+CONSTRAINT)?\s+TRIGGER\s+(\S+).*?ON\s+(\S+)`)
	subMatches := triggerFinder.FindStringSubmatch(t.definition)

	if len(subMatches) < 3 {
		return fmt.Errorf("Can't find a tin %s", t.path)
	}

	t.function = function{codeUnit: codeUnit{definition: t.definition, path: t.path}}
	t.function.parse()

	t.name = subMatches[1]
	t.table = subMatches[2]

	return
}

// drop removes existing trigger from pg.
func (t *trigger) drop() (err error) {
	err = t.codeUnit.drop(`DROP TRIGGER IF EXISTS ` + t.name + ` ON ` + t.table)
	if err != nil {
		return err
	}

	return t.function.drop()
}

// create adds the trigger in pg.
func (t *trigger) create() (err error) {
	return t.codeUnit.create(t.definition)
}
