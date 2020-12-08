package tlsutil_test

import (
	"testing"

	"github.com/go-chassis/foundation/security"
	"github.com/go-chassis/foundation/tlsutil"
	"github.com/stretchr/testify/assert"
)

func TestLoadTLSCertificateFileNotExist(t *testing.T) {
	var cip security.Cipher
	tlsCert, err := tlsutil.LoadTLSCertificate("abc.txt", "abc.txt", "fakepassphase", cip)
	assert.Nil(t, tlsCert)
	assert.Error(t, err)
}
