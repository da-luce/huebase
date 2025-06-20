<h1>
  paletteport
  <img src="./media/ship.png" alt="Description" height="50" style="vertical-align: bottom;" />
</h1>

* CLI tool to view and convert colorschemes for terminals, IDEs, etc.
* A library of popular colorschemes across various tools (think [Repology](https://repology.org/) but for themes!)
* An OS level tool to manage theme profiles concurrently between apps
* Get or make theme from image

## Why?

* Your favorite color schemes should be available everywhere!
* Crafting custom color schemes can be complex and cumbersome, especially for multiple apps
* Many apps have great color schemes that aren't available elsewhere (one of my [long time favorites](https://github.com/muukii/jackhammer-syntax) for the retired code editor [Atom](https://github.com/atom/atom))
* It shouldn't be hard to find and share cool schemes

## Considerations

Obviously, converting color schemes will not be perfect. Some formats contain information that others don't, and thus most conversions will be lossy. Additionally, colors are can have slight to significant differences in interpretation between formats and even themes of the same format. Given all this, `paletteport` simply aims to get you around 90% of the way there with only 10% of the effort of creating a scheme from scratch.

## Conversion Schema

```text
Missing field warnings
for fields expected,
but found not in file
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

## Notes

Add backup conversion fields, i.e. if don't have one-to-one mapping use another
field to fill it's place.

## Other possible names

* palettr
* gotone
* goglow
* PalettePort (put color palette on a pallet?)
