package main

import (
    "fmt"
    "os"
    "seed/exchange"
    "seed/session"
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
