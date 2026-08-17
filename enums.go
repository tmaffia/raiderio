package raiderio

type Region string

const (
	WORLD Region = "world"
	US    Region = "us"
	EU    Region = "eu"
	KR    Region = "kr"
	TW    Region = "tw"
	CN    Region = "cn"
)

type Expansion int

const (
	MIDNIGHT           Expansion = 11
	WAR_WITHIN         Expansion = 10
	DRAGONFLIGHT       Expansion = 9
	SHADOWLANDS        Expansion = 8
	BATTLE_FOR_AZEROTH Expansion = 7
	LEGION             Expansion = 6
)
