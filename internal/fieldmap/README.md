# fieldmap

A Go package for flexible one-to-many field mappings between two structs.

## Example

Suppose we have the following structs

```go
type Source struct {
    A *string
    B *string
}

type Dest struct {
    Group1 struct {
        A *string
        B *string
    }
    Group2 struct {
        C *string
    }
}
```

And want te following mapping:

```yaml
A:
  - Group1.A
  - Group2.C
B:
  - Group1.B
```

This means:

* The value in `Source.A` maps to both `Dest.Group1.A` and `Dest.Group2.C`
* The value in `Source.B` maps to `Dest.Group1.B`

Using this mapping, when you run:

```go
err := ApplySourceToDestMapping(&source, &dest, mapping)
```

The function will:

* Read the value from `source.A`, and set it to both `dest.Group1.A` and `dest.Group2.C`.
* Read the value from `source.B`, and set it to `dest.Group1.B`.

Conversely, running:

```go
err := ApplyDestToSourceMapping(&dest, &source, mapping)
```

will:

* For each destination field mapped to `Source.A` (i.e., `Group1.A` and `Group2.C`), it will pick the first non-nil value found and assign it back to source.A.
* For `Source.B`, it will look at `dest.Group1.B` and assign its value to `source.B` if non-nil.

This mapping structure allows a one-to-many relationship from source to destination fields, and a many-to-one relationship in reverse, giving flexibility in how your data is synchronized across different struct shapes.

## Why store this data in YAML and then convert it to structs?

* Storing mappings in struct tags quickly becomes unreadable, especially when a color field maps to multiple abstract fields. Additionally, in quite a few languages (including Go), struct tags are not compiler checked, so we don't gain much type safety by using tags anyways
* I'd argue that these mappings are more data than logic. Embedding them directly in code tightly couples data with the program’s logic, which reduces flexibility.
* This approach keeps the data language-agnostic—others can easily use or adapt it in different programming languages.
* Using a single, declarative mapping lets you do double duty: generate templates that output Base16-compatible themes and parse Base16 themes back into your abstract scheme representation.

## Verifying a Mapping File

1. Validate scheme struct
   * Ensure all keys in the mapping are valid scheme struct fields
   * Ensure all struct fields have a key in the mapping

2. Validate abstract struct
   * Ensure all mappings in the mapping are valid abstract struct fields
   * Ensure all struct fields have at least one value in the mapping
