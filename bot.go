package haynesbot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/haynesherway/pogo"
	"io"
	"log"
	"os"
	"net/http"
	"strconv"
	"strings"
)

var (
	BotID string
	goBot *discordgo.Session
)

// Error printouts
var (
	ERR_CP_COMMAND           = errors.New("CP command needs to be formatted like this: !maxcp {pokemon} {level} {attack iv} {defense iv} {stamina iv}")
	ERR_IV_COMMAND           = errors.New("IV command needs to be formatted like this: !iv {pokemon} {cp} {level} or !iv {pokemon} {cp}")
	ERR_RAIDCP_COMMAND       = errors.New("Raid CP command needs to be formatted like this: !raidcp {pokemon} or !raidcp {pokemon} {cp}")
	ERR_RAIDCHART_COMMAND    = errors.New("Raid CP Chart command needs to be formatted like this: !raidcpchart {pokemon}")
	ERR_MAXCP_COMMAND        = errors.New("Max CP command needs to be formatted like this: !maxcp {pokemon}")
	ERR_MOVES_COMMAND        = errors.New("Moves command needs to be formatted like this: !moves {pokemon}")
	ERR_TYPES_COMMAND        = errors.New("Types command needs to be formatted like this: !type {pokemon}")
	ERR_TYPECHART_COMMAND    = errors.New("Effect command needs to be formatted like this: !effect {pokemon}")
	ERR_NO_COMBINATIONS      = errors.New("No possible IV combinations for that CP")
	ERR_NO_STATS             = errors.New("Pokemon Master file doesn't have stats for that pokemon yet :(")
	ERR_POKEMON_UNRECOGNIZED = errors.New("Pokemon not recognized.")
	ERR_POKEMON_TYPE_UNRECOGNIZED = errors.New("Pokemon/type not recognized.")
	ERR_COMMAND_UNRECOGNIZED = errors.New("Command not recognized")
	ERR_COORDS_COMMAND = errors.New("Coords command needs to be formatted like this: !coords {lat,long}")

	ERR_NO_CHANNEL = errors.New("Unable to get Channel ID")
	ERR_NO_GUILD = errors.New("Unable to get Guild ID")
	ERR_NO_TEAM = errors.New("No team provided.")
	
	ERR_MISSING_ROLE = errors.New("Missing role.")
	ERR_INVALID_ROLE = errors.New("Invalid role.")
	ERR_ROLE_ADD = errors.New("Unable to add role :(")
	ERR_ROLE_REMOVE = errors.New("Unable to remove role :(")
	
)

//var POKEGO_URL = "https://pokemon.pokego2.com/coords-"
var POKEGO_URL = "https://pokedex100.com/?z="

type botResponse struct {
	s       *discordgo.Session
	m       *discordgo.MessageCreate
	command string
	fields  []string
	err     error
}

// Type Do is a placeholder for the function a command should execute
type Do func(b *botResponse) error 

// BotCommand is a representation of a command the bot can handle
type BotCommand struct {
	Name string
	Format string
	Info string
	Example []string
	Print bool
	Aliases []string
	Do
}

