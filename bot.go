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
    "time"
)

// BotID for discord
var (
	BotID string
	goBot *discordgo.Session
)

// Error printouts
var (
	ERR_CP_COMMAND           = errors.New("CP command needs to be formatted like this: !cp {pokemon} {level} {attack iv} {defense iv} {stamina iv}")
	ERR_IV_COMMAND           = errors.New("IV command needs to be formatted like this: !iv {pokemon} {cp} {hp} {level} or !iv {pokemon} {cp} {hp}")
	ERR_RAIDCP_COMMAND       = errors.New("Raid CP command needs to be formatted like this: !raidcp {pokemon} or !raidcp {pokemon} {cp}")
	ERR_RAIDCHART_COMMAND    = errors.New("Raid CP Chart command needs to be formatted like this: !raidcpchart {pokemon}")
	ERR_MAXCP_COMMAND        = errors.New("Max CP command needs to be formatted like this: !maxcp {pokemon}")
	ERR_MOVES_COMMAND        = errors.New("Moves command needs to be formatted like this: !moves {pokemon}")
	ERR_TYPES_COMMAND        = errors.New("Types command needs to be formatted like this: !type {pokemon}")
	ERR_TYPECHART_COMMAND    = errors.New("Effect command needs to be formatted like this: !effect {pokemon}")
	ERR_PREFIX_COMMAND		 = errors.New("Set the prefix for your guild using !setprefix {prefix}. Max 2 characters.")
	ERR_WELCOME_COMMAND		= errors.New("Set the welcome message for your server using !setwelcome {message}")
	ERR_GOODBYE_COMMAND		= errors.New("Set the goodbye message for your server using !setgoodbye {message}")
	ERR_NO_COMBINATIONS      = errors.New("No possible IV combinations for that CP")
	ERR_NO_STATS             = errors.New("Pokemon Master file doesn't have stats for that pokemon yet :(")
	ERR_POKEMON_UNRECOGNIZED = errors.New("Pokemon not recognized.")
	ERR_POKEMON_TYPE_UNRECOGNIZED = errors.New("Pokemon/type not recognized.")
	ERR_COMMAND_UNRECOGNIZED = errors.New("Command not recognized")
	ERR_COORDS_COMMAND = errors.New("Coords command needs to be formatted like this: !coords {lat,long}")

	ERR_NO_CHANNEL = errors.New("Unable to get Channel ID")
	ERR_NO_GUILD = errors.New("Unable to get Guild ID")
	ERR_NO_TEAM = errors.New("No team provided.")
	ERR_NOT_OWNER = errors.New("Only the server owner can use that command :)")
	
	ERR_MISSING_ROLE = errors.New("Missing role.")
	ERR_INVALID_ROLE = errors.New("Invalid role.")
	ERR_ROLE_ADD = errors.New("Unable to add role :(")
	ERR_ROLE_REMOVE = errors.New("Unable to remove role :(")
	
)

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
var cmdList map[string]BotCommand
var botCommands = []BotCommand{
	{"iv", "!iv [pokemon] [cp] [hp] {level|stardust} {adh}",
		"Get possible IVs of a pokemon", 
		[]string{"!iv machamp 2526 143 33 a", "!iv pikachu 613 56 5000 ad", "!iv raichu 1703 98"}, true,
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
		[]string{"raidcp", "eggcp", "eggiv", "mewcp", "mewiv"},
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
    {"luckydate", "!luckydate", 
        "Returns the date for pokemon to have been caught by for a higher change at luckies.",
        []string{"!luckydate"}, true,
        []string{},
        PrintLuckyDateToDiscord,
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
	/*{"add", "!add", "Add this guild to management",
		[]string{}, false, []string{}, AddGuild,
	},*/
	{"setprefix", "!setprefix {prefix string}", "Change bot prefix for server", 
		[]string{}, false, []string{},
		SetBotPrefix,
	},
	{"setwelcome", "!setwelcome {message}", "Set welcome message for server",
		[]string{}, false, []string{},
		SetWelcome,
	},
	{"setgoodbye", "!setgoodbye {message}", "Set goodbye message for server",
		[]string{}, false, []string{},
		SetGoodbye,
	},
}

// Formatting for info
var INFO_FORMAT = "!cmd [required] [fields|options] {optional}"

// PrintInfo prints the info for a discord command
func (cmd *BotCommand) PrintInfo(prefix string) string {
	examples := Example(strings.Replace(cmd.Format, "!", prefix, 1))
	for _, ex := range cmd.Example {
		examples += Example(strings.Replace(ex, "!", prefix, 1))
	}
	return fmt.Sprintln(cmd.Info, examples)
}

//NewBotResponse creates an instance of a bot interaction
func NewBotResponse(s *discordgo.Session, m *discordgo.MessageCreate, fields []string) *botResponse {
	return &botResponse{s: s, m: m, fields: fields}
}

// GetCommand gets the BotCommand for the input
func (b *botResponse) GetCommand(prefix string) (cmd *BotCommand) {
	if len(b.fields) == 0 {
		b.err = ERR_COMMAND_UNRECOGNIZED
		return cmd
	}
	
	name := strings.ToLower(strings.Replace(b.fields[0], prefix, "", 1))
	if c, ok := cmdMap[name]; ok {
		return &c
	} else {
		b.err = ERR_COMMAND_UNRECOGNIZED
		return cmd
	}
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
		
		log.Println("Adding active guilds...")
		err = initGuilds(goBot.State)
		if err != nil {
			log.Println(err.Error())
		}

		err = goBot.UpdateStatus(0, "!wat")
		if err != nil {
			fmt.Println("Unable to update status: ", err.Error())
		}
		
		if UseImages {
			//Start Image Server
			http.Handle("/img", http.FileServer(http.Dir(ImageServer)))
		}

		log.Println("Bot is running!")
}

func welcomeHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guild, ok := Guilds[m.GuildID];
	if !ok {
		return
	}
	
	if !guild.IsManaged() {
		return
	}
	
	
	err := guild.PrintWelcome(m.User)
	if err != nil {
		log.Println(err)
	}
	
	return
}

func goodbyeHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	guild, ok := Guilds[m.GuildID];
	if !ok {
		return
	}
	
	if !guild.IsManaged() {
		return
	}
	
	err := guild.PrintGoodbye(m.User)
	if err != nil {
		log.Println(err)
	}
	
	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}
	
	guild, ok := Guilds[channel.GuildID]
	if !ok {
		return
	}
	
	guild.Update()
	
	prefix := guild.Settings.BotPrefix
	
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	if m.Author.ID == BotID {
		return
	}
	
	log.Println(m.Content)

	if strings.Contains(m.Content, "mewcp") || strings.Contains(m.Content, "mewiv") {
		m.Content = strings.Replace(m.Content, "mewcp", "mewcp mew", 1)
		m.Content = strings.Replace(m.Content, "mewiv", "mewiv mew", 1)
	}
	bot := NewBotResponse(s, m, strings.Fields(m.Content))
	cmd := bot.GetCommand(prefix)
	if bot.err != nil {
		return
	}
	
	err = cmd.Do(bot)
	if err != nil {
			bot.PrintErrorToDiscord(err)
		}

	return

}

