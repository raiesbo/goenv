# Go env

Lightweight Go package that helps loading the environment variables from an `.env` file that can be stored anywhere in
the project's structure.

This package does not contain any additional dependencies.

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

## License
Copyright (c) 2025 [Raimon Espasa Bou](https://github.com/raiesbo) 

Licensed under [MIT License](./LICENSE)
