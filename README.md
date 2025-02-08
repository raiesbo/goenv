# GO ENV

Lightweight Go package that helps loading the environment variables from an `.env` file that can be saved anywhere in
the project's structure.

## Install

```shell
go get github.com/raiesbo/goenv@latest
```

## Examples

As easy as:

```go
package main

import (
	"goenv"
)

func main() {
	err := goenv.Load()
	if err != nil {
		panic(err)
	}

	// After loading, the environment variables are ready to be used.
	// Example: addr := os.Getenv("ADDR")
}
```

