package facts

import "kpi-drive-test/app/internal/config"

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}
