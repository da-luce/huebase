package objectmap_test

import (
	"reflect"
	"testing"

	"github.com/da-luce/huebase/internal/objectmap"
)

type A struct {
	Name   string `b:"Name"`
	Age    int    `b:"Age"`
	Memory *int   `b:"MemoryB"`
	Inner  Inner  `b:"Inner"`
}

type B struct {
	Name    string
	Age     int
	MemoryB *int
	Inner   Inner
}

type Inner struct {
	Message string
}

func TestTransitiveMap_Property(t *testing.T) {
	var testMem = 50
	src := A{
		Name:   "Alice",
		Age:    30,
		Memory: &testMem,
		Inner:  Inner{Message: "Hello there "},
	}

	var mid B
	err := objectmap.MapInto(&src, &mid, nil, nil, "b")
	if err != nil {
		t.Fatalf("MapInto A → B failed: %v", err)
	}

	var end A
	err = objectmap.MapFrom(&mid, &end, nil, nil, "b")
	if err != nil {
		t.Fatalf("MapFrom B → A failed: %v", err)
	}

	if !reflect.DeepEqual(src, end) {
		t.Errorf("Structs are not equal:\na1=%+v\na2=%+v", src, end)
	}

}
