package structutil_test

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
