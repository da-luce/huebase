package fieldmap

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// getFieldByPath retrieves a reflect.Value for a field given a dotted path.
func getFieldByPath(v reflect.Value, path string) (reflect.Value, error) {
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("non-struct encountered in path: %s", path)
		}
		v = v.FieldByName(part)
		if !v.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %s not found in path: %s", part, path)
		}
	}
	return v, nil
}

// setFieldByPath sets a field value by its dotted path. Assumes types match.
func setFieldByPath(v reflect.Value, path string, val reflect.Value) error {
	field, err := getFieldByPath(v, path)
	if err != nil {
		return err
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set field at path: %s", path)
	}
	if val.Type().AssignableTo(field.Type()) {
		field.Set(val)
		return nil
	}
	return fmt.Errorf("type mismatch for path %s: %s cannot be assigned to %s", path, val.Type(), field.Type())
}

func ApplySourceToDestMapping(source, dest interface{}, mapping map[string][]string) error {
	srcVal := reflect.ValueOf(source)
	destVal := reflect.ValueOf(dest)
	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return errors.New("source and dest must be pointers to structs")
	}

	for srcPath, destPaths := range mapping {
		srcField, err := getFieldByPath(srcVal, srcPath)
		if err != nil {
			return fmt.Errorf("error reading source path %s: %w", srcPath, err)
		}
		for _, destPath := range destPaths {
			err := setFieldByPath(destVal, destPath, srcField)
			if err != nil {
				return fmt.Errorf("error writing dest path %s: %w", destPath, err)
			}
		}
	}
	return nil
}

func ApplyDestToSourceMapping(dest, source interface{}, mapping map[string][]string) error {
	srcVal := reflect.ValueOf(source)
	destVal := reflect.ValueOf(dest)
	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return errors.New("source and dest must be pointers to structs")
	}

	for srcPath, destPaths := range mapping {
		for _, destPath := range destPaths {
			destField, err := getFieldByPath(destVal, destPath)
			if err != nil {
				continue
			}
			// Skip if value is nil or zero
			if (destField.Kind() == reflect.Ptr && destField.IsNil()) || reflect.DeepEqual(destField.Interface(), reflect.Zero(destField.Type()).Interface()) {
				continue
			}
			// Set first non-nil value found
			err = setFieldByPath(srcVal, srcPath, destField)
			if err != nil {
				return fmt.Errorf("error writing source path %s: %w", srcPath, err)
			}
			break
		}
	}
	return nil
}

// hasFieldByPath checks if a struct type t has a nested field path like "SpecialColors.Background".
// Returns true if the full path exists, false otherwise.
func hasFieldByPath(t reflect.Type, path string) bool {
	parts := strings.Split(path, ".")

	curr := t
	for _, p := range parts {
		if curr.Kind() == reflect.Ptr {
			curr = curr.Elem()
		}
		if curr.Kind() != reflect.Struct {
			return false
		}

		field, ok := curr.FieldByName(p)
		if !ok {
			return false
		}
		curr = field.Type
	}
	return true
}

func validateFieldPaths(t reflect.Type, paths []string) (invalid []string) {
	for _, path := range paths {
		if !hasFieldByPath(t, path) {
			invalid = append(invalid, path)
		}
	}
	return
}

func LoadMappingFromString(yamlStr string) (map[string][]string, error) {
	mapping := make(map[string][]string)
	err := yaml.Unmarshal([]byte(yamlStr), &mapping)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}

// verifyMappingString validates keys and values in YAML content given as a string.
// It checks that all keys exist in the source struct, and all values exist in the destination struct.
func VerifyMappingString(sourceType, destType reflect.Type, yamlStr string) error {
	mapping, err := LoadMappingFromString(yamlStr)
	if err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract keys from the mapping map
	var keys []string
	for k := range mapping {
		keys = append(keys, k)
	}

	// Validate that all keys exist in the source struct (e.g., Base16Scheme)
	invalidKeys := validateFieldPaths(sourceType, keys)
	if len(invalidKeys) > 0 {
		return fmt.Errorf("invalid keys in mapping: %v", invalidKeys)
	}

	// Flatten all values for validation against the destination struct (e.g., AbstractScheme)
	var allValues []string
	for _, valList := range mapping {
		allValues = append(allValues, valList...)
	}

	// Validate that all values exist in the destination struct
	invalidValues := validateFieldPaths(destType, allValues)
	if len(invalidValues) > 0 {
		return fmt.Errorf("invalid field paths in mapping values: %v", invalidValues)
	}

	return nil
}
