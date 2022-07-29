package notify

import (
	"github.com/egonzalez49/water-sensor/config"
	"github.com/egonzalez49/water-sensor/logging"
)

type Notifier struct {
	Config *config.Config
	Logger *logging.Logger
}

func NewNotifier(cfg *config.Config, logger *logging.Logger) *Notifier {
	return &Notifier{
		Config: cfg,
		Logger: logger,
	}
}

func (n *Notifier) Notify() {
	n.sendSms()
}
