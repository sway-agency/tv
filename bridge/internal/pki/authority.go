package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"time"
)

const authorityLifetime = 10 * 365 * 24 * time.Hour

type Authority struct {
	Certificate *x509.Certificate
	Key         *ecdsa.PrivateKey
}

func NewAuthority(dir string) (Authority, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return Authority{}, err
	}
	template, err := baseTemplate("Accueil bridge local CA", authorityLifetime)
	if err != nil {
		return Authority{}, err
	}
	template.IsCA = true
	template.BasicConstraintsValid = true
	template.MaxPathLenZero = true
	template.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return Authority{}, err
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		return Authority{}, err
	}
	err = writePair(dir, "ca", der, key)
	return Authority{Certificate: certificate, Key: key}, err
}

func LoadAuthority(dir string) (Authority, error) {
	pair, err := tls.LoadX509KeyPair(dir+"/ca.pem", dir+"/ca-key.pem")
	if err != nil {
		return Authority{}, err
	}
	certificate, err := x509.ParseCertificate(pair.Certificate[0])
	if err != nil {
		return Authority{}, err
	}
	key, ok := pair.PrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return Authority{}, errors.New("ca-key.pem is not an ECDSA key")
	}
	return Authority{Certificate: certificate, Key: key}, nil
}

func baseTemplate(commonName string, lifetime time.Duration) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	now := time.Now().Add(-5 * time.Minute)
	return &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: commonName, Organization: []string{"Accueil"}},
		NotBefore:    now,
		NotAfter:     now.Add(lifetime),
	}, nil
}
