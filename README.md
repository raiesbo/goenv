# Go env

A lightweight Go package that loads environment variables from an .env file, which can be located anywhere within the
project's structure.

This package has no external dependencies.

## Install

```shell
go get github.com/raiesbo/goenv
```

## Examples

As easy as:

```go
package main

import (
	"goenv"
)

func main() {
    if err := goenv.Load(); err != nil {
        panic(err)
    }
    
    // After loading, the environment variables are ready to be used. 
    // Example: addr := os.Getenv("ADDR")
}
```
