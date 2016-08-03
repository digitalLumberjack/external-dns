package utils

import (
	"strings"
	"os"
)

type DnsRecord struct {
	Fqdn    string
	Records []string
	Type    string
	TTL     int
}

type ServiceDnsRecord struct {
	Fqdn        string
	ServiceName string
	StackName   string
}

func ConvertToServiceDnsRecord(dnsRecord DnsRecord) ServiceDnsRecord {
	fqdnSeparator := os.Getenv("FQDN_SEPARATOR")
	if len(fqdnSeparator) == 0 {
		fqdnSeparator = "."
	}
	splitted := strings.Split(dnsRecord.Fqdn, ".")
	if fqdnSeparator != "." {
		splitted = strings.Split(splitted[0], fqdnSeparator)
	}
	serviceRecord := ServiceDnsRecord{dnsRecord.Fqdn, splitted[0], splitted[1]}
	return serviceRecord
}

func ConvertToFqdn(serviceName, stackName, environmentName, rootDomainName string) string {
        fqdnSeparator := os.Getenv("FQDN_SEPARATOR")
        if len(fqdnSeparator) == 0 {
                fqdnSeparator = "."
        }
	labels := []string{serviceName, stackName, environmentName}
	return strings.ToLower(strings.Join(labels, fqdnSeparator)) + "." + rootDomainName
}

// Fqdn ensures that the name is a fqdn adding a trailing dot if necessary.
func Fqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	}
	return name + "."
}

// UnFqdn converts the fqdn into a name removing the trailing dot.
func UnFqdn(name string) string {
	n := len(name)
	if n != 0 && name[n-1] == '.' {
		return name[:n-1]
	}
	return name
}
