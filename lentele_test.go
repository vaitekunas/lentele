package lentele

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io"
	"testing"
)

func recession(v interface{}) interface{} {
	return color.New(color.FgRed).Sprint(v)
}

func round(v interface{}) interface{} {
	vf, ok := v.(float64)
	if !ok {
		return v
	}
	return fmt.Sprintf("%4.2f", vf)
}

func yearFilter(from, to int) func(vals ...interface{}) bool {
	return func(vals ...interface{}) bool {
		if len(vals) != 1 {
			return false
		}
		v := vals[0]
		vi, ok := v.(int)
		if !ok {
			return false
		}
		if vi >= from && vi <= to {
			return true
		}
		return false
	}
}

func badYear(vals ...interface{}) bool {
	if len(vals) != 1 {
		return false
	}
	vf, ok := vals[0].(float64)
	if !ok {
		return false
	}
	if vf < 0 {
		return true
	}
	return false
}

func buildGDPTable(withHeader, header, footer bool) Table {

	var table Table
	if withHeader {
		table = New("Year", "GDP growth", "Inflation")
	} else if header {
		table = New()
		table.AddHeader([]string{"Year", "GDP growth", "Inflation"})
	} else {
		table = New()
	}

	table.AddTitle("Lithuanian GDP growth and inflation time series")
	table.AddTitle("(annual growth rates in percent)")

	table.AddRow("").Insert(1996, 5.1499584261, 2.46181274996)
	table.AddRow("").Insert(1997, 8.2932287195, 8.877890341)
	table.AddRow("").Insert(1998, 7.4671760033, 5.0749342993)
	table.AddRow("").Insert(1999, -1.134642894, 0.7537093994)
	table.AddRow("").Insert(2000, 3.8316671405, 1.0067873303)
	table.AddRow("").Insert(2001, 6.5244308754, 1.3551349535)
	table.AddRow("").Insert(2002, 6.7607495332, 0.2983425414)
	table.AddRow("").Insert(2003, 10.538564772, -1.1457530021)
	table.AddRow("").Insert(2004, 6.5500830275, 1.181321743)
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
	table.AddRow("").Insert(2015, []float64{1.7785761840, 0}, -0.8841084347)
	table.AddRow("typo").Insert(2017, 2.2988432751, 0.9055215537)

	table.AddFootnote("Source: Worldbank")
	table.AddFootnote("Geometric mean was used to calculate the average growth rates")
	table.AddFootnote("Rendered with github.com/vaitekunas/lentele")

	if footer {
		table.AddFooter().Insert("Means:", 4.17, 3.64)
	}

	return table
}

func TestNew(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestNew: this test should never fail/panic")
		}
	}()

	table := buildGDPTable(true, true, true)
	tableNooFoot := buildGDPTable(true, false, false)

	table.Transform(round, "GDP growth", "Inflation")
	table.SetFormat("%-4s", "GDP growth", "Inflation")
	table.SetColumnWidth(15, "GDP growth", "Inflation")

	// Output
	out := bytes.NewBuffer([]byte{})

	table.Render(out, false, true, true, LoadTemplate("classic"))
	table.Render(out, false, true, true, LoadTemplate("smooth"))
	table.Render(out, false, true, true, LoadTemplate("modern"))
	table.Render(out, false, true, true, LoadTemplate("classic"), "Year", "Inflation")
	table.Render(out, true, true, true, LoadTemplate("classic"))
	table.Render(out, false, false, true, LoadTemplate("classic"))
	table.Render(out, false, false, false, LoadTemplate("classic"))
	tableNooFoot.Render(out, false, true, true, LoadTemplate("classic"))
}

