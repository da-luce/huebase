package adapter

import (
	"math/rand"
	"reflect"

	"github.com/da-luce/paletteport/internal/color"
	"github.com/da-luce/paletteport/internal/structutil"
)

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fillDummyScheme(a Adapter) {
	structutil.TraverseStructDFS(a, func(fullPath []string, field reflect.StructField, value reflect.Value) bool {
		// Skip non-pointer fields or already set pointers
		if value.Kind() != reflect.Ptr || !value.IsNil() {
			return true
		}

		fieldType := value.Type().Elem()

		switch fieldType.Name() {
		case "Color":
			randomColor := color.RandomColor()
			ptr := reflect.New(fieldType)
			ptr.Elem().Set(reflect.ValueOf(randomColor))
			value.Set(ptr)

		case "string":
			str := randomString(8)
			ptr := reflect.New(fieldType)
			ptr.Elem().Set(reflect.ValueOf(str))
			value.Set(ptr)

		default:
			// Generic pointer to zero value
			ptr := reflect.New(fieldType)
			value.Set(ptr)
		}

		return true
	})
}
