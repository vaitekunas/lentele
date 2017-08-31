package lentele

import (
	"io"
)

// Table contains all the rows
//
// All methods should be implemented in a thread-safe way.
type Table interface {

	// AddTitle adds an optional title to the table
	AddTitle(title string)

	// AddFootnote adds a footnote to the table
	AddFootnote(footnote string)

	// AddHeader adds a header row (optional, but recommended).
	//
	// Setting the header is required in order to reference columns by name in
	// Row.Modify and Row.Change.
	//
	// The header can later be referenced via Table.GetRowByName("header")
	AddHeader(colnames []string) Row

	// AddFooter adds a footer row (optional).
	//
	// The footer can later be referenced via Table.GetRowByName("footer")
	//
	// The footer row does not have to include values for all columns, i.e.
	// an empty footer can be added to begin with and only the relevant columns
	// changed:
	//  table.AddFooter()
	//  ...
	//  table.GetRowByName("footer").Change("Amount",fmt.Sprintf("Total: %d",total))
	AddFooter() Row

	// Adds a row to the table.
	//
	// If a unique name is provided, then the table can be searched for rows by names.
	// Rownames are case insensitive.
	AddRow(name string) Row

	// SetFormat sets a column's format and returns an error if no such column
	// exists. If no format is specified, then "%v" is going to be used.
	SetFormat(format string, colnames ...string) error

	// SetColumnWidth overrides column width calculations with static values
	SetColumnWidth(width int, colnames ...string) error

	// GetRow returns the nth row from the table or error if no such row exists
	GetRow(nth int) (Row, error)

	// Returns a row
	GetRowByName(name string) (Row, error)

	// Transform a function to all the values in colnames
	Transform(trans func(v interface{}) interface{}, colnames ...string)

	// Filter applies a filter to each row and returns a filtered table.
	// If inplace is set to true, then the filtered-out rows are permanently deleted
	// (references to the rows are removed).
	// Otherwise a new table, *referencing* the relevant rows, is created
	Filter(filter func(values ...interface{}) bool, inplace, keepFooter bool, columns ...string) (Table, error)

	// FilterByRowNames is same as filter, only uses row names instead of column
	// values.
	//
	// Rows without a unique name are treated as having a blank name, i.e. ""
	FilterByRowNames(filter func(rowname string) bool, inplace, keepFooter bool) Table

	// RemoveRows removes a set of rows from the table
	RemoveRows(rowIds ...int) error

	// RemoveRowsByName removes a set of named rows from the table
	RemoveRowsByName(names ...string) error

	// Render renders the table into an io.Writer
	//
	// Setting modified to false will ignore all the applied modifications (Row.Modify).
	// Setting measureModified to true, will use the modified string representation
	// to calculate cell widths.
	Render(dst io.Writer, measureModified, modified, centered bool, template Template, columns ...string)

	// Marshals the table to json including all meta information (row names,
	// modified values, etc.)
	MarshalToRichJSON(io.Writer) (int, error)

	// MarshalToVanillaJSON marshals the table as a simple list of objects,
	// one object per row, i.e. [{col1: val1, col2: val2},{col1: val3, col2: val4}].
	// It does not preserve modifiers, row names and so on.
	MarshalToVanillaJSON(io.Writer) (int, error)
}

// Row represents a single table row.
//
// The implemented methods should return a reference to the same row, so that
// the calls could be chained, e.g.
//  row.Insert(...).Modify(...).
//
// All methods should be implemented in a thread-safe way.
type Row interface {

	// Inserts values into a row
	// Should not fail if some column values are missing or too many values
	// are provided.
	Insert(values ...interface{}) Row

	// Change changes a row cell's value
	// Fails silently if column not available
	Change(colname string, value interface{}) Row

	// Modify modifies (formats) the cell value.
	// This method is lazy, i.e. it only saves the reference to the modifier.
	// The modification is done at render time if modified bool is set to true.
	Modify(modifier func(interface{}) interface{}, colnames ...string) Row
}

// Template handles
type Template interface {

	// SetColumnWidths sets the column widths
	SetColumnWidths([]int)

	// SetDisplayOptions sets some display options
	SetDisplayOptions(center bool)

	// RenderHeader renders the header row
	RenderHeader(mcells, pcells []string) []string

	// RenderRow renders a regular row
	RenderRow(row, rows int, mcells, pcells []string) []string

	// RenderFooter renders the footer row
	RenderFooter(mcells, pcells []string) []string

	// RenderTitles renders table's titles
	RenderTitles(titles []string) []string

	// RenderFootnotes renders table's footnotes
	RenderFootnotes(footnotes []string) []string
}
