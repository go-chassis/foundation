package tlsutil_test

import (
	"testing"

	"github.com/go-chassis/foundation/tlsutil"
	"github.com/stretchr/testify/assert"
)

func TestLoadTLSCertificateFileNotExist(t *testing.T) {
	tlsCert, err := tlsutil.LoadTLSCertificate("abc.txt", "abc.txt", "fakepassphase", func(src string) (s string, err error) {
		return "fakepassphase", nil
	})
	assert.Nil(t, tlsCert)
	assert.Error(t, err)
}
