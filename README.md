# 🎨 HueBase

A CLI app and library designed to enhance the generation and creation of color schemes across a diverse range of applications. Three goals:

* CLI tool to view and convert colorschemes for terminals, IDEs, etc.
* A library of popular colorschemes across various tools (think [Repology](https://repology.org/) but for themes!)
* An OS level tool to manage theme profiles concurrently between apps

## Why?

* Your favorite color schemes should be available everywhere!
* Crafting custom color schemes can be complex and cumbersome, especially for multiple apps
* Many apps have great color schemes that aren't available elsewhere (one of my [long time favorites](https://github.com/muukii/jackhammer-syntax) for the retired code editor [Atom](https://github.com/atom/atom))

## Considerations

Obviously, converting color schemes will not be perfect. Some formats contain information that others don't, and thus most conversions will be lossy. Additionally, colors are can have slight to significant differences in interpretation between formats and even themes of the same format. Given all this, `huebase` simply aims to get you around 90% of the way there with only 10% of the effort of creating a scheme from scratch.

## Conversion Schema

```text
Missing field warnings
for fields expected,
but found in file
(potentially convertible,
but not found)
v
Original theme --> Abstract Theme --> New Theme
                ^                  ^
    Ignored field       Unused field
    warnings            warnings
    for fields never    for fields not
    even converted      used during conversion
    (unnecessary fields) (unconvertible fields)
```
