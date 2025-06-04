package objectmap

import (
	"errors"
	"reflect"
	"strings"
)

// traverseDFS recursively iterates over all fields in a struct, including nested structs.
// If the input is a pointer, it is automatically dereferenced.
//
// Parameters:
//   - value: A reflect.Value representing a struct or a pointer to a struct.
//   - path: The current path in the field hierarchy (used for building full field paths).
//   - visitFunc: A function called for each field. It receives:
//   - fullPath: The dot-separated path to the field (e.g. "Address.Street").
//   - field: The reflect.StructField containing field metadata.
//   - value: The reflect.Value representing the field's value.
func traverseDFS(
	value reflect.Value,
	path []string,
	visitFunc func(fullPath []string, field reflect.StructField, value reflect.Value) bool,
) {

	// Ensure the value is a struct or a pointer to a struct
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return
	}

	// Iterate over the fields of the struct
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldVal := value.Field(i)
		fullPath := append(path, field.Name)

		shouldRecurse := visitFunc(fullPath, field, fieldVal)

		if shouldRecurse && (fieldVal.Kind() == reflect.Struct || fieldVal.Kind() == reflect.Ptr) {
			traverseDFS(fieldVal, fullPath, visitFunc)
		}
	}

}

// HasNestedFieldSlice returns (found, field metadata, field value).
// If not found, found is false and others are zero values.
func hasNestedFieldSlice(structValue reflect.Value, path []string) (bool, reflect.StructField, reflect.Value) {
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}
	if !structValue.IsValid() || structValue.Kind() != reflect.Struct {
		return false, reflect.StructField{}, reflect.Value{}
	}

	currVal := structValue
	currType := structValue.Type()

	for i, part := range path {
		field, ok := currType.FieldByName(part)
		if !ok {
			return false, reflect.StructField{}, reflect.Value{}
		}

		fieldVal := currVal.FieldByIndex(field.Index)
		if !fieldVal.IsValid() {
			return false, reflect.StructField{}, reflect.Value{}
		}

		// Only dereference if not the last element in the path
		if i < len(path)-1 {
			for fieldVal.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					return false, reflect.StructField{}, reflect.Value{}
				}
				fieldVal = fieldVal.Elem()
			}
		}

		currVal = fieldVal
		currType = fieldVal.Type()

		// If not last element, must be struct to continue
		if i < len(path)-1 && currVal.Kind() != reflect.Struct {
			return false, reflect.StructField{}, reflect.Value{}
		}

		// On last element, return the field
		if i == len(path)-1 {
			return true, field, currVal
		}
	}

	return false, reflect.StructField{}, reflect.Value{}
}

// SetNestedField sets the nested field at the given path in structPtr to newVal.
// structPtr must be a reflect.Value of a pointer to a struct (addressable).
// path is a slice of field names, e.g. []string{"Address", "Street", "Name"}.
// newVal must be assignable to the target field type.
func setNestedField(structPtr reflect.Value, path []string, newVal reflect.Value) error {
	if structPtr.Kind() != reflect.Ptr || structPtr.IsNil() {
		return errors.New("input must be a non-nil pointer to a struct")
	}

	currVal := structPtr.Elem()
	if currVal.Kind() != reflect.Struct {
		return errors.New("input must point to a struct")
	}

	for i, part := range path {
		fieldVal := currVal.FieldByName(part)
		if !fieldVal.IsValid() {
			return errors.New("field not found: " + part)
		}

		// If it's a pointer, dereference or allocate if nil (only for intermediate fields)
		if i < len(path)-1 && fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				// Allocate new struct for pointer field
				fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
			}
			fieldVal = fieldVal.Elem()
		}

		// If not last field, descend into nested struct
		if i < len(path)-1 {
			if fieldVal.Kind() != reflect.Struct {
				return errors.New("field " + part + " is not a struct")
			}
			currVal = fieldVal
			continue
		}

		// Last field: set value
		if !fieldVal.CanSet() {
			return errors.New("cannot set field: " + part)
		}

		// Check assignability
		if !newVal.Type().AssignableTo(fieldVal.Type()) {
			return errors.New("cannot assign value of type " + newVal.Type().String() + " to field " + part + " of type " + fieldVal.Type().String())
		}

		fieldVal.Set(newVal)
		return nil
	}

	return errors.New("unexpected error setting nested field")
}

// SplitPath splits a dot-separated string path like "Address.Street.Name"
// into a slice of strings: []string{"Address", "Street", "Name"}.
func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	return strings.Split(path, ".")
}

// JoinPath joins a slice of strings like []string{"Address", "Street", "Name"}
// into a dot-separated string path like "Address.Street.Name".
func joinPath(path []string) string {
	return strings.Join(path, ".")
}

// markPathAndParents marks the given path and all of its parent paths as mapped.
// For example, if path is "N.A.B", it marks "N", "N.A", and "N.A.B" as true in the map.
// This ensures that when nested fields are set, their parent fields are also considered set.
// TODO: I feel like this could be made more generic... think about trees and nodes...
func markPathAndParents(mapped map[string]bool, path string) {
	parts := splitPath(path) // ["N", "A", "B"]
	for i := 1; i <= len(parts); i++ {
		prefix := joinPath(parts[:i])
		mapped[prefix] = true
	}
}

