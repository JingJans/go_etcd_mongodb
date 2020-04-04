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
		err          error
		client       *clientv3.Client
		lease        clientv3.Lease
		leaseResp    *clientv3.LeaseGrantResponse
		leaseId      clientv3.LeaseID
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
		ctx          context.Context
		cancelFun    context.CancelFunc
		kv           clientv3.KV
		txn          clientv3.Txn
		txnResp      *clientv3.TxnResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//上锁
	lease = clientv3.NewLease(client)
	if leaseResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}
	leaseId = leaseResp.ID
	ctx, cancelFun = context.WithCancel(context.TODO())
	defer cancelFun()
	defer lease.Revoke(context.TODO(), leaseId)

	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约失效")
					goto END
				} else {
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()

	//抢锁
	kv = clientv3.NewKV(client)
	//创建事务
	txn = kv.Txn(context.TODO())
	//处理业务

	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9"))
	//
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}
	//判断是否抢到锁
	if !txnResp.Succeeded {
		fmt.Println("抢锁失败")
		fmt.Println(string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)
	//释放锁
	//defer 已经释放掉
}
