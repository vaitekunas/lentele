# lentele

`lentele` is a no-thrills ascii-table builder

`lentele` is thread-safe, i.e. a table can be accessed by several goroutines
simultaneously.

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

# TODO

[] Add `Template.PrintExample` to show an example table
