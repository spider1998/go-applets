package resource

import (
	"github.com/gobuffalo/packr"
)

var MigrationBox packr.Box

func Load() {
	MigrationBox = packr.NewBox("./migrations")
}
