package providers

import (
	"github.com/Sirupsen/logrus"
	rdns "github.com/rancher/external-dns/dns"
	"os"
	"fmt"
	"github.com/miekg/dns"
	"github.com/digitallumberjack/dnsapi"
)

type Bind9Handler struct {
	client *dnsapi.DNSApi
}

func init() {
	server := os.Getenv("BIND9_HOST")
	port := os.Getenv("BIND9_PORT")
	keyname := os.Getenv("BIND9_KEYNAME")
	key := os.Getenv("BIND9_KEY")
	if len(server) == 0 || len(port) == 0 || len(keyname) == 0 || len(key) == 0 {
		logrus.Info("BIND9 environnement not set, skipping init of BIND9 provider")
		return
	}

	bind9Handler := &Bind9Handler{}
	dnsapi := dnsapi.NewDNSApi(server, port, keyname, key, rdns.RootDomainName)
	bind9Handler.client = dnsapi
	if err := RegisterProvider("bind9", bind9Handler); err != nil {
		logrus.Fatal("Could not register bind9 provider")
	}
	logrus.Infof("Configured %s with zone %s and server %s", bind9Handler.GetName(), rdns.RootDomainName, server)

}

func (b *Bind9Handler) AddRecord(record rdns.DnsRecord) error {
	for _, rec := range record.Records {
		err := b.client.Add(record.Fqdn, rec, record.TTL)
		if err != nil {
			return fmt.Errorf("Bind9 API call has failed: %v", err)
		}
	}
	return nil
}
func (b *Bind9Handler) RemoveRecord(record rdns.DnsRecord) error {
	err := b.client.Remove(record.Fqdn)
	if err != nil {
		return fmt.Errorf("Bind9 API call has failed: %v", err)
	}
	return nil
}

func (b *Bind9Handler) UpdateRecord(record rdns.DnsRecord) error {
	err := b.RemoveRecord(record)
	if (err != nil) {
		return err
	} else {
		return b.AddRecord(record)
	}
}

func (b *Bind9Handler) GetRecords() ([]rdns.DnsRecord, error) {
	list, err := b.client.List()
	if (err != nil) {
		return nil, err
	} else {
		records := make([]rdns.DnsRecord, 0)
		for _, rr := range list {
			record := rdns.DnsRecord{}
			record.TTL = int(rr.Header().Ttl)
			record.Fqdn = rr.Header().Name
			record.Type = dns.Type(rr.Header().Rrtype).String()
			names := make([]string, 1)
			names[0] = rr.Header().Name
			record.Records = names
			records = append(records, record)
		}
		return records, nil
	}
}

func (b *Bind9Handler) GetName() string {
	return "bind9"
}