// AddGuild adds a guild to the guild management and checks for requirements
func AddGuild(b *botResponse) error {
	channel, err := b.s.Channel(b.m.ChannelID)
	if err != nil {
		return &botError{ERR_NO_CHANNEL, ""}
	}
	
	guild, ok := Guilds[channel.GuildID];
	if !ok {
		g, err := b.s.Guild(channel.GuildID)
		if err != nil {
			return &botError{ERR_NO_GUILD, ""}
		}
		
		guild = NewGuild(g)
	}
	
	guild.Manage(true)
	
	if len(b.fields) > 1 {
		if b.fields[1] == "teams" {
			// Check Roles
			err = guild.CheckRoles()
			if err != nil {
				return &botError{err, ""}
			}
			guild.ManageTeams(true)
		}
	}
	
	guild.Manage(true)
	
	b.PrintToDiscord("Guild management added!")
	
	return nil
}

// SetBotPrefix sets a bot prefix other than "!" for a certain guild
func SetBotPrefix(b *botResponse) error {
	channel, err := b.s.Channel(b.m.ChannelID)
	if err != nil {
		return &botError{ERR_NO_CHANNEL, ""}
	}
	
	guild, ok := Guilds[channel.GuildID];
	if !ok {
		g, err := b.s.Guild(channel.GuildID)
		if err != nil {
			return &botError{ERR_NO_GUILD, ""}
		}
		
		guild = NewGuild(g)
	}
	
	if !guild.IsOwner(b.m.Author) {
		return &botError{ERR_NOT_OWNER, ""}
	}
	
	if len(b.fields) > 1 {
		prefix := b.fields[1]
		if len(prefix) > 2 {
			return &botError{ERR_PREFIX_COMMAND, ""}
		}
		guild.SetPrefix(prefix)
		
		b.PrintToDiscord("Haynesbot prefix successfully changed to " + prefix)
	} else {
		return &botError{ERR_PREFIX_COMMAND, ""}
	}
	return nil
}

