package lentele

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

// New creates a new table implementing the lentele.Table interface
func New(columns ...string) Table {

	newTable := &table{
		Mutex:       &sync.Mutex{},
		Rows:        []*row{},
		RowNames:    []string{},
		Formats:     []string{},
		Footnotes:   []string{},
		headAndFoot: map[string]*row{},
	}

	if len(columns) > 0 {
		newTable.AddHeader(columns)
	}

	return newTable
}

// NewFromVanillaJSON creates a table from vanilla JSON.
//
// This method expects the source to be a simple JSON representation of a list
// of objects, e.g.:
//
// bytes.NewBuffer([]byte(`[{col1: "val1", col2: "val2"},{col1: "val3", col2: "val4"}]`))
func NewFromVanillaJSON(source io.Reader) (Table, error) {

	return nil, nil
}

// NewFromRichJSON creates a table with meta information (row names, modified
// values, etc.)
func NewFromRichJSON(source io.Reader) (Table, error) {

	return nil, nil
}

// LoadTemplate returns the named template
func LoadTemplate(name string) Template {
	switch strings.ToLower(name) {
	case "smooth":
		return tmplSmooth()
	case "modern":
		return tmplModern()
	default:
		return tmplClassic()
	}
}

// table implements the lentele.Table interface
type table struct {
	*sync.Mutex `json:",omit"`
	Rows        []*row   `json:"rows"`
	RowNames    []string `json:"rownames"`
	Formats     []string `json:"formats"`
	Titles      []string `json:"titles"`
	Footnotes   []string `json:"footnotes"`

	headAndFoot map[string]*row // Map of addresses to header and footer pointers
}

// row implements the lentele.Row interface
type row struct {
	*sync.Mutex `json:",omit"`
	Cells       []*cell `json:"cells"`
	tref        *table  // Parent table reference
}

// value stores individual cell values
type cell struct {
	*sync.Mutex `json:",omit"`
	Value       interface{}                     `json:"value"`
	ModVal      interface{}                     `json:"modified"`
	ModFunc     func(v interface{}) interface{} `json:",omit"`
}

// AddTitle adds a title to the table
// NB: locks t
func (t *table) AddTitle(title string) {
	t.Lock()
	defer t.Unlock()

	if len(title) == 0 {
		return
	}

	t.Titles = append(t.Titles, title)
}

// AddFootnote adds a footnoite to the table
// NB: locks t
func (t *table) AddFootnote(footnote string) {
	t.Lock()
	defer t.Unlock()

	if len(footnote) == 0 {
		return
	}

	t.Footnotes = append(t.Footnotes, footnote)
}

// AddHeader adds a header to the table
// NB: indirectly locks t
func (t *table) AddHeader(colnames []string) Row {
	iface := make([]interface{}, len(colnames))
	for i, v := range colnames {
		iface[i] = v
	}
	return t.AddRow("header").Insert(iface...)
}

// AddFooter adds a footer to the table
// NB: indirectly locks t
func (t *table) AddFooter() Row {
	return t.AddRow("footer")
}

// AddRow adds a new row to the table
// NB: locks t
func (t *table) AddRow(name string) Row {
	t.Lock()
	defer t.Unlock()

	// Normalize name
	name = strings.ToLower(name)

	// Check for header or footer
	switch name {
	case "header", "footer":
		special, ok := t.headAndFoot[name]
		if ok {
			return special
		}
	}

	// Create new row and append it to the table
	newRow := &row{
		Mutex: &sync.Mutex{},
		Cells: []*cell{},
		tref:  t,
	}

	// Add rows and their names
	t.Rows = append(t.Rows, newRow)
	t.RowNames = append(t.RowNames, name)

	// Remember header and footer
	switch name {
	case "header", "footer":
		t.headAndFoot[name] = newRow
	}

	return newRow
}

// SetFormat sets a column's format and returns an error if no such column
// exists. If no format is specified, then "%v" is going to be used.
func (t *table) SetFormat(format string, colnames ...string) error {
	return nil
}

