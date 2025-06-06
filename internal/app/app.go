package app

import (
	"d2runewords/internal/db"
    "d2runewords/internal/ui"
)

func Run() {
	db, err := database.Connect()
	if err != nil {
		return
	}

	window := ui.NewMW(db)
	window.Run()
}