func TestTransforms(t *testing.T) {

	stdCols := []string{"GDP growth", "Inflation"}

	tests := []struct {
		trans                         func(interface{}) interface{}
		transCols                     []string
		format                        string
		formatCols                    []string
		width                         int
		widthCols                     []string
		transErr, formatErr, widthErr bool
	}{
		{round, stdCols, "%-5s", stdCols, 15, stdCols, false, false, false},
		{round, []string{"GDP growth"}, "%-5s", []string{"GDP growth"}, 15, []string{"GDP growth"}, false, false, false},
		{round, []string{"Years"}, "%-5s", stdCols, 15, stdCols, true, false, false},
		{round, stdCols, "%-5s", []string{"Years"}, 15, stdCols, false, true, false},
		{round, stdCols, "%-5s", stdCols, 15, []string{"Years"}, false, false, true},
		{round, []string{}, "%-5s", stdCols, 15, stdCols, true, false, false},
		{round, stdCols, "%-5s", []string{}, 15, stdCols, false, true, false},
		{round, stdCols, "%-5s", stdCols, 15, []string{}, false, false, true},
		{round, stdCols, "static value", stdCols, 15, stdCols, false, false, false},
		{round, stdCols, "%-5s", stdCols, 0, stdCols, false, false, false},
	}

	for i, test := range tests {
		table := buildGDPTable(false, true, true)

		if err := table.Transform(test.trans, test.transCols...); (err != nil) != test.transErr {
			if err != nil {
				t.Errorf("TestTransforms: table.Transform: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestTransforms: table.Transform:  test %d failed", i+1)
			}
		}

		if err := table.SetFormat(test.format, test.formatCols...); (err != nil) != test.formatErr {
			if err != nil {
				t.Errorf("TestTransforms: table.SetFormat: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestTransforms: table.SetFormat: test %d failed", i+1)
			}
		}

		if err := table.SetColumnWidth(test.width, test.widthCols...); (err != nil) != test.widthErr {
			if err != nil {
				t.Errorf("TestTransforms: table.SetColumnWidth: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestTransforms: table.SetColumnWidth: test %d failed", i+1)
			}
		}

	}

}

func TestGetRows(t *testing.T) {

	tests := []struct {
		addHeader, addFooter bool
		rowID                int
		rowName              string
		idErr, nameErr       bool
	}{
		{true, true, 1, "", false, false},
		{true, true, 2, "header", false, false},
		{true, true, 3, "footer", false, false},
		{false, true, 2, "header", false, true},
		{true, false, 3, "footer", false, true},
		{true, true, 4, "typo", false, false},
		{true, true, 5, "not a column", false, true},
		{true, true, 50, "", true, false},
	}

	for i, test := range tests {
		table := buildGDPTable(false, test.addHeader, test.addFooter)

		if _, err := table.GetRow(test.rowID); (err != nil) != test.idErr {
			if err != nil {
				t.Errorf("TestGetRows: table.GetRow: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestGetRows: table.GetRow:  test %d failed", i+1)
			}
		}

		if _, err := table.GetRowByName(test.rowName); (err != nil) != test.nameErr {
			if err != nil {
				t.Errorf("TestGetRows: table.GetRowNames: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestGetRows: table.GetRowNames:  test %d failed", i+1)
			}
		}

	}

}

func TestRemoveRows(t *testing.T) {

	tests := []struct {
		rowIDs []int
		rcount int
		isErr  bool
	}{
		{[]int{1, 2, 20, 21}, 19, false},
		{[]int{0, 1, 2, 20, 21}, 23, true},         // header
		{[]int{1, 2, 20, 22}, 23, true},            // footer
		{[]int{21, 21, 21, 21, 21, 21}, 22, false}, // duplicates get removed
		{[]int{30}, 23, true},
		{[]int{}, 23, true},
	}

	for i, test := range tests {
		table := buildGDPTable(false, true, true)

		if err := table.RemoveRows(test.rowIDs...); (err != nil) != test.isErr {
			if err != nil {
				t.Errorf("TestRemoveRows: table.RemoveRows: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestRemoveRows: table.RemoveRows:  test %d failed", i+1)
			}
		}

		if test.isErr {
			continue
		}

		if rcount := table.GetRowCount(); rcount != test.rcount {
			t.Errorf("TestRemoveRows: table.RemoveRows:  test %d failed. Expected %d rows, got %d", i+1, test.rcount, rcount)
		}

	}

}

