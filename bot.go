package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"haynes-bot/pokemon"
	"strconv"
	"strings"
)

const INFO = 
`**---- HAYNES-BOT ----**
Commands:
	**!cp** {pokemon} {level} {attack iv} {defense iv} {stamina iv}
		Get CP of a pokemon at a specified level with specified IVs
		Example: !cp mewtwo 25 15 14 15
	**!maxcp** {pokemon}
		Get maximum CP of a pokemon with perfect IVs at level 40
		Example: !maxcp latios
	**!raidcp** {pokemon}
		Get range of possible raid CPs for specified pokemon
		Example: !raidcp groudon
	**!raidcp** {pokemon} {cp}
		Get possible IV combinations for specified raid pokemon with specified IV
		Example: !raidcp kyogre 2292
	**!raidchart** {pokemon}
		Get a chart with possible stats for specified pokemon at raid level above 90%
		Example: !raidchart machamp
	**!moves** {pokemon}
		Get a list of fast and charge moves for specified pokemon
		Example: !moves rayquaza
	**!type** {pokemon}
		Get a list of types for a specified pokemon
		Example: !type rayquaza
	**!effect** {pokemon|type}
		Get a list of type relations a specified pokemon or type has
		Example: !effect pikachu or !effect electric`
		
var easterEggs = map[string]string{
    "haynes-bot": "meow",
    "haynes": "The power of haynes is limitless",
    "haynesherway22": "The power of haynes is limitless",
    "haynesherway": "The power of haynes is limitless",
    "alletzhauser": "Alletzhauser is a poopy pants",
}


var (
	BotID string
	goBot *discordgo.Session
	ivCalculator *pokemon.IVCalculator
)

var ivCalcDiscordChannels = map[string]string{
	"405409275410776076": "choppa iv calculator",
	"402885030994509835": "choppa testing",
	"342018683024965635": "raid chat",
}

var (
	ERR_CP_COMMAND = errors.New("CP command needs to be formatted like this: !maxcp {pokemon} {level} {attack iv} {defense iv} {stamina iv}")
	ERR_IV_COMMAND = errors.New("IV command needs to be formatted like this: !iv {pokemon} {cp} {level} or !iv {pokemon} {cp}")
	ERR_RAIDCP_COMMAND = errors.New("Raid CP command needs to be formatted like this: !raidcp {pokemon} or !raidcp {pokemon} {cp}")
	ERR_RAIDCHART_COMMAND = errors.New("Raid CP Chart command needs to be formatted like this: !raidcpchart {pokemon}")
	ERR_MAXCP_COMMAND = errors.New("Max CP command needs to be formatted like this: !maxcp {pokemon}")
	ERR_MOVES_COMMAND = errors.New("Moves command needs to be formatted like this: !moves {pokemon}")
	ERR_TYPES_COMMAND = errors.New("Types command needs to be formatted like this: !type {pokemon}")
	ERR_TYPECHART_COMMAND = errors.New("Effect command needs to be formatted like this: !effect {pokemon}")
	ERR_NO_COMBINATIONS = errors.New("No possible IV combinations for that CP")
	ERR_NO_STATS = errors.New("Pokemon Master file doesn't have stats for that pokemon yet :(")
	ERR_POKEMON_UNRECOGNIZED = errors.New("Pokemon not recognized.")
	ERR_COMMAND_UNRECOGNIZED = errors.New("Command not recognized")
	
	ERR_WRONG_CHANNEL = errors.New("IV Calculator must be used in the designated channel.")
)

type botResponse struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
	command string
	fields []string
	err error
}

type botCommand struct {
	name string
	err error
}

func NewBotResponse(s *discordgo.Session, m *discordgo.MessageCreate, fields []string) *botResponse {
	return &botResponse{s: s, m: m, fields: fields}
}

