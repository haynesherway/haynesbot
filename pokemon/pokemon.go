package pokemon

import (
	"encoding/json"
	//"errors"
	"fmt"
	//"github.com/argandas/gokemon"
	"io/ioutil"
	//"log"
	//"net/http"
	//"github.com/bwmarrin/discordgo"
	"sort"
	"strings"
	//"strconv"
)

const (
    POKEMON_FILE = "pokemon.json"
    MOVES_FILE = "move.json"
)

var PokemonMap map[string]Pokemon

type PokemonList struct {
    Pokemons []*Pokemon
}

type Pokemon struct {
    Name        string `json:"name"`
    ID          string `json:"id"`
    Dex         int `json:"dex"`
    Types       TypeList `json:"types"`
    Stats       PokemonStats `json:"stats"`
    Moves
    MaxCP       int `json:"maxCP"`
    CPData
}

type PokemonStats struct {
    BaseStamina     int `json:"baseStamina"`
    BaseAttack      int `json:"baseAttack"`
    BaseDefense     int `json:"baseDefense"`
}

type CPData struct {
    Max20 int
    Max25 int
    Min20 int
    Min25 int
    Max40 int
    Min40 int
}

func (p *Pokemon) GetRaidCPChart() (string) {
    possibleIVs := []int{15,14,13,12,11,10}

    //ivList := map[int]map[int]map[int]
    ivs := []ivStat{}
    
    str := fmt.Sprintf("CP Chart for **%s**:\n", p.Name)
    str += "`|  %  | Ak | Df | St |  20  |  25  |`\n"
    str += "`|----------------------------------|`\n"
    for _, a := range possibleIVs {
        for _, d := range possibleIVs {
            for _, s := range possibleIVs {
                percent := round(float64((a+d+s)*100)/float64(45))
                if percent < 90 {
                    continue
                }
                cp20 := p.GetCP(20.0, a, d, s)
                cp25 := p.GetCP(25.0, a, d, s)
                iv := ivStat{
                    Attack: a,
                    Defense: d,
                    Stamina: s,
                    CP20: cp20,
                    CP25: cp25,
                    Percent: percent,
                }
                ivs = append(ivs, iv)
            }
        }
    }
    
    ivs = SortChart(ivs)
    chart := []string{}
    for _, iv := range ivs {
        chart = append(chart, iv.PrintChartRow())
    }
    
    return str + strings.Join(chart, "\n")
}

func (p *Pokemon) GetRaidCPRange() (string) {
     min20 := p.GetCP(20.0, 10, 10, 10)
	 max20 := p.GetCP(20.0, 15, 15, 15)
	 min25 := p.GetCP(25.0, 10, 10, 10)
	 max25 := p.GetCP(25.0, 15, 15, 15)
     return fmt.Sprintf("__**CP For %s**__\nLevel 20: %v - **%v**\nLevel 25: %v - **%v**", p.Name, min20, max20, min25, max25)
}

func (p *Pokemon) GetIV(cp int, level float64, stardust int, best string) string {
    ivstat := &ivStat{
        Level: level,
        CP: cp,
        Stardust: stardust,
        Best: best,
    }
    return p.getIV(ivstat)
}