var cmdMap map[string]BotCommand
var botCommands = []BotCommand{
	{"iv", "!iv [pokemon] [cp] {level|stardust} {adh}",
		"Get possible IVs of a pokemon", 
		[]string{"!iv numel 506 33 d", "!iv pikachu 613 500 ad", "!iv raichu 1703"}, true,
		[]string{},
		PrintIVToDiscord,
	},
	{"cp", "!cp [pokemon] [level] [attack iv] [defense iv] [stamina iv]",
		"Get CP of a pokemon at a specified level with specified IVs",
		[]string{"!cp mewtwo 25 15 14 15"}, true,
		[]string{},
		PrintCPToDiscord,
	},
	{"maxcp", "!maxcp [pokemon]",
		"Get maximum CP of a pokemon with perfect IVs at level 40", 
		[]string{"!maxcp latios"}, true,
		[]string{},
		PrintMaxCPToDiscord,
	},
	{"raidiv", "!raidiv [pokemon] {cp}",
		"Get possible IV combinations for specified raid pokemon with specified IV",
		[]string{"!raidcp kyogre 2292", "!raidcp groudon"}, true,
		[]string{"raidcp", "eggcp", "eggiv"},
		PrintRaidCPToDiscord,
	},
	{"raidchart", "!raidchart [pokemon] {'full'}",
		"Get a chart with possible stats for specified pokemon at raid level above 90%",
		[]string{"!raidchart machamp", "!raidchart rayquaza full"}, true,
		[]string{},
		PrintRaidChartToDiscord,
	},
	{"moves", "!moves [pokemon]",
		"Get a list of fast and charge moves for specified pokemon",
		[]string{"!moves rayquaza"}, true,
		[]string{},
		PrintMovesToDiscord,
	},
	{"type", "!type [pokemon]",
		"Get a list of types for a specified pokemon",
		[]string{"!type rayquaza"}, true,
		[]string{},
		PrintTypeToDiscord,
	},
	{"effect", "!effect [pokemon|type]",
		"Get a list of type relations a specified pokemon or type has",
		[]string{"!effect pikachu", "!effect electric"}, true,
		[]string{},
		PrintTypeChartToDiscord,
	},
	{"wat", "!wat {command|'full'}",
		"Get info about commands",
		[]string{"!wat", "!wat full", "!wat raidcp"}, true,
		[]string{"haynes-bot", "haynez-bot"},
		PrintInfoToDiscord,
	},
	{"coords", "!coords {lat, long}",
		"Get a link to google maps for coordinates",
		[]string{"!coords 43.24124,-76.14241"}, false,
		[]string{},
		PrintCoordsToDiscord,
	},
	{"team", "!team {mystic|valor|instinct}",
		"Get assigned to a team",
		[]string{"!team mystic", "!team valor", "!team instinct"}, false,
		[]string{},
		AssignTeam,
	},
	{"add", "!add", "Add this guild to state",
		[]string{}, false, []string{}, AddGuild,
	},
}

var INFO_FORMAT = "!cmd [required] [fields|options] {optional}"

func (cmd *BotCommand) PrintInfo() string {
	examples := Example(cmd.Format)
	for _, ex := range cmd.Example {
		examples += Example(ex)
	}
	return fmt.Sprintln(cmd.Info, examples)
}

//NewBotResponse creates an instance of a bot interaction
func NewBotResponse(s *discordgo.Session, m *discordgo.MessageCreate, fields []string) *botResponse {
	return &botResponse{s: s, m: m, fields: fields}
}

// GetCommand gets the BotCommand for the input
func (b *botResponse) GetCommand() (cmd *BotCommand) {
	if len(b.fields) == 0 {
		b.err = ERR_COMMAND_UNRECOGNIZED
		return cmd
	}
	
	name := strings.ToLower(strings.Replace(b.fields[0], config.BotPrefix, "", 1))
	if c, ok := cmdMap[name]; ok {
		return &c
	} else {
		b.err = ERR_COMMAND_UNRECOGNIZED
		return cmd
	}
}

func (b *botResponse) GetCommandName() string {
	if len(b.fields) == 0 {
		b.err = ERR_COMMAND_UNRECOGNIZED
		return ""
	}
	cmd := strings.ToLower(strings.Replace(b.fields[0], config.BotPrefix, "", 1))
	return cmd
}

func NilFunc(b *botResponse) error {
	return nil
}

// Start starts the bot
func Start() {
	var err error
	goBot, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

		u, err := goBot.User("@me")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		BotID = u.ID

		goBot.AddHandler(messageHandler)
		goBot.AddHandler(welcomeHandler)
		goBot.AddHandler(goodbyeHandler)
		err = goBot.Open()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = goBot.UpdateStatus(0, "!haynes-bot")
		if err != nil {
			fmt.Println("Unable to update status: ", err.Error())
		}
		
		if UseImages {
			//Start Image Server
			http.Handle("/img", http.FileServer(http.Dir(ImageServer)))
		}

		fmt.Println("Bot is running!")
}

func welcomeHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if !IsManaged(m.GuildID) {
		return
	}
	
	g, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}
	
	guild := Guild{g}
	
	err = guild.PrintWelcome(m.User)
	if err != nil {
		log.Println(err)
	}
	
	return
}

func goodbyeHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if !IsManaged(m.GuildID) {
		return
	}
	
	g, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}
	
	guild := Guild{g}
	
	err = guild.PrintGoodbye(m.User)
	if err != nil {
		log.Println(err)
	}
	
	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, config.BotPrefix) {
		return
	}

	if m.Author.ID == BotID {
		return
	}
	
	fmt.Println(m.Content)

	bot := NewBotResponse(s, m, strings.Fields(m.Content))
	cmd := bot.GetCommand()
	if bot.err != nil {
		return
	}
	
	err := cmd.Do(bot)
	if err != nil {
			bot.PrintErrorToDiscord(err)
		}

	return

}

