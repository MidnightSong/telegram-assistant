package test

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_http(t *testing.T) {
	var httpClient = resty.New()
	httpClient.Debug = true
	request := httpClient.R()
	response, err := request.Get("https://www.google.com")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 200, response.StatusCode)
	t.Log(response.String())
}
func TestEnc(t *testing.T) {
	key := "57227176a09c27191875e85ce2ccea571e415fd98038ccb21e892c4d7182bc3e"
	iv := "ff5097cd1d355f6d6f8d9225"
	str := `{"device_id":"test_device","timestamp":1721808181,"uuid":"test_uuid"}`
	decodeKey, _ := hex.DecodeString(key)
	decodeIv, _ := hex.DecodeString(iv)

	encrypt, err := utils.AesGcmEncrypt([]byte(str), decodeKey, decodeIv)

	if err != nil {
		t.Error(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(encrypt))
}

// https://auth.seven-d76.workers.dev/acv
func TestDec(t *testing.T) {
	key := "57227176a09c27191875e85ce2ccea571e415fd98038ccb21e892c4d7182bc3e"
	iv := "ff5097cd1d355f6d6f8d9225"
	decodeKey, _ := hex.DecodeString(key)
	decodeIv, _ := hex.DecodeString(iv)
	str := "0DoLgJMAMV6uacvTUs2YlSPcd1y+ayueIPMWt6mHkpFxVLYC/qmdv3zurLr2oUN38O+fs2RODpVGCV4="
	decodeString, _ := base64.StdEncoding.DecodeString(str)

	plaintext, err := utils.AesGcmDecrypt(decodeString, decodeKey, decodeIv)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(plaintext))
}

func TestGetMachineInfo(t *testing.T) {
	name, err := os.Hostname()
	if err != nil {
		t.Error(err)
	}
	t.Log(name)
}
func TestHttpClient(t *testing.T) {

}