// AssignTeam assigns one of three teams (mystic,valor,instinct)
func AssignTeam(b *botResponse) error {
	if len(b.fields) < 2 {
		return &botError{ERR_NO_TEAM, ""}
	}
	
	team := strings.ToLower(b.fields[1])
	// Make sure this is a valid team
	if !IsValidTeam(team) {
		return &botError{ERR_INVALID_ROLE, b.fields[1]}
	}
	
	// Attempt to get the channel from the state
	// If error, fall back to restapi
	channel, err := b.s.State.Channel(b.m.ChannelID)
	if err != nil {
		channel, err = b.s.Channel(b.m.ChannelID)
		if err != nil {
			return &botError{ERR_NO_CHANNEL, ""}
		}
	}
	
	// Attempt to get the guild from the state
	guild, ok := Guilds[channel.GuildID];
	if !ok {
		g, err := b.s.Guild(channel.GuildID)
		if err != nil {
			return &botError{ERR_NO_GUILD, ""}
		}
		
		guild = NewGuild(g)
	}
	
	if !guild.IsManaged() {
		return &botError{ERR_NOT_MANAGED, ""}
	}
	
	if !guild.TeamsManaged() {
		return &botError{ERR_NOT_MANAGED, ""}
	}

	// Remove all team roles
	err = guild.RemoveAllTeams(b.s, b.m.Author.ID)
	if err != nil {
		return &botError{ERR_ROLE_REMOVE, ""}
	}

	err = guild.AddRole(b.s, b.m.Author.ID, team)
	if err != nil {
		return &botError{err, ""}
	}

	b.PrintToDiscord(fmt.Sprintf("You have been added to team %s!", team))
	
	return nil
}

// SetWelcome allows the server owner to set a welcome message
func SetWelcome(b *botResponse) error {
	channel, err := b.s.Channel(b.m.ChannelID)
	if err != nil {
		return &botError{ERR_NO_CHANNEL, ""}
	}
	
	guild, ok := Guilds[channel.GuildID];
	if !ok {
		g, err := b.s.Guild(channel.GuildID)
		if err != nil {
			return &botError{ERR_NO_GUILD, ""}
		}
		
		guild = NewGuild(g)
	}
	
	if !guild.IsManaged() {
		return &botError{ERR_NOT_MANAGED, ""}
	}
	
	if !guild.IsOwner(b.m.Author) {
		return &botError{ERR_NOT_OWNER, ""}
	}
	
	if len(b.fields) > 1 {
		welcome := strings.Join(b.fields[1:], " ")
		guild.SetWelcome(welcome)
		
		b.PrintToDiscord("Welcome message set!")
	} else {
		return &botError{ERR_WELCOME_COMMAND, ""}
	}

	return nil
}

// SetGoodbye allows the server owner to set a goodbye message
func SetGoodbye(b *botResponse) error {
	channel, err := b.s.Channel(b.m.ChannelID)
	if err != nil {
		return &botError{ERR_NO_CHANNEL, ""}
	}
	
	guild, ok := Guilds[channel.GuildID];
	if !ok {
		g, err := b.s.Guild(channel.GuildID)
		if err != nil {
			return &botError{ERR_NO_GUILD, ""}
		}
		
		guild = NewGuild(g)
	}
	
	if !guild.IsManaged() {
		return &botError{ERR_NOT_MANAGED, ""}
	}
	
	if !guild.IsOwner(b.m.Author) {
		return &botError{ERR_NOT_OWNER, ""}
	}
	
	if len(b.fields) > 1 {
		goodbye := strings.Join(b.fields[1:], " ")
		guild.SetGoodbye(goodbye)
		b.PrintToDiscord("Goodbye message set!")
	} else {
		return &botError{ERR_GOODBYE_COMMAND, ""}
	}

	return nil
}

// PrintInfoToDiscord prints the bot info to discord
func PrintInfoToDiscord(b *botResponse) error {
	channel, err := b.s.Channel(b.m.ChannelID)
	if err != nil {
		return &botError{ERR_NO_CHANNEL, ""}
	}
	
	guild, ok := Guilds[channel.GuildID];
	if !ok {
		return &botError{ERR_NO_GUILD, ""}
	}
	
	prefix := guild.Settings.BotPrefix
	
	emb := NewEmbed().
		//SetTitle("Haynes Bot Commands").
		SetColor(0x00ff00).
		AddField("Commands", Example(strings.Replace(INFO_FORMAT, "!", prefix, 1)))
		
	for _, cmd := range cmdList {
		if !cmd.Print {
			continue
		}
		if len(b.fields) == 1 {
			emb.AddField(prefix+cmd.Name, Example(strings.Replace(cmd.Format, "!", prefix, 1)))
			continue
		} else if len(b.fields) > 1 {
			if strings.ToLower(b.fields[1]) != cmd.Name && strings.ToLower(b.fields[1]) != "full" {
				continue
			}
		}
		emb.AddField(prefix+cmd.Name, cmd.PrintInfo(prefix))
	}
	b.PrintEmbedToDiscord(emb.MessageEmbed)
	return nil
}