func AddGuild(b *botResponse) error {
	channel, err := b.s.Channel(b.m.ChannelID)
	if err != nil {
		return &botError{ERR_NO_CHANNEL, ""}
	}
	
	if !IsManaged(channel.GuildID) {
		return &botError{ERR_NOT_MANAGED, ""}
	}
	
	g, err := b.s.Guild(channel.GuildID)
	if err != nil {
		return &botError{ERR_NO_GUILD, ""}
	}
	
	guild := Guild{g}
	
	// Check Roles
	err = guild.CheckRoles()
	if err != nil {
		return &botError{err, ""}
	}
	b.s.State.GuildAdd(guild.Guild)
	
	return nil
}

//AssignTeam assigns one of three teams (mystic,valor,instinct)
func AssignTeam(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_NO_TEAM, ""}
	}
	
	team := strings.ToLower(b.fields[1])
	// Make sure this is a valid team
	if !IsValidTeam(team) {
		return &botError{ERR_INVALID_ROLE, b.fields[1]}
	}
	
	// Attempt to get the channe from the state
	// If error, fall back to restapi
	channel, err := b.s.State.Channel(b.m.ChannelID)
	if err != nil {
		channel, err = b.s.Channel(b.m.ChannelID)
		if err != nil {
			return &botError{ERR_NO_CHANNEL, ""}
		}
	}
	
	if !IsManaged(channel.GuildID) {
		return &botError{ERR_NOT_MANAGED, ""}
	}
	
	// Attempt to get the guild from the state
	// If error, fall back to restapi
	g, err := b.s.State.Guild(channel.GuildID)
	if err != nil {
		g, err = b.s.Guild(channel.GuildID)
		if err != nil {
			return &botError{ERR_NO_GUILD, ""}
		}
	}
	
	guild := Guild{g}
	
	// Remove all team roles
	err = guild.RemoveAllTeams(b.s, b.m.Author.ID)
	if err != nil {
		return ERR_ROLE_REMOVE
	}
	
	err = guild.AddRole(b.s, b.m.Author.ID, team)
	if err != nil {
		return err
	}
	
	b.PrintToDiscord(fmt.Sprintf("You have been added to team %s!", team))
	
	return nil
}

//PrintInfoToDiscord prints the bot info to discord
func PrintInfoToDiscord(b *botResponse) error {
	emb := NewEmbed().
		//SetTitle("Haynes Bot Commands").
		SetColor(0x00ff00).
		AddField("Commands", Example(INFO_FORMAT))
		
	for _, cmd := range cmdMap {
		if !cmd.Print {
			continue
		}
		if len(b.fields) == 1 {
			emb.AddField("!"+cmd.Name, Example(cmd.Format))
			continue
		} else if len(b.fields) > 1 {
			if strings.ToLower(b.fields[1]) != cmd.Name && strings.ToLower(b.fields[1]) != "full" {
				continue
			}
		}
		emb.AddField("!"+cmd.Name, cmd.PrintInfo())
	}
	b.PrintEmbedToDiscord(emb.MessageEmbed)
	return nil
}

//PrintIVToDiscord prints the IV data to discord
func PrintIVToDiscord(b *botResponse) error {
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
		if strings.Contains(b.fields[4], "h") || strings.Contains(b.fields[4], "s") {
			bestvals += "s"
		}
	}

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
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

