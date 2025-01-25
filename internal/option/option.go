package option

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Option represents a value that may or may not exist.
type Option[T any] struct {
	value T
	isSet bool
}

// Some creates an Option with a value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, isSet: true}
}

// None creates an Option without a value.
func None[T any]() Option[T] {
	var zeroValue T // Zero value for type T
	return Option[T]{value: zeroValue, isSet: false}
}

// IsSome returns true if the Option has a value.
func (o Option[T]) IsSome() bool {
	return o.isSet
}

// IsNone returns true if the Option does not have a value.
func (o Option[T]) IsNone() bool {
	return !o.isSet
}

// Unwrap returns the value if it exists, or panics if the Option is None.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("attempted to unwrap a None value")
	}
	return o.value
}

// UnwrapOr returns the value if it exists, or a default value otherwise.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.IsNone() {
		return defaultValue
	}
	return o.value
}

// Map applies a function to the value if it exists, returning a new Option.
func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.IsSome() {
		return Some(f(o.value))
	}
	return None[T]()
}

// String implements the Stringer interface for debugging purposes.
func (o Option[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

// UnmarshalJSON recursively unmarshals JSON into the Option.
func (o *Option[T]) UnmarshalJSON(data []byte) error {

	// Handle "null" as None
	if string(data) == "null" {
		o.isSet = false
		return nil
	}

	// Unmarshal into the value field
	err := json.Unmarshal(data, &o.value)
	if err != nil {
		return err
	}

	// Mark as set if unmarshaling succeeded
	o.isSet = true
	return nil
}

// UnmarshalYAML recursively unmarshals YAML into the Option.
func (o *Option[T]) UnmarshalYAML(node *yaml.Node) error {

	// Handle null nodes as None
	if node.Tag == "!!null" {
		o.isSet = false
		return nil
	}

	// Unmarshal into the value field
	var value T
	if err := node.Decode(&value); err != nil {
		return err
	}

	o.value = value
	o.isSet = true
	return nil
}

func (o *Option[T]) UnmarshalText(data []byte) error {

	if string(data) == "null" {
		o.isSet = false

		return nil
	}
	// Create a new instance of T
	var value T

	// Check if T implements the UnmarshalText interface
	if unmarshaler, ok := any(&value).(interface{ UnmarshalText([]byte) error }); ok {
		if err := unmarshaler.UnmarshalText(data); err != nil {
			return fmt.Errorf("failed to unmarshal text into Option[%T]: %w", value, err)
		}
	} else {
		// Fallback: Try to unmarshal as JSON or raw string
		fmt.Println("T does not implement UnmarshalText; falling back...")
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("failed to unmarshal text into Option[%T]: %w", value, err)
		}
	}

	// Set the value and mark as set
	o.value = value
	o.isSet = true

	fmt.Printf("Unmarshaled value: %+v\n", o.value)
	return nil
}

// MarshalJSON marshals the Option into JSON.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// MarshalYAML marshals the Option into YAML.
func (o Option[T]) MarshalYAML() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.value, nil
}

// MarshalTOML marshals the Option into TOML.
func (o Option[T]) MarshalTOML() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%v", o.value)), nil
}
