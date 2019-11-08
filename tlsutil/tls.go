package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-chassis/foundation/security"

)

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
func LoadTLSCertificate(certFile, keyFile, passphase string, cipher security.Cipher) ([]tls.Certificate, error) {
	certContent, err := ioutil.ReadFile(certFile)
	if err != nil {
		errorMsg := "read cert file" + certFile + "failed."
		return nil, errors.New(errorMsg)
	}

	keyContent, err := ioutil.ReadFile(keyFile)
	if err != nil {
		errorMsg := "read key file" + keyFile + "failed."
		return nil, errors.New(errorMsg)
	}

	keyBlock, _ := pem.Decode(keyContent)
	if keyBlock == nil {
		errorMsg := "decode key file " + keyFile + " failed"
		return nil, errors.New(errorMsg)
	}

	plainpass, err := cipher.Decrypt(passphase)
	if err != nil {
		return nil, err
	}

	if x509.IsEncryptedPEMBlock(keyBlock) {
		keyData, err := x509.DecryptPEMBlock(keyBlock, []byte(plainpass))
		if err != nil {
			errorMsg := "decrypt key file " + keyFile + " failed."
			return nil, errors.New(errorMsg)
		}

		// 解密成功，重新编码为无加密的PEM格式文件
		plainKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyData,
		}

		keyContent = pem.EncodeToMemory(plainKeyBlock)
	}

	cert, err := tls.X509KeyPair(certContent, keyContent)
	if err != nil {
		errorMsg := "load X509 key pair from cert file " + certFile + " with key file " + keyFile + " failed."
		return nil, errors.New(errorMsg)
	}

	var certs []tls.Certificate
	certs = append(certs, cert)

	return certs, nil
}
