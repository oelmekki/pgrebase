package core

import (
	"gitlab.com/oelmekki/pgrebase/core/function"
	"gitlab.com/oelmekki/pgrebase/core/trigger"
	"gitlab.com/oelmekki/pgrebase/core/view"
)

// Process loads sql code, just once.
func Process() (err error) {
	if err = view.LoadViews(); err != nil {
		return err
	}
	if err = function.LoadFunctions(); err != nil {
		return err
	}
	if err = trigger.LoadTriggers(); err != nil {
		return err
	}

	return
}
