# Battlesnake Rules CLI

This tool allows running a Battlesnake game locally. 

### Installation

Compile from source. Requires Go 1.18 or higher.
```
git clone https://github.com/BattlesnakeOfficial/rules.git
cd rules
go build -o battlesnake ./cli/battlesnake/main.go
```

### Usage

Example command to run a game locally:
```
battlesnake play -W 11 -H 11 --name <SNAKE_NAME> --url <SNAKE_URL> -g solo -v
```

Complete usage documentation:
```
Usage:
  battlesnake play [flags]

Flags:
  -W, --width int                 Width of Board (default 11)
  -H, --height int                Height of Board (default 11)
  -n, --name stringArray          Name of Snake
  -u, --url stringArray           URL of Snake
  -t, --timeout int               Request Timeout (default 500)
  -g, --gametype string           Type of Game Rules (default "standard")
  -m, --map string                Game map to use to populate the board (default "standard")
      --browser                   View the game in the browser using the Battlesnake game board
      --board-url string          Base URL for the game board when using --browser (default "https://board.battlesnake.com")
      --foodSpawnChance int       Percentage chance of spawning a new food every round (default 15)
  -h, --help                      help for play

Global Flags:
      --config string   config file (default is $HOME/.battlesnake.yaml)
      --verbose         Enable debug logging
```

Battlesnake names and URLs will be paired together in sequence, for example:

```
battlesnake play --name Snake1 --url http://snake1-url-whatever --name Snake2 --url http://snake2-url-whatever
```

This will create a game with the following Battlesnakes:
* Snake1, http://snake1-url-whatever
* Snake2, http://snake2-url-whatever

Names are optional, and if you don't provide them UUIDs will be generated instead. However names are way easier to read and highly recommended!

URLs are technically optional too, but your Battlesnake will lose if the server is only sending move requests to http://example.com.

Example creating a 7x7 Standard game with two Battlesnakes:
```
battlesnake play --width 7 --height 7 --name Snake1 --url http://snake1-url-whatever --name Snake2 --url http://snake2-url-whatever
```
