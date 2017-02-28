# echo
a lite log package, provide [zap](https://github.com/uber-go/zap) style API


## install
```
go get -u -v github.com/vizee/echo
```


### usage
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
        echo.Warn("warn message", echo.String("c", c))
        echo.Error("error message", echo.Errors("err", errors.New("oops!")))
        echo.Fatal("fatal message", echo.Stack(true))
}
``
