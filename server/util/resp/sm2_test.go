package resp

import (
	"github.com/SermoDigital/jose/jws"
	"testing"
	"time"
)

func TestSM2(t *testing.T) {

	now := time.Now()
	claims := jws.Claims{}
	claims.SetIssuer("ICBC")  //appid
	claims.SetIssuedAt(now) //timestamp
	claims.SetNotBefore(now)
	claims.SetExpiration(now.Add(time.Duration(10000) * time.Second))
	claims.SetSubject("") //apiurl path
	claims.SetAudience("")
	claims.SetJWTID("") //msgid

	jws.RegisterSigningMethod(SigningMethodSM2)
	signingMethod := jws.GetSigningMethod("SM2")
	j := jws.NewJWT(claims, signingMethod)
	b, err := j.Serialize("")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))

	jwt, err := jws.ParseJWT([]byte(b))
	if err != nil {
		t.Error(err)
	}

	if err := jwt.Validate("", signingMethod); err != nil {
		t.Error(err)
	}

}