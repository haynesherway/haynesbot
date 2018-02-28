package haynesbot

import (
    "encoding/json"
     "errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
   
    "github.com/bwmarrin/discordgo"
)


// Guild Errors
var (
    ERR_NOT_MANAGED = errors.New("Guild is not managed")
    ERR_NO_WELCOME = errors.New("No welcome message set. Set with !setwelcome")
    ERR_NO_GOODBYE = errors.New("No goodbye message set. Set with !setgoodbye")
    
    ERR_MISSING_CHANNEL = errors.New("Channel missing")
    )

var teamRoles = []string{"mystic", "valor", "instinct", "kith"}

// Settings management
var (
    guildSettings *GuildSettings
    Guilds map[string]*Guild
	ManagedGuilds []string
)

// GuildSettings is a json struct to get the slice of Guild Settings
type GuildSettings struct {
	GuildSettings []GuildSetting `json:"GuildSettings"`
}

// GuildSetting is a holder for the settings of a single guild
type GuildSetting struct {
	Name string `json:"Name"`
	ID string `json:"ID"`
	Managed bool `json:"Managed"`
	Teams bool `json:"Teams"`
	BotPrefix string `json:"Prefix,omitempty"`
	Welcome	string `json:"Welcome,omitempty"`
	Goodbye string `json:"Goodbye,omitempty"`
}

// Guild is a representation of a single discord guild
type Guild struct {
    *discordgo.Guild
    Settings GuildSetting
}

func initGuilds(state *discordgo.State) error {
    
    for _, g := range state.Guilds {
        guild, ok := Guilds[g.ID]
        if ok {
            Guilds[g.ID].Guild = g
            if guild.Settings.BotPrefix == "" {
                guild.Settings.BotPrefix = config.BotPrefix
            }
        } else {
            // No previous settings exist, give default
            guild = NewGuild(g)
        }
        
        guildSettings.add(guild).save(config.GuildFile)
    }

    return nil
}

// NewGuild creates a new guild with the default settings
func NewGuild(guild *discordgo.Guild) *Guild {
    botGuild := &Guild{guild, GuildSetting{
                ID: guild.ID,
                Managed: false,
                Teams: false,
                BotPrefix: config.BotPrefix,
            }}
            
            if guild.Name != "" {
                botGuild.Settings.Name = guild.Name
            }
            Guilds[guild.ID] = botGuild
            
            return botGuild
}

// Update update the settings of a guild and updates the json file
func (guild *Guild) Update() error {
    if guild.Settings.Name == "" { 
        return guildSettings.add(guild).save(config.GuildFile)
    }
    return nil
}

