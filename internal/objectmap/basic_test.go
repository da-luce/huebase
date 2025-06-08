package objectmap_test

import (
	"reflect"
	"testing"

	"github.com/da-luce/paletteport/internal/objectmap"
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

	err := objectmap.MapInto(&user, &emp, nil, nil, "map")
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

	err := objectmap.MapInto(&user, &emp, nil, nil, "incorrect tag")
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

	err := objectmap.MapFrom(&emp, &user, nil, nil, "map")
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
	err := objectmap.MapFrom(&emp, &user, nil, nil, "incorrecttag")
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

	err := objectmap.MapInto(src, dst,
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

	err := objectmap.MapInto(src, dst,
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

	err := objectmap.MapFrom(src, dst,
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

	err := objectmap.MapFrom(src, dst,
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
