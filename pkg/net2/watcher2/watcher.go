// Copyright (c) 2018-present,  NebulaChat Studio (https://nebula.chat).
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: Benqi (wubenqi@gmail.com)

package watcher2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/libregram/chatengine/pkg/net2"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// see: /baselib/grpc_util/service_discovery/registry.go
type nodeData struct {
	Addr     string
	Metadata map[string]string
}

// TODO(@benqi): grpc_util/service_discovery集成
type ClientWatcher struct {
	etcCli      *clientv3.Client
	registryDir string
	serviceName string
	// rootPath    string
	client *net2.TcpClientGroupManager
	nodes  map[string]*nodeData
}

func NewClientWatcher(registryDir, serviceName string, cfg clientv3.Config, client *net2.TcpClientGroupManager) (watcher *ClientWatcher, err error) {
	glog.Info("NewClientWatcher enter")
	var etcdClient *clientv3.Client
	if etcdClient, err = clientv3.New(cfg); err != nil {
		glog.Error("Error: cannot connect to etcd:", err)
		glog.Info("Error: cannot connect to etcd:", err)
		return
	}

	watcher = &ClientWatcher{
		etcCli:      etcdClient,
		registryDir: registryDir,
		serviceName: serviceName,
		client:      client,
		nodes:       map[string]*nodeData{},
	}
	glog.Info("NewClientWatcher success")
	return
}

func (m *ClientWatcher) WatchClients(cb func(etype, addr string)) {
	glog.Info("WatchClients enter")
	if m == nil {
		glog.Info("Err m==nil")
		return
	}

	rootPath := fmt.Sprintf("%s/%s", m.registryDir, m.serviceName)

	glog.Info("calling m.etcCli.Get: m.etcCli ", m.etcCli, " context.Background() ", context.Background(), " rootPath ", rootPath, " WithPrefix")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	resp, err := m.etcCli.Get(ctx, rootPath, clientv3.WithPrefix())
	glog.Info("called m.etcCli.Get: resp ", resp, " err ", err)
	cancel()
	if err != nil {
		glog.Info("Error", err)
		glog.Error(err)
		return
	}
	glog.Info("for1 {")
	for _, kv := range resp.Kvs {
		glog.Info("m.addClient")
		m.addClient(kv, cb)
	}
	glog.Info("for1 }")

	glog.Info("calling m.etcCli.Watch")
	rch := m.etcCli.Watch(context.Background(), rootPath, clientv3.WithPrefix())
	glog.Info("called m.etcCli.Watch: rch ", rch)
	glog.Info("for2 {")
	for wresp := range rch {
		glog.Info("for3 {")
		for _, ev := range wresp.Events {
			glog.Info("event ", ev.Type.String())
			if ev.Type.String() == "EXPIRE" {
				// TODO(@benqi): 采用何种策略？？
				// n, ok := m.nodes[string(ev.Kv.Key)]
				// if ok {
				//	 delete(m.nodes, string(ev.Kv.Key))
				// }
				// if cb != nil {
				// 	cb("EXPIRE", string(ev.Kv.Key), string(ev.Kv.Value))
				//}
			} else if ev.Type.String() == "PUT" {
				m.addClient(ev.Kv, cb)
			} else if ev.Type.String() == "DELETE" {
				if n, ok := m.nodes[string(ev.Kv.Key)]; ok {
					m.client.RemoveClient(m.serviceName, n.Addr)
					if cb != nil {
						cb("delete", n.Addr)
					}
					delete(m.nodes, string(ev.Kv.Key))
				}
			}
		}
		glog.Info("for3 }")
	}
	glog.Info("for2 }")
}
func (m *ClientWatcher) addClient(kv *mvccpb.KeyValue, cb func(etype, addr string)) {
	glog.Info("addClient enter")
	node := &nodeData{}
	err := json.Unmarshal(kv.Value, node)
	if err != nil {
		glog.Info("Err", err)
		glog.Error(err)
	}
	if n, ok := m.nodes[string(kv.Key)]; ok {
		if node.Addr != n.Addr {
			m.client.RemoveClient(m.serviceName, n.Addr)
			m.nodes[string(kv.Key)] = node
			if cb != nil {
				cb("delete", n.Addr)
			}
			if cb != nil {
				cb("add", node.Addr)
			}
		}
		glog.Info("m.client.AddClient svc ", m.serviceName, " node ", node.Addr)
		m.client.AddClient(m.serviceName, node.Addr)
	} else {
		m.nodes[string(kv.Key)] = node
		glog.Info("m.client.AddClient svc ", m.serviceName, " node ", node.Addr)
		m.client.AddClient(m.serviceName, node.Addr)
		if cb != nil {
			cb("add", node.Addr)
		}
	}
}
