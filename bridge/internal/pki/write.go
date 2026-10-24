package pki

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

func writePair(dir string, name string, certificate []byte, key *ecdsa.PrivateKey) error {
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	err = writePEM(filepath.Join(dir, name+"-key.pem"), "EC PRIVATE KEY", keyDER, 0o600)
	if err != nil {
		return err
	}
	return writePEM(filepath.Join(dir, name+".pem"), "CERTIFICATE", certificate, 0o644)
}

func writePEM(path string, kind string, der []byte, mode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	err = pem.Encode(file, &pem.Block{Type: kind, Bytes: der})
	if err != nil {
		file.Close()
		return err
	}
	return file.Close()
}