func TestRemoveRowsByName(t *testing.T) {

	tests := []struct {
		rowNames []string
		rcount   int
		isErr    bool
	}{
		{[]string{"typo"}, 22, false},
		{[]string{"typo", ""}, 2, false},
		{[]string{"header"}, 23, true},
		{[]string{"footer"}, 23, true},
		{[]string{"typo", "typo", "typo"}, 22, false}, // Duplicate names get removed
		{[]string{""}, 3, false},
		{[]string{}, 23, true},
	}

	for i, test := range tests {
		table := buildGDPTable(false, true, true)

		if err := table.RemoveRowsByName(test.rowNames...); (err != nil) != test.isErr {
			if err != nil {
				t.Errorf("TestRemoveRowsByName: table.RemoveRowsByName: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestRemoveRowsByName: table.RemoveRowsByName:  test %d failed", i+1)
			}
		}

		if test.isErr {
			continue
		}

		if rcount := table.GetRowCount(); rcount != test.rcount {
			t.Errorf("TestRemoveRowsByName: table.RemoveRowsByName:  test %d failed. Expected %d rows, got %d", i+1, test.rcount, rcount)
		}

	}

}

func TestFilter(t *testing.T) {

	filter := yearFilter(2000, 2009)

	tests := []struct {
		filterFunc          func(vals ...interface{}) bool
		filterCols          []string
		inplace, keepFooter bool
		isErr               bool
		rowCount            int
	}{
		{filter, []string{"Year"}, false, false, false, 11},
		{filter, []string{"Year"}, false, true, false, 12},
		{badYear, []string{"Inflation"}, false, false, false, 3},
		{badYear, []string{"Inflation"}, false, true, false, 4},
		{filter, []string{"Inflation"}, false, false, false, 1},
		{filter, []string{"No such column"}, false, false, true, 1},
		{filter, []string{"Year"}, true, false, false, 11},
		{filter, []string{}, false, false, true, 1},
		{filter, []string{""}, false, false, true, 1},
	}

	for i, test := range tests {
		table := buildGDPTable(false, true, true)

		filtered, err := table.Filter(test.filterFunc, test.inplace, test.keepFooter, test.filterCols...)
		if (err != nil) != test.isErr {
			if err != nil {
				t.Errorf("TestFilter: table.Filter: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestFilter: table.Filter: test %d failed", i+1)
			}
		}

		if test.isErr {
			continue
		}

		if (filtered.GetRowCount() == table.GetRowCount()) != test.inplace {
			t.Errorf("TestFilter: table.Filter: test %d failed. Inplace filtering failed", i+1)
		}

		if rcount := filtered.GetRowCount(); !test.isErr && rcount != test.rowCount {
			t.Errorf("TestFilter: table.Filter: test %d failed: filtered table should return %d rows, got %d", i+1, test.rowCount, rcount)
		}
	}

}

func TestFilterByRowName(t *testing.T) {

	tests := []struct {
		rowNames            []string
		inplace, keepFooter bool
		rowCount            int
	}{
		{[]string{"typo", "footer", "header"}, false, false, 3},
		{[]string{"typo", "header"}, false, true, 3},
		{[]string{"typo"}, false, true, 3},
		{[]string{"typo", ""}, true, true, 23},
		{[]string{""}, false, true, 22},
		{[]string{""}, false, false, 21},
		{[]string{""}, true, false, 21},
	}

	for i, test := range tests {
		table := buildGDPTable(false, true, true)

		filterFunc := func(colname string) bool {
			for _, expected := range test.rowNames {
				if colname == expected {
					return true
				}
			}
			return false
		}

		filtered := table.FilterByRowNames(filterFunc, test.inplace, test.keepFooter)

		if (filtered.GetRowCount() == table.GetRowCount()) != test.inplace {
			t.Errorf("TestFilterByRowName: table.FilterByRowName: test %d failed. Inplace filtering failed", i+1)
		}

		if rcount := filtered.GetRowCount(); rcount != test.rowCount {
			t.Errorf("TestFilterByRowName: table.FilterByRowName: test %d failed: filtered table should return %d rows, got %d", i+1, test.rowCount, rcount)
		}
	}

}

func TestChangeModify(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestChangeModify: this test should never fail/panic")
		}
	}()

	bold := func(v interface{}) interface{} {
		return color.New(color.Bold).Sprint(v)
	}

	table := buildGDPTable(false, true, true)

	row, err := table.GetRow(1)
	if err != nil {
		t.Errorf("TestChangeModify: should not result in error: %s", err.Error())
	}

	defer row.Change("Year", 1999).Modify(bold, "GDP growth", "Inflation")

}

