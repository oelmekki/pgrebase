package core

import (
	"fmt"
)

// codeUnit is the generic representation of any code, be it
// a function, a view, etc.
type codeUnit struct {
	path           string // the absolute path to code file
	name           string // the name of function/view/trigger
	definition     string // the actual code
	previousExists bool   // true if this code unit already exists in database
	parseSignature bool   // true if we need to generate signature (for new functions)
}

// drop is the generic drop function for code units.
func (unit codeUnit) drop(dropQuery string) (err error) {
	rows, err := query(dropQuery)
	if err != nil {
		return err
	}
	rows.Close()

	return
}

// create is the eneric creation function for code units.
func (unit codeUnit) create(definition string) (err error) {
	rows, err := query(definition)
	if err != nil {
		return err
	}
	rows.Close()

	return
}

// codeUnitCreator is the interface that repesents what
// can manipulate code units.
type codeUnitCreator interface {
	load() error
	parse() error
	drop() error
	create() error
}

// downPass performs the Steps used in down pass, when dropping existing code, in dependency
// graph reverse order.
func downPass(unit codeUnitCreator, path string) (err error) {
	errFmt := "  error while loading %s\n  %v\n"

	if err = unit.load(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}
	if err = unit.parse(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}
	if err = unit.drop(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}

	return
}

// upPass performs the steps used in up pass, when creating existing code, in dependency
// graph order
func upPass(unit codeUnitCreator, path string) (err error) {
	errFmt := "  error while creating %s\n  %v\n"
	if err = unit.create(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}

	return
}
