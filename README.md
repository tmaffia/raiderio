# Raider.io API Go wrapper

[![Go Reference](https://pkg.go.dev/badge/github.com/tmaffia/raiderio.svg)](https://pkg.go.dev/github.com/tmaffia/raiderio)
![Go Build & Test](https://github.com/tmaffia/raiderio/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tmaffia/raiderio)](https://goreportcard.com/report/github.com/tmaffia/raiderio)

Wrapper for the raider.io API written in Go 

## Usage

### Get a Character Profile
```go
client, err := raiderio.NewClient()
cq := raiderio.CharacterQuery{
	Region:        region.US,
	Realm:         "illidan",
	Name:          "highervalue",
	TalentLoadout: true,
}

profile, err := client.GetCharacter(&cq)
```

### Get a Guild Profile
```go
gq := raiderio.GuildQuery{
	Region: region.US,
	Realm:  "illidan",
	Name:   "warpath",
	Members: true,
}

profile, err := client.GetGuild(&gq)
```

### Get Raid Rankings for a specific raid
```go
rq := raiderio.RaidQuery{
	Name: 		"aberrus-the-shadowed-crucible",
	Difficulty:	raiderio.MythicRaid,
	Region: 	region.US,
	Limit: 		10,
}

rankings, err := client.GetRaidRankings(&rq)
```

### Get Static Raid data by expansion
```go
raids, err := client.GetRaids(expansion.Dragonflight)
```
