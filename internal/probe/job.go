package probe

import (
	"context"
	"time"
)

func InitJob(p *Prober, intervalSec int, ctx context.Context) {
	tickerDuration := time.Duration(intervalSec) * time.Second
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.ProbeSubscriptions()
		case <-ctx.Done():
			return
		}
	}
}
