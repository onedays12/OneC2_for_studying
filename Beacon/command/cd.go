package command

import (
	"Beacon/utils/common"
	"Beacon/utils/packet"
	"os"
)

func Cd(packer *packet.Parser, ACP int) ([]byte, error) {

	// 解析目录路径
	path := common.ConvertCpToUTF8(packer.ParseString(), ACP)

	// cd
	err := os.Chdir(string(path))
	if err != nil {
		return nil, err
	}

	// 获取当前路径
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 打包返回当前目录
	arr := []interface{}{dir}
	packed, err := packet.PackArray(arr)
	if err != nil {
		return nil, err
	}

	return packed, nil
}
