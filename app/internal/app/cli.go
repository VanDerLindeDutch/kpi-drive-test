package cli

import (
	"kpi-drive-test/app/internal/config"
	"kpi-drive-test/app/internal/domain/arango"
	"log"
)

type Cli struct {
	cfg           *config.Config
	logger        *log.Logger
	arangoService *arango.Service
}

func NewCli(config *config.Config) Cli {
	return Cli{
		cfg:           config,
		arangoService: arango.NewService(config),
		logger:        log.Default(),
	}
}

func (app *Cli) Start() {
	err := app.arangoService.FetchDataAndSave()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
