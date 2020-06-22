package resp

import (
	"crypto"
	"encoding/json"
	"fmt"
	c "github.com/SermoDigital/jose/crypto"
)

type SigningMethodSM struct {
	Name string
	Hash crypto.Hash
	_    struct{}
}

var (
	SigningMethodSM2 = &SigningMethodSM{
		Name: "SM2",
		Hash: crypto.Hash(0),
	}
)

func (m *SigningMethodSM) Alg() string { return m.Name }

func (m *SigningMethodSM) Verify(raw []byte, sig c.Signature, key interface{}) error {
	fmt.Println("call verify")
	return nil
}

func (m *SigningMethodSM) Sign(data []byte, key interface{}) (c.Signature, error) {
	//m.Hash is dig's length
	//rsaKey, ok := key.(*rsa.PrivateKey)
	//if !ok {
	//	return nil, ErrInvalidKey
	//}
	//sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, m.sum(data))
	//if err != nil {
	//	return nil, err
	//}
	//return Signature(sigBytes), nil
	fmt.Println("call Sign")
	return []byte(`sign`), nil
}

//hash func
//func (m *SigningMethodSM) sum(b []byte) []byte {
//	h := m.Hash.New()
//	h.Write(b)
//	return h.Sum(nil)
//}

func (m *SigningMethodSM) Hasher() crypto.Hash { return m.Hash }

func (m *SigningMethodSM) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.Alg() + `"`), nil
}

var _ json.Marshaler = (*SigningMethodSM)(nil)
