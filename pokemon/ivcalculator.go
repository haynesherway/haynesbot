package pokemon

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    STATUS_STARTING int = iota
    STATUS_EXPECTING_POKEMON 
    STATUS_GOT_POKEMON
    STATUS_EXPECTING_CP
    STATUS_GOT_CP
    /*STATUS_EXPECTING_STARDUST
    STATUS_GOT_STARDUST
    STATUS_EXPECTING_BEST
    STATUS_GOT_BEST
    STATUS_EXPECTING_TIER
    STATUS_GOT_TIER*/
    STATUS_EXPECTING_LEVEL
    STATUS_CALCULATING
    STATUS_DONE
)

type IVCalculator struct {
    Session *discordgo.Session
    InputChannel chan interface{}
    RunningCalculations map[string]*IVCalculation
    lock *sync.RWMutex
}

type IVCalculation struct {
    Session *discordgo.Session
    User *discordgo.User
    ChannelID string
    Pokemon *Pokemon
    IV  *ivStat
    Channel chan interface{}
    Status int
}

func StartIVCalculator(s *discordgo.Session) (*IVCalculator) {
    ivCalculator := IVCalculator{
        Session: s,
        InputChannel: make(chan interface{}),
        RunningCalculations: make(map[string]*IVCalculation),
        lock: &sync.RWMutex{},
    }

    go func() {
        for incoming := range ivCalculator.InputChannel {
            if m, ok := incoming.(*discordgo.MessageCreate); ok {
                if !ivCalculator.IsRunning(m.Author.ID) {
                    //Start New Calculation
                    ivCalculator.Start(m)
                } else {
                    ivCalc := ivCalculator.GetCalculation(m.Author.ID)
                    if ivCalc.Status == STATUS_DONE {
                        ivCalculator.Stop(m.Author.ID)
                        ivCalculator.Start(m)
                    } else if ivCalc.ChannelID == m.ChannelID {
                        ivCalc.Channel <- m.Content
                    }
                }
            }
        }
    }()
    return &ivCalculator
}

func (calc *IVCalculator) Start(m *discordgo.MessageCreate) {
    thisCalculation := &IVCalculation{
        Session: calc.Session,
        User: m.Author,
        ChannelID: m.ChannelID,
        Channel: make(chan interface{}),
    }
    calc.lock.Lock()
    calc.RunningCalculations[m.Author.ID] = thisCalculation
    calc.lock.Unlock()
    
    thisCalculation.AskQuestion()
    
    go func() {
        timeout := time.After(1 * time.Minute)
        for {
            select {
                case msg := <- thisCalculation.Channel:
                    if message, ok := msg.(string); ok {
                        thisCalculation.GetResponse(message)
                    }
                    
                case <-timeout:
                    if thisCalculation.Status == STATUS_DONE {
                        //YAY it completed sucessfully
                    } else {
                        thisCalculation.PrintToDiscord("Unable to process your IV Calculation, please try again.")
                    }
                    calc.Stop(m.Author.ID)
                    return
            }
        }
    }()
    
    return
}

func (calc *IVCalculator) Stop(userID string) {
    calc.lock.Lock()
    delete(calc.RunningCalculations, userID)
    calc.lock.Unlock()
}

func (calc *IVCalculator) IsRunning(userID string) bool {
    calc.lock.RLock()
    defer calc.lock.RUnlock()
    if _, ok := calc.RunningCalculations[userID]; ok {
        return true
    }
    return false
}

func (calc *IVCalculator) GetCalculation(userID string) (*IVCalculation) {
    calc.lock.RLock()
    defer calc.lock.RUnlock()
    if ivCalc, ok := calc.RunningCalculations[userID]; ok {
        return ivCalc
    }
    return nil
}

func (ivCalc *IVCalculation) AskQuestion() {
    
    switch ivCalc.Status {
        case STATUS_STARTING:
            ivCalc.PrintToDiscord("Enter pokemon name.")
        case STATUS_GOT_POKEMON:
            ivCalc.PrintToDiscord("Enter CP.")
        case STATUS_GOT_CP:
            ivCalc.PrintToDiscord("Enter level.")
            
    }
    
    ivCalc.Status++
    fmt.Println(ivCalc.Status)
    return
}

func (ivCalc *IVCalculation) GetResponse(m string) {
    
    switch ivCalc.Status {
    case STATUS_EXPECTING_POKEMON:
        if p, ok := PokemonMap[strings.ToLower(m)]; ok {
            ivCalc.Pokemon = &p
            ivCalc.Status++
            ivCalc.AskQuestion()
        } else {
            ivCalc.PrintToDiscord("Unrecognized pokemon. Try again.")
        }
    case STATUS_EXPECTING_CP:
        if cp, err := strconv.Atoi(m); err == nil {
            ivCalc.IV.CP = cp
            ivCalc.Status++
            ivCalc.AskQuestion()
        } else {
            ivCalc.PrintToDiscord(fmt.Sprintf("CP must be an integer, got %s. Try again.", m))
        }
    case STATUS_EXPECTING_LEVEL:
        if lvl, err := strconv.ParseFloat(m, 64); err == nil {
            ivCalc.IV.Level = lvl
            ivCalc.Status++
        } else {
            ivCalc.PrintToDiscord(fmt.Sprintf("Level must be an integer, got %s. Try again.", m))
        }
    }
    
    if ivCalc.Status == STATUS_CALCULATING {
        ivCalc.Calculate()
    }
    
    return
}

func (ivCalc *IVCalculation) Calculate() {
    msg := ivCalc.Pokemon.getIV(ivCalc.IV)
    ivCalc.PrintToDiscord(msg)
    ivCalc.Status++
}

func (ivCalc *IVCalculation) PrintToDiscord(m string) {
    _, _ = ivCalc.Session.ChannelMessageSend(ivCalc.ChannelID, ivCalc.User.Mention() + " " + m)
    return
}