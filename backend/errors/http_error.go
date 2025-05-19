package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type HTTPError struct {
	Message string
	Err     error
	Code    int
	Line    int
	File    string
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s:%d: %d - %s: %v", e.File, e.Line, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s:%d: %d - %s", e.File, e.Line, e.Code, e.Message)
}

func NewHTTPError(err error, msg string, code int) *HTTPError {
	_, file, line, _ := runtime.Caller(1)

	wd, _ := os.Getwd()

	wd = filepath.ToSlash(wd)
	file = filepath.ToSlash(file)

	fmt.Println(wd)
	fmt.Println(file)
	// Преобразуем в относительный путь
	relPath := strings.TrimPrefix(file, fmt.Sprintf("%s/", wd))

	fmt.Println(relPath)

	return &HTTPError{
		Code:    code,
		Message: msg,
		Err:     err,
		File:    relPath,
		Line:    line,
	}
}
