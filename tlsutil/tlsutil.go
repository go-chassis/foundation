package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
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

func ParseSSLCipherSuites(ciphers string, permitTLSCipherSuiteMap map[string]uint16) []uint16 {
	if len(ciphers) == 0 || len(permitTLSCipherSuiteMap) == 0 {
		return nil
	}

	cipherSuiteList := make([]uint16, 0)
	cipherSuiteNameList := strings.Split(ciphers, ",")
	for _, cipherSuiteName := range cipherSuiteNameList {
		cipherSuiteName = strings.TrimSpace(cipherSuiteName)
		if len(cipherSuiteName) == 0 {
			continue
		}

		if cipherSuite, ok := permitTLSCipherSuiteMap[cipherSuiteName]; ok {
			cipherSuiteList = append(cipherSuiteList, cipherSuite)
		}
	}

	return cipherSuiteList
}

func ParseDefaultSSLCipherSuites(ciphers string) []uint16 {
	return ParseSSLCipherSuites(ciphers, TLSCipherSuiteMap)
}

func ParseSSLProtocol(sprotocol string) uint16 {
	var result uint16 = tls.VersionTLS12
	if protocol, ok := TLSVersionMap[sprotocol]; ok {
		result = protocol
	} else {
		panic(fmt.Sprintf("invalid ssl minimal version(%s), use default.", sprotocol))
	}

	return result
}

// GetClientTLSConfig
//  verifyPeer    Whether verify client
//  supplyCert    Whether send certificate
//  verifyCN      Whether verify CommonName
func GetClientTLSConfig(ops ...TLSOption) (tlsConfig *tls.Config, err error) {
	opts := toTLSOptions(ops...)
	var pool *x509.CertPool = nil
	var certs []tls.Certificate
	if opts.VerifyPeer {
		pool, err = GetX509CACertPool(opts.CACertFile)
		if err != nil {
			return nil, err
		}
	}

	if len(opts.CertFile) > 0 {
		certs, err = LoadTLSCertificate(opts.CertFile, opts.KeyFile, opts.KeyPassphase, opts.Decrypt)
		if err != nil {
			return nil, err
		}
	}

	tlsConfig = &tls.Config{
		RootCAs:            pool,
		Certificates:       certs,
		CipherSuites:       opts.CipherSuites,
		InsecureSkipVerify: !opts.VerifyPeer || !opts.VerifyHostName,
		MinVersion:         opts.MinVersion,
		MaxVersion:         opts.MaxVersion,
	}

	return tlsConfig, nil
}

func GetServerTLSConfig(ops ...TLSOption) (tlsConfig *tls.Config, err error) {
	opts := toTLSOptions(ops...)
	clientAuthMode := tls.NoClientCert
	var pool *x509.CertPool = nil
	if opts.VerifyPeer {
		pool, err = GetX509CACertPool(opts.CACertFile)
		if err != nil {
			return nil, err
		}

		clientAuthMode = tls.RequireAndVerifyClientCert
	}

	var certs []tls.Certificate
	if len(opts.CertFile) > 0 {
		certs, err = LoadTLSCertificate(opts.CertFile, opts.KeyFile, opts.KeyPassphase, opts.Decrypt)
		if err != nil {
			return nil, err
		}
	}

	tlsConfig = &tls.Config{
		ClientCAs:                pool,
		Certificates:             certs,
		CipherSuites:             opts.CipherSuites,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		PreferServerCipherSuites: true,
		ClientAuth:               clientAuthMode,
		MinVersion:               opts.MinVersion,
		MaxVersion:               opts.MaxVersion,
		NextProtos:               []string{"h2", "http/1.1"},
	}

	return tlsConfig, nil
}
