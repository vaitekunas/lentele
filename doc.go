// This package can be used to display tabular data in a prettier way.
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
// to color), but increase the size of the string.
//
// Mixing ansi and non-ansi modifiers
//
// Mixing regular modifiers (e.g. insertion of thousands-separators)
// together with ansi-coloring will, generally, have a negative effect, since
// the modifications will not be separable at render time. It is, thus, a better
// idea to insert already (regularly) modified values and then add ansi colors, e.g.::
//  table.AddRow().Insert("ClientX",regMod(8829156.21)).Modify(ansiHigh, "Amount")
//
// String trimmers
//
// lentele.Table does not include any explicit trimming functions to reduce the
// length of very long strings so that cells don't get too wide. Use inline
// modifiers for that, i.e.
//
//  func trim(s string) string {
//  	if utf8.RuneCountInString(s) > 7 {
//  		return fmt.Sprintf("%s...", s[:7])
//  	}
//  	return s
//  }
//
//  table.AddRow().Insert(trim("Iceland"), trim("Eyjafjallajökull"))
//
//  // alternatively:
//  table.AddRow().Insert("Iceland", "Eyjafjallajökull").Modify("Volcano", trim)
//
// Modifier factories
//
// In most cases it makes sense to create a modifier factory in order to make
// the creation of modifiers much easier:
//  import "github.com/fatih/color"
//
//  modFactory := func(c color.Color) func(interface{}) interface{} {
//    return func(v interface{}) interface{}{
//      return c.Sprintf("%v",v)
//    }
//  }
//
//  high := modFactory(color.New(color.FgGreen))
//  low := modFactory(color.New(color.FgRed))
//
// Instead of having several modification functions, one can put the whole
// business logic into the modifier itself and use it for all values:
//
//  revenue := func(v interface{}) interface{}{
//    vf, ok := v.(float)
//    if !ok {
//      return v
//    }
//
//    var c colors.Color
//
//    switch vf {
//      case < 1000:
//        c = color.New(color.FgRed)
//      case < 10000:
//        c = color.New(color.FgYellow)
//      default:
//        c = color.New(color.FgGreen)
//    }
//
//    return c.Sprintf("%6.2f",vf):
//  }
//
//  table.AddRow().Insert("Client Y",2883.12).Modify(revenue, "Amount")
package lentele
