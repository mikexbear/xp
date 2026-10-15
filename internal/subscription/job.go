package subscription

import (
	"context"
	"log/slog"
	"time"
)

func InitJob(subs []*Subscription, intervalSec int, ctx context.Context) {
	tickerDuration := time.Duration(intervalSec) * time.Second
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, sub := range subs {
				if err := sub.UpdateXrayConfigs(); err != nil {
					slog.Error("Failed to update subscription", "err", err)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
