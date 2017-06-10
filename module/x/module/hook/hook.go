package hook

import (
	"github.com/insionng/zenpress/module/hook"

	"qlang.io/spec"
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "github.com/insionng/zenpress/module/hook",

	"CurrentFilterKey": hook.CurrentFilterKey,
	"FiltersMap":       hook.FiltersMap,
	"KueuesMap":        hook.KueuesMap,

	"AddAction":     hook.AddAction,
	"AddFilter":     hook.AddFilter,
	"CurrentFilter": hook.CurrentFilter,
	"DoAction":      hook.DoAction,
	"DoFilter":      hook.DoFilter,
	"HasAction":     hook.HasAction,
	"HasFilter":     hook.HasFilter,
	"RemoveAction":  hook.RemoveAction,
	"RemoveActions": hook.RemoveActions,
	"RemoveFilter":  hook.RemoveFilter,

	"Node": spec.StructOf((*hook.Node)(nil)),
}
