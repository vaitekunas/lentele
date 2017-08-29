package lentele

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
)

// template implements the Template interface
//
// +====+====+====+====+ <- H1: [4]string{"+","=","+","+"}
// |    |    |    |    | <- H2: [1]string{"|"}
// +====+====+====+====+ <- H3: [4]string{"+","=","+","+"}
// ┌────┬────┬────┬────┐ <- C1: [4]string{"┌", "─", "┬", "┐"},
// │    │    │    │    │ <- C2: [1]string{"│"},
// └────┴────┴────┴────┘ <- C3: [4]string{"└", "─", "┴", "┘"},
// +====+====+====+====+ <- F1: [4]string{"+","=","+","+"}
// |    |    |    |    | <- F2: [1]string{"|"}
// +====+====+====+====+ <- F3: [4]string{"+","=","+","+"}
//
type template struct {
	*sync.Mutex
	ColWidths []int
	SkipH1, SkipH3, SkipC1,
	SkipC3, SkipF1, SkipF3 bool
	SkipFirstC1, SkipLastC3 bool
	H1                      [4]string
	H2                      [3]string
	H3                      [4]string
	C1                      [4]string
	C2                      [3]string
	C3                      [4]string
	F1                      [4]string
	F2                      [3]string
	F3                      [4]string
	HR                      string
}

// SetColumnWidths sets the column widths
func (t *template) SetColumnWidths(widths []int) {
	t.Lock()
	defer t.Unlock()

	t.ColWidths = widths
}

// RenderHeader renders the header row
func (t *template) RenderHeader(mcells, pcells []string) []string {
	t.Lock()
	defer t.Unlock()

	L1 := t.H1[0]
	L2 := t.H2[0]
	L3 := t.H3[0]
	for i, width := range t.ColWidths {

		value, sp1, sp2 := measure(i, width, mcells, pcells)

		L1 += strings.Repeat(t.H1[1], width+2)
		L2 += fmt.Sprintf("%s%s%s", sp1, value, sp2)
		L3 += strings.Repeat(t.H3[1], width+2)

		if i != len(t.ColWidths)-1 {
			L1 += t.H1[2]
			L2 += t.H2[1]
			L3 += t.H3[2]
		} else {
			L1 += t.H1[3]
			L2 += t.H2[2]
			L3 += t.H3[3]
		}
	}

	lines := []string{}
	if !t.SkipH1 {
		lines = append(lines, L1)
	}
	lines = append(lines, L2)
	if !t.SkipH3 {
		lines = append(lines, L3)
	}

	return lines

}

// RenderRow renders a regular row
func (t *template) RenderRow(row, rows int, mcells, pcells []string) []string {
	t.Lock()
	defer t.Unlock()

	L1 := t.C1[0]
	L2 := t.C2[0]
	L3 := t.C3[0]
	for i, width := range t.ColWidths {

		value, sp1, sp2 := measure(i, width, mcells, pcells)

		L1 += strings.Repeat(t.C1[1], width+2)
		L2 += fmt.Sprintf("%s%s%s", sp1, value, sp2)
		L3 += strings.Repeat(t.C3[1], width+2)

		if i != len(t.ColWidths)-1 {
			L1 += t.C1[2]
			L2 += t.C2[1]
			L3 += t.C3[2]
		} else {
			L1 += t.C1[3]
			L2 += t.C2[2]
			L3 += t.C3[3]
		}
	}

	lines := []string{}
	if !t.SkipC1 && (row != 1 || !t.SkipFirstC1) {
		lines = append(lines, L1)
	}
	lines = append(lines, L2)

	if !t.SkipC3 && (row != rows || !t.SkipLastC3) {
		lines = append(lines, L3)
	}

	return lines

}

