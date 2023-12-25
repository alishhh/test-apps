package main

import (
	"magnum.kz/services/auth-app/internal/app"
)

func main() {
	application := app.New()
	application.Run()
}