func TestMisc(t *testing.T) {

	table := buildGDPTable(false, false, false)

	//Adding header and footer
	header := table.AddHeader([]string{"Year", "GDP growth", "Inflation"})
	footer := table.AddFooter()

	if header != table.AddRow("header") {
		t.Errorf("TestMisc: repeatedly adding a header should return the same row")
	}

	if footer != table.AddRow("footer") {
		t.Errorf("TestMisc: repeatedly adding a header footer return the same row")
	}

	// Titles and footnotes
	if err := table.AddTitle("some title"); err != nil {
		t.Errorf("TestMisc: adding a title should not result in an error")
	}

	if err := table.AddTitle(""); err == nil {
		t.Errorf("TestMisc: adding an empty title should result in an error")
	}

	if err := table.AddFootnote("some footnote"); err != nil {
		t.Errorf("TestMisc: adding a footnote should not result in an error")
	}

	if err := table.AddFootnote(""); err == nil {
		t.Errorf("TestMisc: adding an empty footnote should result in an error")
	}

	rownames := table.GetRowNames()
	rowcount := table.GetRowCount()
	lastRownames := rownames[rowcount-3:]
	if len(rownames) != rowcount || lastRownames[0] != "typo" || lastRownames[1] != "header" || lastRownames[2] != "footer" {
		t.Errorf("TestMisc: incorrect column names")
	}
}

func TestMarshalToRichJSON(t *testing.T) {

	table := buildGDPTable(true, true, true)
	jsonedRich := bytes.NewBuffer([]byte{})
	jsonedVanilla := bytes.NewBuffer([]byte{})
	jsonedFake := bytes.NewBufferString(`[{"x":12,"y":13},{"y":15},{"x":16,"y":17},{"x":18}]`)
	jsonedFalse := bytes.NewBufferString("[Not JSON")

	if _, err := table.MarshalToRichJSON(jsonedRich); err != nil {
		t.Errorf("TestMarshalToRichJSON: marshalling a table to rich json should not result in error :%s", err.Error())
	}

	if _, err := table.MarshalToVanillaJSON(jsonedVanilla); err != nil {
		t.Errorf("TestMarshalToRichJSON: marshalling a table to vanilla json should not result in error :%s", err.Error())
	}

	NewFromVanillaJSONWrap := func(r io.Reader) (Table, error) {
		return NewFromVanillaJSON(r, "NA")
	}

	tests := []struct {
		source   io.Reader
		loader   func(io.Reader) (Table, error)
		rowCount int
		isErr    bool
	}{
		{jsonedRich, NewFromRichJSON, table.GetRowCount(), false},
		{jsonedVanilla, NewFromVanillaJSONWrap, table.GetRowCount() - 1, false}, // no footer
		{jsonedRich, NewFromVanillaJSONWrap, 0, true},
		{jsonedVanilla, NewFromRichJSON, 0, true},
		{jsonedFake, NewFromVanillaJSONWrap, 5, false},
		{jsonedFake, NewFromRichJSON, 5, true},
		{jsonedFalse, NewFromRichJSON, 0, true},
		{jsonedFalse, NewFromVanillaJSONWrap, 0, true},
	}

	for i, test := range tests {
		table, err := test.loader(test.source)
		if (err != nil) != test.isErr {
			if err != nil {
				t.Errorf("TestMarshalToRichJSON: test %d failed: %s", i+1, err.Error())
			} else {
				t.Errorf("TestMarshalToRichJSON: test %d failed", i+1)
			}
		}
		if test.isErr {
			continue
		}
		if rcount := table.GetRowCount(); rcount != test.rowCount {
			t.Errorf("TestMarshalToRichJSON: test %d failed: expected %d rows, got %d", i+1, test.rowCount, rcount)
		}
	}

}

func TestGoroutines(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestGoroutines: this test should never fail/panic")
		}
	}()

	routines := 1000
	rowsPerRoutine := 10000

	addSomeRows := func(table Table, ready chan<- bool) {
		for i := 1; i <= rowsPerRoutine; i++ {
			table.AddRow("").Insert(i, 0, 0)
		}
		ready <- true
	}

	table := buildGDPTable(true, true, true)

	ready := make(chan bool, routines)
	for j := 1; j <= routines; j++ {
		go addSomeRows(table, ready)
	}

	for k := 1; k <= routines; k++ {
		<-ready
	}

	expected := 23 + routines*rowsPerRoutine
	if rcount := table.GetRowCount(); rcount != expected {
		t.Errorf("TestGoroutines: failed adding some rows: expected %d, got %d:", expected, rcount)
	}

}
