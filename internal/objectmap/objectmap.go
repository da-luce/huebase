package objectmap

import (
	"errors"
	"reflect"
	"strings"

	"github.com/da-luce/paletteport/internal/structutil"
)

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
func MapInto(
	src any,
	dst any,
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
	structutil.TraverseStructDFS(
		src,
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
			dstValid, dstField, dstFieldVal := structutil.HasNestedFieldSlice(dstElem, targetPath)

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
			err := structutil.SetNestedField(dstVal, targetPath, srcValue)
			if err != nil {
				onUnusedSrc(srcPath, srcValue)
				return false
			}
			markPathAndParents(mappedDstPaths, joinPath(targetPath))

			return false
		},
	)

	// Check dst fields that were not mapped
	structutil.TraverseStructDFS(
		dst,
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
func MapFrom(
	src any,
	dst any,
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

	if dstVal.Kind() != reflect.Ptr {
		// WARN: changes won't take effect
	}

	srcElem := srcVal.Elem()
	dstElem := dstVal.Elem()

	if srcElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Struct {
		return errors.New("both src and dst must point to structs")
	}

	mappedDstPaths := make(map[string]bool)

	// Traverse destination to fill it from source
	structutil.TraverseStructDFS(
		dst,
		func(dstPath []string, dstField reflect.StructField, dstValue reflect.Value) bool {
			srcTargetPath := joinPath(dstPath) // default: same name

			if tagVal, ok := dstField.Tag.Lookup(maptag); ok && tagVal != "" {
				srcTargetPath = tagVal
			}
			sourcePath := splitPath(srcTargetPath)

			srcValid, _, srcFieldVal := structutil.HasNestedFieldSlice(srcElem, sourcePath)
			if !srcValid {
				onUnusedDst(dstPath, dstValue)
				return true
			}

			// Normalize for assignment (handle *T → T, T → *T)
			normalizedVal, ok := normalizeForAssignment(srcFieldVal, dstField.Type)
			if !ok {
				onUnusedDst(dstPath, dstValue)
				return true
			}

			if !dstValue.CanSet() {
				onUnusedDst(dstPath, dstValue)
				return true
			}

			err := structutil.SetNestedField(dstVal, dstPath, normalizedVal)
			if err != nil {
				onUnusedDst(dstPath, dstValue)
				return false
			}

			markPathAndParents(mappedDstPaths, joinPath(dstPath))
			return false
		},
	)

	// Track unused source fields
	structutil.TraverseStructDFS(
		src,
		func(srcPath []string, srcField reflect.StructField, srcValue reflect.Value) bool {
			_, mapped := mappedDstPaths[joinPath(srcPath)]
			if !mapped {
				onUnusedSrc(srcPath, srcValue)
			}
			return true
		},
	)

	return nil
}

// normalizeForAssignment attempts to prepare a src value for assignment into dstType.
//
// - Handles pointer dereference (*T → T)
// - Optionally handles value → pointer (T → *T, with auto-allocation)
func normalizeForAssignment(srcVal reflect.Value, dstType reflect.Type) (reflect.Value, bool) {
	// Dereference pointer if necessary
	if srcVal.Kind() == reflect.Ptr && !srcVal.IsNil() {
		srcVal = srcVal.Elem()
	}

	if srcVal.Type().AssignableTo(dstType) {
		return srcVal, true
	}

	// Optional: support T → *T
	if dstType.Kind() == reflect.Ptr && srcVal.Type().AssignableTo(dstType.Elem()) {
		ptr := reflect.New(dstType.Elem())
		ptr.Elem().Set(srcVal)
		return ptr, true
	}

	return reflect.Value{}, false
}
