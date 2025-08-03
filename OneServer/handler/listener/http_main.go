package listener

import (
	"OneServer/utils/request"
	"OneServer/utils/response"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	listenerHTTP *ListenerHTTP
	Listeners    []any
)

func NewBeaconHTTPListener(ts any) any {
	listenerHTTP = &ListenerHTTP{
		ts: ts.(TeamServer),
	}
	return listenerHTTP
}

func (l *ListenerHTTP) ListenerValid(conf request.ConfigDetail) error {

	// 2. 必填项校验
	if conf.HostBind == "" {
		return errors.New("host_bind is required")
	}
	if conf.PortBind < 1 || conf.PortBind > 65535 {
		return errors.New("port_bind must be in the range 1-65535")
	}
	if len(conf.CallbackAddresses) == 0 {
		return errors.New("callback_addresses is required")
	}
	if conf.HBHeader == "" {
		return errors.New("parameter_name is required")
	}
	if conf.UserAgent == "" {
		return errors.New("user_agent is required")
	}
	if !strings.Contains(conf.PagePayload, "<<<PAYLOAD_DATA>>>") {
		return errors.New("web_page_output must contain '<<<PAYLOAD_DATA>>>' template")
	}

	// 3. 回调地址校验
	for _, addr := range conf.CallbackAddresses {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}

		// 拆分 host:port
		host, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("invalid callback address (cannot split host:port): %s", addr)
		}

		// 解析端口
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 || port > 65535 {
			return fmt.Errorf("invalid callback port: %s", addr)
		}

		// 解析 IP 或域名
		if ip := net.ParseIP(host); ip == nil {
			// 域名简单校验
			if len(host) == 0 || len(host) > 253 {
				return fmt.Errorf("invalid callback host: %s", addr)
			}
			for _, part := range strings.Split(host, ".") {
				if len(part) == 0 || len(part) > 63 {
					return fmt.Errorf("invalid callback host: %s", addr)
				}
			}
		}
	}

	// 4. URI 校验
	var uriRegexp = regexp.MustCompile(`^/[a-zA-Z0-9\.\=\-]+(/[a-zA-Z0-9\.\=\-]+)*$`)
	if !uriRegexp.MatchString(conf.URI) {
		return errors.New("uri is invalid")
	}

	return nil
}

func (l *ListenerHTTP) ListenerStart(ListenerName string, listenerConfig request.ConfigDetail) (response.ListenerData, error) {

	listenerData, listener, err := l.HandlerListenerDataAndStart(ListenerName, listenerConfig)
	if err != nil {
		return listenerData, err
	}

	Listeners = append(Listeners, listener)

	return listenerData, nil

}
