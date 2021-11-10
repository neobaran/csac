package tencent

import (
	"time"

	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/neobaran/csac/config"
)

type DNSProvider struct {
	client *dnsPodClient
	config *config.Config
}

func NewDNSProvider(client *dnsPodClient, config *config.Config) *DNSProvider {
	return &DNSProvider{
		client: client,
		config: config,
	}
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	domainData := getDomain(fqdn)

	if err := d.client.CreateRecordData(domainData, value, d.config.TTL); err != nil {
		return err
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, _ := dns01.GetRecord(domain, keyAuth)
	domainData := getDomain(fqdn)

	list, err := d.client.ListRecordData(domainData)
	if err != nil {
		return err
	}

	for _, item := range list {
		_ = d.client.DeleteRecordData(domainData, item)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return time.Duration(d.config.TTL * uint64(time.Second)), dns01.DefaultPollingInterval
}