// mapInto copies matching fields from src to dst (both must be pointers to structs)
// MapFieldsWithTag maps fields from src to dst.
// By default, matches source field name to destination field name.
// If source field has `mapto:"FieldName"` tag, maps to that destination field instead.
// Supports onUnusedDst and onUnusedSrc callbacks.
// The callback functions also add the ability to hook in very helpful behavior
// for testing.
func MapInto[S any, D any](
	src *S,
	dst *D,
	onUnusedSrc func(fieldPath []string, srcVal reflect.Value),
	onUnusedDst func(fieldPath []string, dstVal reflect.Value),
	maptag string,
) error {

	if onUnusedSrc == nil {
		onUnusedSrc = func(_ []string, _ reflect.Value) {}
	}
	if onUnusedDst == nil {
		onUnusedDst = func(_ []string, _ reflect.Value) {}
	}

	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return errors.New("both src and dst must be pointers")
	}

	srcElem := srcVal.Elem()
	dstElem := dstVal.Elem()

	if srcElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Struct {
		return errors.New("both src and dst must point to structs")
	}

	mappedDstPaths := make(map[string]bool)
	traverseDFS(
		srcElem,
		nil,
		func(srcPath []string, srcField reflect.StructField, srcValue reflect.Value) bool {

			// Check tag map
			dstTargetPath := joinPath(srcPath) // default: same name
			// TODO: I may get rid of this default scheme, as what if there happens
			// to be a field in source with the same name that you DON'T want mapping
			if tagVal, ok := srcField.Tag.Lookup(maptag); ok && tagVal != "" {
				dstTargetPath = tagVal
			}
			targetPath := splitPath(dstTargetPath)

			// Try and map the field
			dstValid, dstField, dstFieldVal := hasNestedFieldSlice(dstElem, targetPath)

			if !dstValid {
				onUnusedSrc(srcPath, srcValue)
				return true // Recurse deeper
			}
			// Check type compatibility: must be exactly same type or assignable
			if dstField.Type != srcField.Type && !srcValue.Type().AssignableTo(dstField.Type) {
				onUnusedSrc(srcPath, srcValue)
				return true // Recurse deeper
			}
			// Check that dstFieldVal is settable
			if !dstFieldVal.CanSet() {
				onUnusedSrc(srcPath, srcValue)
				return true // Recurse deeper
			}

			// Now, we can map. Don't recurse any more
			err := setNestedField(dstVal, targetPath, srcValue)
			if err != nil {
				onUnusedSrc(srcPath, srcValue)
				return false
			}
			markPathAndParents(mappedDstPaths, joinPath(targetPath))

			return false
		},
	)

	// Check dst fields that were not mapped
	traverseDFS(
		dstElem,
		nil,
		func(dstPath []string, dstField reflect.StructField, dstValue reflect.Value) bool {
			_, mapped := mappedDstPaths[joinPath(dstPath)]
			if mapped {
				return false // do not recurse more
			} else {
				onUnusedDst(dstPath, dstValue)
				return true // recurse
			}
		},
	)

	return nil
}

// mapFrom copies matching fields from src to dst (both must be pointers to structs).
// MapFieldsWithTag maps fields from src to dst based on tags defined on the destination struct.
//
// By default, matches destination field name to source field name. If a destination field
// has a `mapfrom:"FieldName"` tag, it maps from the specified source field instead.
//
// Fields are matched by name or tag and must have compatible types (either identical or assignable).
// Nested fields are supported via dot-separated paths.
//
// Supports onUnusedSrc and onUnusedDst callbacks, which are invoked for unmapped source or
// destination fields, respectively. These callbacks are useful for debugging, testing,
// or enforcing strict field usage policies.
func MapFrom[S any, D any](
	src *S,
	dst *D,
	onUnusedSrc func(fieldPath []string, srcVal reflect.Value),
	onUnusedDst func(fieldPath []string, dstVal reflect.Value),
	maptag string,
) error {

	if onUnusedSrc == nil {
		onUnusedSrc = func(_ []string, _ reflect.Value) {}
	}
	if onUnusedDst == nil {
		onUnusedDst = func(_ []string, _ reflect.Value) {}
	}

	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return errors.New("both src and dst must be pointers")
	}

	srcElem := srcVal.Elem()
	dstElem := dstVal.Elem()

	if srcElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Struct {
		return errors.New("both src and dst must point to structs")
	}

	mappedDstPaths := make(map[string]bool)
	traverseDFS(
		dstElem,
		nil,
		func(dstPath []string, dstField reflect.StructField, dstValue reflect.Value) bool {
			srcTargetPath := joinPath(dstPath) // default: same name

			if tagVal, ok := dstField.Tag.Lookup(maptag); ok && tagVal != "" {
				srcTargetPath = tagVal
			}
			sourcePath := splitPath(srcTargetPath)

			srcValid, srcField, srcFieldVal := hasNestedFieldSlice(srcElem, sourcePath)

			if !srcValid {
				onUnusedDst(dstPath, dstValue)
				return true // Recurse deeper
			}
			// Check type compatibility
			if dstField.Type != srcField.Type && !srcFieldVal.Type().AssignableTo(dstField.Type) {
				onUnusedDst(dstPath, dstValue)
				return true
			}
			if !dstValue.CanSet() {
				onUnusedDst(dstPath, dstValue)
				return true
			}

			// Perform mapping
			err := setNestedField(dstVal, dstPath, srcFieldVal)
			if err != nil {
				onUnusedDst(dstPath, dstValue)
				return false
			}
			markPathAndParents(mappedDstPaths, joinPath(dstPath))

			return false
		},
	)

	// Check src fields that were not mapped
	traverseDFS(
		srcElem,
		nil,
		func(srcPath []string, srcField reflect.StructField, srcValue reflect.Value) bool {
			_, mapped := mappedDstPaths[joinPath(srcPath)]
			if mapped {
				return false
			} else {
				onUnusedSrc(srcPath, srcValue)
				return true
			}
		},
	)

	return nil
}
