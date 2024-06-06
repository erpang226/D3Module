package socket

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"k8s.io/klog/v2"

	"main/pkg/actor/common"
	"main/pkg/actor/core/socket/broker"
	"main/pkg/actor/core/socket/config"
)

const (
	connectPeriod    = 5 * time.Second
	HandshakeTimeout = 60 * time.Second
)

func getCert(config *config.SocketConfig) (tls.Certificate, error) {
	if config.Key == "" &&
		config.Cert == "" {
		return tls.Certificate{}, nil
	}

	var err error
	var certificate tls.Certificate
	if config.Cert != "" && config.Key != "" {
		certificate, err = tls.LoadX509KeyPair(config.Cert, config.Key)
	} else {
		err = fmt.Errorf("failed to get x509 key pair")
	}
	return certificate, err
}

func GetConnectFunc(moduleType string) broker.ConnectFunc {
	switch moduleType {
	case common.MsgCtxTypeUS:
		return Connect
	}
	klog.Warningf("not supported module type: %v", moduleType)
	return nil
}

func Connect(opts broker.ConnectOptions) (interface{}, error) {
	conn, err := net.Dial(opts.MessageType, opts.Address)
	if err != nil {
		klog.Errorf("failed to dail addrs: %s", opts.Address)
		return nil, err
	}
	return conn, nil
}
