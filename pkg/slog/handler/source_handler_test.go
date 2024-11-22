package handler

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var err2 = errors.New("global error")

func TestSourceCallerError(t *testing.T) {
	slog, buffer := setupSlog()

	slog.Info("info message", "key", "value")
	slog.Error("detail error log", "error", err2)
	slog.Error("detail error log", "error", errors.New("global original error"))
	slog.Error("detail error log", "error", &os.PathError{Op: "open", Path: "/dev/null", Err: errors.New("not a directory")})
	slog.With("key", "value").
		Error("with attributes error log", "error", errors.New("test error"))

	slog.WithGroup("user").Error("with group error log", "error", errors.New("test error"))

	lines := buffer.Logs()

	assert.Len(t, lines, 6, "Log should contain 3 lines")

	assert.Equal(t, `level=INFO msg="info message" key=value`, lines[0])
	assert.Equal(t, `level=ERROR msg="detail error log" error="global error" source=pkg/slog/handler/source_handler_test.go:22`, lines[1])
	assert.Equal(t, `level=ERROR msg="detail error log" error="global original error" source=pkg/slog/handler/source_handler_test.go:23`, lines[2])
	assert.Equal(t, `level=ERROR msg="detail error log" error="open /dev/null: not a directory" source=pkg/slog/handler/source_handler_test.go:24`, lines[3])
	assert.Equal(t, `level=ERROR msg="with attributes error log" key=value error="test error" source=pkg/slog/handler/source_handler_test.go:26`, lines[4])
	assert.Equal(t, `level=ERROR msg="with group error log" user.error="test error" user.source=pkg/slog/handler/source_handler_test.go:28`, lines[5])
}

func TestNewSourceHandler(t *testing.T) {
	h := NewSourceHandler(slog.NewTextHandler(os.Stdout, nil), nil)
	assert.NotNil(t, h)
}

func BenchmarkHandler(b *testing.B) {
	slog, _ := setupSlog()

	var table = []struct {
		input int
	}{
		{input: 100},
		{input: 1000},
		{input: 10000},
		{input: 100000},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.input), func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					for i := 0; i < v.input; i++ {
						slog.Error("error message", "error", errors.New("test error"))
					}
				}
			})
		})
	}
}

// Test helpers

type BufferedHandler struct {
	buffer *bytes.Buffer
}

func (h BufferedHandler) Logs() []string {
	return slices.DeleteFunc(strings.Split(h.buffer.String(), "\n"), func(e string) bool {
		return e == ""
	})
}

func setupSlog() (*slog.Logger, BufferedHandler) {
	var buf bytes.Buffer

	opts := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
		AddSource: false,
	}
	h := slog.NewTextHandler(&buf, &opts)
	baseHandler := NewSourceHandler(h, &HandlerOptions{StripFilePath: getRelativePath, AddSource: false})
	slogger := slog.New(baseHandler)
	slog.SetDefault(slogger)

	return slogger, BufferedHandler{
		buffer: &buf,
	}
}

func getRelativePath(absPath string) string {
	projectRoot := "pkg"
	index := strings.Index(absPath, projectRoot)
	if index == -1 {
		return absPath
	}
	return absPath[index:]
}
