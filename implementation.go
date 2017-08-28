package lentele

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode/utf8"
)

// New creates a new table implementing the lentele.Table interface
func New(columns ...string) Table {

	newTable := &table{
		Mutex:    &sync.Mutex{},
		Rows:     []*row{},
		RowNames: []string{},
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
	case "plain":
		return tmplClassic()
	case "smooth":
		return tmplClassic()
	case "modern":
		return tmplClassic()
	default:
		return tmplClassic()
	}
}

// table implements the lentele.Table interface
type table struct {
	*sync.Mutex `json:",omit"`
	Rows        []*row   `json:"rows"`
	RowNames    []string `json:"rownames"`

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

// AddHeader adds a header to the table
// NB: indirectly locks t
func (t *table) AddHeader(colnames []string) Row {
	return t.AddRow("header").Insert(colnames)
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
	if name == "header" {
		header, ok := t.headAndFoot["header"]
		if ok {
			return header
		}
	} else if name == "footer" {
		footer, ok := t.headAndFoot["footer"]
		if ok {
			return footer
		}
	}

	// Create new row and append it to the table
	newRow := &row{
		Mutex: &sync.Mutex{},
		Cells: []*cell{},
		tref:  t,
	}

	t.Rows = append(t.Rows, newRow)

	return newRow
}

// SetFormat sets a column's format and returns an error if no such column
// exists. If no format is specified, then "%v" is going to be used.
func (t *table) SetFormat(column, format string) error {
	return nil
}

// GetRow returns the nth row from the table or error if no such row exists
func (t *table) GetRow(nth int) (Row, error) {
	return nil, nil
}

// Returns a row
func (t *table) GetRowByName(name string) (Row, error) {
	return nil, nil
}

// Filter applies a filter to each row and returns a filtered table.
// If inplace is set to true, then the filtered-out rows are permanently deleted
// (references to the rows are removed).
// Otherwise a new table, *referencing* the relevant rows, is created
func (t *table) Filter(filter func(columns ...string) bool, inplace bool, columns ...string) Table {

	return nil
}

// FilterByRowNames is same as filter, only uses row names instead of column
// values.
//
// Rows without a unique name are treated as having a blank name, i.e. ""
func (t *table) FilterByRowNames(filter func(rowname string) bool, inplace bool) Table {

	return nil
}

// Removes the nth row from the table
func (t *table) RemoveRow(nth int) error {
	return nil
}

// Removes the named row from the table
func (t *table) RemoveRowByName(name string) error {
	return nil
}

// Render writes a rendered table into an io.Writer
// NB: locks t
func (t *table) Render(dst io.Writer, measureModified, modified bool, template Template, columns ...string) {
	t.Lock()
	defer t.Unlock()
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

	// Find relevant column
	index := r.getColnameIndex(colname)
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

	for _, colname := range colnames {

		index := r.getColnameIndex(colname)
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

// getHeadOrFoot returns a header or a footer if they exist
// NB: locks parent table
func (r *row) getHeadOrFoot(name string) (*row, bool) {
	r.tref.Lock()
	defer r.tref.Unlock()

	header, ok := r.tref.headAndFoot[strings.ToLower(name)]
	return header, ok
}

// Returns the position of a header's columns "colname"
// NB: locks header
// NB: indirectly locks parent table
func (r *row) getColnameIndex(colname string) int {

	// Default: no such column
	index := -1

	// Check the header
	header, ok := r.getHeadOrFoot("header")
	if !ok {
		return index
	}

	// Find colname index
	header.Lock()
	for i, hcell := range header.Cells {
		vs, ok := hcell.Value.(string)
		if !ok {
			continue
		}
		if strings.ToLower(vs) == strings.ToLower(colname) {
			return i
		}
	}
	header.Unlock()

	return index
}

// template implements the Template interface
//
// +====+====+====+====+ <- H1: [2]string{"+","="}
// |    |    |    |    | <- H2: [1]string{"|"}
// +====+====+====+====+ <- H3: [2]string{"+","="}
// +----+----+----+----+ <- C1: [2]string{"+","-"}
// |    |    |    |    | <- C2: [1]string{"|"}
// +----+----+----+----+ <- C3: [2]string{"+","-"}
// +====+====+====+====+ <- F1: [2]string{"+","="}
// |    |    |    |    | <- F2: [1]string{"|"}
// +====+====+====+====+ <- F3: [2]string{"+","="}
//
type template struct {
	*sync.Mutex
	ColWidths []int
	H1        [2]string
	H2        [1]string
	H3        [2]string
	C1        [2]string
	C2        [1]string
	C3        [2]string
	F1        [2]string
	F2        [1]string
	F3        [2]string
}

// SetColumnWidths sets the column widths
func (t *template) SetColumnWidths(widths []int) {
	t.Lock()
	defer t.Unlock()

	t.ColWidths = widths
}

// RenderHeader renders the header row
func (t *template) RenderHeader(cells []string) []string {
	t.Lock()
	defer t.Unlock()

	L1 := t.H1[0]
	L2 := t.H2[0]
	L3 := t.H3[0]
	for i, width := range t.ColWidths {

		vlen := utf8.RuneCountInString(cells[i])
		reps := int((width + 2 - vlen) / 2)
		sp1 := strings.Repeat(" ", reps)
		sp2 := strings.Repeat(" ", width+2-vlen-reps)

		L1 += strings.Repeat(t.H1[1], width+2) + t.H1[0]
		L2 += fmt.Sprintf("%s%s%s", sp1, cells[i], sp2) + t.H2[0]
		L3 += strings.Repeat(t.H3[1], width+2) + t.H3[0]
	}

	return []string{L1, L2, L3}

}

// RenderRow renders a regular row
func (t *template) RenderRow(row int, cells []string) []string {
	t.Lock()
	defer t.Unlock()

	L1 := t.C1[0]
	L2 := t.C2[0]
	L3 := t.C3[0]
	for i, width := range t.ColWidths {

		vlen := utf8.RuneCountInString(cells[i])
		reps := int((width + 2 - vlen) / 2)
		sp1 := strings.Repeat(" ", reps)
		sp2 := strings.Repeat(" ", width+2-vlen-reps)

		L1 += strings.Repeat(t.C1[1], width+2) + t.C1[0]
		L2 += fmt.Sprintf("%s%s%s", sp1, cells[i], sp2) + t.C2[0]
		L3 += strings.Repeat(t.C3[1], width+2) + t.C3[0]
	}

	return []string{L1, L2, L3}
}

// RenderFooter renders the footer row
func (t *template) RenderFooter(cells []string) []string {
	t.Lock()
	defer t.Unlock()

	L1 := ""
	L2 := ""
	L3 := ""
	prev := false
	for i, width := range t.ColWidths {

		vlen := utf8.RuneCountInString(cells[i])
		reps := int((width + 2 - vlen) / 2)
		sp1 := strings.Repeat(" ", reps)
		sp2 := strings.Repeat(" ", width+2-vlen-reps)

		if vlen > 0 {
			if !prev {
				L1 += t.F1[0]
				L2 += t.F2[0]
				L3 += t.F3[0]
			}
			L1 += strings.Repeat(t.F1[1], width+2) + t.F1[0]
			L2 += fmt.Sprintf("%s%s%s", sp1, cells[i], sp2) + t.F2[0]
			L3 += strings.Repeat(t.F3[1], width+2) + t.F3[0]
			prev = true
		} else {
			if !prev {
				L1 += " "
				L2 += " "
				L3 += " "
			}
			L1 += strings.Repeat(" ", width+2)
			L2 += strings.Repeat(" ", width+2)
			L3 += strings.Repeat(" ", width+2)
			prev = false
		}
	}

	return []string{L1, L2, L3}
}

// Classic template
func tmplClassic() *template {

	return &template{
		Mutex:     &sync.Mutex{},
		ColWidths: []int{},
		H1:        [2]string{"+", "="},
		H2:        [1]string{"|"},
		H3:        [2]string{"+", ""},
		C1:        [2]string{"+", "-"},
		C2:        [1]string{"|"},
		C3:        [2]string{"+", "-"},
		F1:        [2]string{"+", "="},
		F2:        [1]string{"|"},
		F3:        [2]string{"+", "="},
	}
}
