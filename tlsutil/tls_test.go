package tlsutil_test

import (
	"github.com/go-chassis/foundation/tlsutil"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var rsaCertPEM = `-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJANq58YD5coE2MA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMjEwMzAxMDQzMDAzWhcNMjEwMzMxMDQzMDAzWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAySqpnHOj42/LcGguAIj/ZYDb96ccdAjhuSgScWnOiF2obMVuBxAkaYRa
IcBvphv2N6dCs/AiOzwhyi7d4X82NQ5ftMBjzEHZaRCXQtS2JLHmNi4iuA5GHN0Z
EjinwXeT8ZsJP1wIHtnqF7D8PZdhS8V/SYimx4ejYG3J/+AIDU4YYyb14/3jjVzy
X4UnMy1igPbPtx6CbjNxUaVCmy4RUbrLwYdY1k+QbGguhfk4YERiV0P5W2pZzVqn
9rjvrEdFn0lgyRjNqvsRVneEcd7Y+OqgXvB69wiFrEeoEq/qbsDYQNFEm30Bx0wi
cbqMMYsuZTDRcXPz8gveyjNw0E2zDQIDAQABo1AwTjAdBgNVHQ4EFgQUGL3vXIio
1B7uvqGmzTXNaDZ/u7YwHwYDVR0jBBgwFoAUGL3vXIio1B7uvqGmzTXNaDZ/u7Yw
DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAAjBY2o/Jyh6MLo4DletW
PyS/8f46HMWo8kPTeZ77oh7iJNnpzbI4pgJ5yVTR4RqAj25ibSE0UuOrRRAgEWzT
5Y4C0r+XZghxyt9XET2RSC+BJxm4rC+bvsIIE0fNgX21o5hhSfSBpIl5NZOdVIbw
3VkvFi2hOtViVRxkk28SdvymDgDYU6djixf3qGYlvE+YSMUfDNFflxLWNCRXOyK0
9YLSOLZyaX8VkENVSZb3OmSDQoCTpnmrVEKHp4OcjbqfKB2o1bvREJ4CRBHmGMgd
pkM+Xeu/qofei+ekZsGLIaceM+cSV4w42vaHdwC5HkTpV4fIaGHEwsI8FD7aWKrU
eA==
-----END CERTIFICATE-----
`

var rsaKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAySqpnHOj42/LcGguAIj/ZYDb96ccdAjhuSgScWnOiF2obMVu
BxAkaYRaIcBvphv2N6dCs/AiOzwhyi7d4X82NQ5ftMBjzEHZaRCXQtS2JLHmNi4i
uA5GHN0ZEjinwXeT8ZsJP1wIHtnqF7D8PZdhS8V/SYimx4ejYG3J/+AIDU4YYyb1
4/3jjVzyX4UnMy1igPbPtx6CbjNxUaVCmy4RUbrLwYdY1k+QbGguhfk4YERiV0P5
W2pZzVqn9rjvrEdFn0lgyRjNqvsRVneEcd7Y+OqgXvB69wiFrEeoEq/qbsDYQNFE
m30Bx0wicbqMMYsuZTDRcXPz8gveyjNw0E2zDQIDAQABAoIBAQC3Yoz8cu8UhvWO
o2pMUpeAkNf2DAGERhSAFme5vBrrdXX0soZ7KdwH1P/VhPhDFXp/gZrtLhwGo+qp
xc+/oZhpBZF51Wkk62KmxNkfs4nYKdUTzzsXTuvbpDMWyU8krz3PIuZrPBqrBTzC
HDXWcAniaUiAYHKpspzdaziaakDs3obT5grtg6UNB/s/mPBrvHH9ftmolOBOxLhH
rvNxoCyP/tAB8AW1AJ6g3TSK1SANPgqzROe0X0ZHzW1T6/bZAC3OPA2PdhhnkdjN
zH2OLu3cU1Ybgtku/Ce7/AO1yKX2qdTn3v9iO0xO8swqAetepAqk8RKhfx0w61yc
Nugyya6BAoGBAOW5IzpB9USbTTzDuvVBBTS1e4OJRmkkR9KtZTfLNjVY7aY1mhrl
/Fi3JX+qPT9OsLyiRq9/UKuNcVeBYiVamBV4eBtmYqivOnnMe6P8V20RCWmhgZqw
2nyQfgIVFgSnTOlVdg5YD+GXptnfl/pqO2f2Pl4vQdM/nemyIk/HodgPAoGBAOAt
UZ8HTTzVnzyjPq1VGe2mQBmwC1fMkpmD2/d56TIDT4vVbAsu9ge106mdB52ZvlNc
nkl6KvsBcCXyJT6S/bpOINFViFAvigRpmez1rmFMqBf2OvTnELtlZ3kTn8/X1uIX
bFIX2G1XAKAQ6MBCYLzvCLMTf+EkebENA4kQr0cjAoGBAJrdfY8nqgYvQBmHxgDS
bYUEF5ksMQhuifDQLh035HqAUe2r0xDxHHZeOWxgQtvr25+/MjHbfXG5b8BTG+wc
r8xBo46tLjOTtbMok+2QDwwa4SKR24KCWTiCXEBhIK/QbTwb/fNbkJE/oB7e6mDJ
vvSt/4uVBiY4i+dgzFrGNSgnAoGBAJ/RbzkSuYu/N/DA6LQl0YBNX7FwggWsAG+V
Q8JglVFkbtdf5dDrP9crV6S6IG3I55kClI4JnI6p7cv/n3HG1UB25oqWkcGowpp2
tpfqZtFTFxtOHaXu/Uy79FKrHOnOFJHG5SB5g4Af4IA8zdITAGhxeSBBrI9Ts7X3
cyfKT0tFAoGBAJvS+AZK93SLYvsRdyytevKuYMAqIO+YRBOfZ+NN9+vosGlzbmuj
8zHkQQ8yeMJASbTpEAakz2ybJeN5zNELMKGUrsGVefiW0Z/5skO2toeF1PWtF0G8
FF2uvhCC6xWFT/PzYycBYMzq6gUus3aJQ+wjFnEflZYj7wGzf4tGzA5I
-----END RSA PRIVATE KEY-----

`)

func TestLoadTLSCertificate(t *testing.T) {
	td := os.TempDir()
	f1, err := os.Create(filepath.Join(td, "key.pem"))
	assert.NoError(t, err)
	defer f1.Close()
	_, err = io.WriteString(f1, string(rsaKeyPEM))
	assert.NoError(t, err)

	f2, err := os.Create(filepath.Join(td, "cert.pem"))
	assert.NoError(t, err)
	defer f1.Close()
	_, err = io.WriteString(f2, rsaCertPEM)
	assert.NoError(t, err)

	t.Run("given no file,should failed", func(t *testing.T) {
		tlsCert, err := tlsutil.LoadTLSCertificate("abc.txt", "abc.txt", "fakepassphase", func(src string) (s string, err error) {
			return "fakepassphase", nil
		})
		assert.Nil(t, tlsCert)
		assert.Error(t, err)
	})
	t.Run("given rsa key pair,should success", func(t *testing.T) {
		tlsCert, err := tlsutil.LoadTLSCertificate(f2.Name(), f1.Name(), "", nil)
		assert.NotNil(t, tlsCert)
		assert.NoError(t, err)
	})
}
