package listener

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

const (
	defaultProtocol = "http"
	encryptKeyLen   = 16
)

// ---------- 主函数 ----------

func (l *ListenerHTTP) HandlerListenerDataAndStart(listenerName string, conf request.ConfigDetail) (response.ListenerData, *HTTP, error) {

	if conf.RequestHeaders == nil {
		conf.RequestHeaders = make(map[string]string)
	}
	if conf.ResponseHeaders == nil {
		conf.ResponseHeaders = make(map[string]string)
	}

	if conf.HostHeader == "" {
		conf.RequestHeaders["Host"] = conf.HostHeader
	}

	if conf.UserAgent == "" {
		conf.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	}

	if conf.HBHeader == "" {
		conf.HBHeader = "X-Session-ID"
	}

	if conf.HBPrefix == "" {
		conf.HBPrefix = "SESSIONID="
	}
	if conf.EncryptKey == "" {
		key := make([]byte, encryptKeyLen)
		if _, err := rand.Read(key); err != nil {
			return response.ListenerData{}, nil, fmt.Errorf("generate key: %w", err)
		}
		conf.EncryptKey = string(key)
	}

	conf.Protocol = defaultProtocol

	httpListener := &HTTP{
		GinEngine: gin.New(),
		Name:      listenerName,
		Config:    conf,
		Active:    false,
	}
	if err := httpListener.Start(l.ts); err != nil {
		return response.ListenerData{}, nil, fmt.Errorf("start listener: %w", err)
	}

	listenerData := response.ListenerData{
		BindHost:   conf.HostBind,
		BindPort:   strconv.Itoa(conf.PortBind),
		BeaconAddr: strings.Join(conf.CallbackAddresses, ","),
		Status:     "Listen",
		Data:       httpListener.Config,
	}
	if !httpListener.Active {
		listenerData.Status = "Closed"
	}

	return listenerData, httpListener, nil
}
