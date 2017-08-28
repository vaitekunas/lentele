package lentele

import (
	"fmt"
	"os"
	"testing"
	"unicode/utf8"
)

func trim(s string) string {
	if utf8.RuneCountInString(s) > 7 {
		return fmt.Sprintf("%s...", s[:7])
	}
	return s
}

func modify(v interface{}) interface{} {
	return fmt.Sprintf("[%v]", v)
}

func TestNew(t *testing.T) {

	classic := LoadTemplate("classic")

	table := New("ID", "Client", "Amount")
	table.AddRow("").Insert(1, "Acme", 100)
	table.AddRow("").Insert(2, "IronMountain", 200)
	table.AddRow("").Insert(3, "E corp", 300).Modify(modify, "ID", "Client")

	table.Render(os.Stdout, true, true, classic)
}
