package lentele

import (
	"fmt"
	"github.com/fatih/color"
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

	evil := func(v interface{}) interface{} {
		return color.New(color.FgRed).Add(color.Bold).Sprint(v)
	}

	table := New("ID", "Client", "Amount")
	table.AddRow("").Insert(1, fmt.Sprintf("%-12s", "Acme"), 100)
	table.AddRow("").Insert(2, fmt.Sprintf("%-12s", "IronMountain"), 200)
	table.AddRow("").Insert(3, fmt.Sprintf("%-12s", "E corp"), 300).Modify(evil, "Client")
	table.AddFooter().Insert("", fmt.Sprintf("%12s", "Total"), "600")

	table.Render(os.Stdout, false, true, LoadTemplate("classic"))
	//table.Render(os.Stdout, false, true, LoadTemplate("smooth"))
	//table.Render(os.Stdout, false, true, LoadTemplate("modern"))
}
