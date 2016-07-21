package influxcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := influxcloudProvider()
	a.Equal(provider.ClientKey, "testkey")
	a.Equal(provider.Secret, "testsecret")
	a.Equal(provider.CallbackURL, "/callback")
}

func TestNewConfigDefaults(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	config := influxcloudProvider().Config
	a.NotNil(config)
	a.Equal("testkey", config.ClientID)
	a.Equal("testsecret", config.ClientSecret)
	a.Equal(authURL, config.Endpoint.AuthURL)
	a.Equal(tokenURL, config.Endpoint.TokenURL)
	a.Equal("/callback", config.RedirectURL)

}

func Test_Implements_Provider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Implements((*goth.Provider)(nil), influxcloudProvider())
}

func Test_BeginAuth(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := influxcloudProvider()
	session, err := provider.BeginAuth("test_state")
	s := session.(*Session)
	a.NoError(err)
	//FIXME: we really need to be able to run this against the acceptance server, too.
	// How should we do this? Maybe a test envvar switch?
	a.Contains(s.AuthURL, "cloud.influxdata.com/oauth/authorize")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", os.Getenv("INFLUXCLOUD_KEY")))
	a.Contains(s.AuthURL, "state=test_state")
	a.Contains(s.AuthURL, "scope=user")
}
func Test_SessionFromJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := influxcloudProvider()

	//FIXME: What is this testing exactly?
	s, err := provider.UnmarshalSession(`{"AuthURL":"http://github.com/auth_url","AccessToken":"1234567890"}`)
	a.NoError(err)
	session := s.(*Session)
	a.Equal(session.AuthURL, "http://github.com/auth_url")
	a.Equal(session.AccessToken, "1234567890")
}

func influxcloudProvider() *Provider {
	return New("testkey", "testsecret", "/callback", "user")
}
