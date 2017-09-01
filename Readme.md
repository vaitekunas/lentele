# lentele

[![Build Status](https://travis-ci.org/vaitekunas/lentele.svg?branch=master)](https://travis-ci.org/vaitekunas/lentele)

`lentele` is a mix between a very lite DataFrame interface and a
no-thrills ascii-table builder. Its primary use is in cli-applications to display
structured data (logs, cluster information and so on).

See [the journal repo](https://github.com/vaitekunas/journal) for an implemented
example.

![vortex_cluster](https://user-images.githubusercontent.com/3492398/29931653-1b68a308-8e72-11e7-9309-fbd092286945.png)

`lentele` is thread-safe, i.e. a table can be accessed by several goroutines
simultaneously.

`Go` documentation is available on [godoc](https://godoc.org/github.com/vaitekunas/lentele).

# Usage

```go
// Define a new table
table := New()
table.AddHeaders([]string{"ID", "Client", "Amount", "On time"}, []string{"%.3d","%s","%7.2f€"})

// Insert data
table.AddRow("").Insert(1, "Acme", 29123.23, true)
table.AddRow("").Insert(2, "", 12211.12, true)
table.AddRow("Typo").Insert(3, "", 7781.2, false)

// Make some corrections
table.GetRow(2).Value("Client").Set("Ecorp")
table.GetRowByName("Typo").Change(3, "", 77812.2, true)

// Render to stdout (different)
table.Render(os.Stdout, true, []string{"Client","On time"})

// Render to stdout (all columns)
table.Render(os.Stdout, true, []string{})
```
```shell

```

## Modifications

Modifications can be chained:

```Go
table := New("Col X", "Col Y", "Col Z")

table.AddRow("chain").
      Insert("xcol","ycol","zcol").
      Modify(high, "Col X","Col Z").
      Modify(low, "Col Y")
```

## Table templates

Currently the library provides following table templates:

```go
for _, tmpl := range []func() Template{tmplClassic, tmplSmooth, tmplModern} {
  template := tmpl()
  template.PrintExample(os.Stdout)
}
```

A good source of characters that can be used in creating database designs can be
found [here](https://en.wikipedia.org/wiki/Box-drawing_character).

# TODO

- [ ] Increase test coverage
- [ ] Add `Template.PrintExample` to show an example table
