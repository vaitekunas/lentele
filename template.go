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
func (t *template) SetDisplayOptions(center bool) {
	t.Lock()
	defer t.Unlock()

	t.Center = center
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
		} else {
			lines = append(lines, title)
		}
	}

	lines = append(lines, "")

	return lines
}

// renderL1L2L3 renders a template line
func renderL1L2L3(T1 [4]string, T2 [3]string, T3 [4]string, widths []int, mcells, pcells []string, center bool) (L1 string, L2 string, L3 string, isEmpty bool) {

	var tlsum int
	lines := newLines(pcells)
	L2Slice := []string{}
	for line := 1; line <= lines; line++ {

		L1 = T1[0]
		L2 = T2[0]
		L3 = T3[0]

		tlsum = 1
		for i, width := range widths {

			// Skip irrelevant columns
			if width == 0 {
				continue
			}

			// Cell values and spacing
			value, sp1, sp2, tl := measure(i, width, mcells, pcells)

			// Cell lines and prelines (empty lines)
			ilines := strings.Count(value, "\n") + 1
			prelines := 0
			if lines > 1 {
				prelines = int((lines - ilines) / 2)
			}

			// Upper border
			L1 += strings.Repeat(T1[1], width+2)

			// Prelines, lines, postlines
			if line <= prelines || line > prelines+ilines {
				L2 += fmt.Sprintf("%s",strings.Repeat(" ",width+2))
			}else{
				valueParts := strings.Split(value, "\n")
				sp1Parts := strings.Split(sp1, "\n")
				sp2Parts := strings.Split(sp2, "\n")
				iline := line-prelines-1
				L2 += fmt.Sprintf("%s%s%s", sp1Parts[iline], valueParts[iline], sp2Parts[iline])
			}

			// Bottom border
			L3 += strings.Repeat(T3[1], width+2)

			// Cell walls to the right
			if i != len(widths)-1 {
				L1 += T1[2]
				L2 += T2[1]
				L3 += T3[2]
			} else {
				L1 += T1[3]
				L2 += T2[2]
				L3 += T3[3]
			}

			if len(value) > 0 {
				isEmpty = false
			}

			// Calculate row width
			tlsum += tl + 1
		}

		if line <= lines {
			if center {
			L2Slice = append(L2Slice, fmt.Sprintf("%s%s", strings.Repeat(" ", getOffset(tlsum)), L2))
		}else{
			L2Slice = append(L2Slice, L2)
		}
			L2 = T2[0]
		}

	}

	if center {
		L1 = centerStr(L1)
		L2 = strings.Join(L2Slice,"\n")
		L3 = centerStr(L3)
	}

	return L1, L2, L3, isEmpty

}

// mesure mesaures string widths and returns printable strings
func measure(i, width int, mcells, pcells []string) (string, string, string, int) {

	pvalue := ""
	mvalue := ""
	if i < len(pcells) {
		pvalue = pcells[i]
		mvalue = mcells[i]
	}

	sp1Slice := []string{}
	sp2Slice := []string{}

	totalLen := 0
	mvalueParts := strings.Split(mvalue, "\n")
	for _, mpart := range mvalueParts {
		vlen := utf8.RuneCountInString(mpart)
		reps := int((width + 2 - vlen) / 2)
		sp1Slice = append(sp1Slice, strings.Repeat(" ", reps))
		sp2Slice = append(sp2Slice, strings.Repeat(" ", width+2-vlen-reps))

		if width+2 > totalLen {
			totalLen = width + 2
		}
	}

	return pvalue, strings.Join(sp1Slice, "\n"), strings.Join(sp2Slice, "\n"), totalLen
}

// newLines returns the max number of linebreaks (measured as \n) in a row
func newLines(pcells []string) int {
	lines := 1
	for _, cell := range pcells {
		if n := strings.Count(cell, "\n"); n+1 > lines {
			lines = n + 1
		}
	}
	return lines
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
