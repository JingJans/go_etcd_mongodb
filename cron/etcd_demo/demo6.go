package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		opResp clientv3.OpResponse
		getOp  clientv3.Op

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
	putOp = clientv3.OpPut("cron/jobs/job1", "11223")
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(opResp.Put().Header.Revision)

	getOp = clientv3.OpGet("cron/jobs/job1")
	if opResp,err = kv.Do(context.TODO(),getOp);err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(opResp.Get().Kvs[0].ModRevision)
	fmt.Println(string(opResp.Get().Kvs[0].Value))

}
