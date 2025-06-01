package objectmap

import (
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------
// Basic MapField
// -----------------------------------------------------------------------------

type User struct {
	Name     string
	Email    string `mapto:"ContactEmail"`
	Age      int
	Location string `mapto:"Office"`
}

type Employee struct {
	Name         string
	ContactEmail string
	ID           string
	Office       string
}

func TestMapFields_CorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{ID: "E123"}

	err := MapFields(&user, &emp, nil, nil, "mapto")
	if err != nil {
		t.Fatalf("MapFields returned error: %v", err)
	}

	fields := []struct {
		got, want, label string
	}{
		{emp.Name, user.Name, "Name"},
		{emp.ContactEmail, user.Email, "ContactEmail"},
		{emp.Office, user.Location, "Office"},
		{emp.ID, "E123", "ID"},
	}

	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("%s mismatch: got %q, want %q", f.label, f.got, f.want)
		}
	}

}

func TestMapFields_IncorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{ID: "E123"}

	err := MapFields(&user, &emp, nil, nil, "incorrecttag")
	if err != nil {
		t.Fatalf("MapFields returned error: %v", err)
	}

	fields := []struct {
		got, want, label string
	}{
		{emp.Name, user.Name, "Name"},
		{emp.ContactEmail, emp.ContactEmail, "ContactEmail"}, // this will not map
		{emp.Office, emp.Office, "Office"},                   // this will not map
		{emp.ID, "E123", "ID"},                               // this should be the same
	}

	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("%s mismatch: got %q, want %q", f.label, f.got, f.want)
		}
	}
}

// -----------------------------------------------------------------------------
// Nested MapField
// -----------------------------------------------------------------------------

type Source struct {
	A string `mapto:"N.A"`
	B string `mapto:"N.B"`
	C string
	D string
}

type Nested struct {
	A string
	B string
}

type Dest struct {
	N Nested
	C string
	D string
}

func TestMapFields_BasicMapping(t *testing.T) {
	src := &Source{
		A: "foo",
		B: "bar",
		C: "baz",
		D: "qux",
	}
	dst := &Dest{}

	var unusedSrcFields [][]string
	var unusedDstFields [][]string

	err := MapFields(src, dst,
		func(path []string, val reflect.Value) {
			unusedSrcFields = append(unusedSrcFields, path)
		},
		func(path []string, val reflect.Value) {
			unusedDstFields = append(unusedDstFields, path)
		},
		"mapto",
	)
	if err != nil {
		t.Fatalf("MapFields returned error: %v", err)
	}

	if dst.N.A != "foo" || dst.N.B != "bar" || dst.C != "baz" || dst.D != "qux" {
		t.Errorf("unexpected destination values: %+v", dst)
	}

	if len(unusedSrcFields) != 0 {
		t.Errorf("expected no unused src fields, got: %+v", unusedSrcFields)
	}

	if len(unusedDstFields) != 0 {
		t.Errorf("expected no unused dst fields, got: %+v", unusedDstFields)
	}
}

func TestMapFields_UnusedFields(t *testing.T) {
	type Src struct {
		Used   string `mapto:"Used"`
		Unused string
	}
	type Dst struct {
		Used string
		Free string
	}

	src := &Src{
		Used:   "hello",
		Unused: "skipme",
	}
	dst := &Dst{}

	var unusedSrc [][]string
	var unusedDst [][]string

	err := MapFields(src, dst,
		func(path []string, val reflect.Value) {
			unusedSrc = append(unusedSrc, path)
		},
		func(path []string, val reflect.Value) {
			unusedDst = append(unusedDst, path)
		},
		"mapto",
	)
	if err != nil {
		t.Fatalf("MapFields returned error: %v", err)
	}

	if dst.Used != "hello" {
		t.Errorf("expected Used to be 'hello', got %q", dst.Used)
	}

	expectedUnusedSrc := [][]string{{"Unused"}}
	expectedUnusedDst := [][]string{{"Free"}}

	if !reflect.DeepEqual(unusedSrc, expectedUnusedSrc) {
		t.Errorf("unexpected unused src fields: got %v, want %v", unusedSrc, expectedUnusedSrc)
	}

	if !reflect.DeepEqual(unusedDst, expectedUnusedDst) {
		t.Errorf("unexpected unused dst fields: got %v, want %v", unusedDst, expectedUnusedDst)
	}
}

// -----------------------------------------------------------------------------
// Generic helper struct
// -----------------------------------------------------------------------------

type Struct struct {
	A string
	B string
	C string
	D NestedStruct
}

type NestedStruct struct {
	E string
	F string
	G SubNestedStruct
}

type SubNestedStruct struct {
	H string
	I SubSubNestedStruct
}

type SubSubNestedStruct struct {
	J SubSubSubNestedStruct
}

type SubSubSubNestedStruct struct {
	I string
}

func newTestStruct() Struct {
	return Struct{
		A: "a",
		B: "b",
		C: "c",
		D: NestedStruct{
			E: "e",
			F: "f",
			G: SubNestedStruct{
				H: "h",
				I: SubSubNestedStruct{
					J: SubSubSubNestedStruct{
						I: "deep",
					},
				},
			},
		},
	}
}

