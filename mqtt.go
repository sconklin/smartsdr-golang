package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type mqttPubInfo struct {
	Value      string
	LastChange time.Time
}

func MqttHandler(errc chan<- error, mqttcmdc chan<- string) {

	var lastAz float64
	var deltap float64

	timeLast := time.Now()

	log.Info("MQTT New")
	mqttClient := client.New(&client.Options{
		ErrorHandler: func(err error) {
			errc <- err
		},
	})

	defer mqttClient.Terminate()

	log.Info("MQTT Cert: %s", conf.MqttI.BrokerCert)
	// Read the certificate file.
	b, err := ioutil.ReadFile(conf.MqttI.BrokerCert)
	if err != nil {
		errc <- err
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(b); !ok {
		errc <- errors.New("failed to parse root certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: roots,
	}

	ipstr := conf.MqttI.Broker + ":" + conf.MqttI.BrokerPort
	log.Infof(" MQTT Server: %s", ipstr)

	log.Info("MQTT Connect")
	// Connect to the MQTT Server.
	err = mqttClient.Connect(&client.ConnectOptions{
		Network:   "tcp",
		Address:   ipstr,
		ClientID:  []byte("netrotor"),
		UserName:  []byte(conf.MqttI.BrokerUser),
		Password:  []byte(conf.MqttI.BrokerPass),
		TLSConfig: tlsConfig,
	})
	if err != nil {
		errc <- err
	}

	// Subscribe to an MQTT topic. Accept commands on that and pass to the radio

	log.Infof("MQTT Subscribe to <%s>", conf.MqttI.TopicRead)
	err = mqttClient.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(conf.MqttI.TopicRead),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					log.Infof("MQTT RX <%s><%s>", string(topicName), string(message))
					mqttcmdc <- message
				},
			},
		},
	})
	if err != nil {
		errc <- err
	}
}

func MqttPublish(topicSuffix string, data string) error {
	/* Publish */
	topic := conf.MqttI.PubTopicPrefix + topicSuffix
	log.Infof(" MQTT Server: publishing %s = %s", topic, data)
	return mqttClient.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS1,
		TopicName: []byte(topic),
		Message:   []byte(data),
	})
}

/*
type mqttPubInfo struct {
	Value      string
	LastChange Time
}
*/
var topicmap map[string]mqttPubInfo

func MqttEnqueue(topicSuffix string, value string) err {
	oldinfo, present := topicmap[topicSuffix]
	var publish bool

	if present {
		/* we've sent this before */
		if oldinfo.Value != value {
			publish = true
		} else if time.Now().Sub(oldinfo.LastChange) > 60*time.Second {
			/* send it */
		}
	} else {
		/* first time we've seen this, save it */
	}

}
