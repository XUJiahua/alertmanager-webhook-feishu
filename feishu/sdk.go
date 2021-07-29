package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Sdk struct {
	appID     string
	appSecret string
	token     string
	client    http.Client
}

func NewSDK(appID string, appSecret string) *Sdk {
	s := &Sdk{
		appID:     appID,
		appSecret: appSecret,
		client:    http.Client{},
	}

	if appID != "" && appSecret != "" {
		s.refreshToken()
	}

	return s
}

type batchGetIDResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		EmailUsers map[string][]struct {
			OpenId string `json:"open_id"`
			UserId string `json:"user_id"`
		} `json:"email_users"`
		EmailsNotExist []string `json:"emails_not_exist"`
		MobileUsers    map[string][]struct {
			OpenId string `json:"open_id"`
			UserId string `json:"user_id"`
		} `json:"mobile_users"`
		MobilesNotExist []string `json:"mobiles_not_exist"`
	} `json:"data"`
}

// BatchGetID https://open.feishu.cn/document/ukTMukTMukTM/uUzMyUjL1MjM14SNzITN
func (s Sdk) BatchGetID(emails []string) (map[string]string, error) {
	if len(emails) == 0 {
		return nil, errors.New("at least 1 email")
	}
	if len(emails) > 50 {
		return nil, errors.New("at most 50 emails")
	}

	api := "https://open.feishu.cn/open-apis/user/v1/batch_get_id?"
	for _, email := range emails {
		api += "emails=" + email + "&"
	}
	api = api[:len(api)-1]
	var response batchGetIDResponse
	err := s.get(api, s.token, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, errors.New(fmt.Sprintf("code: %d, err: %s", response.Code, response.Msg))
	}

	res := make(map[string]string)
	for k, vv := range response.Data.EmailUsers {
		for _, v := range vv {
			res[k] = v.OpenId
		}
	}
	return res, nil
}

type tokenRequest struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type tokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

// TenantAccessToken https://open.feishu.cn/document/ukTMukTMukTM/uIjNz4iM2MjLyYzM
func (s Sdk) TenantAccessToken() (*tokenResponse, error) {
	request := tokenRequest{
		AppID:     s.appID,
		AppSecret: s.appSecret,
	}
	var response tokenResponse
	err := s.post("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/", s.token, request, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, errors.New(fmt.Sprintf("code: %d, err: %s", response.Code, response.Msg))
	}

	return &response, nil
}

type webhookV2Response struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
}

func (s Sdk) WebhookV2(webhook string, body io.Reader) error {
	req, err := http.NewRequest("POST", webhook, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	do, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()

	var resp webhookV2Response
	err = json.NewDecoder(do.Body).Decode(&resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != 0 {
		return errors.New(fmt.Sprintf("code: %d, err: %s", resp.StatusCode, resp.StatusMessage))
	}

	return nil
}

func (s Sdk) get(url string, auth string, responseBody interface{}) error {
	return s.call("GET", url, auth, nil, responseBody)
}

func (s Sdk) post(url string, auth string, requestBody, responseBody interface{}) error {
	return s.call("POST", url, auth, requestBody, responseBody)
}

func (s Sdk) call(method string, url string, auth string, requestBody, responseBody interface{}) error {
	logrus.Debugf("%s %s with %v", method, url, requestBody)
	var body io.Reader
	if requestBody != nil {
		bs, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(bs)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if auth != "" {
		req.Header.Set("Authorization", "Bearer "+auth)
	}

	do, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()

	err = json.NewDecoder(do.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	logrus.Debug(responseBody)

	return nil
}

func (s *Sdk) refreshToken() {
	response, err := s.TenantAccessToken()
	if err != nil {
		logrus.Errorf("refresh token failed, %v", err)

		// sleep and try again
		time.Sleep(time.Second * 1)
		s.refreshToken()
		return
	}
	s.token = response.TenantAccessToken

	// https://open.feishu.cn/document/ukTMukTMukTM/uIjNz4iM2MjLyYzM
	// Token 有效期为 2 小时，在此期间调用该接口 token 不会改变。当 token 有效期小于 30 分的时候，再次请求获取 token 的时候，会生成一个新的 token，与此同时老的 token 依然有效。
	// 在过期前 1 分钟刷新
	time.AfterFunc(time.Second*time.Duration(response.Expire-60), func() {
		s.refreshToken()
	})
}
