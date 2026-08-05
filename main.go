package main

import (
	"fmt"
	"os"
	"seed/exchange"
	"seed/session"
)

func main() {
	fmt.Println("=========== Welcome to Seed Toolkit ===========")

	if len(os.Args) < 2 {
		fmt.Println("Usage: seed <command>")
		fmt.Println()
		fmt.Println("Available commands:")
		fmt.Println("• seed exchange       -  Start exchange process")
		fmt.Println("• seed session <peer> -  Start session with saved peer")
		fmt.Println()
		fmt.Println("Seed is offline set of utilities that allows you to exchange sensitive data over untrusted channels.")
		fmt.Println("It is designed to be the last emergency channel that is used to find more sustainable communication channels.")
		fmt.Println("Hence, no app, no servers, just CLI for mathematical algorithms.")
		fmt.Println()
		fmt.Println("Settings directory: ~/.seed")
		fmt.Println("Upstream: https://github.com/y9san9/seed-cli")
		os.Exit(0)
	}
	switch os.Args[1] {
	case "exchange":
		exchange.Run()
	case "session":
		if len(os.Args) < 3 {
			fmt.Println("Usage: seed session <peer>")
			os.Exit(1)
		}
		session.Run(os.Args[2])
	default:
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
