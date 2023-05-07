# Library for logging

## Using

```go
package main

import "github.com/iooojik/go-logger"

func main() {
	logger.WriteLogsToFile(true, "./log.txt")
	logger.LogInfo("hello world!")
	logger.LogDebug("hello world!")
	logger.LogError("hello world!")
}
```

## Methods

#### Write logs to a file

```go
logger.WriteLogsToFile(logPath string)
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