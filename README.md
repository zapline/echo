# echo
A lite log package, provide [zap](https://github.com/uber-go/zap) style API.

## Install
```
go get -u -v github.com/vizee/echo
```

### Usage
```go
package main

import (
	"errors"

	"github.com/vizee/echo"
)

func main() {
	a := 0
	b := true
	c := "c"
	echo.Debug("debug message", echo.Int("a", a)) // default InfoLevel
	echo.Info("info message", echo.Bool("b", b))
	echo.Warn("warn message", echo.String("c", c), echo.Var("var", map[string]int{"a": 1, "b": 2}))
	echo.Error("error message", echo.Errors("err", errors.New("oops!")))
	echo.Fatal("fatal message", echo.Stack(true))
}
```
