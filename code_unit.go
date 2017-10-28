package main

import (
	"fmt"
)

type CodeUnit struct {
	Path           string
	Name           string
	Definition     string
	previousExists bool
	parseSignature bool
}

/*
 * Generic drop function for code units
 */
func (unit CodeUnit) Drop(dropQuery string) (err error) {
	rows, err := Query(dropQuery)
	if err != nil {
		return err
	}
	rows.Close()

	return
}

/*
 * Generic creation function for code units
 */
func (unit CodeUnit) Create(definition string) (err error) {
	rows, err := Query(definition)
	if err != nil {
		return err
	}
	rows.Close()

	return
}

type CodeUnitCreator interface {
	Load() error
	Parse() error
	Drop() error
	Create() error
}

/*
 * Steps used in down pass, when dropping existing code, in dependency
 * graph reverse order
 */
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

/*
 * Steps used in up pass, when creating existing code, in dependency
 * graph order
 */
func UpPass(unit CodeUnitCreator, path string) (err error) {
	errFmt := "  error while creating %s\n  %v\n"
	if err = unit.Create(); err != nil {
		return fmt.Errorf(errFmt, path, err)
	}

	return
}
