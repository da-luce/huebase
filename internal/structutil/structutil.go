package structutil

import (
	"errors"
	"reflect"
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

// HasNestedFieldSlice returns (found, field metadata, field value).
// If not found, found is false and others are zero values.
func HasNestedFieldSlice(structValue reflect.Value, path []string) (bool, reflect.StructField, reflect.Value) {
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
func SetNestedField(structPtr reflect.Value, path []string, newVal reflect.Value) error {
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
