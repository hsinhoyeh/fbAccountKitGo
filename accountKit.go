package fbAccountKitGo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Profile is a user profile returned from facebook's account Kit
type Profile struct {
	ID    string `json:"id"`
	Phone struct {
		Number         string `json:"number"`
		CountryPrefix  string `json:"country_prefix"`
		NationalNumber string `json:"national_number"`
	} `json:"phone"`
}

// CodeToToken is a structure for capturing the response when we send auth code to fbaccount kit.
type CodeToToken struct {
	ID                   *string `json:"id"`
	AccessToken          *string `json:"access_token"`
	TokenRefreshInterval *int    `json:"token_refresh_interval_sec"`
}

type AccountKit struct {
	AppID     string
	AppSecret string
	Debug     bool

	client *http.Client
}

func NewAccountKit(appID, appSecret string) *AccountKit {
	return &AccountKit{
		AppID:     appID,
		AppSecret: appSecret,
		Debug:     true,
	}
}

func (a *AccountKit) VerifyByToken(token string) (*Profile, error) {
	queryParams := url.Values{}
	queryParams.Add("access_token", token)

	// hash hmac sha256 token with appSecret
	appsecret_proof := computeHmac256(token, a.AppSecret)
	queryParams.Add("appsecret_proof", appsecret_proof)

	request, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://graph.accountkit.com/v1.0/me/?%s", queryParams.Encode()),
		nil)
	if err != nil {
		return nil, err
	}
	response, err := a.do(request)
	if err != nil {
		return nil, err
	}

	profile := Profile{}
	err = a.readAndClose(response, &profile)
	return &profile, err
}

func (a *AccountKit) VerifyByCode(code string) (*CodeToToken, error) {

	queryParams := url.Values{}
	queryParams.Add("grant_type", "authorization_code")
	queryParams.Add("code", code)
	queryParams.Add("access_token", fmt.Sprintf("AA|%s|%s", a.AppID, a.AppSecret))
	request, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://graph.accountkit.com/v1.0/access_token?%s", queryParams.Encode()),
		nil)
	if err != nil {
		return nil, err
	}
	response, err := a.do(request)
	if err != nil {
		return nil, err
	}

	cot := &CodeToToken{}
	err = a.readAndClose(response, cot)
	return cot, err
}

func (a *AccountKit) readAndClose(resp *http.Response, picker interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, picker); err != nil {
		return err
	}
	return nil
}

func (a *AccountKit) do(req *http.Request) (*http.Response, error) {
	client := a.client
	if client == nil {
		client = http.DefaultClient
	}

	if a.Debug {
		rr, _ := httputil.DumpRequest(req, true)
		fmt.Printf("<----:%v\n", string(rr))
	}
	resp, err := client.Do(req)

	if a.Debug {
		rr, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("---->:%v\n", string(rr))
	}

	return resp, err
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
