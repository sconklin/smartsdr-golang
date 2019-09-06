package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type mqttPubInfo struct {
	Value       string
	LastChange  time.Time
	WhiteListed bool
}

var mqttClient *client.Client
var topicmap map[string]mqttPubInfo
var topicDumpFileName string

func ReadPubTopics() error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Errorf("Topics Dir Error: %v", err)
		return err
	}

	fpath := filepath.Join(dir, conf.MqttI.TopicsFile)
	if err != nil {
		log.Errorf("Topics path Error: %v", err)
		return err
	}

	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Errorf("Unable to read Topics File %s", fpath)
		return err
	}

	lines := strings.Split(string(content), "\n")

	// TODO remove blank lines and ones which are commented with '#'
	for _, line := range lines {
		val := strings.TrimSpace(line)
		if val != "" && !strings.HasPrefix(val, "#") {
			// Init the topic map with all the ones which are whitelisted
			ts := mqttPubInfo{"", time.Now(), true}
			topicmap[val] = ts
		}
	}
	log.Debugf("Topics Pub List: %v", topicmap)
	return nil
}

func MqttInit(errc chan<- error) error {
	log.Info("MQTTInit")
	// if there's a publish topic list file then get those topics into a list
	// otherwise, publish everything
	err := ReadPubTopics()
	if err != nil {
		return err
	}

	// If there's a dump file named, then assemble the filename
	if conf.MqttI.TopicDumpFile != "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Errorf("Topic Dump Dir Error: %v", err)
			return err
		}

		topicDumpFileName = filepath.Join(dir, conf.Radio.SubsFile)
		if err != nil {
			log.Errorf("Topic Dump path Error: %v", err)
			return err
		}

	} else {
		topicDumpFileName = ""
	}

	// open the mqtt client
	mqttClient = client.New(&client.Options{
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

	return nil
}

func MqttHandler(errc chan<- error, mqttcmdc chan<- string) {

	log.Info("MQTTHandler")

	// Subscribe to an MQTT topic. Accept commands on that and pass to the radio

	log.Infof("MQTT Subscribe to <%s>", conf.MqttI.TopicRead)
	err := mqttClient.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(conf.MqttI.TopicRead),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					log.Infof("MQTT RX <%s><%s>", string(topicName), string(message))
					mqttcmdc <- string(message)
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
	err := mqttClient.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS1,
		TopicName: []byte(topic),
		Message:   []byte(data),
	})
	return err
}

func MqttEnqueue(topicSuffix string, value string) error {
	var publish bool

	oldinfo, present := topicmap[topicSuffix]
	if present {
		/* we've sent this before */
		if oldinfo.WhiteListed {
			/* it's ok to publish if we need to */
			if oldinfo.Value != value {
				publish = true
			} else if time.Now().Sub(oldinfo.LastChange) > 60*time.Second {
				publish = true
			} else {
				publish = false
			}
		} else {
			publish = false
		}
	} else {
		/* first time we've seen this, by definition not whitelisted */
		publish = false
		/* Add it to the list with whitelist = false */
		ts := mqttPubInfo{value, time.Now(), false}
		topicmap[topicSuffix] = ts
		/* rewrite the dump file of topics */
		if topicDumpFileName != "" {
			f, err := os.Create("/tmp/dat2")
			if err != nil {
				return err
			}
			defer f.Close()

			w := bufio.NewWriter(f)

			for k, v := range topicmap {
				if !v.WhiteListed {
					fmt.Fprintf(w, "#")
				}
				fmt.Fprintf(w, "%s/n", k)
			}
			w.Flush()
		}

	}

	if publish {
		err := MqttPublish(topicSuffix, value)
		if err != nil {
			return err
		}
		// update the pubinfo
		oldinfo.LastChange = time.Now()
		oldinfo.Value = value
		topicmap[topicSuffix] = oldinfo
	}
	return nil
}
