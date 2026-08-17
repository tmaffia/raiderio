# Raider.io API Go wrapper

[![Go Reference](https://pkg.go.dev/badge/github.com/tmaffia/raiderio.svg)](https://pkg.go.dev/github.com/tmaffia/raiderio)
![Go Build & Test](https://github.com/tmaffia/raiderio/actions/workflows/go.yml/badge.svg)

Wrapper for the raider.io API written in Go

## Usage

### Add the module

```
go get github.com/tmaffia/raiderio
```

### Authentication

The Raider.IO API provides higher rate limits for authenticated requests. You can generate an API key by registering your application on the [Raider.IO Application Settings](https://raider.io/settings/apps) page.

```go
package main

import "github.com/tmaffia/raiderio"

func main() {
	_ = raiderio.NewClient("YOUR_API_KEY")
}
```

`NewClient()` with no argument works for unauthenticated requests.

### Get a Character Profile

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	profile, err := client.GetCharacter(context.Background(), &raiderio.CharacterQuery{
		Region:        raiderio.US,
		Realm:         "illidan",
		Name:          "thehighvalue",
		TalentLoadout: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(profile.Class) // Mage
}
```

### Get a Guild Profile

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	profile, err := client.GetGuild(context.Background(), &raiderio.GuildQuery{
		Region:  raiderio.US,
		Realm:   "illidan",
		Name:    "warpath",
		Members: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(profile.Name)
}
```

### Get Raid Rankings for a specific raid

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	rankings, err := client.GetRaidRankings(context.Background(), &raiderio.RaidQuery{
		Slug:       "nerubar-palace",
		Difficulty: raiderio.MYTHIC_RAID,
		Region:     raiderio.US,
		Limit:      10,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(len(rankings.RaidRanking))
}
```

### Get Static Raid data by expansion

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	raids, err := client.GetRaids(context.Background(), raiderio.WAR_WITHIN)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(raids.Raids))
}
```

### Get the current Mythic+ affixes

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	affixes, err := client.GetMythicPlusAffixes(context.Background(), &raiderio.AffixesQuery{
		Region: raiderio.US,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(affixes.Title)
}
```

### Get a character's Mythic+ score

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	profile, err := client.GetCharacter(context.Background(), &raiderio.CharacterQuery{
		Region:                    raiderio.US,
		Realm:                     "illidan",
		Name:                      "thehighvalue",
		MythicPlusScoresBySeason: []string{"current"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(profile.MythicPlusScoresBySeason[0].Scores.All)
}
```

### Get the Mythic+ runs leaderboard

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmaffia/raiderio"
)

func main() {
	client := raiderio.NewClient()

	runs, err := client.GetMythicPlusRuns(context.Background(), &raiderio.MythicPlusRunsQuery{
		Region:  raiderio.US,
		Season:  "season-tww-3",
		Dungeon: "all",
		Affixes: "all",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(len(runs.Rankings))
}
```
