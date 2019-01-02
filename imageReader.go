package haynesbot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract"
)

var (
	cpRegexp = regexp.MustCompile(`CP[0-9]{2,4}`)
	hpRegexp = regexp.MustCompile(`[0-9]{2,3}HP`)
)

func ReadImage(att *discordgo.MessageAttachment) {
	client := gosseract.NewClient()
	defer client.Close()

	path, err := ProcessAttachment(att)
	if err != nil {
		fmt.Println(err)
	}
	client.SetImage(path)
	client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZ,/0123456789")
	text, _ := client.Text()

	cp := strings.Replace(cpRegexp.FindString(text), "CP", "", -1)
	hp := strings.Replace(hpRegexp.FindString(text), "HP", "", -1)
	fmt.Println(text)
	fmt.Println(cp, hp)
	return
}

func ProcessAttachment(att *discordgo.MessageAttachment) (string, error) {
	response, _ := http.Get(att.ProxyURL)
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("/tmp/%s.png", att.ID))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
