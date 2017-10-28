package main

import (
	"fmt"
)

// CodeUnit is the generic representation of any code, be it
// a function, a view, etc.
type CodeUnit struct {
	Path           string // the absolute path to code file
	Name           string // the name of function/view/trigger
	Definition     string // the actual code
	previousExists bool   // true if this code unit already exists in database
	parseSignature bool   // true if we need to generate signature (for new functions)
}

// Drop is the generic drop function for code units.
func (unit CodeUnit) Drop(dropQuery string) (err error) {
	rows, err := Query(dropQuery)
	if err != nil {
		return err
	}
	rows.Close()

	return
}

// Create is the eneric creation function for code units.
func (unit CodeUnit) Create(definition string) (err error) {
	rows, err := Query(definition)
	if err != nil {
		return err
	}
	rows.Close()

	return
}

// CodeUnitCreator is the interface that repesents what
// can manipulate code units.
type CodeUnitCreator interface {
	Load() error
	Parse() error
	Drop() error
	Create() error
}

// DownPass performs the Steps used in down pass, when dropping existing code, in dependency
// graph reverse order.
func DownPass(unit CodeUnitCreator, path string) (err error) {
	errFmt := "  error while loading %s\n  %v\n"

	if err = unit.Load(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}
	if err = unit.Parse(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}
	if err = unit.Drop(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}

	return
}

// UpPass performs the steps used in up pass, when creating existing code, in dependency
// graph order
func UpPass(unit CodeUnitCreator, path string) (err error) {
	errFmt := "  error while creating %s\n  %v\n"
	if err = unit.Create(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}

	return
}