// SetPrefix sets the bot prefix for a guild
func (guild *Guild) SetPrefix(pre string) error {
    guild.Settings.BotPrefix = pre
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

// SetWelcome sets the welcome messages for a guild
func(guild *Guild) SetWelcome(msg string) error {
    guild.Settings.Welcome = msg
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

// SetGoodbye sets the goodbye message for a guild
func(guild *Guild) SetGoodbye(msg string) error {
    guild.Settings.Goodbye = msg
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

// Manage adds the guild into the guilds managed by the bot
func(guild *Guild) Manage(manage bool) error {
    guild.Settings.Managed = manage
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

// ManageTeams allows the guild to have teams (valor, instinct, mystic) managed by the bot
func(guild *Guild) ManageTeams(manage bool) error {
    guild.Settings.Teams = manage
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

// CheckRoles verfies the necessary team roles exist in the guild
func (guild *Guild) CheckRoles() error {
    roleCheck := make(map[string]bool)
    for _, role := range teamRoles {
        roleCheck[role] = false
    }
    
    for _, role := range guild.Guild.Roles {
        if _, ok := roleCheck[role.Name]; ok {
            roleCheck[role.Name] = true
        }
    }
    
    for role, exists := range roleCheck {
        if !exists {
            return &botError{ERR_MISSING_ROLE, role}
        }
    }
    
    return nil
}

// GetChannelID gets the channel id for a channel name in a guild
func (guild Guild) GetChannelID(c string) (string, error) {
    for _, channel := range guild.Channels {
        if channel.Name == c {
            return channel.ID, nil
        }
    }
    
    return "", &botError{ERR_MISSING_CHANNEL, c}
}

// GetRoleID gets the role id for a named role in a guild
func (guild *Guild) GetRoleID(r string) (string, error) {
    for _, role := range guild.Guild.Roles {
        if role.Name == r {
            return role.ID, nil
        }
    }
    
    return "", &botError{ERR_MISSING_ROLE, r}
}

// AddRold adds a role to the given user for a guild
func (guild *Guild) AddRole(session *discordgo.Session, userID string, roleName string) error {
    roleID, err := guild.GetRoleID(roleName)
    if err != nil {
        return err
    }
    
    err = session.GuildMemberRoleAdd(guild.ID, userID, roleID)
	if err != nil {
		return ERR_ROLE_ADD
	}
	
	return nil
}

// RemoveRole removes a role from the given user for a guild
func (guild *Guild) RemoveRole(session *discordgo.Session, userID string, roleName string) error {
    roleID, err := guild.GetRoleID(roleName)
    if err != nil {
        return err
    }
    
    err = session.GuildMemberRoleRemove(guild.ID, userID, roleID)
	if err != nil {
		return ERR_ROLE_REMOVE
	}
	
	return nil
}

// RemoveAllTeams removes all team roles from the given user for a guild
func (guild *Guild) RemoveAllTeams(session *discordgo.Session, userID string) error {
    for _, t := range teamRoles {
        err := guild.RemoveRole(session, userID, t)
        if err != nil {
            return err
        }
    }

    return nil
}

// PrintWelcome prints the stored welcome message in a welcome channel for a guild
func (guild *Guild) PrintWelcome(user *discordgo.User) error {
    if !guild.IsManaged() {
		return ERR_NOT_MANAGED
	}
    
    if guild.Settings.Welcome == "" {
        return ERR_NO_WELCOME
    }
    
    welcomeChannel, err := guild.GetChannelID("welcome")
    if err != nil {
        return &botError{ERR_MISSING_CHANNEL, "welcome"}
    }
    
    var messageStrReplace = map[string]string{
	    "{mention}": user.Mention(),
	    "{guild}": guild.Name,
	    "{user}": user.Username,
	}

	message := guild.Settings.Welcome
	for str, rep := range messageStrReplace {
	    message = strings.Replace(message, str, rep, -1)
	}
    
    _, _ = goBot.ChannelMessageSend(welcomeChannel, message)
    
    return nil
}

// PrintGoodbye prints the stored goodbye message in a welcome channel for a guild
func (guild *Guild) PrintGoodbye(user *discordgo.User) error {
    if !guild.IsManaged() {
		return ERR_NOT_MANAGED
	}
    
    if guild.Settings.Goodbye == "" {
        return ERR_NO_GOODBYE
    }
    
    welcomeChannel, err := guild.GetChannelID("welcome")
    if err != nil {
        return &botError{ERR_MISSING_CHANNEL, "welcome"}
    }
    
    var messageStrReplace = map[string]string{
	    "{mention}": user.Mention(),
	    "{guild}": guild.Name,
	    "{user}": user.Username,
	}

	message := guild.Settings.Goodbye
	for str, rep := range messageStrReplace {
	    message = strings.Replace(message, str, rep, -1)
	}
    
    
    _, _ = goBot.ChannelMessageSend(welcomeChannel, message)
    
    return nil
}

// IsOwner returns true if the given user is the owner of the guild
func (guild *Guild) IsOwner(user *discordgo.User) bool {
    if user.ID == guild.OwnerID {
        return true
    } 
    return false
}

// IsManaged returns true if the guild is managed by the bot
func (guild *Guild) IsManaged() bool {
   return guild.Settings.Managed
}

// TeamsManaged returns true if the teams are managed in this guild by the bot
func (guild *Guild) TeamsManaged() bool {
   return guild.Settings.Teams
}

// IsValidTeam returns true if a team is a valid Pokemon Go team
func IsValidTeam(s string) bool {
    for _, t := range teamRoles {
        if t == s {
            return true
        }
    }
    return false
}

func readGuildSettings(f string) error {
    Guilds = make(map[string]*Guild)
    guildSettings = &GuildSettings{}
    
    file, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	
	err = json.Unmarshal(file, &guildSettings)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	
	for _, gs := range guildSettings.GuildSettings {
	    Guilds[gs.ID] = &Guild{Settings: gs}
	}
	
	return nil
}

func (gs *GuildSettings) add(g *Guild) *GuildSettings {
    if g.Name != "" {
        g.Settings.Name = g.Name
    }
    for i, s := range gs.GuildSettings {
        if s.ID == g.ID {
            gs.GuildSettings[i] = g.Settings
            return gs
        }
    }
    gs.GuildSettings = append(gs.GuildSettings, g.Settings)
    return gs
}

func (gs *GuildSettings) save(file string) error {
    out, err := json.MarshalIndent(gs, "", "  ")
    if err != nil {
       return err 
    }
    
    log.Println("Writing to guild settings file...")
    
    return ioutil.WriteFile(file, out, 0600)
}