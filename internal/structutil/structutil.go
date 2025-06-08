package structutil

import "reflect"

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

// Provide wrapper for traverseDFS
func TraverseStructDFS(data any,
	visitFunc func(fullPath []string, field reflect.StructField, value reflect.Value) bool,
) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	traverseDFS(val, nil, visitFunc)
}
