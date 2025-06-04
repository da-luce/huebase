package objectmap

import (
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------
// Basic mapping
// -----------------------------------------------------------------------------

type User struct {
	Name     string
	Email    string `map:"ContactEmail"`
	Age      int
	Location string `map:"Office"`
}

type Employee struct {
	Name         string
	ContactEmail string
	ID           string
	Office       string
}

func TestMapInto_CorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{ID: "E123"}

	err := MapInto(&user, &emp, nil, nil, "map")
	if err != nil {
		t.Fatalf("MapInto returned error: %v", err)
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

func TestMapInto_IncorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{ID: "E123"}

	err := MapInto(&user, &emp, nil, nil, "incorrect tag")
	if err != nil {
		t.Fatalf("MapInto returned error: %v", err)
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

// Run in opposite direction
func TestMapFrom_CorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{Name: "Bob", ContactEmail: "bob@example.com", ID: "E123", Office: "Remote"}

	err := MapFrom(&emp, &user, nil, nil, "map")
	if err != nil {
		t.Fatalf("mapFrom returned error: %v", err)
	}

	fields := []struct {
		got, want interface{}
		label     string
	}{
		{user.Name, emp.Name, "Name"},
		{user.Email, emp.ContactEmail, "Email"},
		{user.Location, emp.Office, "Location"},
		{user.Age, 30, "Age"},
	}

	for _, f := range fields {
		if f.got != f.want {
			t.Errorf("%s mismatch: got %q, want %q", f.label, f.got, f.want)
		}
	}
}

func TestMapFrom_IncorrectTag(t *testing.T) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30, Location: "HQ"}
	emp := Employee{Name: "Bob", ContactEmail: "bob@example.com", ID: "E123", Office: "Remote"}

	// Using a wrong tag name means no mapping happens for tagged fields
	err := MapFrom(&emp, &user, nil, nil, "incorrecttag")
	if err != nil {
		t.Fatalf("mapFrom returned error: %v", err)
	}

	fields := []struct {
		got, want interface{}
		label     string
	}{
		{user.Name, emp.Name, "Name"},
		{user.Email, "alice@example.com", "Email"}, // no mapping, remains original
		{user.Location, "HQ", "Location"},          // no mapping, remains original
		{user.Age, 30, "Age"},                      // no mapping, remains original
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

type NestedDst struct {
	A string `mapfrom:"A"`
	B string `mapfrom:"B"`
}

type Dest struct {
	N NestedDst
	C string
	D string
}

func TestMapInto_BasicMapping(t *testing.T) {
	src := &Source{
		A: "foo",
		B: "bar",
		C: "baz",
		D: "qux",
	}
	dst := &Dest{}

	var unusedSrcFields [][]string
	var unusedDstFields [][]string

	err := MapInto(src, dst,
		func(path []string, val reflect.Value) {
			unusedSrcFields = append(unusedSrcFields, path)
		},
		func(path []string, val reflect.Value) {
			unusedDstFields = append(unusedDstFields, path)
		},
		"mapto",
	)
	if err != nil {
		t.Fatalf("MapInto returned error: %v", err)
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

func TestMapInto_UnusedFields(t *testing.T) {
	type Src struct {
		Used   string `map:"Used"`
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

	err := MapInto(src, dst,
		func(path []string, val reflect.Value) {
			unusedSrc = append(unusedSrc, path)
		},
		func(path []string, val reflect.Value) {
			unusedDst = append(unusedDst, path)
		},
		"map",
	)
	if err != nil {
		t.Fatalf("MapInto returned error: %v", err)
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

func TestMapFrom_BasicMapping(t *testing.T) {
	src := &Source{
		A: "foo",
		B: "bar",
		C: "baz",
		D: "qux",
	}
	dst := &Dest{}

	var unusedSrcFields [][]string
	var unusedDstFields [][]string

	err := MapFrom(src, dst,
		func(path []string, val reflect.Value) {
			unusedSrcFields = append(unusedSrcFields, path)
		},
		func(path []string, val reflect.Value) {
			unusedDstFields = append(unusedDstFields, path)
		},
		"mapfrom",
	)
	if err != nil {
		t.Fatalf("mapFrom returned error: %v", err)
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

func TestMapFrom_UnusedFields(t *testing.T) {
	type Src struct {
		Used   string
		Unused string
	}
	type Dst struct {
		Used string `map:"Used"`
		Free string
	}

	src := &Src{
		Used:   "hello",
		Unused: "skipme",
	}
	dst := &Dst{}

	var unusedSrc [][]string
	var unusedDst [][]string

	err := MapFrom(src, dst,
		func(path []string, val reflect.Value) {
			unusedSrc = append(unusedSrc, path)
		},
		func(path []string, val reflect.Value) {
			unusedDst = append(unusedDst, path)
		},
		"map",
	)
	if err != nil {
		t.Fatalf("mapFrom returned error: %v", err)
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
	H *string
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
				H: nil, // Testing nil is important!
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
		// Important regression test!
		{
			name:      "Nested field D.G.H",
			path:      []string{"D", "G", "H"},
			wantFound: true,
			wantValue: nil,
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
			if found {
				if val.Kind() == reflect.Ptr && val.IsNil() {
					// ok
				} else if val.Interface() != tc.wantValue {
					t.Errorf("expected value=%v, got %v", tc.wantValue, val.Interface())
				}
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

	traverseDFS(reflect.ValueOf(testStruct), nil, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
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

	traverseDFS(reflect.ValueOf(s), nil, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
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

	traverseDFS(reflect.ValueOf("not a struct"), nil, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
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
