# MapFields

`MapFields` is a generic Go function to copy matching fields from a source struct to a destination struct. Both `src` and `dst` must be pointers to structs. Because this package uses dynamic reflection, it must be tested water tight to catch any mapping errors or unexpected behaviors.

## Features

- By default, matches fields by name.
- Supports mapping fields using struct tags (`mapto:"FieldName"`) on the source struct to rename destination fields.
- Provides callbacks for handling unused source or destination fields.
- Supports nested fields via dot-separated paths.

## Usage

```go
err := MapFields(&srcStruct, &dstStruct,
    func(srcPath []string, srcVal reflect.Value) {
        // Called for each source field not mapped to destination
    },
    func(dstPath []string, dstVal reflect.Value) {
        // Called for each destination field not mapped from source
    },
    "mapto", // name of struct tag to use for field mapping
)
if err != nil {
    // handle error
}
