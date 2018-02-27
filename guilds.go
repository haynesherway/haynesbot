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

type GuildSettings struct {
	GuildSettings []GuildSetting `json:"GuildSettings"`
}

type GuildSetting struct {
	Name string `json:"Name"`
	ID string `json:"ID"`
	Managed bool `json:"Managed"`
	Teams bool `json:"Teams"`
	BotPrefix string `json:"Prefix,omitempty"`
	Welcome	string `json:"Welcome,omitempty"`
	Goodbye string `json:"Goodbye,omitempty"`
}

type Guild struct {
    *discordgo.Guild
    Settings GuildSetting
}

func InitGuilds(state *discordgo.State) error {
    
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
        if guild.Settings.Managed {
            ManagedGuilds = append(ManagedGuilds, guild.ID)
        }
    }

    return nil
}

func NewGuild(guild *discordgo.Guild) *Guild {
    botGuild := &Guild{guild, GuildSetting{
                Name: guild.Name,
                ID: guild.ID,
                Managed: false,
                Teams: false,
                BotPrefix: config.BotPrefix,
            }}
            Guilds[guild.ID] = botGuild
            
            return botGuild
}

func (guild *Guild) SetPrefix(pre string) error {
    guild.Settings.BotPrefix = pre
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

func(guild *Guild) SetWelcome(msg string) error {
    guild.Settings.Welcome = msg
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

func(guild *Guild) SetGoodbye(msg string) error {
    guild.Settings.Goodbye = msg
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

func(guild *Guild) Manage(manage bool) error {
    guild.Settings.Managed = manage
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

func(guild *Guild) ManageTeams(manage bool) error {
    guild.Settings.Teams = manage
    guildSettings.add(guild).save(config.GuildFile)
    return nil
}

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

func (guild Guild) GetChannelID(c string) (string, error) {
    for _, channel := range guild.Channels {
        if channel.Name == c {
            return channel.ID, nil
        }
    }
    
    return "", &botError{ERR_MISSING_CHANNEL, c}
}

func (guild *Guild) GetRoleID(r string) (string, error) {
    for _, role := range guild.Guild.Roles {
        if role.Name == r {
            return role.ID, nil
        }
    }
    
    return "", &botError{ERR_MISSING_ROLE, r}
}

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

func (guild *Guild) RemoveAllTeams(session *discordgo.Session, userID string) error {
    for _, t := range teamRoles {
        err := guild.RemoveRole(session, userID, t)
        if err != nil {
            return err
        }
    }

    return nil
}

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

func (guild *Guild) IsOwner(user *discordgo.User) bool {
    if user.ID == guild.OwnerID {
        return true
    } 
    return false
}

func (guild *Guild) IsManaged() bool {
   return guild.Settings.Managed
}

func (guild *Guild) TeamsManaged() bool {
   return guild.Settings.Teams
}

func IsValidTeam(s string) bool {
    for _, t := range teamRoles {
        if t == s {
            return true
        }
    }
    return false
}

func ReadGuildSettings(f string) error {
    Guilds = make(map[string]*Guild)
    
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
    g.Settings.Name = g.Name
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
    
    return ioutil.WriteFile(file, out, 0600)
}