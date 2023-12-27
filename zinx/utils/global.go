package utils

import (
	"encoding/json"
	"gameTcp/zinx/iface"
	"os"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type GlobalConfig struct {
	//server
	TcpServer iface.Server //全局Server对象
	Host      string       //当前服务器主机IP
	TcpPort   int          //当前服务器主机监听端口号
	Name      string       //当前服务器名称

	//zInx
	Version          string //版本号
	IPVersion        string //
	MaxPacketSize    uint32 //数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
}

/*
定义一个全局的对象
*/
var Global *GlobalConfig

/*
提供init方法，默认加载
*/
func init() {
	data, err := os.ReadFile("../../zinx/conf/config.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &Global)
	if err != nil {
		panic(err)
	}

}
