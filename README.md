# Library for logging

## Using

```go
package main

import "github.com/iooojik-dev/go-logger"

func main() {
	logger.WriteLogsToFile(true, "./log.txt")
	logger.LogInfo("hello world!")
	logger.LogDebug("hello world!")
	logger.LogError("hello world!")
	logger.MakeError(nil, errors.New("smth went wrong"))
}
```

## Methods

#### Write logs to a file

```go
logger.WriteLogsToFile(write bool, logPath string)
```

#### Log without highlighting

```go
logger.LogInfo(msg any)
```

#### Log with green highlighting

```go
logger.LogPositive(msg any)
```

#### Log with yellow highlighting

```go
logger.LogDebug(msg any)
```

#### Log errors

```go
logger.LogError(msg any)
```

#### Make new instance of custom error

```go
logger.MakeError(msg *string, err error) *CustomError
```