package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"net"
	"time"
)

const leafLifetime = 397 * 24 * time.Hour

func (a Authority) IssueServer(dir string, hosts []string) error {
	template, err := baseTemplate("Accueil bridge", leafLifetime)
	if err != nil {
		return err
	}
	for _, host := range hosts {
		ip := net.ParseIP(host)
		if ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	return a.issue(dir, "server", template)
}

func (a Authority) IssueClient(dir string, name string) error {
	template, err := baseTemplate(name, leafLifetime)
	if err != nil {
		return err
	}
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	return a.issue(dir, "client-"+name, template)
}

func (a Authority) issue(dir string, name string, template *x509.Certificate) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	template.KeyUsage = x509.KeyUsageDigitalSignature
	der, err := x509.CreateCertificate(rand.Reader, template, a.Certificate, &key.PublicKey, a.Key)
	if err != nil {
		return err
	}
	return writePair(dir, name, der, key)
}
