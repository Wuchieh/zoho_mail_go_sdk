package zoho_mail_go_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Client struct {
	// 創建時取得
	ClientId string `json:"client_id"`

	// 創建時取得
	ClientSecret string `json:"client_secret"`

	// Scope ZohoMail.accounts.READ,ZohoMail.messages.CREATE
	Code string `json:"code"`

	// 需透過API取得
	//AccountID string `json:"account_id"`

	Account *Account
	Auth    *Auth

	mux sync.Mutex
}

type Account struct {
	Mail string `json:"mail"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (a *Account) GetMail() string {
	if a.Name != "" {
		return fmt.Sprintf("%s <%s>", a.Name, a.Mail)
	}
	return a.Mail
}

type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ApiDomain    string `json:"api_domain"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// New 初始化
//
//	scope: "ZohoMail.accounts.READ,ZohoMail.messages.CREATE"
func New(clientID, clientSecret, code string) (*Client, error) {
	c := &Client{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		Code:         code,
	}

	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) init() error {
	if err := c.initAuth(); err != nil {
		return err
	}

	if err := c.getAccountID(); err != nil {
		return err
	}

	return nil
}

func (c *Client) initAuth() error {
	v := url.Values{
		"client_id":     {c.ClientId},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {c.Code},
	}

	request, err := http.NewRequest("POST", UrlGetToken, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var a *Auth
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}

	c.Auth = a

	return nil
}

func (c *Client) getAccountID() error {
	req, err := http.NewRequest("GET", UrlGetAccountID, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.Auth.AccessToken)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var a accountResp
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}

	if len(a.Data) > 0 {
		var account Account
		account.ID = a.Data[0].AccountId
		account.Name = a.Data[0].Name
		account.Mail = a.Data[0].Email
		c.Account = &account
	} else {
		return errors.New("not get account id")
	}

	return nil
}

func (c *Client) RefreshToken() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	v := url.Values{
		"client_id":     {c.ClientId},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.Auth.RefreshToken},
	}

	req, err := http.NewRequest("POST", UrlGetToken, strings.NewReader(v.Encode()))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	defer req.Body.Close()

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var a *Auth
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}

	c.Auth.AccessToken = a.AccessToken

	return nil
}

func (c *Client) SendMail(m *Mail) error {
	var from string

	if m.FromAddress != "" {
		from = m.FromAddress
	} else {
		from = c.Account.GetMail()
	}

	return c.SendMailC(from, m.ToAddress, m.Subject, m.Content, m.AskReceipt)
}

func (c *Client) SendMailC(from, to, subject, content string, askReceipt bool) error {
	if err := c.RefreshToken(); err != nil {
		return errors.Join(errors.New("refresh token Error"), err)
	}

	m := map[string]string{
		"fromAddress": from,
		"toAddress":   to,
		"subject":     subject,
		"content":     content,
	}

	if askReceipt {
		m["askReceipt"] = "yes"
	}

	mBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(UrlSendMail, c.Account.ID), bytes.NewBuffer(mBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var e emailResp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &e)
	if err != nil {
		return err
	}

	if e.Status.Code != 200 {
		return errors.New(e.Status.Description)
	}

	return nil
}
