package haynesbot

import (
    "strconv"
    "testing"
    "github.com/haynesherway/pngtable")

func TestTableColors(t *testing.T) {
    table := pngtable.New()
    for perc := 100; perc >= 0; perc-- {
        table.AddRow([]string{strconv.Itoa(perc)}).SetBackground(PercentColor(perc)).SetColor(BLACK)
    }
    table.Draw()
    return
}