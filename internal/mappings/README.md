# Decoding the Madness

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
