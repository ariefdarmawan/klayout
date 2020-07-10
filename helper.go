package klayout

import (
	"reflect"

	"github.com/eaciit/toolkit"
)

var DefaultDateFormat = "DD-MMM-YYYY"

func tagDefault(tag reflect.StructTag, name string, def string) string {
	if tmp := tag.Get(name); tmp == "" {
		return def
	} else {
		return tmp
	}
}

func tagMemberDefault(tag reflect.StructTag, name string, alloweds []string, def string) string {
	if tmp := tag.Get(name); toolkit.HasMember(alloweds, tmp) {
		return tmp
	} else {
		return def
	}
}
