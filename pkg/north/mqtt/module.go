package mqtt

import (
	"encoding/json"
	"k8s.io/klog/v2"
	core "main/pkg/actor/core"
	actorContext "main/pkg/actor/core/context"
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"
	"main/pkg/common/modules"
	"main/pkg/config"
	"main/pkg/north/adapter"
	"time"
)

type MqttModule struct {
	enable          bool
	name            string
	dataUploadTopic string
	mqttClient      adapter.Adapter
	stopChan        chan struct{}
}

func NewMqttModuleFromDB(m *module.Module, property *[]moduleproperty.ModuleProperty) *MqttModule {
	mqttClientConfig := config.MqttModuleConfig{
		Enable: true,
		Name:   m.Name,
	}
	for _, moduleProperty := range *property {
		if moduleProperty.Name == "host" {
			mqttClientConfig.Host = moduleProperty.Value
			continue
		} else if moduleProperty.Name == "port" {
			mqttClientConfig.Port = moduleProperty.Value
			continue
		} else if moduleProperty.Name == "username" {
			mqttClientConfig.User = moduleProperty.Value
			continue
		} else if moduleProperty.Name == "password" {
			mqttClientConfig.Passwd = moduleProperty.Value
			continue
		} else if moduleProperty.Name == "clientId" {
			mqttClientConfig.ClientId = moduleProperty.Value
			continue
		}
	}
	newMqttModule := NewMqttModule(mqttClientConfig)
	return newMqttModule
}

func NewMqttModule(config config.MqttModuleConfig) *MqttModule {
	mqttClient := NewMqttClient(config)
	return &MqttModule{
		enable:          config.Enable,
		name:            config.Name,
		dataUploadTopic: config.DataUploadTopic,
		mqttClient:      &mqttClient,
		stopChan:        make(chan struct{}),
	}
}

func Register(m *MqttModule) {
	core.Register(m)
}

func (m *MqttModule) Name() string {
	return m.name
}

func (m *MqttModule) Group() string {
	return modules.NorthGroup
}

func (m *MqttModule) Enable() bool {
	return m.enable
}

func (m *MqttModule) SetEnable(enable bool) {
	m.enable = enable
}

func (m *MqttModule) Start() {
	err := m.mqttClient.Init()
	if err != nil {
		klog.Error("init mqtt client error %s", err)
		m.stopChan <- struct{}{}
	}
	for {
		select {
		case <-actorContext.Done():
			return
		default:
		}

		go m.sendToNorth()
		go m.sendToMiddle()
		//// go sendToSouth()
		//go m.keepalive("/heartbeat")

		<-m.stopChan
		klog.V(4).Infof("mqtt module %s stop...", m.Name())
		m.mqttClient.UnInit()
		m.enable = false
	}

}

func (m *MqttModule) sendToNorth() {
	for {
		select {
		case <-actorContext.Done():
			return
		default:
		}
		if !m.enable {
			break
		}
		message, err := actorContext.Receive(m.Name())
		if err != nil {
			klog.Errorf("receive message from %s channel error %s", m.Name(), err)
			continue
		}
		klog.V(4).Infof("send message to mqtt module %s.message %v", m.Name(), message)
		d, err := json.Marshal(message.Content)
		if err != nil {
			klog.Error("json error %s", err)
		}
		err = m.mqttClient.SendToNorth(d, m.dataUploadTopic)
		if err != nil {
			klog.Error("send to north error %s", err)
		}
	}
}

func (m *MqttModule) sendToMiddle() {
	for {
		select {
		case <-actorContext.Done():
			return
		default:
		}
		if !m.enable {
			break
		}
		//message := model.NewMessage("").SetRoute(modules.MqttModuleName, modules.MiddleGroup).
		//	SetResourceOperation("test", model.DeleteOperation).FillBody("test-message-" + time.Now().String())
		//actorContext.Send(modules.DeviceTwinModuleName, *message)
		klog.V(5).Infof("send message to middle module")
		time.Sleep(5 * time.Second)
	}
}

func (m *MqttModule) keepalive(topic string) {
	for {
		select {
		case <-actorContext.Done():
			return
		default:
		}
		if !m.enable {
			break
		}
		klog.V(5).Infof("mqtt module %s heartbeat...", m.Name())
		err := m.mqttClient.SendToNorth("ping-"+m.name, topic)
		if err != nil {
			klog.Warningf("mqtt module %s heartbeat error %v", m.Name(), err)
		}
		time.Sleep(10 * time.Second)
	}
}

func (m *MqttModule) Stop() {
	m.stopChan <- struct{}{}
	actorContext.Cleanup(m.Name())
}
