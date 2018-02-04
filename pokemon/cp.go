package pokemon

import (
	"fmt"
	"math"
	"sort"
	
	"github.com/bwmarrin/discordgo"
)

type ivStat struct {
	Stardust int
	Attack int
	Defense int
	Stamina int
	Percent int
	CP int
	CP20 int
	CP25 int
	Level float64
	Best string
}

var StardustMap = map[int][]float64{
    200: {1.0,1.5,2.0,2.5},
    400: {3.0,3.5,4.0,4.5},
	600: {5.0,5.5,6.0,6.5},
	800: {7.0,7.5,8.0,8.5},
	1000: {9.0,9.5,10.0,10.5},
	1300: {11.0,11.5,12.0,12.5},
	1600: {13.0,13.5,14.0,14.5},
	1900: {15.0,15.5,16.0,16.5},
	2200: {17.0,17.5,18.0,18.5},
	2500: {19.0,19.5,20.0,20.5},
	3000: {21.0,21.5,22.0,22.5},
	3500: {23.0,23.5,24.0,24.5},
	4000: {25.0,25.5,26.0,26.5},
	4500: {27.0,27.5,28.0,28.5},
	5000: {29.0,29.5,30.0,30.5},
	6000: {31.0,31.5,32.0,32.5},
	7000: {33.0,33.5,34.0,34.5},
	8000: {35.0,35.5,36.0,36.5},
	9000: {37.0,37.5,38.0,38.5},
	10000: {39.0,39.5},
}

var NewMultiplierMap = map[float64]float64{
	1.0:	0.094,
	1.5:    0.135137432,
	2.0:	0.16639787,
	2.5:	0.192650919,
	3.0:	0.21573247,
	3.5:	0.236572661,
	4.0:	0.25572005,
	4.5:	0.273530381,
	5.0:	0.29024988,
	5.5:	0.335445036,
	6.0:	0.3210876,
	6.5:	0.335445036,
	7.0:	0.34921268,
	7.5:	0.362457751,
	8.0:	0.37523559,
	8.5:	0.387592406,
	9.0:	0.39956728,
	9.5:	0.411193551,
	10.0:	0.42250001,
	10.5:	0.432926419,
	11.0:	0.44310755,
	11.5:	0.4530599578,
	12.0:	0.46279839,
	12.5:	0.472336083,
	13.0:	0.48168495,
	13.5:	0.4908558,
	14.0:	0.49985844,
	14.5:	0.508701765,
	15.0:	0.51739395,
	15.5:	0.525942511,
	16.0:	0.53435433,
	16.5:	0.542635767,
	17.0:	0.55079269,
	17.5:	0.558830576,
	18.0:	0.56675452,
	18.5:	0.574569153,
	19.0:	0.58227891,
	19.5:	0.589887917,
	20.0:	0.59740001,
	20.5:	0.604818814,
	21.0:	0.61215729,
	21.5:	0.619399365,
	22.0:	0.62656713,
	22.5:	0.633644533,
	23.0:	0.64065295,
	23.5:	0.647576426,
	24.0:	0.65443563,
	24.5:	0.661214806,
	25.0:	0.667934,
	25.5:	0.674577537,
	26.0:	0.68116492,
	26.5:	0.687680648,
	27.0:	0.69414365,
	27.5:	0.700538673,
	28.0:	0.70688421,
	28.5:	0.713164996,
	29.0:	0.71939909,
	29.5:	0.725571552,
	30.0:	0.7317,
	30.5:	0.734741009,
	31.0:	0.73776948,
	31.5:	0.740785574,
	32.0:	0.74378943,
	32.5:	0.746781211,
	33.0:	0.74976104,
	33.5:   0.752729087,
	34.0:	0.75568551,
	34.5:	0.758630378,
	35.0:	0.76156384,
	35.5:	0.764486065,
	36.0:	0.76739717,
	36.5:	0.770297266,
	37.0:	0.7731865,
	37.5:	0.776064962,
	38.0:	0.77893275,
	38.5:	0.781790055,
	39.0:	0.78463697,
	39.5:	0.787473578,
	40.0:	0.79030001,
}

