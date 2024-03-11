package grpc

import (
	"fmt"
	"github.com/As1433223/framework_one/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"log"
	"net"
)

type GrpcConf struct {
	Ip   string `yaml:"Ip"`
	Port int    `yaml:"Port"`
}

func getConfig(servername string) (GrpcConf, error) {
	conf, err := config.GetConfig(servername)
	if err != nil {
		return GrpcConf{}, err
	}
	var GrpcConfig GrpcConf
	var data map[interface{}]interface{}
	yaml.Unmarshal([]byte(conf), &data)
	GrpcConfig.Ip = data["GrpcConf"].(map[interface{}]interface{})["Ip"].(string)
	GrpcConfig.Port = data["GrpcConf"].(map[interface{}]interface{})["Port"].(int)
	return GrpcConfig, nil
}
func RegisterGrpc(servername string, f func(server *grpc.Server)) error {
	conf, err := getConfig(servername)
	if err != nil {
		return err
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Println("Listening failed:", conf.Port)
		return err
	}
	err = config.RegisterConsul(servername, conf.Port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	reflection.Register(s)
	f(s)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	log.Println()
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve:%v", err)
		return err
	}
	return nil

}