func (p *Pokemon) getIV(stats *ivStat) (string) {
    possibleIVs := []int{15,14,13,12,11,10,9,8,7,6,5,4,3,2,1,0}
    
    possibleLevels := []float64{}
    if stats.Level != 0.0 {
        possibleLevels = append(possibleLevels, stats.Level)
    } else if stats.Stardust != 0 {
        if _, ok := StardustMap[stats.Stardust]; ok {
            possibleLevels = StardustMap[stats.Stardust]
        } else {
           for k := range NewMultiplierMap {
                possibleLevels = append(possibleLevels, k)
            } 
        }
    } else {
        for k := range NewMultiplierMap {
            possibleLevels = append(possibleLevels, k)
        }
    }
    cp := stats.CP
    
    ivList := []ivStat{}
    
    message := fmt.Sprintf("Possible IVs for **%s** with CP of **%d**:\n", p.Name, cp)
    for _, l := range possibleLevels {
        for _, a := range possibleIVs {
            for _, d := range possibleIVs {
                for _, s := range possibleIVs {
                    calccp := p.GetCP(l, a, d, s)
                    
                    if stats.Best != "" {
                        beststr := ""
                        vals := []int{a,d,s}
                        sort.Ints(vals)
                        highest := vals[2]
                        if a == highest {
                            beststr += "a"
                        }
                        if d == highest {
                            beststr += "d"
                        }
                        if s == highest {
                            beststr += "s"
                        }
                        if beststr != stats.Best {
                            continue
                        }
                    }
                    if cp == calccp {
                        perc := round(float64((a+d+s)*100)/float64(45))
                        stat := ivStat{
                            Level: l,
                            Attack: a, 
                            Defense: d,
                            Stamina: s,
                            //CP20: cp20,
                            //CP25: cp25,
                            Percent: perc,
                        }
                        ivList = append(ivList, stat)
                    }
                }
            }
        }
    }
    
    if ivList == nil || len(ivList) == 0  {
        return ""
    }
    
    afterMessage := ""
    ivList = SortChart(ivList)
    chart := []string{}
    if len(ivList) > 50 {
        ivList = ivList[0:50]
        afterMessage = "\nFull Chart too long to display, refine results by adding more info if possible :("
    }
    for _, s := range ivList {
       chart = append(chart, s.PrintIVRow())
    }
    
    return message + strings.Join(chart, "\n") + afterMessage
}

func (p *Pokemon) GetRaidIV(raidcp int) (string) {
    possibleIVs := []int{15,14,13,12,11,10}

    ivList := []ivStat{}
    
    message := fmt.Sprintf("Possible IVs for **%s** with CP of **%d**:\n", p.Name, raidcp)
    for _, a := range possibleIVs {
        for _, d := range possibleIVs {
            for _, s := range possibleIVs {
                cp20 := p.GetCP(20.0, a, d, s)
                cp25 := p.GetCP(25.0, a, d, s)
                if raidcp == cp20 || raidcp == cp25 {
                    perc := round(float64((a+d+s)*100)/float64(45))
                    stat := ivStat{
                        Attack: a, 
                        Defense: d,
                        Stamina: s,
                        //CP20: cp20,
                        //CP25: cp25,
                        Percent: perc,
                    }
                    ivList = append(ivList, stat)
                }
            }
        }
    }
    
    if ivList == nil || len(ivList) == 0  {
        return ""
    }
    
    ivList = SortChart(ivList)
    chart := []string{}
    for _, s := range ivList {
       chart = append(chart, s.PrintIVRow())
    }
    
    return message + strings.Join(chart, "\n")
}

/*func (p *Pokemon) GetCP(level int, ivAttack int, ivDefense int, ivStamina int) (cp int){
    attack := GetStatValue(p.Stats.BaseAttack, ivAttack, level)
    defense := GetStatValue(p.Stats.BaseDefense, ivDefense, level)
    stamina := GetStatValue(p.Stats.BaseStamina, ivStamina, level) 
    
    cp = CalculateCP(attack, defense, stamina, level)
    return
}*/

func (p *Pokemon) GetCP(level float64, ivAttack int, ivDefense int, ivStamina int) (cp int){
    attack := GetStatValue(p.Stats.BaseAttack, ivAttack, level)
    defense := GetStatValue(p.Stats.BaseDefense, ivDefense, level)
    stamina := GetStatValue(p.Stats.BaseStamina, ivStamina, level) 
    
    cp = CalculateCP(attack, defense, stamina, level)
    return
}

func (p *Pokemon) GetMaxCP() (cp int) {
    if p.Stats.BaseAttack == 1 && p.Stats.BaseDefense == 1 && p.Stats.BaseStamina == 1 {
        return 0
    }
    return p.GetCP(40.0, 15, 15, 15)
}

func (p *Pokemon) GetTypeRelations() (relations map[string]map[string]float64) {
    relations = make(map[string]map[string]float64)
    relations["attack"] = make(map[string]float64)
    relations["defense"] = make(map[string]float64)
    
    for _, pt := range p.Types {
        attackScalars := GetAttackTypeScalars(pt.ID)
        for tName, tScalar := range attackScalars {
            if _, ok := relations["attack"][tName]; !ok {
                relations["attack"][tName] = 1
            }
            relations["attack"][tName] = relations["attack"][tName] * tScalar
        }
        
        defenseScalars := GetDefenseTypeScalars(pt.ID)
        for tName, tScalar := range defenseScalars {
            if _, ok := relations["defense"][tName]; !ok {
                relations["defense"][tName] = 1
            }
            relations["defense"][tName] = relations["defense"][tName] * tScalar
        }
    }

    return
}
    