// PrintIVToDiscord prints the IV data to discord
func PrintIVToDiscord(b *botResponse) error {
	if len(b.fields) < 3 {
		return &botError{ERR_IV_COMMAND, ""}
	}

	pokemonName := strings.ToLower(b.fields[1])

	cp, err := strconv.Atoi(b.fields[2])
	if err != nil {
		return &botError{ERR_IV_COMMAND, ""}
	}

    hp, err := strconv.Atoi(b.fields[3])
    if err != nil {
        return &botError{ERR_IV_COMMAND, ""}
    }

	level := 0.0
	stardust := 0
	if len(b.fields) > 4 {

		val, err := strconv.ParseFloat(b.fields[4], 64)
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
	if len(b.fields) > 5 {
		if strings.Contains(b.fields[5], "a") {
			bestvals += "a"
		}
		if strings.Contains(b.fields[5], "d") {
			bestvals += "d"
		}
		if strings.Contains(b.fields[5], "h") || strings.Contains(b.fields[5], "s") {
			bestvals += "s"
		}
	}

	if p, err := pogo.GetPokemon(pokemonName); err == nil {
		stats, ivChart := p.GetIV(cp, hp, level, stardust, bestvals)
		if len(ivChart) == 0 {
			return &botError{ERR_NO_COMBINATIONS, p.Name}
		} else {
			emb := NewEmbed().
					SetColor(0x9013FE).
					AddField(fmt.Sprintf("CP: %d", cp), Example(ivChart)).
					SetAuthor(p.Name, p.API.Sprites.Front)
					//SetImage(p.API.Sprites.Front).MessageEmbed
			if len(stats) > 30 {
				emb.SetDescription("Full chart too long to display, displaying first 30 rows. \nAdd more data to limit results.")
			}
			b.PrintEmbedToDiscord(emb.MessageEmbed)
			//b.PrintToDiscord(ivChart)
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
			if ImageExists(imgName) {
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

// PrintLuckyDateToDiscord prints the lucky date to discord
func PrintLuckyDateToDiscord(b *botResponse) error {
    now := time.Now()
    luckydate := now.AddDate(0, 0, -780)
    msg := fmt.Sprintf("Any Pokémon older than **%s** has the highest chance to become lucky.", luckydate.Format("01/02/2006"))
    _, _ = b.s.ChannelMessageSend(b.m.ChannelID, msg)

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
	var pokego_url = "https://pokedex100.com/?z="
	
	if len(b.fields) < 2 {
		return &botError{ERR_COORDS_COMMAND, ""}
	}
	
	c := config.BotPrefix + "coords"
	coords := strings.Replace(strings.Replace(strings.Join(b.fields, ""), c, "", 1), " ", "", -1)
	fmt.Println(coords)
	
	link := pokego_url + coords
	
	b.PrintToDiscord(link)
	
	return nil
}

// SendImageToDiscord sends an image as a file attachment to discord
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

var silentErrors = []error{ERR_NO_CHANNEL, ERR_NOT_MANAGED, ERR_NO_GUILD}

// Error formats the error message for printing
func (e *botError) Error() string {
	for _, se := range silentErrors {
		if e.err == se {
			log.Println(e.err)
			return ""
		}
	}
	
	if e.err == ERR_POKEMON_UNRECOGNIZED && e.value != "" {
		return fmt.Sprintf("Pokemon unrecognized: %s", e.value)
	} else if e.err == ERR_POKEMON_TYPE_UNRECOGNIZED && e.value != "" {
		return fmt.Sprintf("Pokemon/type unrecognized: %s", e.value)
	} else if e.err == ERR_NO_COMBINATIONS && e.value != "" {
		return fmt.Sprintf("No possible IV combinations for that CP for %s", e.value)
	} else if e.err == ERR_NO_STATS && e.value != "" {
		return fmt.Sprintf("No stats available for %s in the Pokemon Go Master file yet :(", e.value)
	} else if e.err == ERR_MISSING_ROLE && e.value != "" {
		return fmt.Sprintf("Missing role: %s", e.value)
	} else if e.err == ERR_INVALID_ROLE && e.value != "" {
		return fmt.Sprintf("Invalid role: %s", e.value)
	}
	return e.err.Error()
}

// ImageExists checks if an image exists in the image server folder
func ImageExists(name string) bool {
	if _, err := os.Stat(ImageServer + "/" + name); err != nil {
	  return false
	}
	return true
}

func init() {
	cmdList = make(map[string]BotCommand)
	cmdMap = make(map[string]BotCommand)
	for _, cmd := range botCommands {
		cmdList[cmd.Name] = cmd
		cmdMap[cmd.Name] = cmd
		if len(cmd.Aliases) > 0 {
			for _, alias := range cmd.Aliases {
				cmdMap[alias] = cmd
			}
		}
	}
}
