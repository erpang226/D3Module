package socket

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"k8s.io/klog/v2"

	"main/pkg/actor/common"
	"main/pkg/actor/core/model"
	"main/pkg/actor/core/socket/broker"
	"main/pkg/actor/core/socket/config"
	"main/pkg/actor/core/socket/store"
	"main/pkg/actor/core/socket/wrapper"
)

type Context struct {
	contexts map[string]*context
	groups   map[string]*context
	sync.RWMutex
}

var globalSocketContext Context
var once = sync.Once{}

func InitSocketContext() *Context {
	once.Do(func() {
		globalSocketContext.contexts = make(map[string]*context)
		globalSocketContext.groups = make(map[string]*context)
	})
	return &globalSocketContext
}

func (s *Context) AddModule(info *common.ModuleInfo) {
	name := info.ModuleName
	s.setContext(name)
	context := s.getContext(name)
	if !info.IsRemote {
		context.AddModule(name, info.Connection)
	} else {
		context.AddModuleRemote(name)
	}
}

func (s *Context) AddModuleGroup(module, group string) {
	s.Lock()
	s.groups[module] = s.contexts[module]
	s.Unlock()

	s.getContext(module).AddModuleGroup(module, group)
}

func (s *Context) Cleanup(module string) {
	s.getContext(module).Cleanup(module)
}

func (s *Context) Send(module string, message model.Message) {
	s.getContext(module).Send(module, message)
}

func (s *Context) Receive(module string) (model.Message, error) {
	return s.getContext(module).Receive(module)
}

func (s *Context) SendSync(module string, message model.Message, timeout time.Duration) (model.Message, error) {
	return s.getContext(module).SendSync(module, message, timeout)
}

func (s *Context) SendResp(message model.Message) {
	module := message.GetSource()
	s.getContext(module).SendResp(message)
}

func (s *Context) SendToGroup(group string, message model.Message) {
	s.getGroupContext(group).SendToGroup(group, message)
}

func (s *Context) SendToGroupSync(module string, message model.Message, timeout time.Duration) error {
	return s.getContext(module).SendToGroupSync(module, message, timeout)
}

func (s *Context) getGroupContext(group string) *context {
	s.RLock()
	defer s.RUnlock()
	return s.groups[group]
}

func (s *Context) getContext(module string) *context {
	s.RLock()
	defer s.RUnlock()
	return s.contexts[module]
}

func (s *Context) setContext(module string) {
	s.Lock()
	defer s.Unlock()
	s.contexts[module] = newContext(module)
}

type context struct {
	name       string
	address    string
	moduleType string
	bufferSize int

	certificate tls.Certificate
	store       *store.PipeStore
	broker      *broker.RemoteBroker
}

func newContext(module string) *context {
	sConfig, err := config.GetClientSocketConfig(module)
	if err != nil {
		klog.Errorf("failed to get config with error %+v", err)
		return nil
	}

	certificate, err := getCert(&sConfig)
	if err != nil {
		klog.Errorf("failed to get cert with error %+v", err)
	}

	remoteBroker := broker.NewRemoteBroker()

	return &context{
		name:        sConfig.ModuleName,
		moduleType:  sConfig.SocketType,
		address:     sConfig.Address,
		bufferSize:  sConfig.BufferSize,
		certificate: certificate,
		broker:      remoteBroker,
		store:       store.NewPipeStore(),
	}
}

func (m *context) AddModuleRemote(module string) {
	klog.Infof("add remote module: %s", module)
	conn := m.Connect(module, GetConnectFunc(m.moduleType))
	if conn == nil {
		klog.Errorf("failed to connect")
	}
}

func (m *context) AddModule(module string, usConn interface{}) {
	klog.Infof("add module: %v", module)
	conn, ok := usConn.(wrapper.Conn)
	if !ok {
		klog.Errorf("failed to add module, bad us conn")
		return
	}
	m.store.Add(module, conn)
}

func (m *context) AddModuleGroup(module, group string) {
	klog.Infof("add module(%v) to group(%v)", module, group)
	pipeInfo, err := m.store.Get(module)
	if err != nil {
		klog.Warningf("bad module name %s", module)
		return
	}

	conn := pipeInfo.Wrapper()
	if conn != nil {
		m.store.AddGroup(module, group, conn)
	}
}

func (m *context) Cleanup(module string) {
	klog.Infof("clean up module: %s", module)
	pipeInfo, err := m.store.Get(module)
	if err != nil {
	}

	conn := pipeInfo.Wrapper()
	if conn != nil {
		err = conn.Close()
	}
	m.store.Delete(module)
}

