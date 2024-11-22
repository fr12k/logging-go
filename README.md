# `go/errlog` module

The `errlog` library provides enhanced error handling and logging utilities for Go applications. It focuses on enriching error logs with contextual information such as function names, file paths, and line numbers to help developers debug effectively.

## Features

1. **Custom Errors with Stack Traces**:
   - Implements `CallerError`, an error type that includes a message and stack frame details.
   - Provides methods to extract the caller's function name and file details.

2. **Custom Log Handler**:
   - A `CustomHandler` for the `slog` logging framework that enriches log records with error context.
   - Automatically adds the function name and file location to log entries containing `CallerError`.

3. **Unit Tests**:
   - Comprehensive test suite to validate the behavior of `CallerError` and `CustomHandler`.

---

## Installation

Add the library to your project using:

```bash
go get github.com/goflink/go/errlog
```

---

## Usage

### 1. Creating and Using `CallerError`

The `CallerError` type provides detailed error context, including caller function and file information.

```go
package main

import (
	"fmt"
	"github.com/goflink/go/errlog/pkg/errors"
)

func main() {
	err := errors.New("something went wrong")
	funcName, fileDetails := err.ExtractCallerInfo()

	fmt.Printf("Error: %s\nFunction: %s\nFile: %s\n", err.Error(), funcName, fileDetails)
}
```

**Output**:
```
Error: something went wrong
Function: main.main
File: /path/to/main.go:10
```

---

### 2. Custom Log Handler with `slog`

Integrate the `CustomHandler` to enrich your logs with error context automatically.

```go
package main

import (
	"log/slog"
  "os"
	"github.com/goflink/go/errlog/pkg/handler"
	"github.com/goflink/go/errlog/pkg/errors"
)

func main() {
	// Set up the custom handler
	baseHandler := slog.NewTextHandler(os.StdOut, nil)
	customHandler := handler.New(baseHandler)
	logger := slog.New(customHandler)
	slog.SetDefault(logger)

  err := errors.New("something went wrong") // some function that returns an error
  if err != nil {
	  slog.Error("An error occurred", "error", errors.Wrap(err))
  }

  slog.Error("just an error message")
}
```

**Sample Log**:
```
level=ERROR msg="An error occurred" error="something went wrong" func=main.main file=/path/to/main.go:16
level=ERROR msg="just an error message"
```

---

## Testing

To run the tests:

```bash
go test ./... -v
```

---

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.