// Transform transforms the values of colnames using function 'trans'
// NB: locks t
func (t *table) Transform(trans func(v interface{}) interface{}, colnames ...string) {
	t.Lock()
	defer t.Unlock()

	if len(colnames) == 0 {
		return
	}

	// Get indexes
	colIdx := t.getColIdx(false, colnames...)

	// Get header and footer
	header := t.headAndFoot["header"]
	footer := t.headAndFoot["footer"]

	// Transform values
	for i := range t.Rows {
		if t.Rows[i] == header || t.Rows[i] == footer {
			continue
		}

		for _, j := range colIdx {
			if j <= len(t.Rows[i].Cells)-1 {

				switch reflect.TypeOf(t.Rows[i].Cells[j].Value).Kind() {

				case reflect.Slice:
					slice := reflect.ValueOf(t.Rows[i].Cells[j].Value)
					result := make([]interface{}, slice.Len(), slice.Len())
					for k := 0; k < slice.Len(); k++ {
						result[k] = trans(slice.Index(k).Interface())
					}
					t.Rows[i].Cells[j].Value = result

				default:
					t.Rows[i].Cells[j].Value = trans(t.Rows[i].Cells[j].Value)
				}

			}
		}
	}

}

// GetRow returns the nth row from the table or error if no such row exists
func (t *table) GetRow(nth int) (Row, error) {
	return nil, nil
}

// Returns a row
func (t *table) GetRowByName(name string) (Row, error) {
	t.Lock()
	defer t.Unlock()

	name = strings.ToLower(name)

	switch name {
	case "header", "footer":
		special, ok := t.headAndFoot[name]
		if !ok {
			return nil, fmt.Errorf("GetRowByName: no such rowname '%s'", name)
		}
		return special, nil

	default:
		for i, rowname := range t.RowNames {
			if rowname == name {
				return t.Rows[i], nil
			}
		}
		return nil, fmt.Errorf("GetRowByName: no such rowname '%s'", name)
	}

}

// Filter applies a filter to each row and returns a filtered table.
// If inplace is set to true, then the filtered-out rows are permanently deleted
// (references to the rows are removed).
// Otherwise a new table, *referencing* the relevant rows, is created
func (t *table) Filter(filter func(values ...interface{}) bool, inplace, keepFooter bool, columns ...string) (Table, error) {
	t.Lock()
	defer t.Unlock()

	// Validate columns
	if len(columns) == 0 {
		return nil, fmt.Errorf("Filter: at least one column must be provided")
	}

	colIdx := t.getColIdx(false, columns...)
	if len(colIdx) == 0 {
		return nil, fmt.Errorf("Filter: unknown columns")
	}

	// Header and footer
	header := t.headAndFoot["header"]
	footer := t.headAndFoot["footer"]

	// Fields for the new table
	headAndFoot := map[string]*row{}
	rowNames := []string{}
	rows := []*row{}

	// Find relevant rows
	for i, row := range t.Rows {

		// Always keep the header
		if row == header {
			rows = append(rows, t.Rows[i])
			rowNames = append(rowNames, t.RowNames[i])
			headAndFoot["header"] = header
		}

		// Optionally keep the footer
		if row == footer && keepFooter {
			rows = append(rows, t.Rows[i])
			rowNames = append(rowNames, t.RowNames[i])
			headAndFoot["footer"] = footer
		}

		// Run the filter
		values := make([]interface{}, len(colIdx), len(colIdx))
		for j, col := range colIdx {
			if len(row.Cells)-1 < col {
				continue
			}
			values[j] = row.Cells[col].Value
		}
		if filter(values...) {
			rows = append(rows, t.Rows[i])
			rowNames = append(rowNames, t.RowNames[i])
		}
	}

	fTable := t.tableFromRows(false, inplace, rows, rowNames, headAndFoot)

	return fTable, nil

}

