package structutil_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/da-luce/paletteport/internal/structutil"
)

type BasicTestStruct struct {
	A string
	B int
	C BasicNestedStruct
}

type BasicNestedStruct struct {
	D bool
}

func TestTraverseStructDFS_SimpleStruct(t *testing.T) {
	s := BasicTestStruct{
		A: "hello",
		B: 42,
		C: BasicNestedStruct{D: true},
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

	structutil.TraverseStructDFS(s, visitFunc)

	// Expect all top-level fields visited
	if !visited["A"] || !visited["B"] || !visited["C"] {
		t.Errorf("Expected fields A, B, C to be visited, got %v", visited)
	}
}

func TestTraverseStructDFS_PointerToStruct(t *testing.T) {
	s := &BasicTestStruct{
		A: "world",
		B: 123,
		C: BasicNestedStruct{D: false},
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

	structutil.TraverseStructDFS(s, visitFunc)

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

	structutil.TraverseStructDFS(data, visitFunc)

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

// -----------------------------------------------------------------------------
// traverseFields
// -----------------------------------------------------------------------------

func TestTraverseFields_AllFields(t *testing.T) {
	var visited []string

	testStruct := newTestStruct()

	structutil.TraverseStructDFS(testStruct, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		visited = append(visited, strings.Join(fullPath, "."))
		return true // always recurse
	})

	expected := []string{
		"A", "B", "C",
		"D", "D.E", "D.F",
		"D.G", "D.G.H",
		"D.G.I", "D.G.I.J", "D.G.I.J.I",
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf("visited paths = %v\nexpected = %v", visited, expected)
	}
}

func TestTraverseFields_SkipNested(t *testing.T) {
	var visited []string

	s := Struct{}

	structutil.TraverseStructDFS(s, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		visited = append(visited, strings.Join(fullPath, "."))
		// only recurse into top-level field "D"
		return strings.Join(fullPath, ".") == "D"
	})

	expected := []string{"A", "B", "C", "D", "D.E", "D.F", "D.G"}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf("visited paths = %v\nexpected = %v", visited, expected)
	}
}

func TestTraverseFields_NonStructInput(t *testing.T) {
	called := false

	structutil.TraverseStructDFS("not a struct", func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		called = true
		return false
	})

	if called {
		t.Error("visitFunc should not be called for non-struct input")
	}
}
