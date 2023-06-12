/**
 * @Author: lenovo
 * @Description:
 * @File:  register
 * @Version: 1.0.0
 * @Date: 2023/06/12 19:04
 */

package discovery

import (
	"context"
	"errors"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// 服务注册的通用接口
type Registrar interface {
	//注册服务
	Register(service Service) error
	//注销服务
	DeRegister() error
}

const leaseTTL = 3

// 基于etcd的服务发现中间件,实现Registrar
type RegistrarEtcd struct {
	//客户端信息
	cli *clientv3.Client
	//租约信息,基于租约来做健康检测
	leaseID  clientv3.LeaseID
	leaseTTL int64
	//续约响应channel
	leaseKeepAliveRespCh <-chan *clientv3.LeaseKeepAliveResponse
}

// 实现接口方法
func (re *RegistrarEtcd) Register(service Service) error {
	//1.租约的申请 lease grant(授权)
	grantResp, err := re.cli.Grant(context.Background(), re.leaseTTL)
	if err != nil {
		return err
	}
	//grantResp.ID 就是租约ID
	re.leaseID = grantResp.ID
	//2.将服务标识: 服务地址 put到etcd中,同时绑定租约
	key := service.Name() + "-" + uuid.New().String()
	_, err = re.cli.Put(context.Background(), key, service.Addr(), clientv3.WithLease(grantResp.ID))
	if err != nil {
		return err
	}
	//3.租约续约
	re.leaseKeepAliveRespCh, err = re.cli.KeepAlive(context.Background(), re.leaseID)
	if err != nil {
		return err
	}
	//从channel中读取响应内容,做处理,需要使用独立的goroutine并发接收
	go re.HandleKeepAliveResp(re.leaseKeepAliveRespCh)
	log.Printf("service %s was registered", service.Name())

	return nil
}

func (re *RegistrarEtcd) HandleKeepAliveResp(ch <-chan *clientv3.LeaseKeepAliveResponse) {
	for rsp := range ch {
		log.Println(rsp.ID)
	}
}

// DeRegister 撤销租约
func (re *RegistrarEtcd) DeRegister() error {
	//lease revoke

	if _, err := re.cli.Revoke(context.Background(), re.leaseID); err != nil {
		return err
	}
	if err := re.cli.Close(); err != nil {
		return err
	}
	return nil
}

// 初始化的方法
func NewRegistrarEtcd(endpoints []string) (*RegistrarEtcd, error) {
	//没有etcd的地址的话,直接返回
	if len(endpoints) == 0 {
		return nil, errors.New("endpoints is empty")
	}
	//连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	//初始化RegistrarEtcd
	//设置必要属性
	re := &RegistrarEtcd{
		cli:      cli,
		leaseTTL: leaseTTL,
	}
	//leaseID是动态获得的,必须先要申请租约,才能有值

	//返回对象
	return re, nil
}
