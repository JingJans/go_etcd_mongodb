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
		config             clientv3.Config
		client             *clientv3.Client
		err                error
		kv                 clientv3.KV
		watcher            clientv3.Watcher
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watchChan          <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
		cancelFun          context.CancelFunc
		ctx                context.Context
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
	go func() {
		for {
			kv.Put(context.TODO(), "cron/jobs/job7", "demojob7")
			kv.Delete(context.TODO(), "cron/jobs/job7")
			time.Sleep(1 * time.Second)
		}
	}()

	if getResp, err = kv.Get(context.TODO(), "cron/jobs/job7", ); err != nil {
		fmt.Println(err)
		return
	}
	if len(getResp.Kvs) != 0 {
		fmt.Println(string(getResp.Kvs[0].Value))
	}
	watchStartRevision = getResp.Header.Revision + 1

	//创建监听

	watcher = clientv3.NewWatcher(client)
	fmt.Println("开始版本", watchStartRevision)

	ctx, cancelFun = context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFun()
	})
	watchChan = watcher.Watch(ctx, "cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	for watchResp = range watchChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改：", string(event.Kv.Value), "Revision:", event.Kv.ModRevision, event.Kv.CreateRevision)
			case mvccpb.DELETE:
				fmt.Println("删除：", event.Kv.ModRevision)
			}
		}
	}

}