// RenderFooter renders the footer row
func (t *template) RenderFooter(mcells, pcells []string) []string {
	t.Lock()
	defer t.Unlock()

	any := false
	L1 := t.F1[0]
	L2 := t.F2[0]
	L3 := t.F3[0]
	for i, width := range t.ColWidths {

		value, sp1, sp2 := measure(i, width, mcells, pcells)

		L1 += strings.Repeat(t.F1[1], width+2)
		L2 += fmt.Sprintf("%s%s%s", sp1, value, sp2)
		L3 += strings.Repeat(t.F3[1], width+2)

		if i != len(t.ColWidths)-1 {
			L1 += t.F1[2]
			L2 += t.F2[1]
			L3 += t.F3[2]
		} else {
			L1 += t.F1[3]
			L2 += t.F2[2]
			L3 += t.F3[3]
		}
		if len(value) > 0 {
			any = true
		}
	}

	lines := []string{}
	if !t.SkipF1 {
		lines = append(lines, L1)
	}
	if any {
		lines = append(lines, L2)
	}

	if !t.SkipF3 && any {
		lines = append(lines, L3)
	}

	return lines
}

// RRenderTitle renders the title
func (t *template) RenderTitle(title string) []string {
	return []string{"", title, ""}
}

// RenderFootnotes renders footnotes
func (t *template) RenderFootnotes(footnotes []string) []string {

	lines := []string{"", "<HR>"}

	longest := 0
	for i, note := range footnotes {
		formatted := fmt.Sprintf("%d. %s", i+1, note)
		if length := utf8.RuneCountInString(formatted); length > longest {
			longest = length
		}
		lines = append(lines, formatted)
	}

	lines[1] = strings.Repeat(t.HR, longest)
	lines = append(lines, "", "")

	return lines
}

// mesure mesaures string widths and returns printable strings
func measure(i, width int, mcells, pcells []string) (string, string, string) {

	pvalue := ""
	mvalue := ""
	if i < len(pcells) {
		pvalue = pcells[i]
		mvalue = mcells[i]
	}

	vlen := utf8.RuneCountInString(mvalue)
	reps := int((width + 2 - vlen) / 2)
	sp1 := strings.Repeat(" ", reps)
	sp2 := strings.Repeat(" ", width+2-vlen-reps)

	return pvalue, sp1, sp2
}

// Classic template
func tmplClassic() *template {

	return &template{
		Mutex:      &sync.Mutex{},
		ColWidths:  []int{},
		SkipC1:     true,
		SkipLastC3: true,
		SkipF3:     true,
		H1:         [4]string{"╔", "═", "╦", "╗"},
		H2:         [3]string{"║", "║", "║"},
		H3:         [4]string{"╠", "═", "╩", "╣"},
		C1:         [4]string{"╟", "─", "┼", "╢"},
		C2:         [3]string{"║", "│", "║"},
		C3:         [4]string{"╟", "─", "┼", "╢"},
		F1:         [4]string{"╚", "═", "╧", "╝"},
		F2:         [3]string{" ", " ", " "},
		F3:         [4]string{"", "", "", ""},
		HR:         "─",
	}
}

// Smooth template
func tmplSmooth() *template {

	return &template{
		Mutex:      &sync.Mutex{},
		ColWidths:  []int{},
		SkipC1:     true,
		SkipLastC3: true,
		H1:         [4]string{"┌", "─", "┬", "┐"},
		H2:         [3]string{"│", "│", "│"},
		H3:         [4]string{"├", "─", "┼", "┤"},
		C1:         [4]string{"├", "─", "┼", "┤"},
		C2:         [3]string{"│", "│", "│"},
		C3:         [4]string{"├", "─", "┼", "┤"},
		F1:         [4]string{"├", "─", "┴", "┤"},
		F2:         [3]string{"│", " ", "│"},
		F3:         [4]string{"└", "─", "─", "┘"},
		HR:         "─",
	}
}

// Modern template
func tmplModern() *template {

	return &template{
		Mutex:      &sync.Mutex{},
		ColWidths:  []int{},
		SkipH1:     true,
		SkipC1:     true,
		SkipLastC3: true,
		H1:         [4]string{" ", " ", " ", " "},
		H2:         [3]string{" ", " ", " "},
		H3:         [4]string{"━", "━", "━", "━"},
		C1:         [4]string{" ", " ", " ", " "},
		C2:         [3]string{" ", " ", " "},
		C3:         [4]string{" ", " ", " ", " "},
		F1:         [4]string{"━", "━", "━", "━"},
		F2:         [3]string{" ", " ", " "},
		F3:         [4]string{" ", " ", " ", " "},
		HR:         "─",
	}
}
