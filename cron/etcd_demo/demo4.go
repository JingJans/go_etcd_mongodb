package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config       clientv3.Config
		client       *clientv3.Client
		err          error
		lease        clientv3.Lease
		leaseResp    *clientv3.LeaseGrantResponse
		leaseId      clientv3.LeaseID
		putResp      *clientv3.PutResponse
		kv           clientv3.KV
		getResp      *clientv3.GetResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	lease = clientv3.NewLease(client)
	if leaseResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	leaseId = leaseResp.ID

	/*ctx,_ := context.WithTimeout(context.TODO(),5*time.Second)*/

	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("失效")
					goto END
				} else {
					fmt.Println("收到续租", keepResp.ID)
				}
			}
		}
		END:
	}()

	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
	}
	fmt.Println(putResp.Header.Revision)

	for {
		if getResp, err = kv.Get(context.TODO(), "cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv 过期")
			return
		}
		fmt.Println("还没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}
