package objectmap_test

import (
	"reflect"
	"testing"

	"github.com/da-luce/paletteport/internal/objectmap"
)

type TestStruct struct {
	A string
	B int
	C NestedStruct
}

type NestedStruct struct {
	D bool
}

func TestTraverseStructDFS_SimpleStruct(t *testing.T) {
	s := TestStruct{
		A: "hello",
		B: 42,
		C: NestedStruct{D: true},
	}

	visited := map[string]bool{}
	visitFunc := func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		// Join path by dot
		pathStr := ""
		for i, p := range fullPath {
			if i > 0 {
				pathStr += "."
			}
			pathStr += p
		}
		visited[pathStr] = true
		return true // continue traversal
	}

	objectmap.TraverseStructDFS(s, visitFunc)

	// Expect all top-level fields visited
	if !visited["A"] || !visited["B"] || !visited["C"] {
		t.Errorf("Expected fields A, B, C to be visited, got %v", visited)
	}
}

func TestTraverseStructDFS_PointerToStruct(t *testing.T) {
	s := &TestStruct{
		A: "world",
		B: 123,
		C: NestedStruct{D: false},
	}

	visited := map[string]bool{}
	visitFunc := func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		pathStr := ""
		for i, p := range fullPath {
			if i > 0 {
				pathStr += "."
			}
			pathStr += p
		}
		visited[pathStr] = true
		return true
	}

	objectmap.TraverseStructDFS(s, visitFunc)

	if !visited["A"] || !visited["B"] || !visited["C"] {
		t.Errorf("Expected fields A, B, C to be visited, got %v", visited)
	}
}

func TestTraverseStructDFS_Complex(t *testing.T) {

	type Inner struct {
		X int
		Y *string
	}

	type Outer struct {
		A string
		B *Inner
		C *Inner
		D *int
	}

	strVal := "hello"
	intVal := 99

	data := Outer{
		A: "outer",
		B: &Inner{
			X: 10,
			Y: &strVal,
		},
		C: nil,
		D: &intVal,
	}

	visited := map[string]bool{}

	visitFunc := func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		// Build dot-separated path
		pathStr := ""
		for i, p := range fullPath {
			if i > 0 {
				pathStr += "."
			}
			pathStr += p
		}

		visited[pathStr] = true

		// Skip traversing into pointer field C if nil
		if pathStr == "C" && value.IsNil() {
			return false
		}

		// Example: skip traversing into Outer.D pointer to int (primitive)
		if pathStr == "D" {
			return false
		}

		return true
	}

	objectmap.TraverseStructDFS(data, visitFunc)

	expectedPaths := []string{
		"A",
		"B",
		"B.X",
		"B.Y",
		"C", // visited but skipped since nil
		"D", // visited but skipped
	}

	for _, p := range expectedPaths {
		if !visited[p] {
			t.Errorf("Expected to visit %q but did not", p)
		}
	}
}