// FilterByRowNames is same as filter, only uses row names instead of column
// values.
//
// Rows without a unique name are treated as having a blank name, i.e. ""
func (t *table) FilterByRowNames(filter func(rowname string) bool, inplace, keepFooter bool) Table {
	t.Lock()
	defer t.Unlock()

	// Header and footer
	header := t.headAndFoot["header"]
	footer := t.headAndFoot["footer"]

	// Fields for the new table
	headAndFoot := map[string]*row{}
	rowNames := []string{}
	rows := []*row{}

	// Find relevant rows
	for i, row := range t.Rows {
		if relevant := filter(t.RowNames[i]); relevant || row == header || (row == footer && keepFooter) {

			if row == header {
				headAndFoot["header"] = row
			}

			if row == footer {
				headAndFoot["footer"] = row
			}

			rows = append(rows, row)
			rowNames = append(rowNames, t.RowNames[i])
		}

	}

	fTable := t.tableFromRows(false, inplace, rows, rowNames, headAndFoot)

	return fTable
}

// RemoveRows removes a set of rows from the table
// NB: locks t
func (t *table) RemoveRows(rowIDs ...int) error {
	t.Lock()
	defer t.Unlock()

	if len(rowIDs) == 0 {
		return fmt.Errorf("RemoveRows: no row IDs provided")
	}

	// Header and footer
	header := t.headAndFoot["header"]
	footer := t.headAndFoot["footer"]

	selected := []int{}

	// Validate rowIDs
	for _, nth := range rowIDs {
		if nth < 0 || nth > len(t.Rows) {
			return fmt.Errorf("RemoveRow: no such row")
		}

		srow := t.Rows[nth]
		if srow == header || srow == footer {
			return fmt.Errorf("RemoveRow: cannot remove header/footer")
		}
		selected = append(selected, nth)
	}

	// Remove rows
	t.removeRows(false, selected)

	return nil
}

// RemoveRowsByName removes a set of named row from the table
// NB: locks t
func (t *table) RemoveRowsByName(names ...string) error {
	t.Lock()
	defer t.Unlock()

	if len(names) == 0 {
		return fmt.Errorf("RemoveRows: no row IDs provided")
	}

	// Header and footer
	header := t.headAndFoot["header"]
	footer := t.headAndFoot["footer"]

	// Gather RowIds
	selected := []int{}
	for _, rmname := range names {
		for i, name := range t.RowNames {
			if strings.ToLower(rmname) == name && t.Rows[i] != header && t.Rows[i] != footer {
				selected = append(selected, i)
			}
		}
	}

	// Remove rows
	t.removeRows(false, selected)

	return nil
}

