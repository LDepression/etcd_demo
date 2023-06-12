/**
 * @Author: lenovo
 * @Description:
 * @File:  discovery
 * @Version: 1.0.0
 * @Date: 2023/06/12 20:57
 */

package discovery

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"time"
)

type Discovery interface {
	GetServiceAddr(serviceName string) (string, error)
	WatchService(serviceName string) error
}
type DiscoveryEtcd struct {
	cli *clientv3.Client
}

func NewDiscoveryEtcd(endpoints []string) (*DiscoveryEtcd, error) {
	if len(endpoints) == 0 {
		return nil, errors.New("endpoints must be specified")
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	de := &DiscoveryEtcd{cli: cli}
	return de, nil
}
func (de *DiscoveryEtcd) GetServiceAddr(serviceName string) (string, error) {
	getResp, err := de.cli.Get(context.Background(), serviceName, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(getResp.Kvs) == 0 {
		return "", errors.New("service not found")
	}
	//load balance
	randIndex := rand.Intn(len(getResp.Kvs)) // [0.n)
	addr := string(getResp.Kvs[randIndex].Value)
	return addr, nil
}