func (b *botResponse) GetCommandName() string {
	if len(b.fields) == 0 {
		b.err = ERR_COMMAND_UNRECOGNIZED
		return ""
	}
	cmd := strings.Replace(strings.Replace(b.fields[0], "!", "", 1), "?", "", 1)
	return cmd
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	err = goBot.UpdateStatus(0, "!haynes-bot")
	if err != nil {
		fmt.Println("Unable to update status: ", err.Error())
	}
	
	ivCalculator = pokemon.StartIVCalculator(goBot)

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, config.BotPrefix) {
		if ivCalculator.IsRunning(m.Author.ID) {
			ivCalculator.InputChannel <- m
		}
		return
	}

	if m.Author.ID == BotID {
		return
	}

	fmt.Println(m.Content)
	
	bot := NewBotResponse(s, m, strings.Fields(m.Content))
	cmd := bot.GetCommandName()
	if bot.err != nil {
		return
	}
	
	if len(bot.fields) > 1 {
		if egg, ok := easterEggs[strings.ToLower(bot.fields[1])]; ok {
			bot.PrintToDiscord(egg)
			return
		}
	}

	switch cmd {
	case "cp":
		err := bot.PrintCPToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "ivcalc":
		if _, ok := ivCalcDiscordChannels[m.ChannelID]; !ok {
			bot.PrintErrorToDiscord(ERR_WRONG_CHANNEL)
		} else {
			err := bot.StartIVCalculation(ivCalculator.InputChannel)
			if err != nil {
				bot.PrintErrorToDiscord(err)
			}
		}
	case "iv":
		err := bot.PrintIVToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "raidcp":
		err := bot.PrintRaidCPToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "raidchart": 
		err := bot.PrintRaidChartToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "maxcp":
		err := bot.PrintMaxCPToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "moves":
		err := bot.PrintMovesToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "types", "type":
		err := bot.PrintTypeToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "effect":
		err := bot.PrintTypeChartToDiscord()
		if err != nil {
			bot.PrintErrorToDiscord(err)
		}
	case "haynes-bot", "haynez-bot":
		bot.PrintInfoToDiscord()
	}
	
	return
}

func (b *botResponse) PrintInfoToDiscord() error {
	b.PrintToDiscord(INFO)
	return nil
}

func (b *botResponse) StartIVCalculation(calcChan chan interface{}) (error) {
	calcChan <- b.m
	
	return nil
}

func (b *botResponse) PrintIVToDiscord() error {
	if len(b.fields) < 3 {
		return &botError{ERR_IV_COMMAND, ""}
	}
	
	pokemonName := strings.ToLower(b.fields[1])
	
	cp, err := strconv.Atoi(b.fields[2])
	if err != nil {
		return &botError{ERR_IV_COMMAND, ""}
	}
	
	level := 0.0
	stardust := 0
	if len(b.fields) > 3 {
			
			val, err := strconv.ParseFloat(b.fields[3], 64)
			if err != nil {
				return &botError{ERR_IV_COMMAND, ""}
			}
			if val <= 40.0 {
				level = val
			} else {
				stardust = int(val)
			}
	}
	
	bestvals := ""
	if len(b.fields) > 4 {
		if strings.Contains(b.fields[4], "a") {
			bestvals += "a"
		}
		if strings.Contains(b.fields[4], "d") {
			bestvals += "d"
		}
		if strings.Contains(b.fields[4], "h") {
			bestvals += "h"
		}
	}
	
	if p, ok := pokemon.PokemonMap[pokemonName]; ok {
	    ivChart := p.GetIV(cp, level, stardust, bestvals)
		if len(ivChart) == 0 {
	            return &botError{ERR_NO_COMBINATIONS, p.Name}
	        } else {
	            b.PrintToDiscord(ivChart)
	        }
	} else {
	    return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}
	
	return nil
}

