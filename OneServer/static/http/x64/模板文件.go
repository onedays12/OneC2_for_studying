package main  
  
import (  
    "encoding/binary"    
    "encoding/json"    
    "fmt")  
  
// CONFIG_MARKER_2024 + 5120 字节空洞  
var placeholder = [5120]byte{  
    'C', 'O', 'N', 'F', 'I', 'G', '_', 'M', 'A', 'R',  
    'K', 'E', 'R', '_', '2', '0', '2', '4',  
}  
  
// 运行时读取并反序列化  
func loadConfig() (map[string]any, error) {  
  
    // 检查  
    if len(placeholder) < 4 {  
       return nil, fmt.Errorf("buffer too small")  
    }  
  
    // 读取json长度  
    lengthBytes := placeholder[:4]  
    length := binary.LittleEndian.Uint32(lengthBytes)  
    if length == 0 {  
       return nil, fmt.Errorf("no config embedded")  
    }  
  
    // 检查  
    end := 4 + int(length)  
    if end > len(placeholder) {  
       return nil, fmt.Errorf("invalid length")  
    }  
  
    // 反序列化json  
    var cfg map[string]any  
    if err := json.Unmarshal(placeholder[4:end], &cfg); err != nil {  
       return nil, err  
    }  
    println(string(placeholder[4:end]))  
    return cfg, nil  
}  
  
func main() {  
    // demo 行为：打印配置  
    cfg, err := loadConfig()  
    if err != nil {  
       fmt.Println("load config failed:", err)  
       return  
    }  
    fmt.Printf("running with cfg=%s\n", cfg)  
}