package fbAccountKitGo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleAccountKit(t *testing.T) {
	const (
		facebookAppID = "<FB-APP-ID>"
		appSecret     = "<FB-APP-ACCOUNTKIT-SECRET"
		code          = "<AUTH-CODE>"
	)

	accountKit := NewAccountKit(facebookAppID, appSecret)
	c2t, err := accountKit.VerifyByCode(code)
	assert.NoError(t, err)
	assert.NotNil(t, c2t.ID)
	assert.NotNil(t, c2t.AccessToken)
	assert.NotNil(t, c2t.TokenRefreshInterval)

	profile, err := accountKit.VerifyByToken(*c2t.AccessToken)
	assert.NoError(t, err)
	assert.True(t, profile.ID != "")
	assert.True(t, profile.Phone.CountryPrefix != "")
	assert.True(t, profile.Phone.NationalNumber != "")
	assert.True(t, profile.Phone.Number != "")
	fmt.Printf("profile:%v\n", profile)
}

func TestUnmarshalProfile(t *testing.T) {
	payload := []byte(`{"id":"3433294898","access_token":"EMAWfh3zcrRDP0BhvA8ZA0djfePsbYWzoE1IZD","token_refresh_interval_sec":2592000}`)

	c2t := CodeToToken{}
	err := json.Unmarshal(payload, &c2t)
	assert.NoError(t, err)
	assert.NotNil(t, c2t.ID)
	assert.NotNil(t, c2t.AccessToken)
	assert.NotNil(t, c2t.TokenRefreshInterval)
	assert.Equal(t, "3433294898", *c2t.ID)
	assert.Equal(t, "EMAWfh3zcrRDP0BhvA8ZA0djfePsbYWzoE1IZD", *c2t.AccessToken)
	assert.Equal(t, 2592000, *c2t.TokenRefreshInterval)
}
