package objectmap_test

import (
	"reflect"
	"testing"

	"github.com/da-luce/paletteport/internal/objectmap"
)

type SrcNested struct {
	X int
	Y string
}

type Src struct {
	A string
	B int
	N SrcNested
}

type DstNested struct {
	Y string `mapfrom:"N.Y"`
}

type Dst struct {
	A     string
	B     int
	Ren   string `mapfrom:"A"`
	Nest  DstNested
	Extra string // will be unused
}

func TestMapFrom_BasicAndTagged(t *testing.T) {
	src := &Src{
		A: "hello",
		B: 42,
		N: SrcNested{
			X: 99,
			Y: "nested",
		},
	}

	dst := &Dst{}

	var unusedSrcFields [][]string
	var unusedDstFields [][]string

	err := objectmap.MapFrom(
		src, dst,
		func(path []string, val reflect.Value) {
			unusedSrcFields = append(unusedSrcFields, path)
		},
		func(path []string, val reflect.Value) {
			unusedDstFields = append(unusedDstFields, path)
		},
		"mapfrom",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Basic and tagged fields should map correctly
	if dst.A != "hello" {
		t.Errorf("dst.A expected 'hello', got %q", dst.A)
	}
	if dst.B != 42 {
		t.Errorf("dst.B expected 42, got %d", dst.B)
	}
	if dst.Ren != "hello" {
		t.Errorf("dst.Ren expected 'hello' from mapfrom tag, got %q", dst.Ren)
	}
	if dst.Nest.Y != "nested" {
		t.Errorf("dst.Nest.Y expected 'nested' from nested mapfrom tag, got %q", dst.Nest.Y)
	}

	// Check that Extra was unused
	expectedUnusedDst := [][]string{{"Extra"}}
	if !equalPathSlices(unusedDstFields, expectedUnusedDst) {
		t.Errorf("unexpected unusedDstFields: got %v, want %v", unusedDstFields, expectedUnusedDst)
	}

	// Check that N.X in src was unused
	expectedUnusedSrc := [][]string{{"N", "X"}}
	if !equalPathSlices(unusedSrcFields, expectedUnusedSrc) {
		t.Errorf("unexpected unusedSrcFields: got %v, want %v", unusedSrcFields, expectedUnusedSrc)
	}
}

// Helper to compare [][]string slices
func equalPathSlices(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// -----------------------------------------------------------------------------
// pointer tests
// -----------------------------------------------------------------------------

type SrcPtrNested struct {
	Message *string
	Count   *int
}

type SrcPtr struct {
	Name  *string
	Age   *int
	Inner *SrcPtrNested
}

type DstPtrNested struct {
	Message string
	Count   int
}

type DstPtr struct {
	Name  string
	Age   int
	Inner DstPtrNested
	Note  *string // unused
}

func TestMapFrom_PointerFields(t *testing.T) {
	msg := "hi"
	count := 7
	name := "Alice"
	age := 30

	src := &SrcPtr{
		Name: &name,
		Age:  &age,
		Inner: &SrcPtrNested{
			Message: &msg,
			Count:   &count,
		},
	}

	var unusedSrc [][]string
	var unusedDst [][]string

	dst := &DstPtr{}

	err := objectmap.MapFrom(
		src, dst,
		func(path []string, val reflect.Value) {
			unusedSrc = append(unusedSrc, path)
		},
		func(path []string, val reflect.Value) {
			unusedDst = append(unusedDst, path)
		},
		"mapfrom",
	)

	if err != nil {
		t.Fatalf("MapFrom error: %v", err)
	}

	if dst.Name != "Alice" {
		t.Errorf("expected dst.Name = 'Alice', got %q", dst.Name)
	}
	if dst.Age != 30 {
		t.Errorf("expected dst.Age = 30, got %d", dst.Age)
	}
	if dst.Inner.Message != "hi" {
		t.Errorf("expected dst.Inner.Message = 'hi', got %q", dst.Inner.Message)
	}
	if dst.Inner.Count != 7 {
		t.Errorf("expected dst.Inner.Count = 7, got %d", dst.Inner.Count)
	}

	// Should detect unused DstPtr.Note and nothing unused in src
	expectedUnusedDst := [][]string{{"Note"}}
	if !equalPathSlices(unusedDst, expectedUnusedDst) {
		t.Errorf("unexpected unusedDst: got %v, want %v", unusedDst, expectedUnusedDst)
	}

	if len(unusedSrc) != 0 {
		t.Errorf("expected no unused src fields, got: %v", unusedSrc)
	}
}
