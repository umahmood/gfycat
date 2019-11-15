# Gyfcat

A Go library which generates Gfycat strings in the format *AdjectiveAdjectiveAnimal*.

*"Most randomly generated URLs look like this: G1XeD4SwlHReDA. We thought it would be fun to do it differently. Our URLs follow the nomenclature: AdjectiveAdjectiveAnimal This is enough to give us a namespace of billions, while also letting humans write them easier."*. [https://gfycat.com/about]()

# Installation

> $ go get github.com/umahmood/gfycat

# Usage

```
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
    fmt.Println(g.GenerateNameOrder(gfycat.AnimalSecond))
}
```
Output:
```
similardeardowitcher
nastyminibeastagonizing
```

# Documentation

- http://godoc.org/github.com/umahmood/gfycat

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

