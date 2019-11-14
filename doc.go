/*
Package gfycat generates a random adjective + adjective + animal

Usage:

    package main

    import (
        "fmt"

        "github.com/umahmood/gfycat"
    )

    func main() {
        g, err := gfycat.New()
        if err != nil {
            // handle error
        }
        fmt.Println(g.GenerateName())
        fmt.Println(g.GenerateNameOrder(gfycat.AnimalFirst))
    }
*/
package gfycat
