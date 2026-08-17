package raiderio

// Region is a region slug accepted by the Raider.IO API
type Region string

// Regions available in the Raider.IO API
const (
	WORLD Region = "world"
	US    Region = "us"
	EU    Region = "eu"
	KR    Region = "kr"
	TW    Region = "tw"
	CN    Region = "cn"
)

// Expansion identifies a World of Warcraft expansion by its Raider.IO API id
type Expansion int

// Expansions available in the Raider.IO API
const (
	MIDNIGHT           Expansion = 11
	WAR_WITHIN         Expansion = 10
	DRAGONFLIGHT       Expansion = 9
	SHADOWLANDS        Expansion = 8
	BATTLE_FOR_AZEROTH Expansion = 7
	LEGION             Expansion = 6
)
