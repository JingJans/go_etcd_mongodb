package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		//putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		idx     int
		kvpair  *mvccpb.KeyValue
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
	if _, err = kv.Put(context.TODO(), "cron/jobs/job1", "this Test", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	}
	if _, err = kv.Put(context.TODO(), "cron/jobs/job2", "this Test2", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	}
	if _, err = kv.Put(context.TODO(), "cron/jobs/job3", "this Test3", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	}
	if getResp, err = kv.Get(context.TODO(), "cron/jobs", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	}
	result := getResp.Kvs
	fmt.Println(result)

	if delResp, err = kv.Delete(context.TODO(), "cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	}
	if len(delResp.PrevKvs) != 0 {
		for idx, kvpair = range delResp.PrevKvs {
			fmt.Print(string(idx))
			fmt.Print(string(kvpair.Key),string(kvpair.Value))
		}
	}

}