func (b *botResponse) PrintCPToDiscord() error {
	if len(b.fields) < 6 {
		return &botError{ERR_CP_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])
	
	level, err := strconv.ParseFloat(b.fields[2], 64)
	ivA, err := strconv.Atoi(b.fields[3])
	ivD, err := strconv.Atoi(b.fields[4])
	ivS, err := strconv.Atoi(b.fields[5])
	if err != nil {
		return &botError{ERR_CP_COMMAND, ""}
	}

	if p, ok := pokemon.PokemonMap[pokemonName]; ok {
	    cp := p.GetCP(level, ivA, ivD, ivS)
	    message := fmt.Sprintf("CP for **%s** at level %d with IVs %d/%d/%d is %d", p.Name, level, ivA, ivD, ivS, cp)
	    b.PrintToDiscord(message)
	} else {
	    return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

func (b *botResponse) PrintMaxCPToDiscord() error {
    if len(b.fields) < 2 {
		return &botError{ERR_MAXCP_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])
	
	if pokemon, ok := pokemon.PokemonMap[pokemonName]; ok {
	    maxcp := pokemon.GetMaxCP()
	    if maxcp == 0 {
	    	return &botError{ERR_NO_STATS, pokemon.Name}
	    }
	    b.PrintToDiscord(fmt.Sprintf("Max CP for %s is %d", pokemon.Name, pokemon.GetMaxCP()))
	} else {
	    return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}
	
    return nil
}

func (b *botResponse) PrintRaidChartToDiscord() error {
    if len(b.fields) < 2 {
		return &botError{ERR_RAIDCHART_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])
	
	if p, ok := pokemon.PokemonMap[pokemonName]; ok {
		b.PrintToDiscord(p.GetRaidCPChart())
	} else {
	    return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

func (b *botResponse) PrintRaidCPToDiscord() error {
	if len(b.fields) < 2 {
		return &botError{ERR_RAIDCP_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	if p, ok := pokemon.PokemonMap[pokemonName]; ok {
	    if len(b.fields) == 2 {
	    	b.PrintToDiscord(p.GetRaidCPRange())
	    } else {
	        cp, err := strconv.Atoi(b.fields[2])
	        if err != nil {
	            return &botError{ERR_RAIDCP_COMMAND, ""}
	        }
	        ivChart := p.GetRaidIV(cp)
	        if len(ivChart) == 0 {
	            return &botError{ERR_NO_COMBINATIONS, p.Name}
	        } else {
	            b.PrintToDiscord(ivChart)
	        }
	    }
	} else {
	    return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

func (b *botResponse) PrintMovesToDiscord() error {
    if len(b.fields) < 2 {
		return &botError{ERR_MOVES_COMMAND, ""}
	}
	
	pokemonName := strings.ToLower(b.fields[1])
	
	if p, ok := pokemon.PokemonMap[pokemonName]; ok {
	    message := fmt.Sprintf("Moves for **%s**:\nFast: %s\nCharge: %s", p.Name, p.Moves.Fast.Print(), p.Moves.Charge.Print())
	    b.PrintToDiscord(message)
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}
	
	return nil
}

func (b *botResponse) PrintTypeToDiscord() error {
	if len(b.fields) < 2 {
		return &botError{ERR_TYPES_COMMAND, ""}
	}
	
	pokemonName := strings.ToLower(b.fields[1])
	
	if p, ok := pokemon.PokemonMap[pokemonName]; ok {
		message := fmt.Sprintf("Type for **%s**: %s", p.Name, p.Types.Print())
		b.PrintToDiscord(message)
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}
	
	return nil
}

func (b *botResponse) PrintTypeChartToDiscord() error {
	if len(b.fields) < 2 {
		return &botError{ERR_TYPECHART_COMMAND, ""}
	}
	
	typeValue := strings.ToLower(b.fields[1])
	
	if p, ok := pokemon.PokemonMap[typeValue]; ok {
		b.PrintToDiscord(p.PrintTypeChart())
	} else if t, ok := pokemon.TypeMap[pokemon.Type2ID[typeValue]]; ok {
		b.PrintToDiscord(t.PrintTypeChart())
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}
	return nil
}

func (b *botResponse) PrintToDiscord(msg string) {
	_, _ = b.s.ChannelMessageSend(b.m.ChannelID, msg)
	return
}

func (b *botResponse) PrintErrorToDiscord(err error) {
	if berr, ok := err.(*botError); ok {
		_, _ = b.s.ChannelMessageSend(b.m.ChannelID, berr.Error())
	} else {
		_, _ = b.s.ChannelMessageSend(b.m.ChannelID, err.Error())
	}
	return
}

type botError struct {
	err error
	pokemon string
}

func (e *botError) Error() string {
	if e.err == ERR_POKEMON_UNRECOGNIZED && e.pokemon != "" {
		return fmt.Sprintf("Pokemon unrecognized: %s", e.pokemon)
	} else if e.err == ERR_NO_COMBINATIONS && e.pokemon != "" {
		return fmt.Sprintf("No possible IV combinations for that CP for %s", e.pokemon)
	} else if e.err == ERR_NO_STATS && e.pokemon != "" {
		return fmt.Sprintf("No stats available for %s in the Pokemon Go Master file yet :(", e.pokemon)
	}
	return e.err.Error()
}

func (e *botError) pokemonUnrecognized() bool {
	if _, ok := pokemon.PokemonMap[strings.ToLower(e.pokemon)]; !ok {
		return false
	}
	return true
}


