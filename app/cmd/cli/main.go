package main

import (
	cli "kpi-drive-test/app/internal/app"
	"kpi-drive-test/app/internal/config"
)

func main() {
	cfg := config.GetConfig()
	app := cli.NewCli(cfg)
	app.Start()
}
