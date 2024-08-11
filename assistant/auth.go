package assistant

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/go-resty/resty/v2"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var restyCli = resty.New()
var key, _ = hex.DecodeString("57227176a09c27191875e85ce2ccea571e415fd98038ccb21e892c4d7182bc3e")
var iv, _ = hex.DecodeString("ff5097cd1d355f6d6f8d9225")
var config = dao.Config{}

type AuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UUID      string `json:"uuid"`
		Exp       int    `json:"exp"`
		Duration  int    `json:"duration"`
		Timestamp int    `json:"timestamp"`
	} `json:"data"`
}

func Auth() (*AuthResponse, error) {
	var uid string
	if fyne.CurrentDevice().IsMobile() {
		id, err := utils.GetAndroidID()
		if err != nil {
			return nil, errors.New("获取Android设备信息失败，请联系客服处理" + err.Error())
		}
		uid = id
	} else {
		address, err := getMACAddress()
		if err != nil {
			return nil, errors.New("获取设备信息失败，请联系客服处理" + err.Error())
		}
		uid = address
	}

	params := make(map[string]interface{})
	params["device_id"] = uid
	params["uuid"] = config.Get("authCode")
	params["timestamp"] = time.Now().Unix()
	result := &AuthResponse{}
	if config.Get("socksOpen") == "true" {
		err := setSocks5(true, config.Get("socksAddr"), config.Get("socksPort"))
		if err != nil {
			return nil, err
		}
	} else {
		_ = setSocks5(false, "", "")
	}
	err := post("https://auth.seven-d76.workers.dev/acv", params, result)
	if err != nil {
		return nil, errors.New("获取设备信息失败，请联系客服处理" + err.Error())
	}
	return result, nil
}

func setSocks5(b bool, address string, port string) error {
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

func post(toUrl string, params map[string]interface{}, result any) error {
	restyCli.Debug = true

	request := restyCli.R()
	request.SetHeader("Content-Type", "text/plain; charset=UTF-8")
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return err
	}
	//LogInfo(context.Background(), "序列化："+string(jsonStr))
	encrypt, err := utils.AesGcmEncrypt(jsonStr, key, iv)
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
	plaintext, err := utils.AesGcmDecrypt(decodeBytes, key, iv)
	if err != nil {
		return errors.New("DecError")
	}
	err = json.Unmarshal(plaintext, result)
	if err != nil {
		return errors.New("json Unmarshal err:" + err.Error())
	}
	return nil
}
func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, inter := range interfaces {
		if inter.Flags&net.FlagUp != 0 && len(inter.HardwareAddr) > 0 {
			return inter.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("no valid MAC address found")
}
