package haynesbot

import (
    "errors"
    "fmt"
    "github.com/bwmarrin/discordgo"
)

// Messages
const WELCOME_MESSAGE = `%s, welcome to **%s**! The Syracuse spoofing discord... we accept you... you are safe here...
Set your team by typing **!team valor**, **!team mystic** or **!team instinct**`
const GOODBYE_MESSAGE = `**%s** just left **%s**. Bye bye **%s**. What a big mistake...`

// Guild Errors
var (
    ERR_NOT_MANAGED = errors.New("Guild is not managed")
    
    ERR_MISSING_CHANNEL = errors.New("Channel missing")
    )

var teamRoles = []string{"mystic", "valor", "instinct", "kith"}

type Guild struct {
    *discordgo.Guild
}

func IsManaged(guildid string) bool {
	for _, id := range ManagedGuilds {
		if id == guildid {
			return true
		} 
	}
	return false
}

func (guild *Guild) CheckRoles() error {
    if !IsManaged(guild.ID) {
        return ERR_NOT_MANAGED
    }
    
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
    if !IsManaged(guild.ID) {
        return ERR_NOT_MANAGED
    }
    
    welcomeChannel, err := guild.GetChannelID("welcome")
    if err != nil {
        return &botError{ERR_MISSING_CHANNEL, "welcome"}
    }
    
    
    _, _ = goBot.ChannelMessageSend(welcomeChannel, fmt.Sprintf(WELCOME_MESSAGE, user.Mention(), guild.Name))
    
    return nil
}

func (guild *Guild) PrintGoodbye(user *discordgo.User) error {
    if !IsManaged(guild.ID) {
        return ERR_NOT_MANAGED
    }
    
    welcomeChannel, err := guild.GetChannelID("welcome")
    if err != nil {
        return &botError{ERR_MISSING_CHANNEL, "welcome"}
    }
    
    
    _, _ = goBot.ChannelMessageSend(welcomeChannel, fmt.Sprintf(GOODBYE_MESSAGE, user.Username, guild.Name, user.Username))
    
    return nil
}

func IsValidTeam(s string) bool {
    for _, t := range teamRoles {
        if t == s {
            return true
        }
    }
    return false
}