package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	nacosClient config_client.IConfigClient
)

func NaCosInit(ip string, port int) error {
	var err error
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   uint64(port),
		},
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	nacosClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	return err
}
func GetConfig(id string) (string, error) {
	content, err := nacosClient.GetConfig(vo.ConfigParam{
		DataId: id,
		Group:  "DEFAULT_GROUP",
	})
	return content, err
}
