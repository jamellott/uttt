package main

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// ErrValidationFailed is returned if the path provided was not valid
var ErrValidationFailed = errors.New("path validation failed")

// ValidatePath checks that a user is allowed to access the given path.
// Note that, because of implementation details, it will deny any files
// which contain successive periods
func ValidatePath(path string) (string, error) {
	validated := filepath.Join("ui/", path)

	if strings.Contains(validated, "..") ||
		(strings.Index(validated, "ui/") != 0 && validated != "ui") {
		return "", ErrValidationFailed
	}

	return "./" + validated, nil
}

func main() {
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/", func(e echo.Context) error {
		return e.Redirect(301, "/ui/index.html")
	})

	server.GET("/ui/*", func(e echo.Context) error {
		path := e.Request().URL.Path[len("/ui/"):]

		validatedPath, err := ValidatePath(path)
		if err != nil {
			return e.String(403, "Forbidden")
		}

		return e.File(validatedPath)
	})

	server.Start("0.0.0.0:8081")
}
