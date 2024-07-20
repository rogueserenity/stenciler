// Package main provides the entry point for the Stenciler application.
package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/rogueserenity/stenciler/cmd"
)

// these values are automagically populated by Goreleaser.
var (
	version string
	commit  string
	date    string
)

func main() {
	logLevel := os.Getenv("STENCILER_LOG_LEVEL")
	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case slog.LevelDebug.String():
		level = slog.LevelDebug
	case slog.LevelWarn.String():
		level = slog.LevelWarn
	case slog.LevelError.String():
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	cmd.SetReleaseVersion(version)
	cmd.Execute()
}
