package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"strconv"
)

type httpClient struct{}

var HttpClient httpClient
var restyCli = resty.New()
var key, _ = hex.DecodeString("57227176a09c27191875e85ce2ccea571e415fd98038ccb21e892c4d7182bc3e")
var iv, _ = hex.DecodeString("ff5097cd1d355f6d6f8d9225")

func (cli *httpClient) SetSocks5(b bool, address string, port string) error {
	if b {
		proxyURL, _ := url.Parse(fmt.Sprintf("socks5://%s:%s", address, port))
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return errors.New("socks5 proxy init error: " + err.Error())
		}
		transport := &http.Transport{Dial: dialer.Dial}
		restyCli.SetTransport(transport)
	} else {
		restyCli.RemoveProxy()
	}
	return nil
}

func (cli *httpClient) Post(toUrl string, params map[string]interface{}, result any) error {
	restyCli.Debug = true

	request := restyCli.R()
	request.SetHeader("Content-Type", "text/plain; charset=UTF-8")
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return err
	}
	//LogInfo(context.Background(), "序列化："+string(jsonStr))
	encrypt, err := AesGcmEncrypt(jsonStr, key, iv)
	if err != nil {
		return errors.New("EncError")
	}
	toString := base64.StdEncoding.EncodeToString(encrypt)

	response, err := request.SetBody(toString).Post(toUrl)
	if err != nil {
		return errors.New("请求错误：请检查网络连接或代理设置")
	}
	if response.StatusCode() != http.StatusOK {
		return errors.New("http status code:" + strconv.Itoa(response.StatusCode()))
	}
	body := response.Body()
	decodeBytes, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return errors.New("DecError")
	}
	plaintext, err := AesGcmDecrypt(decodeBytes, key, iv)
	if err != nil {
		return errors.New("DecError")
	}
	err = json.Unmarshal(plaintext, result)
	if err != nil {
		return errors.New("json Unmarshal err:" + err.Error())
	}
	return nil
}
