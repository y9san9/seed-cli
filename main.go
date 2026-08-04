package main

import (
    "fmt"
    "os"
    "seed/exchange"
)

func main() {
    if len(os.Args) < 2 {
        // TODO: proper help
        fmt.Println("Usage: seed <command>")
        os.Exit(1)
    }
    switch os.Args[1] {
    case "exchange":
        exchange.Run()
    default:
        fmt.Println("Unknown subcommand")
    }
}