// Render writes a rendered table into an io.Writer
// NB: locks t
func (t *table) Render(dst io.Writer, measureModified, modified, centered bool, template Template, columns ...string) {
	t.Lock()
	defer t.Unlock()

	// Header and footer info
	headRow, footRow := -1, -1
	header, _ := t.headAndFoot["header"]
	footer, _ := t.headAndFoot["footer"]

	// Get relevant columns
	colIdx := t.getColIdx(false, columns...)

	// Final rows
	measureRows := [][]string{}
	printRows := [][]string{}

	// Walk through rows
	rowCount := 0
	widths := []int{}
	for i, row := range t.Rows {

		if t.Rows[i] == header {
			headRow = i
		} else if t.Rows[i] == footer {
			footRow = i
		} else {
			rowCount++
		}

		// Relevant columns
		rangeVar := []int{}
		if len(colIdx) != 0 {
			rangeVar = colIdx
		} else {
			for k := range row.Cells {
				rangeVar = append(rangeVar, k)
			}
		}

		measureRow := []string{}
		printRow := []string{}

		// Walk through all the columns of a row
		for j, jcol := range rangeVar {

			// Ignore overflowing cells
			if jcol >= len(row.Cells) {
				continue
			}

			// Select cell
			jcell := row.Cells[jcol]

			// Determine formats and modifiers
			jcell.Lock()
			if len(widths) < j+1 {
				widths = append(widths, 0)
				t.Formats = append(t.Formats, "%v")
			}

			format := t.Formats[j]
			if jcell.ModFunc == nil {
				jcell.ModFunc = func(v interface{}) interface{} { return v }
			}

			// Prepare formated and modified values
			var valueNorm string
			var valueMod string

			switch reflect.TypeOf(jcell.Value).Kind() {

			case reflect.String:
				valueNormSlice := []string{}
				valueModSlice := []string{}
				value, _ := jcell.Value.(string)
				for _, part := range strings.Split(value, "\n") {
					valueNormSlice = append(valueNormSlice, fmt.Sprintf(format, part))
					valueModSlice = append(valueModSlice, fmt.Sprintf(format, jcell.ModFunc(part)))
				}
				valueNorm = strings.Join(valueNormSlice, "\n")
				valueMod = strings.Join(valueModSlice, "\n")

			case reflect.Slice:
				slice := reflect.ValueOf(jcell.Value)
				valueNormSlice := []string{}
				valueModSlice := []string{}
				for i := 0; i < slice.Len(); i++ {
					valueNormSlice = append(valueNormSlice, fmt.Sprintf(format, slice.Index(i).Interface()))
					valueModSlice = append(valueModSlice, fmt.Sprintf(format, jcell.ModFunc(slice.Index(i).Interface())))
				}
				valueNorm = strings.Join(valueNormSlice, "\n")
				valueMod = strings.Join(valueModSlice, "\n")

			default:
				valueNorm = fmt.Sprintf(format, jcell.Value)
				valueMod = fmt.Sprintf(format, jcell.ModFunc(jcell.Value))
			}

			jcell.ModVal = valueMod

			if measureModified {
				measureRow = append(measureRow, valueMod)
			} else {
				measureRow = append(measureRow, valueNorm)
			}
			if modified {
				printRow = append(printRow, valueMod)
			} else {
				printRow = append(printRow, valueNorm)
			}

			// Remember column widths
			if measureModified {
				if length := utf8.RuneCountInString(valueMod); length > widths[j] {
					widths[j] = length
				}
			} else {
				if length := utf8.RuneCountInString(valueNorm); length > widths[j] {
					widths[j] = length
				}
			}

			jcell.Unlock()
		}

		measureRows = append(measureRows, measureRow)
		printRows = append(printRows, printRow)
	}

	// Set template widths
	template.SetColumnWidths(widths)
	template.SetDisplayOptions(centered)

	// Prepare table slice
	lines := []string{""}

	// Title
	if len(t.Titles) > 0 {
		lines = append(lines, template.RenderTitles(t.Titles)...)
	}

	// Render header
	if headRow != -1 {
		lines = append(lines, template.RenderHeader(measureRows[headRow], printRows[headRow])...)
	}

	// Render rows
	rnr := 1
	for i := range measureRows {
		if i == headRow || i == footRow {
			continue
		}
		lines = append(lines, template.RenderRow(rnr, rowCount, measureRows[i], printRows[i])...)
		rnr++
	}

	// Render footer
	if footRow != -1 {
		lines = append(lines, template.RenderFooter(measureRows[footRow], printRows[footRow])...)
	} else {
		lines = append(lines, template.RenderFooter([]string{}, []string{})...)
	}

	// Render Footnotes
	if len(t.Footnotes) > 0 {
		lines = append(lines, template.RenderFootnotes(t.Footnotes)...)
	}

	// Write to destination
	dst.Write([]byte(strings.Join(lines, "\n")))

}

// Marshals the table to json including all meta information (row names,
// modified values, etc.)
// NB: locks t
func (t *table) MarshalToRichJSON(dst io.Writer) (int, error) {
	return 0, nil
}

// MarshalToVanillaJSON marshals the table as a simple list of objects,
// one object per row, i.e. [{col1: val1, col2: val2},{col1: val3, col2: val4}].
// It does not preserve modifiers, row names and so on.
// NB: locks t
func (t *table) MarshalToVanillaJSON(dst io.Writer) (int, error) {
	return 0, nil
}

// removeRows removes a set of rows from the trable
func (t *table) removeRows(lock bool, selected []int) {
	if lock {
		t.Lock()
		defer t.Unlock()
	}

	// Sort entries
	sort.Ints(selected)

	// remove rows
	for i := len(selected) - 1; i >= 0; i-- {
		nth := selected[i]
		if nth != len(t.Rows)-1 {
			t.Rows = append(t.Rows[:nth], t.Rows[nth+1:]...)
			t.RowNames = append(t.RowNames[:nth], t.RowNames[nth+1:]...)
		} else {
			t.Rows = t.Rows[:nth]
			t.RowNames = t.RowNames[:nth]
		}
	}

}

