package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
	"main/pkg/config"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"testing"
)

// TestNewMqttClient tests the NewMqttClient function that creates the MqttClient object
func TestNewMqttClient(t *testing.T) {
}

// TestInit tests the procurement of the MqttClient
func TestInit(t *testing.T) {

}

// TestSend checks send function by sending message to server
func TestSend(t *testing.T) {

}

// TestReceive sends the message through send function then calls receive function to see same message is received or not
func TestReceive(t *testing.T) {

}

func TestSubscribe(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	mqttClient := MqttClient{
		config: config.MqttModuleConfig{
			Enable:   true,
			Name:     "test",
			Host:     "localhost",
			Port:     "1883",
			ClientId: "my-test",
			User:     "admin",
			Passwd:   "123456",
		},
		client: nil,
	}

	err := mqttClient.Init()
	if err != nil {
		t.Error("init error")
	}
	err = mqttClient.subscribeTopic("/test", func(client mqtt.Client, message mqtt.Message) {
		t.Logf("client receive message %v", string(message.Payload()))
	})
	if err != nil {
		t.Error("subscribe error", err)
	}

	// 阻塞等待接收信号
	<-sigs
	fmt.Println("\nReceived an interrupt, stopping service...")
	mqttClient.UnInit()
	// 优雅退出程序
	os.Exit(0)

}

// 1.2w/s
var payload1200b = "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234" +
	"5678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
	"012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123" +
	"4567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678" +
	"901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012" +
	"345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456" +
	"7890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901" +
	"2345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456" +
	"7890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901" +
	"2345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456" +
	"7890123456789012345678901234567890123456789012345678901234567890"

// 3.4w/s
var payload400b = "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234" +
	"5678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
	"012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123" +
	"4567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678"

func TestMqttClients(t *testing.T) {
	var wg sync.WaitGroup

	// 启动多个goroutine来模拟并发任务
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		mqttClient := MqttClient{
			config: config.MqttModuleConfig{
				Enable:   true,
				Name:     "test",
				Host:     "192.168.21.43",
				Port:     "1883",
				ClientId: "test-client" + strconv.Itoa(i),
				User:     "admin",
				Passwd:   "123456",
			},
			client: nil,
		}
		err := mqttClient.Init()
		if err != nil {
			t.Error("init error", err)
		}
		go simulateTask(i, mqttClient, payload400b, 1000000, &wg)
	}

	// 等待所有任务完成
	wg.Wait()
	klog.Info("All tasks completed")
}

func simulateTask(taskId int, client MqttClient, payload string, total int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer client.UnInit()

	for i := 0; i < total; i++ {
		err := client.publishTopic("test/topic", payload)
		if err != nil {
			klog.Error("publish topic error", err)
		}
	}
	klog.Infof("tast %d total completed...", taskId)
}
