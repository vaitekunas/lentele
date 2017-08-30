package lentele

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
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

	Center bool

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

// SetDisplayOptions sets some display options
func (t *template) SetDisplayOptions(center bool){
	t.Lock()
	defer t.Unlock()

	t.Center = center
}

// renderL1L2L3 renders a template line
func renderL1L2L3(T1[4]string, T2[3]string, T3[4]string, widths []int, mcells, pcells []string, center bool) (L1 string, L2 string, L3 string, isEmpty bool) {

	L1 = T1[0]
	L2 = T2[0]
	L3 = T3[0]

	tlsum := 1
	for i, width := range widths {

		value, sp1, sp2, tl := measure(i, width, mcells, pcells)

		L1 += strings.Repeat(T1[1], width+2)
		L2 += fmt.Sprintf("%s%s%s", sp1, value, sp2)
		L3 += strings.Repeat(T3[1], width+2)

		if i != len(widths)-1 {
			L1 += T1[2]
			L2 += T2[1]
			L3 += T3[2]
		} else {
			L1 += T1[3]
			L2 += T2[2]
			L3 += T3[3]
		}
		tlsum += tl + 1

		if len(value) > 0 {
			isEmpty = false
		}
	}

	if center {
		L1 = centerStr(L1)
		L2 = fmt.Sprintf("%s%s", strings.Repeat(" ", getOffset(tlsum)), L2)
		L3 = centerStr(L3)
	}

	return L1, L2, L3, isEmpty

}

// RenderHeader renders the header row
func (t *template) RenderHeader(mcells, pcells []string) []string {
	t.Lock()
	defer t.Unlock()

	// Render lines
	L1, L2, L3, _ := renderL1L2L3(t.H1, t.H2, t.H3, t.ColWidths, mcells, pcells, t.Center)

	// Append or skip
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

	// Render lines
	L1, L2, L3, _ := renderL1L2L3(t.C1, t.C2, t.C3, t.ColWidths, mcells, pcells, t.Center)

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

	// Render lines
	L1, L2, L3, isEmpty := renderL1L2L3(t.F1, t.F2, t.F3, t.ColWidths, mcells, pcells, t.Center)

	lines := []string{}
	if !t.SkipF1 {
		lines = append(lines, L1)
	}
	if !isEmpty {
		lines = append(lines, L2)
	}

	if !t.SkipF3 && !isEmpty {
		lines = append(lines, L3)
	}

	return lines
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

// RenderTitles renders the title
func (t *template) RenderTitles(titles []string) []string {
	t.Lock()
	defer t.Unlock()

	lines := []string{""}

	for _, title := range titles {
		if t.Center {
			lines = append(lines, centerStr(title))
		}else{
			lines = append(lines, title)
		}
	}

	lines = append(lines, "")

	return lines
}

// mesure mesaures string widths and returns printable strings
func measure(i, width int, mcells, pcells []string) (string, string, string, int) {

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

	totalLen := width + 2

	return pvalue, sp1, sp2, totalLen
}

// getOffset returns the available tty space
func getOffset(width int) int {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	offset := int((w - width) / 2)
	if offset < 0 {
		return 0
	}
	return offset
}

// centerStr centers a string
func centerStr(value string) string {
	width := utf8.RuneCountInString(value)
	offset := getOffset(width)

	return fmt.Sprintf("%s%s", strings.Repeat(" ", offset), value)
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
		F3:         [4]string{" ", " ", " ", " "},
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