// tableFromRows builds a table from rows
func (t *table) tableFromRows(lock, inplace bool, rows []*row, rowNames []string, hf map[string]*row) *table {

	if lock {
		t.Lock()
		defer t.Unlock()
	}

	var fTable *table

	// Create a new table or replace table rows
	if inplace {
		t.Rows = rows
		t.RowNames = rowNames
		t.headAndFoot = hf
		fTable = t
	} else {
		fTable = &table{
			Mutex:       &sync.Mutex{},
			Rows:        rows,
			RowNames:    rowNames,
			Formats:     t.Formats,
			Titles:      t.Titles,
			Footnotes:   t.Footnotes,
			headAndFoot: hf,
		}
	}

	return fTable
}

// getColIdx returns indexes of selected columns
func (t *table) getColIdx(lock bool, columns ...string) []int {
	if lock {
		t.Lock()
		defer t.Unlock()
	}

	colIdx := []int{}
	for _, col := range columns {
		idx := t.getColnameIndex(col, false, false)
		if idx != -1 {
			colIdx = append(colIdx, idx)
		}
	}

	return colIdx
}

// getColnameIndex returns the position of a header's columns "colname"
// NB: locks t
func (t *table) getColnameIndex(colname string, lockTable, lockHead bool) int {
	if lockTable {
		t.Lock()
		defer t.Unlock()
	}

	// Default: no such column
	index := -1

	// Check the header
	header, ok := t.headAndFoot["header"]
	if !ok {
		return index
	}

	// Find colname index
	if lockHead {
		header.Lock()
		defer header.Unlock()
	}

	for i, hcell := range header.Cells {

		vs, ok := hcell.Value.(string)
		if !ok {
			continue
		}
		if strings.ToLower(vs) == strings.ToLower(colname) {
			return i
		}
	}

	return index
}

// Insert inserts some values into the row
func (r *row) Insert(vals ...interface{}) Row {
	r.Lock()
	defer r.Unlock()

	if len(vals) == 0 {
		return r
	}

	// Insert cells
	//
	// NB: If no header has been set, then all the values are going to be shown
	// in separate generically names columns (COL1 - COLK)
	// Setting the header later will hide the overflowing cells.
	for i := 0; i < len(vals); i++ {
		r.Cells = append(r.Cells, &cell{
			Mutex: &sync.Mutex{},
			Value: vals[i],
		})
	}

	return r
}

// Change changes a row cell's value
func (r *row) Change(colname string, value interface{}) Row {
	r.Lock()
	defer r.Unlock()

	colname = strings.ToLower(colname)

	// Header and footer rows
	header := r.tref.headAndFoot["header"]
	footer := r.tref.headAndFoot["footer"]

	// Find relevant column
	var index int
	switch r {
	case header, footer:
		index = r.tref.getColnameIndex(colname, true, false)
	default:
		index = r.tref.getColnameIndex(colname, true, true)
	}

	if index == -1 {
		return r
	}

	// Change the value
	rcell := r.Cells[index]
	rcell.Lock()
	rcell.Value = value
	rcell.ModFunc = func(v interface{}) interface{} { return v }
	rcell.Unlock()

	return r
}

// Modify modifies an entry using a modifier
func (r *row) Modify(modifier func(interface{}) interface{}, colnames ...string) Row {
	r.Lock()
	defer r.Unlock()

	if len(colnames) == 0 {
		return r
	}

	// Header and footer rows
	header := r.tref.headAndFoot["header"]
	footer := r.tref.headAndFoot["footer"]

	// Modify all relevant cells
	for _, colname := range colnames {

		colname = strings.ToLower(colname)

		var index int

		switch r {
		case header, footer:
			index = r.tref.getColnameIndex(colname, true, false)
		default:
			index = r.tref.getColnameIndex(colname, true, true)
		}

		if index == -1 {
			continue
		}

		rcell := r.Cells[index]
		rcell.Lock()
		rcell.ModFunc = modifier
		rcell.Unlock()
	}

	return r
}