// PrintCPToDiscord prints CP info based on input to discord
func PrintCPToDiscord(b *botResponse) error {
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

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		cp := p.GetCP(level, ivA, ivD, ivS)
		emb := NewEmbed().
		SetColor(0x9013FE).
		AddField(p.Name, fmt.Sprintf("CP at level %v with IVs %d/%d/%d: %d", level, ivA, ivD, ivS, cp)).
		SetThumbnail(p.API.Sprites.Front).MessageEmbed
		b.PrintEmbedToDiscord(emb)
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

// PrintMaxCPToDiscord prints an embed with the max cp to discord
func PrintMaxCPToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_MAXCP_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		maxcp := p.GetMaxCP()
		if maxcp == 0 {
			return &botError{ERR_NO_STATS, p.Name}
		}
		emb := NewEmbed().
		SetColor(0x9013FE).
		AddField(p.Name, fmt.Sprintf("Max CP: %v", maxcp)).
		SetThumbnail(p.API.Sprites.Front).MessageEmbed
		b.PrintEmbedToDiscord(emb)
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

// PrintRaidChartToDiscord prints a chart with CP/IVs to discord
func PrintRaidChartToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_RAIDCHART_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		ivList, chart := p.GetRaidCPChart()
		if UseImages {
			imgName := fmt.Sprintf("RAIDCHART-%s.png", p.Name)
			if ImageExists("poop") {
					f, err := os.Open(ImageServer + "/" + imgName)
					if err != nil {
						fmt.Println(err.Error())
					}
					b.SendImageToDiscord(imgName, f)
				} else {
					fmt.Println("Getting table")
					GetTable(p, ivList, imgName)
					
					f, err := os.Open(ImageServer + "/" + imgName)
					if err != nil {
						fmt.Println(err.Error())
					}
					
					b.SendImageToDiscord(imgName, f)
				} 
			
		} else {
			rows := strings.Split(chart, "\n")
			emb := NewEmbed().
				SetColor(0x9013FE).
				AddField("Raid Chart", Example(strings.Join(rows[:40], "\n")))
				
				if len(b.fields) > 2 {
					if strings.ToLower(b.fields[2]) == "full" {
						rowCount := len(rows)
						st, en := 41, 80
						for {
							if en > rowCount {
								en = rowCount
							}
							emb.AddField("Continued", Example(strings.Join(rows[st:en], "\n")))
							
							if en == rowCount || en > 200 {
								break
							}
							st+=40
							en+=40
						}
					}
				}
				emb.SetAuthor(p.Name, p.API.Sprites.Front)
				b.PrintEmbedToDiscord(emb.MessageEmbed)
		}
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

// PrintRaidCPToDiscord prints either a range or a list of possible CPs for a raid pokemon
func PrintRaidCPToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_RAIDCP_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		if len(b.fields) == 2 {
			emb := NewEmbed().
			SetColor(0x9013FE).
			AddField(p.Name + " Raid CP", p.GetRaidCPRange()).
			SetThumbnail(p.API.Sprites.Front).MessageEmbed
			b.PrintEmbedToDiscord(emb)
		} else {
			cp, err := strconv.Atoi(b.fields[2])
			if err != nil {
				return &botError{ERR_RAIDCP_COMMAND, ""}
			}
			//imgName := fmt.Sprintf("RAID-%s-%d.png", p.Name, cp)
			_, ivChart := p.GetRaidIV(cp)
			/*if UseImages {
				//imgName := fmt.Sprintf("RAID-%s-%d", p.Name, cp)
				//imgName := "draw.png"
				if ImageExists(imgName) {
					f, err := os.Open(ImageServer + "/" + imgName)
					if err != nil {
						fmt.Println(err.Error())
					}
					b.SendImageToDiscord(imgName, f)
				} else {
					file := GetTable(ivList, imgName)
					
					b.SendImageToDiscord(imgName, file)
				} 
			} else {*/
				if len(ivChart) == 0 {
					return &botError{ERR_NO_COMBINATIONS, p.Name}
				} else {
					emb := NewEmbed().
					SetColor(0x9013FE).
					AddField(fmt.Sprintf("CP: %d", cp), Example(ivChart)).
					SetAuthor(p.Name, p.API.Sprites.Front).MessageEmbed
					//SetImage(p.API.Sprites.Front).MessageEmbed
					b.PrintEmbedToDiscord(emb)
				}
			//}
		}
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

// PrintMovesToDiscord prints an embed with moves to discord
func PrintMovesToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_MOVES_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		emb := NewEmbed().
		SetTitle(fmt.Sprintf("Moves for %s", p.Name)).
		SetColor(0x0B9EFF).
		AddField("Fast", p.Moves.Fast.Print()).
		AddField("Charge", p.Moves.Charge.Print()).
		SetThumbnail(p.API.Sprites.Front).MessageEmbed
		b.PrintEmbedToDiscord(emb)
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

// PrintTypeToDiscord prints an embed with type info to discord
func PrintTypeToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_TYPES_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		emb := NewEmbed().
		SetColor(0x9013FE).
		AddField(fmt.Sprintf("Type for %s", p.Name), p.Types.Print()).
		SetThumbnail(p.API.Sprites.Front).MessageEmbed
		b.PrintEmbedToDiscord(emb)
	} else {
		return &botError{ERR_POKEMON_UNRECOGNIZED, b.fields[1]}
	}

	return nil
}

