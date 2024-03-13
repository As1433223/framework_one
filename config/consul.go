package config

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
	"net"
)

type ConsulConf struct {
	Ip   string `yaml:"Ip"`
	Port int    `yaml:"Port"`
}

func getConfig(servername string) (ConsulConf, error) {
	conf, err := GetConfig(servername)
	if err != nil {
		return ConsulConf{}, nil
	}
	var ConsulConfig ConsulConf
	var nacosConf map[interface{}]interface{}
	err = yaml.Unmarshal([]byte(conf), &nacosConf)
	ConsulConfig.Ip = nacosConf["ConsulConf"].(map[interface{}]interface{})["Ip"].(string)
	ConsulConfig.Port = nacosConf["ConsulConf"].(map[interface{}]interface{})["Port"].(int)
	if err != nil {
		return ConsulConf{}, nil
	}
	return ConsulConfig, nil
}
func GetIp() (ip []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, addr := range addrs {
		ipNet, isVailIpNet := addr.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = append(ip, ipNet.IP.String())
			}
		}

	}
	return ip
}
func RegisterConsul(servername string, port int, Name string) error {
	conf, err := getConfig(servername)
	if err != nil {
		return err
	}
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", conf.Ip, conf.Port),
	})
	if err != nil {
		return err
	}
	ip := GetIp()
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    Name,
		Tags:    []string{"Grpc"},
		Port:    port,
		Address: ip[0],
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			GRPC:                           fmt.Sprintf("%s:%d", ip[0], port),
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		return err
	}
	return nil
}
func AgentHealthService(servername string) (string, error) {
	conf, err := getConfig(servername)
	if err != nil {
		return "", err
	}
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", conf.Ip, conf.Port),
	})
	if err != nil {
		return "", err
	}
	name, i, err := client.Agent().AgentHealthServiceByName(servername)
	if err != nil {
		return "", err
	}
	if name != "passing" {
		return "", errors.New("is not health service")
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", i[0].Service.Address, i[0].Service.Port), nil
}
