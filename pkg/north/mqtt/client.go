package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
	"main/pkg/config"
	"time"
)

type MqttClient struct {
	config config.MqttModuleConfig
	client mqtt.Client
}

// NewMqttClient initializes a new mqtt client instance
func NewMqttClient(conf config.MqttModuleConfig) MqttClient {
	return MqttClient{config: conf}
}

func (m *MqttClient) Init() error {
	klog.V(4).Infof("init mqttclient...")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", m.config.Host, m.config.Port))
	if len(m.config.ClientId) > 0 {
		opts.SetClientID(m.config.ClientId)
	} else {
		opts.SetClientID(fmt.Sprintf("%s-%d", m.config.Name, time.Now().Unix()))
	}
	opts.SetUsername(m.config.User)
	opts.SetPassword(m.config.Passwd)
	opts.SetConnectionLostHandler(nil)
	opts.SetOnConnectHandler(nil)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetCleanSession(true)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetReconnectingHandler(nil)
	mqttClient := mqtt.NewClient(opts)
	token := mqttClient.Connect()
	if token.Wait() && token.Error() == nil {
		m.client = mqttClient
		return nil
	}
	return token.Error()
}

func (m *MqttClient) UnInit() {
	klog.V(4).Info("un-init mqtt client...", m.config.Name)
	m.client.Disconnect(2000)
}

// SendToNorth topic
func (m *MqttClient) SendToNorth(payload interface{}, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("topic is not set")
	}
	return m.publishTopic(params[0], payload)
}

func (m *MqttClient) ReceiveFromNorth(params ...interface{}) error {
	var topic = ""
	if len(params) > 0 {
		topic = params[0].(string)
	} else {
		return fmt.Errorf("topic is not set")
	}
	return m.subscribeTopic(topic, params[1].(mqtt.MessageHandler))
}

func (m *MqttClient) subscribeTopic(topic string, handler mqtt.MessageHandler) error {
	token := m.client.Subscribe(topic, 0, handler)
	if token.Wait() && token.Error() != nil {
		klog.Errorf("subscribe topic %s error %v", topic, token.Error())
		return token.Error()
	}
	return nil
}

func (m *MqttClient) publishTopic(topic string, payload interface{}) error {
	token := m.client.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		klog.Error("publish error ", token.Error())
		return token.Error()
	}
	return nil
}
