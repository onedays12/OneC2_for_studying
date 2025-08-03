package sysinfo

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var Kernel32 = syscall.NewLazyDLL("Kernel32.dll")

// 获取本机第一个非回环IPv4地址
func GetInternalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

// 获取域名（Windows下为WORKGROUP或域名，Linux下为主机名）
func GetDomain() string {
	if runtime.GOOS == "windows" {
		// Windows下尝试获取USERDOMAIN
		return os.Getenv("USERDOMAIN")
	}
	// Linux下用主机名
	name, _ := os.Hostname()
	return name
}

func GetCodePageANSI() (int32, error) {
	fnGetACP := Kernel32.NewProc("GetACP")
	if fnGetACP.Find() != nil {
		return 0, errors.New("not found GetACP")
	}
	acp, _, _ := fnGetACP.Call()
	return int32(acp), nil
}

// 获取计算机名
func GetComputerName() string {
	name, _ := os.Hostname()
	return name
}

// 获取用户名
func GetUsername() string {
	if u, err := user.Current(); err == nil {
		// Windows下user.Current().Username可能是"DOMAIN\\User"
		parts := strings.Split(u.Username, "\\")
		return parts[len(parts)-1]
	}
	return ""
}

// 获取当前进程名
func GetProcessName() string {
	return filepath.Base(os.Args[0])
}

// 获取当前进程PID
func GetPid() int16 {
	return int16(os.Getpid())
}

func GetTID() int16 {
	GetCurrentThreadId := Kernel32.NewProc("GetCurrentThreadId")
	tid, _, _ := GetCurrentThreadId.Call()
	return int16(tid)
}

func GetCodePageOEM() (int32, error) {
	procGetOEMCP := Kernel32.NewProc("GetOEMCP")

	if err := procGetOEMCP.Find(); err != nil {
		return 0, errors.New("not found GetOEMCP")
	}

	oemacp, _, _ := procGetOEMCP.Call()
	return int32(oemacp), nil
}

func GetGmtOffset() int32 {
	_, offset := time.Now().Zone()
	// offset 是相对于 UTC 的秒数
	// Windows Bias 是分钟，且正值代表西区（负时区），Go 的 offset 正值代表东区（正时区）
	// 所以这里直接用 offset/3600 得到小时数
	return int32(offset / 3600)
}

type OSVERSIONINFOEXW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        byte
	wReserved           byte
}

func getOSVersionInfo() (*OSVERSIONINFOEXW, error) {
	ntdll := syscall.NewLazyDLL("ntdll.dll")
	rtlGetVersion := ntdll.NewProc("RtlGetVersion")

	var osvi OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))

	ret, _, _ := rtlGetVersion.Call(uintptr(unsafe.Pointer(&osvi)))
	if ret != 0 {
		return nil, fmt.Errorf("RtlGetVersion failed: %d", ret)
	}
	return &osvi, nil
}

func GetWindowsMajorVersion() (int8, error) {
	osvi, err := getOSVersionInfo()
	if err != nil {
		return 0, err
	}
	return int8(osvi.dwMajorVersion), nil
}

func GetWindowsMinorVersion() (int8, error) {
	osvi, err := getOSVersionInfo()
	if err != nil {
		return 0, err
	}
	return int8(osvi.dwMinorVersion), nil
}

func GetWindowsBuildNumber() (int32, error) {
	osvi, err := getOSVersionInfo()
	if err != nil {
		return 0, err
	}
	return int32(osvi.dwBuildNumber), nil
}
