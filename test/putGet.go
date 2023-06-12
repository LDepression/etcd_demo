/**
 * @Author: lenovo
 * @Description:
 * @File:  putGet
 * @Version: 1.0.0
 * @Date: 2023/06/12 17:59
 */

package test

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func PutGet() {
	//1.连接etcd
	//实例化etcd客户端
	//给定连接信息
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"http://localhost:2379"},
		//连接时长
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	//put测试
	putResp, err := cli.Put(context.Background(), "testKey", "testValue")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(putResp)

	//get测试
	getResp, err := cli.Get(context.Background(), "testKey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(getResp.Kvs[0].Key), string(getResp.Kvs[0].Value))
}
