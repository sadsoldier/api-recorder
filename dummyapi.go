/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "fmt"
    "dummyapi/server"
)

func main() {
    server := server.New()
    server.Configure()
    err := server.Run()
    if err != nil {
        fmt.Println("exit on error: ", err)
    }
}
