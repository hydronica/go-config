package encode

import (
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/jbsmith7741/go-tools/appenderr"
)

// CheckRequired reports an aggregated error for every field tagged req:"true" that is
// still the zero value. Call it after all config sources (defaults, env, files, flags)
// have loaded, so a struct default satisfies the requirement but an explicit zero value
// (for example 0 for a required number) does not.
//
// Each error names the flag name for the field, the same name the flag package
// registers (kebab-case, dot-joined for nested structs), for example:
//
//	db.password is required
//
// A struct-typed field may also carry req:"true"; it is reported missing when every
// field in the struct is zero. CheckRequired still recurses into nested structs to
// check their own req tags regardless of the parent's tag.
//
// Returns nil if i is not a pointer to a struct.
func CheckRequired(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return nil
	}
	errs := appenderr.New()
	checkRequired(errs, "", v.Elem())
	return errs.ErrOrNil()
}

func checkRequired(errs *appenderr.AppendErr, prefix string, vStruct reflect.Value) {
	for i := 0; i < vStruct.NumField(); i++ {
		sf := vStruct.Type().Field(i)
		field := vStruct.Field(i)
		if !ShouldProcessField(sf, field) {
			continue
		}
		if sf.Tag.Get(ConfigTag) == "ignore" {
			continue
		}
		raw := sf.Tag.Get(FlagTag)
		if raw == "-" {
			continue
		}
		name, omitPrefix, _ := ParseTag(raw)

		if IsRecursableStruct(field) {
			if sf.Tag.Get(RequiredTag) == "true" && field.IsZero() {
				errs.Addf("%s is required", requiredKey(prefix, sf, name, omitPrefix))
			}
			childPrefix := prefix
			switch {
			case name == "omitprefix":
				childPrefix = ""
			case name != "":
				childPrefix = JoinPrefix(prefix, name, ".")
			}
			child, ok := StructValue(field)
			if !ok {
				continue
			}
			checkRequired(errs, childPrefix, child)
			continue
		}

		if sf.Tag.Get(RequiredTag) != "true" {
			continue
		}
		if !field.IsZero() {
			continue
		}
		errs.Addf("%s is required", requiredKey(prefix, sf, name, omitPrefix))
	}
}

// requiredKey builds the flag-style name used in a required-field error, matching the
// name the flag package registers for the same field.
func requiredKey(prefix string, sf reflect.StructField, name string, omitPrefix bool) string {
	segment := strcase.ToKebab(sf.Name)
	if name != "" && name != "omitprefix" {
		segment = name
	}
	key := JoinPrefix(prefix, segment, ".")
	if omitPrefix {
		key = segment
	}
	return key
}
