/**
 * @Author: lenovo
 * @Description:
 * @File:  services
 * @Version: 1.0.0
 * @Date: 2023/06/12 19:43
 */

package discovery

import (
	"log"
	"os"
	"os/signal"
)

type Service interface {
	Name() string
	Addr() string
}

// orderService
type OrderService struct {
	name string
	addr string
}

func (o *OrderService) Name() string {
	return o.name
}
func (o *OrderService) Addr() string {
	return o.addr
}

func ServiceOrder(addr string) {
	//初始化orderService
	orderService := &OrderService{
		name: "order",
		addr: addr,
	}
	//获取RegistrarEtcd
	re, err := NewRegistrarEtcd([]string{"localhost:2379"})
	if err != nil {
		log.Fatalln(err)
	}
	//将初始化orderService注册到RegistrarEtcd中去
	if err := re.Register(orderService); err != nil {

	}
	//阻塞执行
	log.Printf("service %s is running", orderService.Name())
	chInt := make(chan os.Signal, 1)
	signal.Notify(chInt, os.Interrupt) //监控终止ctrl + c
	select {
	case <-chInt:
		if err := re.DeRegister(); err != nil {
			log.Fatalf("service %s(%s) was deregistered", orderService.Name(), orderService.addr)
		}
	}
}

func ServiceDriver() {
	de, err := NewDiscoveryEtcd([]string{"localhost:2379"})
	if err != nil {
		log.Fatalln(err)
	}
	serviceName := "order"
	addr, err := de.GetServiceAddr(serviceName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("service: %s was discoveried on %s\n", serviceName, addr)

	//然后连接目标服务,使用grpc
}
