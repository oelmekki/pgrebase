package trigger

import (
	"gitlab.com/oelmekki/pgrebase/core/config"
)

var conf *config.Config

// Init stores configuration for further usage.
func Init(cfg *config.Config) {
	conf = cfg
}
