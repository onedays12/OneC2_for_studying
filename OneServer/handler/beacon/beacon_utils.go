package beacon

import (
	"crypto/rc4"
	"errors"
	"golang.org/x/text/transform"
	"io"
	"net"
	"strings"
)

func ConvertCpToUTF8(input string, codePage int) string {
	enc, exists := codePageMapping[codePage]
	if !exists {
		return input
	}

	reader := transform.NewReader(strings.NewReader(input), enc.NewDecoder())
	utf8Text, err := io.ReadAll(reader)
	if err != nil {
		return input
	}

	return string(utf8Text)
}

func GetOsVersion(majorVersion uint8, minorVersion uint8, buildNumber uint, isServer bool, systemArch string) (int, string) {
	var (
		desc string
		os   = OS_UNKNOWN
	)

	osVersion := "unknown"
	if majorVersion == 10 && minorVersion == 0 && isServer && buildNumber >= 19045 {
		osVersion = "Win 2022 Serv"
	} else if majorVersion == 10 && minorVersion == 0 && isServer && buildNumber >= 17763 {
		osVersion = "Win 2019 Serv"
	} else if majorVersion == 10 && minorVersion == 0 && !isServer && buildNumber >= 22000 {
		osVersion = "Win 11"
	} else if majorVersion == 10 && minorVersion == 0 && isServer {
		osVersion = "Win 2016 Serv"
	} else if majorVersion == 10 && minorVersion == 0 {
		osVersion = "Win 10"
	} else if majorVersion == 6 && minorVersion == 3 && isServer {
		osVersion = "Win Serv 2012 R2"
	} else if majorVersion == 6 && minorVersion == 3 {
		osVersion = "Win 8.1"
	} else if majorVersion == 6 && minorVersion == 2 && isServer {
		osVersion = "Win Serv 2012"
	} else if majorVersion == 6 && minorVersion == 2 {
		osVersion = "Win 8"
	} else if majorVersion == 6 && minorVersion == 1 && isServer {
		osVersion = "Win Serv 2008 R2"
	} else if majorVersion == 6 && minorVersion == 1 {
		osVersion = "Win 7"
	}

	desc = osVersion + " " + systemArch
	if strings.Contains(osVersion, "Win") {
		os = OS_WINDOWS
	}
	return os, desc
}

func int32ToIPv4(ip uint) string {
	bytes := []byte{
		byte(ip),
		byte(ip >> 8),
		byte(ip >> 16),
		byte(ip >> 24),
	}
	return net.IP(bytes).String()
}

func ConvertUTF8toCp(input string, codePage int) string {
	enc, exists := codePageMapping[codePage]
	if !exists {
		return input
	}

	transform.NewWriter(io.Discard, enc.NewEncoder())
	encodedText, err := io.ReadAll(transform.NewReader(strings.NewReader(input), enc.NewEncoder()))
	if err != nil {
		return input
	}

	return string(encodedText)
}

func RC4Crypt(data []byte, key []byte) ([]byte, error) {
	rc4crypt, err := rc4.NewCipher(key)
	if err != nil {
		return nil, errors.New("rc4 crypt error")
	}
	decryptData := make([]byte, len(data))
	rc4crypt.XORKeyStream(decryptData, data)
	return decryptData, nil
}
