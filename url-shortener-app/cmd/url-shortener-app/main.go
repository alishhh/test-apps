package main

import (
	"github.com/alishhh/url-shortener-app/internal/app"
)

func main() {
	application := app.New()
	application.Run()
}
