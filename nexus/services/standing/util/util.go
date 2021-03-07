package util

var LECTeamsAbvToName map[string]string
var LECTeamsNametoAbv map[string]string

const LECGamesPerTeam = 18

func init() {
	LECTeamsAbvToName = map[string]string{
		"AST": "Astralis",
		"XL":  "Excel",
		"S04": "Schalke 04",
		"FNC": "Fnatic",
		"G2":  "G2 Esports",
		"MAD": "MAD Lions",
		"MSF": "Misfits Gaming",
		"RGE": "Rogue",
		"SK":  "SK Gaming",
		"VIT": "Team Vitality",
	}

	LECTeamsNametoAbv = map[string]string{
		"Astralis":       "AST",
		"Excel":          "XL",
		"Schalke 04":     "S04",
		"Fnatic":         "FNC",
		"G2 Esports":     "G2",
		"MAD Lions":      "MAD",
		"Misfits Gaming": "MSF",
		"Rogue":          "RGE",
		"SK Gaming":      "SK",
		"Team Vitality":  "VIT",
	}
}
