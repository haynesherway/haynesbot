package haynesbot

import (
    "fmt"
    "io"
    "image"
    "image/color"
    "image/png"
    //"errors"
    "net/http"
    "os"
    "strconv"
    
    "github.com/haynesherway/pngtable"
    "github.com/haynesherway/pogo"
)

// Default colors
var (
    BLACK = image.Black
    WHITE = image.White
)

// PercentColor returns a color based on a percentage
func PercentColor(x int) color.Color {
    m := 12
    green := (m * (x * 255/100)) - (255 * (m-2))
    red := m * ((100 - x) * 255/100)

    blue := uint8(0)
    
    tp := 95
    if green > 255 && x > tp {
        green = 255 + (m * (tp - x))
    } else if green > 255 {
        green = 255
    } else if green < 0 {
        green = 0
    }
    if red > 255 {
        red = 255
    }
    
    return color.RGBA{uint8(red), uint8(green), blue, uint8(255)}
}

// GetTable gets a png table based on the pokemon stats
func GetTable(p *pogo.Pokemon, data interface{}, fileName string) *os.File {
    table := pngtable.New()
    
    //Get image
    f, err := Download(p.API.Sprites.Front, p.Name)
    img, _, err := image.Decode(f)
    if err != nil {
        table.SetTitle(fmt.Sprintf("%s - Raid CP Chart", p.Name)).SetColor(WHITE).SetBackground(BLACK)
    } else {
        table.SetTitlePicture(img).SetColor(WHITE).SetBackground(BLACK).SetHeight(100)
    }
   
    table.Options.SetRowHeight(15)
    table.Options.SetFontSize(10)
    table.Options.SetBorderColor(BLACK)
    
    if ivList, ok := data.([]pogo.IVStat); ok {
        table.SetHeaders([]string{"IV%", "A", "D", "S", "CP@20", "CP@25"}).SetColor(WHITE)
        for _, iv := range ivList {
            if iv.Percent < 88 {
                continue
            }
            a := strconv.Itoa(iv.Attack)
            d := strconv.Itoa(iv.Defense)
            s := strconv.Itoa(iv.Stamina)
            p := strconv.Itoa(iv.Percent) + "%"
            cp20 := strconv.Itoa(iv.CP20)
            cp25 := strconv.Itoa(iv.CP25)
            table.AddRow([]string{p, a, d, s, cp20, cp25}).SetBackground(PercentColor(iv.Percent)).SetColor(BLACK)
        }
    }
    table.Options.SetColWidths([]int{35, 20, 20, 20, 50, 50})
    table.Draw()
    f, err = os.Create(ImageServer + "/" + fileName)
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

// Download downloads a sprite image to a file to be added to the png
func Download(s string, n string) (f *os.File, err error) {
    n += ".png"
    
    if !ImageExists(n) {
    
        response, err := http.Get(s)
        if err != nil {
            return f, err
        }
        defer response.Body.Close()
        
        f, err = os.Create(ImageServer + "/" + n)
        if err != nil {
            return f, err
        }
        _, err = io.Copy(f, response.Body)
        if err != nil {
            return f, err
        }
        f.Close()
    }
    
    if f, err = os.Open(ImageServer + "/" + n); err == nil{
        return f, err
    }

    
    return f, nil
}