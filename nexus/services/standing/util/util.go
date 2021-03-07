package util

var LECTeamsNametoAbv map[string]string

const LECGamesPerTeam = 18

func init() {
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
