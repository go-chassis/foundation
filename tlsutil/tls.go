package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type Decrypt func(src string) (string, error)

//GetX509CACertPool is a function used to get certificate
func GetX509CACertPool(caCertFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("read ca cert file %s failed", caCertFile)
	}

	pool.AppendCertsFromPEM(caCert)
	return pool, nil
}

//LoadTLSCertificate is a function used to load a certificate
// RFC 1423 is insecure, password and decrypt is not required
func LoadTLSCertificate(certFile, keyFile, password string, decrypt Decrypt) ([]tls.Certificate, error) {
	certContent, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	keyContent, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyContent)
	if keyBlock == nil {
		return nil, fmt.Errorf("decrypt key file "+keyFile+" failed: %w", err)
	}
	if password != "" {
		var plainpass string
		if decrypt != nil {
			plainpass, err = decrypt(password)
			if err != nil {
				return nil, err
			}
		} else {
			plainpass = password
		}
		if x509.IsEncryptedPEMBlock(keyBlock) {
			keyData, err := x509.DecryptPEMBlock(keyBlock, []byte(plainpass))
			if err != nil {
				return nil, fmt.Errorf("decrypt key file "+keyFile+" failed. err: %w", err)
			}

			// 解密成功，重新编码为无加密的PEM格式文件
			keyBlock = &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: keyData,
			}

		}
	}
	keyContent = pem.EncodeToMemory(keyBlock)
	cert, err := tls.X509KeyPair(certContent, keyContent)
	if err != nil {
		return nil, fmt.Errorf("load cert failed:%w", err)
	}

	var certs []tls.Certificate
	certs = append(certs, cert)

	return certs, nil
}