func (m *context) Send(module string, message model.Message) {
	pipeInfo, err := m.store.Get(module)
	if err != nil {
		klog.Warningf("failed to get module %s", module)
		return
	}
	message.SetType(m.moduleType)
	message.SetDestination(module)
	conn := pipeInfo.Wrapper()
	if conn != nil {
		err = m.broker.Send(conn, message)
		return
	}
	klog.Warningf("bad module name %s", module)
}

func (m *context) Receive(module string) (model.Message, error) {
	pipeInfo, err := m.store.Get(module)
	if err != nil {
		klog.Warningf("failed to get module pipe: %s", module)
		return model.Message{}, fmt.Errorf("failed to get module pipe: %v", err)
	}

	conn := pipeInfo.Wrapper()
	if conn != nil {
		return m.broker.Receive(conn)
	}

	klog.Warningf("bad module name: %s", module)
	return model.Message{}, fmt.Errorf("bad module name(%s)", module)
}

func (m *context) SendSync(module string, message model.Message, timeout time.Duration) (model.Message, error) {
	pipeInfo, err := m.store.Get(module)
	if err != nil {
		klog.Warningf("failed to get module pipe: %s", module)
		return model.Message{}, fmt.Errorf("failed to get module pipe: %v", err)
	}

	conn := pipeInfo.Wrapper()
	if conn == nil {
		klog.Warningf("bad module name: %s", module)
		return model.Message{}, fmt.Errorf("bad module name(%s)", module)
	}
	message.SetType(m.moduleType)
	message.SetDestination(module)
	return m.broker.SendSync(conn, message, timeout)
}

func (m *context) SendResp(message model.Message) {
	pipeInfo, err := m.store.Get(message.GetSource())
	if err != nil {
		klog.Warningf("failed to get module:%s", message.GetSource())
		return
	}

	conn := pipeInfo.Wrapper()
	if conn == nil {
		klog.Warningf("bad module name:%s", message.GetSource())
		return
	}
	message.SetDestination(message.GetSource())
	err = m.broker.Send(conn, message)
}

func (m *context) SendToGroup(group string, message model.Message) {
	var err error
	walkFunc := func(module string, pipe store.PipeInfo) error {
		conn := pipe.Wrapper()
		if conn == nil {
			klog.Warningf("bad pipe")
			return nil
		}
		message.SetDestination(module)
		err = m.broker.Send(conn, message)
		if err != nil {
			return err
		}
		return nil
	}

	err = m.store.WalkGroup(group, walkFunc)
}

func (*context) SendToGroupSync(moduleType string, message model.Message, timeout time.Duration) error {
	return fmt.Errorf("not supported now")
}

type ModuleExchange struct {
	Modules []string            `json:"modules"`
	Groups  map[string][]string `json:"groups"`
}

func (m *context) exchangeModuleInfo(conn wrapper.Conn, module string) error {
	moduleMsg := model.NewMessage("").
		BuildRouter(module, "", common.ResourceTypeModule, common.OperationTypeModule).
		SetType(m.moduleType).
		FillBody("")
	resp, err := m.broker.SendSyncInternal(conn, *moduleMsg, 0)
	if err != nil {
		klog.Errorf("failed to send module message with error %+v", err)
		return fmt.Errorf("failed to send module message, response:%+v, error: %+v", resp, err)
	}

	var exchange ModuleExchange
	bytes, err := json.Marshal(resp.GetContent())
	if err != nil {
		klog.Errorf("failed to marshal response with error %+v", err)
		return fmt.Errorf("failed to marshal response, error: %+v", err)
	}

	err = json.Unmarshal(bytes, &exchange)
	if err != nil {
		klog.Errorf("bad modules info from remote with error %+v", err)
		return fmt.Errorf("bad modules info from remote %+v with error %s", resp, err.Error())
	}

	for _, name := range exchange.Modules {
		if name == module {
			continue
		}
		klog.Infof("socket module: %s", name)
		m.store.Add(name, conn)
	}

	for group, modules := range exchange.Groups {
		for _, module := range modules {
			m.store.AddGroup(module, group, conn)
		}
	}

	klog.Infof("success to send module message")
	return nil
}

func (m *context) Connect(module string, connect broker.ConnectFunc) wrapper.Conn {
	opts := broker.ConnectOptions{
		Address:     m.address,
		MessageType: m.moduleType,
		BufferSize:  m.bufferSize,
		Cert:        m.certificate,
	}

	for {
		conn := m.broker.Connect(opts, connect)
		if conn == nil {
			time.Sleep(connectPeriod)
			continue
		}

		m.AddModule(module, conn)

		err := m.exchangeModuleInfo(conn, module)
		if err == nil {
			return conn
		}
		klog.Errorf("error to connect with %+v", err)

		err = conn.Close()
		time.Sleep(connectPeriod)
	}
}
