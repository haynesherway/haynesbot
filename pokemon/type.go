package pokemon

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "strings"
)

const (
    TYPES_FILE = "type.json"
)

var TypeMap = map[string]Type{}
var Type2ID= map[string]string{}

type Type struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Damage []*TypeDamage `json:"damage"`
}

type TypeList []*PokemonType

type PokemonType struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
}

type TypeDamage struct {
    ID string `json:"id"`
    Scalar float64 `json:"attackScalar"`
}

func (typeList TypeList) Print() string {
    types := []string{}
    for _, t := range typeList {
        types = append(types, t.Name)
    }
    return strings.Join(types, ", ")
}
func (t *Type) PrintTypeChart() string {
        typeRelations := make(map[string]map[string]float64)
        typeRelations["attack"] = GetAttackTypeScalars(t.ID)
        typeRelations["defense"] = GetDefenseTypeScalars(t.ID)
        
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
        
        msg := fmt.Sprintf("Type Effects for **%s**:\n", t.Name)
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

func GetAttackTypeScalars(id string) (map[string]float64) {
    typeScalars := map[string]float64{}
    if ty, ok := TypeMap[id]; ok {
    
        for _, damage := range ty.Damage {
            typeScalars[TypeMap[damage.ID].Name] = damage.Scalar
        }
    }
    
    return typeScalars
}

func GetDefenseTypeScalars(id string) (map[string]float64) {
    if _, ok := TypeMap[id]; !ok {
        return nil
    } 
    
    typeScalars := map[string]float64{}
    for _, ty := range TypeMap {
        for _, typeDamage := range ty.Damage {
            if typeDamage.ID == id {
                typeScalars[ty.Name] = typeDamage.Scalar
            }
        }
    }
    
    return typeScalars
}

func init() {
    TypeMap = make(map[string]Type)
    
    //Types
	file, err := ioutil.ReadFile(TYPES_FILE)
	if err != nil {
	    fmt.Println(err.Error())
	    return
	}
	
	typeList := []Type{}
	err = json.Unmarshal(file, &typeList)
	if err != nil {
	    fmt.Println(err.Error())
	    return
	}
	
	for _, ty := range typeList {
	    TypeMap[ty.ID] = ty
	    Type2ID[strings.ToLower(ty.Name)] = ty.ID
	}
}