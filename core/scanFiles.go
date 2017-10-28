package core

import (
	"github.com/oelmekki/pgrebase/core/config"
)

// scanFiles scans sql directory for sql files.
func scanFiles(cfg *config.Config) {
	cfg.FunctionFiles = make([]string, 0)
	cfg.TriggerFiles = make([]string, 0)
	cfg.TypeFiles = make([]string, 0)
	cfg.ViewFiles = make([]string, 0)

	walker := sourceWalker{Config: cfg}
	walker.Process()
}
