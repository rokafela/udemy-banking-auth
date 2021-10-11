package main

import (
	"github.com/rokafela/udemy-banking-auth/app"
	"github.com/rokafela/udemy-banking-auth/logger"
)

func main() {
	logger.Info("application starting")
	app.Start()
}
