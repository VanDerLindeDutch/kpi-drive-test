package cli

import (
	"kpi-drive-test/app/internal/config"
	"kpi-drive-test/app/internal/domain/arango"
	"kpi-drive-test/app/internal/domain/facts"
	"log"
)

type Cli struct {
	cfg           *config.Config
	logger        *log.Logger
	factService   *facts.Service
	arangoService *arango.Service
}

func NewCli(config *config.Config) Cli {
	return Cli{
		cfg:           config,
		factService:   facts.NewService(config),
		arangoService: arango.NewService(config),
		logger:        log.Default(),
	}
}

func (app *Cli) Start() {
	err := app.arangoService.Process()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
