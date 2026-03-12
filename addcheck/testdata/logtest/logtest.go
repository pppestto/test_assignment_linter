// Package logtest holds slog calls for analysistest; line-ending comments use backtick regex only.
package logtest

import (
	"log/slog"
)

var token string

func BadLowercase() {
	slog.Error("Failed to connect to database") // want `LowerCase`
}

func BadEnglish() {
	slog.Error("ошибка подключения к базе данных") // want `English`
}

func BadSensitive() {
	slog.Info("token: " + token) // want `sensitive`
}

func BadSymbols() {
	slog.Info("server started!✅") // want `Emoji` `English`
	slog.Error("connection failed!!!") // want `Emoji`
	slog.Warn("bad!") // want `Emoji`
}

func Good() {
	slog.Info("token validated")
	slog.Info("server started")
}
