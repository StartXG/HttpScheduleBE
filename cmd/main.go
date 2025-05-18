package main

import (
	"HttpScheduleBE/cmd/commands"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.NewSource(time.Now().UnixNano())

	rootCmd := commands.HttpScheduleBeCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
