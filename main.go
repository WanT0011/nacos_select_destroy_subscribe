package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

const (
	ServerName = `test_service`
)

func main() {
	c := newClient()

	// 正常情况
	normal(c)

	// select破坏情况
	//destroy(c)

	// select不破环情况
	//noDestroyAfterSelect(c)

	<-context.Background().Done()
}

// normal 正常的Subscribe
func normal(c naming_client.INamingClient) {
	randomRegister()

	err := c.Subscribe(&vo.SubscribeParam{
		ServiceName: ServerName,
		Clusters:    []string{"DEFAULT"},
		GroupName:   constant.DEFAULT_GROUP,
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Println("service change:", services)
		},
	})
	if err != nil {
		panic(err)
	}
}

// destroy SelectInstances 破坏了 Subscribe
func destroy(c naming_client.INamingClient) {
	randomRegister()

	err := c.Subscribe(&vo.SubscribeParam{
		ServiceName: ServerName,
		Clusters:    []string{"DEFAULT"},
		GroupName:   constant.DEFAULT_GROUP,
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Println("service change:", services)
		},
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
	instances, err := c.SelectInstances(vo.SelectInstancesParam{
		ServiceName: ServerName,
		HealthyOnly: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("selected instances:", instances)

}

// noDestroyAfterSelect SelectInstances 未破坏 Subscribe
func noDestroyAfterSelect(c naming_client.INamingClient) {
	randomRegister()

	err := c.Subscribe(&vo.SubscribeParam{
		ServiceName: ServerName,
		Clusters:    []string{"DEFAULT"},
		GroupName:   constant.DEFAULT_GROUP,
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Println("service change:", services)
		},
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
	instances, err := c.SelectInstances(vo.SelectInstancesParam{
		Clusters:    []string{"DEFAULT"},
		ServiceName: ServerName,
		HealthyOnly: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("selected instances:", instances)

}

// randomRegister random将注册一个名为ServerName的服务；有一个常驻的实例，一个实例会以5s的频率进行上下线；
func randomRegister() {
	c1 := newClient()
	success, err := c1.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8080,
		Weight:      100,
		Enable:      true,
		Healthy:     true,
		ServiceName: ServerName,
		Ephemeral:   true,
	})
	if err != nil {
		panic(err)
	}
	if !success {
		panic("register instance failed")
	}
	fmt.Println("register instance success")

	c2 := newClient()

	register := false
	go func() {
		for {
			if register {
				register = false
				success, err := c2.DeregisterInstance(vo.DeregisterInstanceParam{
					Ip:          "127.0.0.1",
					Port:        8082,
					ServiceName: ServerName,
					Ephemeral:   true,
				})
				if err != nil {
					panic(err)
				}
				if !success {
					panic("deregister instance failed")
				}
				fmt.Println("deregister 127.0.0.1:8082 instance success")
			} else {
				register = true
				success, err = c2.RegisterInstance(vo.RegisterInstanceParam{
					Ip:          "127.0.0.1",
					Port:        8082,
					Weight:      100,
					Enable:      true,
					Healthy:     true,
					ServiceName: ServerName,
					Ephemeral:   true,
				})
				if err != nil {
					panic(err)
				}
				if !success {
					panic("register instance failed")
				}
				fmt.Println("register 127.0.0.1:8082 instance success")
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func newClient() naming_client.INamingClient {
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NotLoadCacheAtStart: true,
		},
		ServerConfigs: []constant.ServerConfig{
			{
				Scheme:      "http",
				ContextPath: "/nacos",
				IpAddr:      "localhost",
				Port:        8848,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return client
}
