package adapter

import (
	"fmt"
	"reflect"
	"testing"
)

// End-to-end adapter transformation test
// FIXME: this test is bad!
func TestAdapterRoundTripCompatibility(t *testing.T) {
	for _, srcAdapter := range Adapters {
		t.Run(fmt.Sprintf("From_%s", srcAdapter.Name()), func(t *testing.T) {
			fillDummyScheme(srcAdapter)

			for _, dstAdapter := range Adapters {
				t.Run(fmt.Sprintf("From_%s_To_%s", srcAdapter.Name(), dstAdapter.Name()), func(t *testing.T) {

					// Adapt from src → dst
					err := adaptScheme(srcAdapter, dstAdapter)
					if err != nil {
						t.Fatalf("adaptScheme failed: %v", err)
					}

					// Render destination to string
					rendered, err := RenderAdapterToString(dstAdapter)
					if err != nil {
						t.Fatalf("RenderAdapterToString failed: %v", err)
					}
					if rendered == "" {
						t.Fatal("rendered string is empty")
					}

					// Parse into new destination instance
					dstParsed := newAdapterInstance(dstAdapter)
					err = dstParsed.FromString(rendered)
					if err != nil {
						t.Fatalf("FromString failed: %v", err)
					}

					// Similarity check (can be relaxed if needed)
					sim := FieldSimilarity(dstAdapter, dstParsed)
					if sim < 1.0 {
						t.Errorf("Round-trip mismatch: similarity %.2f < 1.0", sim)
						fmt.Println("Original:")
						printSchemeContents("DST", dstAdapter)
						fmt.Println("Parsed:")
						printSchemeContents("DST_PARSED", dstParsed)
					}
				})
			}
		})
	}
}

// Helper to allocate new adapter instance
func newAdapterInstance(ad Adapter) Adapter {
	typ := reflect.TypeOf(ad).Elem()
	return reflect.New(typ).Interface().(Adapter)
}
