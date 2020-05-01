package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	cfg, err := initializeConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config")
		panic(err)
	}

	server := buildServer(cfg)
	listenAddress := fmt.Sprintf("%v:%v", cfg.Host, strconv.Itoa(cfg.Port))
	if cfg.AcmeTLS {
		server.StartAutoTLS(listenAddress)
	} else {
		server.Start(listenAddress)
	}
}

// buildServer constructs an echo instance with the routes setup
// according to the configuration given
func buildServer(cfg *config) *echo.Echo {
	server := echo.New()

	server.Use(middleware.Recover())
	if cfg.RequestLogs {
		server.Use(middleware.Logger())
	}

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

	return server
}
