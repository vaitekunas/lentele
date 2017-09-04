# lentele [![godoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/vaitekunas/lentele) [![Build Status](https://travis-ci.org/vaitekunas/lentele.svg?branch=master)](https://travis-ci.org/vaitekunas/lentele) [![Coverage Status](https://coveralls.io/repos/github/vaitekunas/lentele/badge.svg?branch=master)](https://coveralls.io/github/vaitekunas/lentele?branch=master)

`lentele` is a mix between a very lite DataFrame interface and a no-thrills
ascii-table builder. Its primary use is to display structured data (logs,
cluster information, etc) in cli-applications.

The distributed logging facility [journal](https://github.com/vaitekunas/journal)
uses `lentele` to display log and service statistics and is a good example of
using this library.

![vortex_cluster](https://user-images.githubusercontent.com/3492398/29931653-1b68a308-8e72-11e7-9309-fbd092286945.png)

Feature overview:

* Light data manipulation capabilities - filtering, transformation and visual modification
* Dual-view rendering - a table can be rendered with/without visual modifications.
e.g. the table above, in its non-modified state, contains slices of floats, however
a modifier function turns them into utilization indicators.
* Lazy evaluation of modifications - the modification function is not executed until
render time
* Custom look and feel - `table.Render` accepts anything implementing the `lentele.Template` interface
* Thread-safe execution
* Unmarshalling of JSON objects straight into a `lentele.Table`



# Usage

A new table can be created using one of the following methods:

* `New()` - not providing any column names
* `New("COL1", "COL2")` - providing column names at creation time
* `NewFromRichJSON(src)` - loading a previously marshaled table
* `NewFromVanillaJSON(src)` - loading a json-encoded list `[{"x": 1, "y": 2},{"x": 3, "y": 4},...]`

```go
// GDP/Inflation Table
table := lentele.New("Year", "GDP growth", "Inflation")

table.AddTitle("Lithuanian GDP growth and inflation time series")
table.AddTitle("(annual growth rates in percent)")

table.AddFootnote("Source: Worldbank")

table.AddRow("").Insert(1996, 5.1499584261, 24.6181274996)
table.AddRow("").Insert(1997, 8.2932287195, 8.877890341)
table.AddRow("provide").Insert(1998, 7.4671760033, 5.0749342993)
table.AddRow("row").Insert(1999, -1.134642894, 0.7537093994)
table.AddRow("names").Insert(2000, 3.8316671405, 1.0067873303)
table.AddRow("").Insert(2001, 6.5244308754, 1.3551349535)
table.AddRow("for").Insert(2002, 6.7607495332, 0.2983425414)
table.AddRow("easy").Insert(2003, 10.538564772, -1.1457530021)
table.AddRow("filtering").Insert(2004, 6.5500830275, 1.181321743)
table.AddRow("").Insert(2005, 7.7274079181, 2.6434629364)
table.AddRow("").Insert(2006, 7.4064443553, 3.7450370211)
table.AddRow("").Insert(2007, 11.086954387, 5.7302441043)
table.AddRow("").Insert(2008, 2.6280779606, 10.9274114655)
table.AddRow("").Insert(2009, -14.81416331, 4.4515124791)
table.AddRow("").Insert(2010, 1.6398196491, 1.3191844446)
table.AddRow("").Insert(2011, 6.0493534967, 4.1303006884)
table.AddRow("").Insert(2012, 3.8349018688, 3.0899832694)
table.AddRow("").Insert(2013, 3.5068256833, 1.0474666208)
table.AddRow("").Insert(2014, 3.4950161647, 0.1037899145)
table.AddRow("").Insert(2015, 1.7785761840, -0.8841084347)
table.AddRow("").Insert(2016)

// Sequential inserts
if row, err := table.GetRow(table.GetRowCount() - 1); err == nil {
  row.Insert(2.2988432751)
  row.Insert(0.9055215537)
}

table.AddFooter().Insert("Means:", 4.17, 3.64)

// Render table
table.Render(os.Stdout, false, true, true, lentele.LoadTemplate("classic"))
```

This code snippet results in the following table:
```

                                          Lithuanian GDP growth and inflation time series
                                                 (annual growth rates in percent)

                                             ╔════════╦══════════════╦═══════════════╗
                                             ║  Year  ║  GDP growth  ║   Inflation   ║
                                             ╠════════╩══════════════╩═══════════════╣
                                             ║  1996  │ 5.1499584261 │ 24.6181274996 ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  1997  │ 8.2932287195 │  8.877890341  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  1998  │ 7.4671760033 │ 5.0749342993  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  1999  │ -1.134642894 │ 0.7537093994  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2000  │ 3.8316671405 │ 1.0067873303  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2001  │ 6.5244308754 │ 1.3551349535  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2002  │ 6.7607495332 │ 0.2983425414  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2003  │ 10.538564772 │ -1.1457530021 ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2004  │ 6.5500830275 │  1.181321743  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2005  │ 7.7274079181 │ 2.6434629364  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2006  │ 7.4064443553 │ 3.7450370211  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2007  │ 11.086954387 │ 5.7302441043  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2008  │ 2.6280779606 │ 10.9274114655 ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2009  │ -14.81416331 │ 4.4515124791  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2010  │ 1.6398196491 │ 1.3191844446  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2011  │ 6.0493534967 │ 4.1303006884  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2012  │ 3.8349018688 │ 3.0899832694  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2013  │ 3.5068256833 │ 1.0474666208  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2014  │ 3.4950161647 │ 0.1037899145  ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2015  │ 1.778576184  │ -0.8841084347 ║
                                             ╟────────┼──────────────┼───────────────╢
                                             ║  2016  │ 2.2988432751 │ 0.9055215537  ║
                                             ╚════════╧══════════════╧═══════════════╝
                                               Means:       4.17           3.64       

────────────────────
1. Source: Worldbank

```

## Modify

Table cells can be visually modified by providing a `func(interface{})interface{}`
modifier function. A good example of such modifiers are functions that add
ansi color escape sequences, e.g.:

```Go
// Transform floats
table.Transform(round, "GDP growth", "Inflation")

// Apply modifications
for i := 1; i <= table.GetRowCount(); i++ {
  if row, err := table.GetRow(i); err == nil {
    row.Modify(gdp, "GDP growth").Modify(infl, "Inflation").Modify(bold, "Year")
  }
}

// Bold header/footer
if header, err := table.GetRowByName("header"); err == nil {
  header.Modify(bold, "Year", "GDP growth", "Inflation")
}
if footer, err := table.GetRowByName("footer"); err == nil {
  footer.Modify(bold, "Year", "GDP growth", "Inflation")
}

// Render table
table.Render(os.Stdout, false, true, true, lentele.LoadTemplate("classic"))
```

here the float64 values are transformed into rounded strings and then colored
using the `github.com/fatih/color` lib:
```

                                          Lithuanian GDP growth and inflation time series
                                                 (annual growth rates in percent)

                                                ╔════════╦════════════╦═══════════╗
                                                ║  Year  ║ GDP growth ║ Inflation ║
                                                ╠════════╩════════════╩═══════════╣
                                                ║  1996  │    5.15    │   24.62   ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  1997  │    8.29    │   8.88    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  1998  │    7.47    │   5.07    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  1999  │   -1.13    │   0.75    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2000  │    3.83    │   1.01    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2001  │    6.52    │   1.36    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2002  │    6.76    │   0.30    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2003  │   10.54    │   -1.15   ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2004  │    6.55    │   1.18    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2005  │    7.73    │   2.64    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2006  │    7.41    │   3.75    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2007  │   11.09    │   5.73    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2008  │    2.63    │   10.93   ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2009  │   -14.81   │   4.45    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2010  │    1.64    │   1.32    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2011  │    6.05    │   4.13    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2012  │    3.83    │   3.09    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2013  │    3.51    │   1.05    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2014  │    3.50    │   0.10    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2015  │    1.78    │   -0.88   ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2016  │    2.30    │   0.91    ║
                                                ╚════════╧════════════╧═══════════╝
                                                  Means:      4.17        3.64     

────────────────────
1. Source: Worldbank
```

## Filter

Tables can be filtered by providing a `func(...interface{}) bool` filter function.
If the parameter `inplace` is set to true, the filtered-out rows are removed from the
table. Otherwise a new table containing the filtered rows is created.

**NB**: filtered
rows are references to the original rows, i.e. modifying them is going to change the
original table too.


Tables can also be filtered by their row names (`Table.FilterByRowNames`). In this
case a `func(string)bool` filter function must be provided.

```Go

// Filter bad years
filtered, err := table.Filter(badYear, false, false, "GDP growth", "Inflation")
if err != nil {
	log.Fatal(err.Error())
}
filtered.AddFootnote("Years with negative GDP growth and/or inflation")

// Change some value
if row, err := filtered.GetRow(2); err == nil {
	row.Change("GDP growth", "-10.54").Modify(gdp, "GDP growth")
}
filtered.AddFootnote("GDP growth value for 2003 has been overwritten")

// Render the filtered table
filtered.Render(os.Stdout, false, true, true, lentele.LoadTemplate("classic"))

```
```

                                          Lithuanian GDP growth and inflation time series
                                                 (annual growth rates in percent)

                                                 ╔══════╦════════════╦═══════════╗
                                                 ║ Year ║ GDP growth ║ Inflation ║
                                                 ╠══════╩════════════╩═══════════╣
                                                 ║ 1999 │   -1.13    │   0.75    ║
                                                 ╟──────┼────────────┼───────────╢
                                                 ║ 2003 │   -10.54   │   -1.15   ║
                                                 ╟──────┼────────────┼───────────╢
                                                 ║ 2009 │   -14.81   │   4.45    ║
                                                 ╟──────┼────────────┼───────────╢
                                                 ║ 2015 │    1.78    │   -0.88   ║
                                                 ╚══════╧════════════╧═══════════╝


──────────────────────────────────────────────────
1. Source: Worldbank
2. Years with negative GDP growth and/or inflation
3. GDP growth value for 2003 has been overwritten
```

## Remove rows

Rows can also be removed manually by providing their rowID (line) or row name.
Headers/footers cannot be removed using the `Table.RemoveRows` method.

Rows (including headers/footers) can also be removed using `Table.RemoveRowsByName`.

```Go
rrows := func() []int {v:=make([]int,10);for i:=range v{v[i]=(i+1)*2};return v}
table.RemoveRows(rrows()...)
table.Render(os.Stdout, false, true, true, lentele.LoadTemplate("classic"))
```

```
                                          Lithuanian GDP growth and inflation time series
                                                 (annual growth rates in percent)

                                                ╔════════╦════════════╦═══════════╗
                                                ║  Year  ║ GDP growth ║ Inflation ║
                                                ╠════════╩════════════╩═══════════╣
                                                ║  1996  │    5.15    │   24.62   ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  1998  │    7.47    │   5.07    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2000  │    3.83    │   1.01    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2002  │    6.76    │   0.30    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2004  │    6.55    │   1.18    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2006  │    7.41    │   3.75    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2008  │    2.63    │   10.93   ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2010  │    1.64    │   1.32    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2012  │    3.83    │   3.09    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2014  │    3.50    │   0.10    ║
                                                ╟────────┼────────────┼───────────╢
                                                ║  2016  │    2.30    │   0.91    ║
                                                ╚════════╧════════════╧═══════════╝
                                                  Means:      4.17        3.64     

────────────────────
1. Source: Worldbank
```

## (Un)marshalling

`lentele.Table` can be marshaled into rich (containing modified values, titles, etc)
or vanilla (just non-modified values and table header) JSON:

```
buf := bytes.NewBuffer([]byte{})

// Marshal to rich JSON
if _, err := table.MarshalToRichJSON(buf); err != nil {
  log.Fatal(err.Error())
}
buf.WriteString("\n\n")

// Marshal to vanilla JSON
if _, err := table.MarshalToVanillaJSON(buf); err != nil {
  log.Fatal(err.Error())
}
buf.WriteString("\n\n")

fmt.Println(buf.String())
```
```
{"rows":[{"cells":[{"value":"Year","modified":"\u001b[1mYear\u001b[0m"},{"value":"GDP growth","modified":"\u001b[1mGDP growth\u001b[0m"},{"value":"Inflation","modified":"\u001b[1mInflation\u001b[0m"}]},{"cells":[{"value":1996,"modified":"\u001b[1m1996\u001b[0m"},{"value":"5.15","modified":"\u001b[33m5.15\u001b[0m"},{"value":"24.62","modified":"\u001b[33m24.62\u001b[0m"}]},{"cells":[{"value":1998,"modified":"\u001b[1m1998\u001b[0m"},{"value":"7.47","modified":"\u001b[33m7.47\u001b[0m"},{"value":"5.07","modified":"\u001b[33m5.07\u001b[0m"}]},{"cells":[{"value":2000,"modified":"\u001b[1m2000\u001b[0m"},{"value":"3.83","modified":"\u001b[32m3.83\u001b[0m"},{"value":"1.01","modified":"\u001b[32m1.01\u001b[0m"}]},{"cells":[{"value":2002,"modified":"\u001b[1m2002\u001b[0m"},{"value":"6.76","modified":"\u001b[33m6.76\u001b[0m"},{"value":"0.30","modified":"\u001b[32m0.30\u001b[0m"}]},{"cells":[{"value":2004,"modified":"\u001b[1m2004\u001b[0m"},{"value":"6.55","modified":"\u001b[33m6.55\u001b[0m"},{"value":"1.18","modified":"\u001b[32m1.18\u001b[0m"}]},{"cells":[{"value":2006,"modified":"\u001b[1m2006\u001b[0m"},{"value":"7.41","modified":"\u001b[33m7.41\u001b[0m"},{"value":"3.75","modified":"\u001b[32m3.75\u001b[0m"}]},{"cells":[{"value":2008,"modified":"\u001b[1m2008\u001b[0m"},{"value":"2.63","modified":"\u001b[32m2.63\u001b[0m"},{"value":"10.93","modified":"\u001b[33m10.93\u001b[0m"}]},{"cells":[{"value":2010,"modified":"\u001b[1m2010\u001b[0m"},{"value":"1.64","modified":"\u001b[32m1.64\u001b[0m"},{"value":"1.32","modified":"\u001b[32m1.32\u001b[0m"}]},{"cells":[{"value":2012,"modified":"\u001b[1m2012\u001b[0m"},{"value":"3.83","modified":"\u001b[32m3.83\u001b[0m"},{"value":"3.09","modified":"\u001b[32m3.09\u001b[0m"}]},{"cells":[{"value":2014,"modified":"\u001b[1m2014\u001b[0m"},{"value":"3.50","modified":"\u001b[32m3.50\u001b[0m"},{"value":"0.10","modified":"\u001b[32m0.10\u001b[0m"}]},{"cells":[{"value":2016,"modified":"\u001b[1m2016\u001b[0m"},{"value":"2.30","modified":"\u001b[32m2.30\u001b[0m"},{"value":"0.91","modified":"\u001b[32m0.91\u001b[0m"}]},{"cells":[{"value":"Means:","modified":"\u001b[1mMeans:\u001b[0m"},{"value":4.17,"modified":"\u001b[1m4.17\u001b[0m"},{"value":3.64,"modified":"\u001b[1m3.64\u001b[0m"}]}],"rownames":["header","","","","","","","","","","","","footer"],"formats":{"0":"%v","1":"%v","2":"%v"},"titles":["Lithuanian GDP growth and inflation time series","(annual growth rates in percent)"],"footnotes":["Source: Worldbank"],"width":{}}

[{"GDP growth":"5.15","Inflation":"24.62","Year":1996},{"GDP growth":"7.47","Inflation":"5.07","Year":1998},{"GDP growth":"3.83","Inflation":"1.01","Year":2000},{"GDP growth":"6.76","Inflation":"0.30","Year":2002},{"GDP growth":"6.55","Inflation":"1.18","Year":2004},{"GDP growth":"7.41","Inflation":"3.75","Year":2006},{"GDP growth":"2.63","Inflation":"10.93","Year":2008},{"GDP growth":"1.64","Inflation":"1.32","Year":2010},{"GDP growth":"3.83","Inflation":"3.09","Year":2012},{"GDP growth":"3.50","Inflation":"0.10","Year":2014},{"GDP growth":"2.30","Inflation":"0.91","Year":2016}]

```

A JSON object can be unmarshaled into a table by using either the `NewFromRichJSON(src io.Reader)`
or the `NewFromVanillaJSON(src io.Reader)` methods. The later method can also me
used to create tables from unrelated JSON objects (e.g. a marshalled slice of maps or structs)

```Go
jsoned := bytes.NewBufferString(`[{"x":1, "y":2},{"y":4},{"x":5, "y":6},{"x":7}]`)
newTable,err := lentele.NewFromVanillaJSON(jsoned, "N/A")
if err != nil {
  log.Fatal(err.Error())
}
newTable.AddTitle("Some JSON object")
newTable.AddFootnote("Unmarshaled using NewFromVanillaJSON")
```
```

                                                         Some JSON object

                                                           ╔═════╦═════╗
                                                           ║  x  ║  y  ║
                                                           ╠═════╩═════╣
                                                           ║  1  │  2  ║
                                                           ╟─────┼─────╢
                                                           ║ N/A │  4  ║
                                                           ╟─────┼─────╢
                                                           ║  5  │  6  ║
                                                           ╟─────┼─────╢
                                                           ║  7  │ N/A ║
                                                           ╚═════╧═════╝


───────────────────────────────────────
1. Unmarshaled using NewFromVanillaJSON

```

# Table templates

Currently the library provides following templates for table rendering:

```go
for _, tmpl := range []func() Template{tmplClassic, tmplSmooth, tmplModern} {  
  tmpl().PrintExample(os.Stdout)
}
```
```
╔════╦═════════════════════════════╦══════════════════════╗
║ ID ║ Client                      ║               Amount ║
╠════╩═════════════════════════════╩══════════════════════╣
║ 1  │ Dunder Mifflin              │              172,341 ║
╟────┼─────────────────────────────┼──────────────────────╢
║ 2  │ Acme Corporation            │               43,223 ║
╟────┼─────────────────────────────┼──────────────────────╢
║ 3  │ Monsters, Inc               │              666,666 ║
╟────┼─────────────────────────────┼──────────────────────╢
║ 4  │ Advanced Idea Mechanics     │              469,218 ║
╟────┼─────────────────────────────┼──────────────────────╢
║ 5  │ Michael Scott Paper Company │                9,288 ║
╟────┼─────────────────────────────┼──────────────────────╢
║ 6  │ Weyland-Yutani Corporation  │          982,283,767 ║
╚════╧═════════════════════════════╧══════════════════════╝
       Total:                        983,644,505 (983.6m)  


╭────┬─────────────────────────────┬──────────────────────╮
│ ID │ Client                      │               Amount │
├────┼─────────────────────────────┼──────────────────────┤
│ 1  │ Dunder Mifflin              │              172,341 │
├────┼─────────────────────────────┼──────────────────────┤
│ 2  │ Acme Corporation            │               43,223 │
├────┼─────────────────────────────┼──────────────────────┤
│ 3  │ Monsters, Inc               │              666,666 │
├────┼─────────────────────────────┼──────────────────────┤
│ 4  │ Advanced Idea Mechanics     │              469,218 │
├────┼─────────────────────────────┼──────────────────────┤
│ 5  │ Michael Scott Paper Company │                9,288 │
├────┼─────────────────────────────┼──────────────────────┤
│ 6  │ Weyland-Yutani Corporation  │          982,283,767 │
├────┴─────────────────────────────┴──────────────────────┤
│      Total:                        983,644,505 (983.6m) │
╰─────────────────────────────────────────────────────────╯


  ID   Client                                      Amount  
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  1    Dunder Mifflin                             172,341  

  2    Acme Corporation                            43,223  

  3    Monsters, Inc                              666,666  

  4    Advanced Idea Mechanics                    469,218  

  5    Michael Scott Paper Company                  9,288  

  6    Weyland-Yutani Corporation             982,283,767  
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       Total:                        983,644,505 (983.6m)  
```

You can create your own templates by implementing the `template.Template` interface.
A good source of characters that can be used in creating table designs can be
found [here](https://en.wikipedia.org/wiki/Box-drawing_character).

# TODO

- [x] Increase test coverage
- [x] Add `Template.PrintExample` to show an example table
