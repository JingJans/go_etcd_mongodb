package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		//ctx     context.Context
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "cron/jobs/job1", "go",clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(putResp.Header)
	if putResp.PrevKv!=nil {
		fmt.Println(string(putResp.PrevKv.Value))
	}
	if getResp, err = kv.Get(context.TODO(), "cron/jobs/job1");err!=nil {
		fmt.Println(err)
	}
	fmt.Println(getResp)
}
