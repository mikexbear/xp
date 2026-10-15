package probe

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
	"xp/internal/config"
	"xp/internal/subscription"
	"xp/internal/xray"
)

type Prober struct {
	subscriptions   []*subscription.Subscription
	probeUrls       []config.Url
	probeTimeoutSec int
}

func NewProber(subscriptions []*subscription.Subscription, probeUrls []config.Url, probeTimeoutSec int) *Prober {
	return &Prober{subscriptions: subscriptions, probeUrls: probeUrls, probeTimeoutSec: probeTimeoutSec}
}

func (p *Prober) ProbeSubscriptions() {
	for _, sub := range p.subscriptions {
		xrayConfigs := sub.GetXrayConfigs()

		for _, xrayConfig := range xrayConfigs {
			p.probeOneXrayConfig(xrayConfig)
		}
	}
}

func (p *Prober) probeOneXrayConfig(xrayConfig *xray.Config) {
	host, port, err := xrayConfig.GetSocksInboundHostAndPort()
	if err != nil {
		slog.Error(
			fmt.Sprintf("Failed to get the host and port of the SOCKS inbound (remarks: %s)", xrayConfig.Remarks),
			"err", err,
		)
		// TODO: save metrics
		return
	}

	proxyUrl, err := url.Parse(fmt.Sprintf("socks5://%s:%d", host, port))
	if err != nil {
		slog.Error(
			fmt.Sprintf("Invalid SOCKS inbound (remarks: %s, host: %s, port: %d)", xrayConfig.Remarks, host, port),
			"err", err,
		)
		// TODO: save metrics
		return
	}

	xrayServer := xray.NewServer(xrayConfig)
	if err := xrayServer.Start(); err != nil {
		slog.Error(fmt.Sprintf("Failed to start Xray server (remarks: %s)", xrayConfig.Remarks), "err", err)
		// TODO: save metrics
		return
	}

	defer func() {
		if err := xrayServer.Stop(); err != nil {
			slog.Error(fmt.Sprintf("Failed to stop Xray server (remarks: %s)", xrayConfig.Remarks), "err", err)
		}
	}()

	httpClient := http.Client{
		Timeout: time.Duration(p.probeTimeoutSec) * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	for _, probeUrl := range p.probeUrls {
		err = p.probeOneUrl(&httpClient, probeUrl)
		success := err == nil

		// TODO: save metrics

		slog.Info(
			fmt.Sprintf(
				"Remote node %q probe status: %t (probe URL: %s)",
				xrayConfig.Remarks,
				success,
				probeUrl.String(),
			),
			"err", err,
		)
	}
}

func (p *Prober) probeOneUrl(client *http.Client, probeUrl config.Url) error {
	probeUrlString := probeUrl.String()

	response, err := client.Get(probeUrlString)
	if err != nil {
		return fmt.Errorf("probe %q failed: %w", probeUrlString, err)
	}

	defer func() { _ = response.Body.Close() }()
	_, _ = io.Copy(io.Discard, response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("probe %q failed: unexpected status: %s", probeUrlString, response.Status)
	}

	return nil
}
