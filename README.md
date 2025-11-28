# Raider.io API Go wrapper

[![Go Reference](https://pkg.go.dev/badge/github.com/tmaffia/raiderio.svg)](https://pkg.go.dev/github.com/tmaffia/raiderio)
![Go Build & Test](https://github.com/tmaffia/raiderio/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tmaffia/raiderio)](https://goreportcard.com/report/github.com/tmaffia/raiderio)

Wrapper for the raider.io API written in Go

## Usage

### Include module in your go.mod

```
include github.com/tmaffia/raiderio v0.4.0
```

### Get a Character Profile

```go
client, err := raiderio.NewClient()

profile, err := client.GetCharacter(&CharacterQuery{
 Region: regions.US,
 Realm:  "illidan",
 Name:   "thehighvalue",
 TalentLoadout: true,
})

fmt.Println(profile.Class) // Mage
```

### Get a Guild Profile

```go
gq := raiderio.GuildQuery{
 Region: regions.US,
 Realm:  "illidan",
 Name:   "warpath",
 Members: true,
}

profile, err := client.GetGuild(&gq)
```

### Get Raid Rankings for a specific raid

```go
rq := raiderio.RaidQuery{
 Name:   "nerubar-palace",
 Difficulty: raiderio.MYTHIC_RAID,
 Region:  regions.US,
 Limit:   10,
}

rankings, err := client.GetRaidRankings(&rq)
```

### Get Static Raid data by expansion

```go
raids, err := client.GetRaids(expansions.WAR_WITHIN)
```
