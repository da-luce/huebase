package adapter

import (
	"os"
	"testing"

	"reflect"

	"github.com/da-luce/huebase/internal/fieldmap"
)

func TestVerifyAllSchemeMappings(t *testing.T) {
	for _, adapter := range Adapters {
		data, err := os.ReadFile(adapter.MappingPath())
		if err != nil {
			t.Errorf("Failed to read %s: %v", adapter.MappingPath(), err)
			continue
		}
		err = fieldmap.VerifyMappingString(
			reflect.TypeOf(adapter).Elem(),
			reflect.TypeOf(AbstractScheme{}),
			string(data),
		)
		if err != nil {
			t.Errorf("Verification failed for %s: %v", adapter.MappingPath(), err)
		}
	}

}
