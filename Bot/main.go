package main

import (
	bot "github.com/Discipuls/LexiLens/Bot/Bot"
)

func main() {
	bot.Wg.Add(10000)
	bot.Main()
}