var MultiplierMap = map[int]float64{
1:	0.094,
2:	0.16639787,
3:	0.21573247,
4:	0.25572005,
5:	0.29024988,
6:	0.3210876,
7:	0.34921268,
8:	0.37523559,
9:	0.39956728,
10:	0.42250001,
11:	0.44310755,
12:	0.46279839,
13:	0.48168495,
14:	0.49985844,
15:	0.51739395,
16:	0.53435433,
17:	0.55079269,
18:	0.56675452,
19:	0.58227891,
20:	0.59740001,
21:	0.61215729,
22:	0.62656713,
23:	0.64065295,
24:	0.65443563,
25:	0.667934,
26:	0.68116492,
27:	0.69414365,
28:	0.70688421,
29:	0.71939909,
30:	0.7317,
31:	0.73776948,
32:	0.74378943,
33:	0.74976104,
34:	0.75568551,
35:	0.76156384,
36:	0.76739717,
37:	0.7731865,
38:	0.77893275,
39:	0.78463697,
40:	0.79030001,
}

func GetStatValue(base int, iv int, level float64) (value float64) {
		value = (float64(base) + float64(iv))
	
	return
}

func CalculateCP(attack float64, defense float64, stamina float64, level float64) (cp int) {
	if multiplier, ok := NewMultiplierMap[level]; ok {
		cp = int((attack * math.Pow(defense, 0.5) * math.Pow(stamina, 0.5) * math.Pow(multiplier, 2)) / 10)
	}
	if cp < 10 {
	    cp = 10
	}
	return
}

func PrintToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, fields []string) error {
	if len(fields) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "CP command should be in the following format: !cp mewtwo")

		return nil
	}

	//pokemon := fields[1]

	return nil

}

type By func(s1, s2 *ivStat) bool

func (by By) Sort(stats []ivStat) {
	ss := &statSorter{
		stats: stats,
		by: by,
	}
	sort.Sort(ss)
}

type statSorter struct {
	stats	[]ivStat
	by 	func(s1, s2 *ivStat) bool
}

func (s *statSorter) Len() int {
	return len(s.stats)
}

func (s *statSorter) Swap(i, j int) {
	s.stats[i], s.stats[j] = s.stats[j], s.stats[i]
} 

func (s *statSorter) Less(i, j int) bool {
	return s.by(&s.stats[i], &s.stats[j])
}

func SortChart(stats []ivStat) []ivStat {
	sort := func(s1, s2 *ivStat) bool {
		if s1.Level == s2.Level {
			if s1.CP25 == s2.CP25 {
				if s1.CP20 == s2.CP20 {
					if s1.Percent == s2.Percent {
						if s1.Attack == s2.Attack {
							if s1.Defense == s2.Defense {
								return s1.Stamina > s2.Stamina
							}
							return s1.Defense > s2.Defense
						}
						return s1.Attack > s2.Attack
					}
					return s1.Percent > s2.Percent
				}
				return s1.CP20 > s2.CP20
			}
			return s1.CP25 > s2.CP25
		}
		return s1.Level > s2.Level
	}
	
	By(sort).Sort(stats)
	return stats
}

func (s *ivStat) PrintChartRow() string {
	if s.Percent == 100 {
		return fmt.Sprintf("`| %d | %d | %d | %d | %d | %d |`", s.Percent, s.Attack, s.Defense, s.Stamina, s.CP20, s.CP25)
	}
	return fmt.Sprintf("`| %d%% | %d | %d | %d | %d | %d |`", s.Percent, s.Attack, s.Defense, s.Stamina, s.CP20, s.CP25)
}

func (s *ivStat) PrintIVRow() string {
	return fmt.Sprintf("`%4.1f | %2d | %2d | %2d ---> %d%%`", s.Level, s.Attack, s.Defense, s.Stamina, s.Percent)
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}