// -----------------------------------------------------------------------------
// hasNestedFields
// -----------------------------------------------------------------------------

func TestHasNestedFieldSlice(t *testing.T) {

	testStruct := newTestStruct()

	tests := []struct {
		name      string
		path      []string
		wantFound bool
		wantValue any
	}{
		{
			name:      "Top level field A",
			path:      []string{"A"},
			wantFound: true,
			wantValue: "a",
		},
		{
			name:      "Nested field D.E",
			path:      []string{"D", "E"},
			wantFound: true,
			wantValue: "e",
		},
		{
			name:      "Deeply nested field D.G.I.J.I",
			path:      []string{"D", "G", "I", "J", "I"},
			wantFound: true,
			wantValue: "deep",
		},
		{
			name:      "Non-existent field D.G.X",
			path:      []string{"D", "G", "X"},
			wantFound: false,
		},
		{
			name:      "Invalid middle field D.X.J",
			path:      []string{"D", "X", "J"},
			wantFound: false,
		},
		{
			name:      "Partially correct path D.G.I.J.X",
			path:      []string{"D", "G", "I", "J", "X"},
			wantFound: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			found, _, val := hasNestedFieldSlice(reflect.ValueOf(testStruct), tc.path)
			if found != tc.wantFound {
				t.Errorf("expected found=%v, got %v", tc.wantFound, found)
			}
			if found && val.Interface() != tc.wantValue {
				t.Errorf("expected value=%v, got %v", tc.wantValue, val.Interface())
			}
		})
	}
}

// -----------------------------------------------------------------------------
// setNestedFields
// -----------------------------------------------------------------------------

func TestSetNestedField(t *testing.T) {
	tests := []struct {
		name      string
		path      []string
		newValue  any
		wantErr   bool
		verifyVal func(s Struct) bool
	}{
		{
			name:     "Set top-level field A",
			path:     []string{"A"},
			newValue: "updated A",
			wantErr:  false,
			verifyVal: func(s Struct) bool {
				return s.A == "updated A"
			},
		},
		{
			name:     "Set nested field D.E",
			path:     []string{"D", "E"},
			newValue: "updated E",
			wantErr:  false,
			verifyVal: func(s Struct) bool {
				return s.D.E == "updated E"
			},
		},
		{
			name:     "Set deeply nested field D.G.I.J.I",
			path:     []string{"D", "G", "I", "J", "I"},
			newValue: "deep value",
			wantErr:  false,
			verifyVal: func(s Struct) bool {
				return s.D.G.I.J.I == "deep value"
			},
		},
		{
			name:     "Set with invalid path",
			path:     []string{"D", "X", "Z"},
			newValue: "won't work",
			wantErr:  true,
			verifyVal: func(s Struct) bool {
				return true // doesn't matter, should error
			},
		},
		{
			name:     "Set with incompatible type",
			path:     []string{"A"},
			newValue: 123, // int, not string
			wantErr:  true,
			verifyVal: func(s Struct) bool {
				return s.A != "123" // should not change
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := Struct{}
			err := setNestedField(reflect.ValueOf(&s), tc.path, reflect.ValueOf(tc.newValue))
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
			if err == nil && !tc.verifyVal(s) {
				t.Errorf("field was not correctly set for path %v", tc.path)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// traverseFields
// -----------------------------------------------------------------------------

func TestTraverseFields_AllFields(t *testing.T) {
	var visited []string

	testStruct := newTestStruct()

	traverseFields(reflect.ValueOf(testStruct), nil, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		visited = append(visited, joinPath(fullPath))
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

	traverseFields(reflect.ValueOf(s), nil, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		visited = append(visited, joinPath(fullPath))
		// only recurse into top-level field "D"
		return joinPath(fullPath) == "D"
	})

	expected := []string{"A", "B", "C", "D", "D.E", "D.F", "D.G"}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf("visited paths = %v\nexpected = %v", visited, expected)
	}
}

func TestTraverseFields_NonStructInput(t *testing.T) {
	called := false

	traverseFields(reflect.ValueOf("not a struct"), nil, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		called = true
		return false
	})

	if called {
		t.Error("visitFunc should not be called for non-struct input")
	}
}

// -----------------------------------------------------------------------------
// splitPath
// -----------------------------------------------------------------------------

func TestSplitPath(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"A", []string{"A"}},
		{"A.B", []string{"A", "B"}},
		{"A.B.C", []string{"A", "B", "C"}},
		{"..", []string{"", "", ""}},
		{"A..B", []string{"A", "", "B"}},
	}

	for _, tc := range tests {
		result := splitPath(tc.input)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("splitPath(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

// -----------------------------------------------------------------------------
// joinPath
// -----------------------------------------------------------------------------

func TestJoinPath(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"A"}, "A"},
		{[]string{"A", "B"}, "A.B"},
		{[]string{"A", "B", "C"}, "A.B.C"},
		{[]string{"", "", ""}, ".."},
		{[]string{"A", "", "B"}, "A..B"},
	}

	for _, tc := range tests {
		result := joinPath(tc.input)
		if result != tc.expected {
			t.Errorf("joinPath(%v) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}