func (p *Pokemon) PrintTypeChart() string {
        typeRelations := p.GetTypeRelations()
        
        superEffective := []string{}
        notEffective := []string{}
        weakness := []string{}
        resistance := []string{}
        
        //Attack
        for ty, sc := range typeRelations["attack"] {
            if sc > 1.9 {
                superEffective = append(superEffective, ty + "(x2)")  
            } else if sc >= 1.4 {
                superEffective = append(superEffective, ty)
            } else if sc <= .6 {
                notEffective = append(notEffective, ty + "(x2)")
            } else if sc <= .8 {
                notEffective = append(notEffective, ty)
            }
        }
        
        //Defense 
        for ty, sc := range typeRelations["defense"] {
            if sc > 1.9 {
                weakness = append(weakness, ty + "(x2)")
            } else if sc >= 1.4 {
                weakness = append(weakness, ty)
            } else if sc <= .6 {
                resistance = append(resistance, ty + "(x2)")
            } else if sc <= .8 {
                resistance = append(resistance, ty)
            }
        }
        
        msg := fmt.Sprintf("Type Effects for **%s** (%s):\n", p.Name, p.Types.Print())
        /*if len(doubleEffective) > 0 {
            msg += fmt.Sprintf("Double Effective Against: %s\n", strings.Join(doubleEffective, ", "))
        }*/
        if len(superEffective) > 0 {
            msg += fmt.Sprintf("Super Effective Against: %s\n", strings.Join(superEffective, ", "))
        }
        if len(notEffective) > 0 {
            msg += fmt.Sprintf("Not Very Effective Against: %s\n", strings.Join(notEffective, ", "))
        }
        if len(weakness) > 0 {
            msg += fmt.Sprintf("Weak To: %s\n", strings.Join(weakness, ", "))
        }
        if len(resistance) > 0 {
            msg += fmt.Sprintf("Resistant To: %s\n", strings.Join(resistance, ", "))
        }
        
        return msg
    }
    
   /* for _, pt := range Pokemon.Types {
        if ty, ok := TypeMap[pt.ID]; ok {
            for _, t := range ty {
                scalars := p.GetTypeScalars()
                for tName, tScalar := range scalars {
                    if _, ok := relations[tName]; !ok {
                        relations[tName] = 1
                    }
                    relations[t.Name] = relations[t.Name] * 
                }
                
            }
        }
    }*/

/*func PrintPokemonToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, fields []string) error {
    if len(fields) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pokemon command should be in the following format: !poke mewtwo")

		return nil
	}

	pokemonName := strings.ToLower(fields[1])
	
	p, err := gokemon.GetPokemon(pokemonName)
	fmt.Println(err)
	fmt.Println(p.String())
	if err != nil {
	    _, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Pokemon %s not recognized", fields[1]))
	}
	
	_, _ = s.ChannelMessageSend(m.ChannelID, p.String())
	
	return nil
}

func PrintWeaknessToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, fields []string) error {
    /*if len(fields) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Weakness command should be in the following format: !weakness mewtwo")

		return nil
	}

	pokemonName := strings.ToLower(fields[1])
	
	if p, err := gokemon.GetPokemon(pokemonName); err != nil {
	    types := p.Types
	    message := fmt.Sprintf("Weaknesses for **%s:** %s", p.Name, weaknessString)
	    _, _ = s.ChannelMessageSend(m.ChannelID, message)
	} else {
	    _, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Pokemon %s not recognized", fields[1]))
	}
	
	return nil
} */

func init() {
    PokemonMap = make(map[string]Pokemon)
    
    //Pokemon
    file, err := ioutil.ReadFile(POKEMON_FILE)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	pokemonList := []Pokemon{}
	err = json.Unmarshal(file, &pokemonList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	for _, poke := range pokemonList {
	    PokemonMap[strings.ToLower(poke.Name)] = poke
	}
	
    return
}