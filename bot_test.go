package haynesbot

import (
    "fmt"
    //"strings"
    //"testing"
    
    "github.com/bwmarrin/discordgo"
)

const (
    TEST_CHANNEL_ID = "402885030994509835"
    )

func init() {
    test = true
    
    err := ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()
	
	goBot.State.ChannelAdd(&discordgo.Channel{
	    Name: "#haynes-bot-testing",
	    GuildID: "341978052869226496",
	    ID: "402885030994509835",
	})
}

/*func TestBotResponse_PrintInfoToDiscord(t *testing.T) {
    message := &discordgo.Message{
        Content: "!wat",
        ChannelID: TEST_CHANNEL_ID,
    }
    m := &discordgo.MessageCreate{message}
    bot := NewBotResponse(goBot, m, strings.Fields(m.Content))
    bot.PrintInfoToDiscord()
    return
}*/