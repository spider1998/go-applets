package applets

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewWechat(appID, secret, templateID string) Wechat {
	return Wechat{
		appID:      appID,
		secret:     secret,
		templateID: templateID,
	}
}

type Wechat struct {
	appID       string
	secret      string
	templateID  string
	accessToken *AccessToken
	sync.Mutex
}

type OpenIDResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	WxErr
}

type CheckTokenRequest struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Echostr   string `json:"echostr"`
}

type WxErr struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	WxErr
}

type SendRequest struct {
	Touser          string  `json:"touser"`
	FormID          string  `json:"form_id"`
	Page            string  `json:"page"`
	Data            Message `json:"data"`
	EmphasisKeyword string  `json:"emphasis_keyword"`
}

type Message map[string]interface{}

func (w *Wechat) SendMsg(req SendRequest) (err error) {
	token, err := w.GetAccessToken()
	if err != nil {
		return
	}
	api, err := TokenAPI("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send", token)
	if err != nil {
		return
	}
	for key := range req.Data {
		req.Data[key] = Message{"value": req.Data[key]}
	}
	body := map[string]interface{}{
		"touser":           req.Touser,
		"template_id":      w.templateID,
		"page":             req.Page,
		"form_id":          req.FormID,
		"data":             req.Data,
		"emphasis_keyword": req.EmphasisKeyword,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err := http.Post(api, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = errors.New("WECHAT_SERVER_ERROR")
		return err
	}

	var resp WxErr
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Errcode != 0 {
		return errors.New(resp.Errmsg)
	}
	return nil

}

func (w *Wechat) GetAccessToken() (token string, err error) {
	w.Lock()
	defer w.Unlock()
	if w.accessToken == nil || w.accessToken.ExpiresIn < time.Now().Unix() {
		for i := 0; i < 3; i++ {
			err = w.getAccessToken()
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			return
		}
	}
	token = w.accessToken.AccessToken
	return
}

func (w *Wechat) CheckSignature(req CheckTokenRequest) (err error) {
	if sig := w.sortSha1(req.Timestamp, req.Nonce, req.Echostr); sig != req.Signature {
		err = errors.New("check signature failed.")
		return
	}
	return
}

func (w *Wechat) GetOpenID(authCode string) (openID string, err error) {
	urlStr := "https://api.weixin.qq.com/sns/jscode2session?appid=" + w.appID + "&secret=" + w.secret + "&js_code=" + authCode + "&grant_type=authorization_code"
	logrus.WithField("url", urlStr).Info("get open id.")
	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("wechat internal server error.")
		return
	}

	var result OpenIDResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	if result.Errcode != 0 {
		err = errors.New(result.Errmsg)
		return
	}
	openID = result.Openid
	return
}

func TokenAPI(api, token string) (string, error) {
	u, err := url.Parse(api)
	if err != nil {
		return "", err
	}
	query := u.Query()
	query.Set("access_token", token)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (w *Wechat) getAccessToken() (err error) {
	urls, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		return
	}
	query := urls.Query()
	query.Set("appid", w.appID)
	query.Set("secret", w.secret)
	query.Set("grant_type", "client_credential")

	urls.RawQuery = query.Encode()

	res, err := http.Get(urls.String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("wechat internal server error.")
	}

	var token AccessToken
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return
	}

	if token.Errcode != 0 {
		return errors.New(token.Errmsg)
	}
	w.accessToken.AccessToken = token.AccessToken
	w.accessToken.ExpiresIn = token.ExpiresIn
	return
}

func (w *Wechat) sortSha1(s ...string) string {
	sort.Strings(s)
	h := sha1.New()
	h.Write([]byte(strings.Join(s, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}
