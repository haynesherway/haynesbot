package haynesbot

import (
    //"fmt"
    "image"
    "image/color"
    "image/png"
    //"errors"
    "os"
    "strconv"
    
    "github.com/haynesherway/pngtable"
    "github.com/haynesherway/pogo"
)

var (
    BLACK = image.Black
    GREEN = color.RGBA{0,0x99,0,0xff}
    GREEN1 = color.RGBA{0x66, 0x99, 0, 0xff}
    GREEN2 = color.RGBA{0x99, 0x99, 0, 0xff}
    GREEN3 = color.RGBA{0xcc, 0xaa, 0, 0xff}
    GREEN4 = color.RGBA{0xcc, 0x66, 0, 0xff}
)

func PercentColor(perc int) color.Color {
    if perc == 100 {
        return GREEN
    } else if perc >= 98 {
        return GREEN1
    } else if perc >= 96 {
        return GREEN2
    } else if perc >= 93 {
        return GREEN3
    } else if perc >= 91 {
        return GREEN4
    } else {
        return BLACK
    }
}

func GetTable(data interface{}, fileName string) *os.File {
    /*table := pngtable.New()
    table.Options.SetColWidths([]int{35,35,35,45})
    table.Options.SetRowHeight(25)
    table.Options.SetFontSize(15)
    if ivList, ok := data.([]pogo.IVStat); ok {
        table.SetHeaders([]string{"A", "D", "S", "%"})
        for _, iv := range ivList {
            a := strconv.Itoa(iv.Attack)
            d := strconv.Itoa(iv.Defense)
            s := strconv.Itoa(iv.Stamina)
            p := strconv.Itoa(iv.Percent)
            table.AddRow([]string{a, d, s, p}).SetBackground(PercentColor(iv.Percent))
        }
    }*/
    
    table := pngtable.New()
    table.Options.SetRowHeight(15)
    table.Options.SetFontSize(10)
    
    if ivList, ok := data.([]pogo.IVStat); ok {
        table.SetHeaders([]string{"IV%", "A", "D", "S", "CP@20", "CP@25"})
        for _, iv := range ivList {
            if iv.Percent < 90 {
                continue
            }
            a := strconv.Itoa(iv.Attack)
            d := strconv.Itoa(iv.Defense)
            s := strconv.Itoa(iv.Stamina)
            p := strconv.Itoa(iv.Percent) + "%"
            cp20 := strconv.Itoa(iv.CP20)
            cp25 := strconv.Itoa(iv.CP25)
            table.AddRow([]string{p, a, d, s, cp20, cp25}).SetBackground(PercentColor(iv.Percent))
        }
    }
    table.Options.SetColWidths([]int{35, 20, 20, 20, 50, 50})
    table.Draw()
    f, err := os.Create(ImageServer + "/" + fileName)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    err = png.Encode(f, table.Image)
    if err != nil {
        panic(err)
    }
    return f
}