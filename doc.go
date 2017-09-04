// Package lentele can be used to display structured data (log entries, cluster
// information, etc.) as a table. It also permits light data manipulations (filtering
// transforming, etc.)
//
// Lazy evaluation
//
// lentele.Table is lazy when it comes to modifications. Modified values are
// created only at render or export time. Until then only the reference to the
// modifier is kept.
//
// Modifiers
//
// A cell value can be modified by a modifier function with the signature
//  func(v interface{}) interface{}
//
// Usually such modifications are used to format the value in a more complicated
// way than is possible with fmt.Sprintf, e.g. trimming strings and adding ellipsis,
// putting in the thousands separator, etc.
// One special modification is the addition of ansi escape sequences, e.g.:
//
//  func important(v interface{}) interface{} {
//	 return fmt.Sprintf("\033[1;33m%v\033[0m", v)
//  }
//
// Such a modifier pads the string with characters that are invisible to the
// console (they generate colors), but influence the calculation of the length
// of the string.
//
// Rendering ansi colors
//
// If ansi escape sequences are used to modify cell values, then calculation
// of cell widths will lead to peculiar results when printing to os.Stdout,
// since the escape sequences are not visible in the console (get converted
// to color), but increase the size of the string. Use "measureModified=false"
// when rendering tables that include ansi colors.
//
// Mixing ansi and non-ansi modifiers
//
// Mixing regular modifiers (e.g. insertion of thousands-separators)
// together with ansi-coloring will, generally, have a negative effect, since
// the modifications will not be separable at render time. It is, thus, a better
// idea to insert already (regularly) modified values and then add ansi colors, e.g.:
//  table.AddRow().Insert("ClientX",regMod(8829156.21)).Modify(ansiHigh, "Amount")
//
// Transformations
//
// Sometimes it makes sense to transform a whole column (e.g. trim strings and add
// ellipsis) instead of modifying individual cells. Columns can be transformed by
//  Table.Transform()
//
// Chaining
//
// All the row operations can be chained, e.g.:
//  table.AddRow("").Insert(1).Insert("E corp").Insert(425567.32).Modify(evil, "Client").Modify(mm, "Amount").Change("ID", 0)
package lentele