// PrintTypeToDiscord prints an embed with a type chart to discord
func PrintTypeChartToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_TYPECHART_COMMAND, ""}
	}

	typeValue := strings.ToLower(b.fields[1])

	if p, err := pogo.GetPokemon(typeValue); err == nil {
		emb := NewEmbed().
		SetColor(0x9013FE).
		SetTitle(fmt.Sprintf("Type Effects for %s", p.Name)).
		AddField("Super Effective", p.SuperEffective.Print()).
		AddField("Not Effective", p.NotEffective.Print()).
		AddField("Weaknesses", p.Weakness.Print()).
		AddField("Resistance", p.Resistance.Print()).
		SetThumbnail(p.API.Sprites.Front).MessageEmbed
		b.PrintEmbedToDiscord(emb)
	} else if t, err := pogo.GetType(typeValue); err == nil {
		emb := NewEmbed().
		SetColor(0x9013FE).
		SetTitle(fmt.Sprintf("Type Effects for %s", t.Name)).
		AddField("Super Effective", t.SuperEffective.Print()).
		AddField("Not Effective", t.NotEffective.Print()).
		AddField("Weaknesses", t.Weakness.Print()).
		AddField("Resistance", t.Resistance.Print()).
		SetThumbnail(t.Thumbnail).MessageEmbed
		b.PrintEmbedToDiscord(emb)
	} else {
		return &botError{ERR_POKEMON_TYPE_UNRECOGNIZED, b.fields[1]}
	}
	return nil
}

// PrintCoordsToDiscord prints a Pokego++ link to the coords
func PrintCoordsToDiscord(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_COORDS_COMMAND, ""}
	}
	
	c := config.BotPrefix + "coords"
	coords := strings.Replace(strings.Replace(strings.Join(b.fields, ""), c, "", 1), " ", "", -1)
	fmt.Println(coords)
	
	link := POKEGO_URL + coords
	
	b.PrintToDiscord(link)
	
	return nil
}

func (b *botResponse) SendImageToDiscord(fileName string, r io.Reader) {
	_, _ = b.s.ChannelFileSend(b.m.ChannelID, fileName, r)
	return
}
// PrintToDiscord prints the message string to discord
func (b *botResponse) PrintToDiscord(msg string) {
	_, _ = b.s.ChannelMessageSend(b.m.ChannelID, msg)
	return
}

// Print embed to discord prints an embed to discord
func (b *botResponse) PrintEmbedToDiscord(e *discordgo.MessageEmbed) {
	_, _ = b.s.ChannelMessageSendEmbed(b.m.ChannelID, e)
}

// PrintErrorToDiscord prints the error to discord
func (b *botResponse) PrintErrorToDiscord(err error) {
	if berr, ok := err.(*botError); ok {
		if berr.Error() == "" {
			return
		}
		_, _ = b.s.ChannelMessageSend(b.m.ChannelID, berr.Error())
	} else {
		_, _ = b.s.ChannelMessageSend(b.m.ChannelID, err.Error())
	}
	return
}

type botError struct {
	err     error
	value string
}

// Error formats the error message for printing
func (e *botError) Error() string {
	if e.err == ERR_POKEMON_UNRECOGNIZED && e.value != "" {
		return fmt.Sprintf("Pokemon unrecognized: %s", e.value)
	} else if e.err == ERR_POKEMON_TYPE_UNRECOGNIZED && e.value != "" {
		return fmt.Sprintf("Pokemon/type unrecognized: %s", e.value)
	} else if e.err == ERR_NO_COMBINATIONS && e.value != "" {
		return fmt.Sprintf("No possible IV combinations for that CP for %s", e.value)
	} else if e.err == ERR_NO_STATS && e.value != "" {
		return fmt.Sprintf("No stats available for %s in the Pokemon Go Master file yet :(", e.value)
	} else if e.err == ERR_NOT_MANAGED {
		return ""
	} else if e.err == ERR_MISSING_ROLE && e.value != "" {
		return fmt.Sprintf("Missing role: %s", e.value)
	} else if e.err == ERR_INVALID_ROLE && e.value != "" {
		return fmt.Sprintf("Invalid role: %s", e.value)
	}
	return e.err.Error()
}

func ImageExists(name string) bool {
	if _, err := os.Stat(ImageServer + "/" + name); err != nil {
	  return false
	}
	return true
}

func init() {
	cmdMap = make(map[string]BotCommand)
	for _, cmd := range botCommands {
		cmdMap[cmd.Name] = cmd
		if len(cmd.Aliases) > 0 {
			for _, alias := range cmd.Aliases {
				cmdMap[alias] = cmd
			}
		}
	}
}