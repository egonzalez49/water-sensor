package notify

import "github.com/egonzalez49/water-sensor/config"

type Notifier struct {
	Config *config.Config
}

func NewNotifier(cfg *config.Config) *Notifier {
	return &Notifier{
		Config: cfg,
	}
}

func (n *Notifier) Notify() {
	n.sendSms()
}
