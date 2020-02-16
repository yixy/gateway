package handler

import (
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	u, err := url.Parse("https://user:passwd@github.com:443/xxx/yyy/zzz?a==b&&c=d?e=f#123")
	if err != nil {
		t.Error(err)
	}
	t.Log(u.String())
	t.Log(u.Scheme)
	t.Log(u.User.String())
	t.Log(u.Host)
	t.Log(u.RequestURI())
	t.Log(u.Fragment)
}
