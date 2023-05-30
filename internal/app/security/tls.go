package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"time"
)

// CertificateSerialNumber серийный номер сертификата.
var CertificateSerialNumber = big.NewInt(1)

// CreateCertificate создает TLS сертификат для веб-сервера.
func CreateCertificate() (cert tls.Certificate, err error) {
	x509Cert := &x509.Certificate{
		SerialNumber: CertificateSerialNumber,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return tls.Certificate{}, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, x509Cert, x509Cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert.Certificate = append(cert.Certificate, certBytes)
	cert.PrivateKey = privateKey
	cert.Leaf = x509Cert

	return cert, nil
}
