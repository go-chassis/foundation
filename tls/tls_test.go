package tlsutil_test

import (
	"github.com/go-chassis/foundation/security"
	"github.com/go-chassis/foundation/tls"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadTLSCertificateFileNotExist(t *testing.T) {
	var cip security.Cipher
	tlsCert, err := tlsutil.LoadTLSCertificate("abc.txt", "abc.txt", "fakepassphase", cip)
	assert.Nil(t, tlsCert)
	assert.Error(t, err)
}
