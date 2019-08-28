package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Radio related represents the radio information we can subscribe to
type Rad struct {
	SubsFile        string `json:"radio_sub_list_file"`
	IncludeMessages string `json:"include_messages"`
}

// Mqtt Represents the MQTT items
type Mqtt struct {
	PubTopicPrefix  string `json:"pub_topic_prefix"`
	TopicRead       string `json:"topic_read"`
	Broker          string `json:"broker"`
	BrokerUser      string `json:"broker_user"`
	BrokerPass      string `json:"broker_pass"`
	BrokerPort      string `json:"broker_port"`
	PubOnChange     string `json:"only_publish_on_change"`
	FreqErrorChange string `json:"freq_error_change"`
	TopicsFile      string `json:"topic_pub_list_file"`
	TopicDumpFile   string `json:"topic_dump_file"`
}

// Config Represents the top-level config structure
type Config struct {
	Radio Rad  `json:"Radio"`
	MqttI Mqtt `json:"Mqtt"`
}

// ReadConfig reads the config json file
func ReadConfig(jsonFileName string) (*Config, error) {
	file, err := os.Open(jsonFileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}

	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// DumpConfig prints the config information
func DumpConfig(conf *Config) {
	fmt.Printf("Radio:\n")
	fmt.Printf("    Radio Sub List File:  %s\n", conf.Radio.SubsFile)
	fmt.Printf("    Include messages?:    %s\n", conf.Radio.IncludeMessages)
	fmt.Printf("Mqtt:\n")
	fmt.Printf("    Pub Topic Prefix:     %s\n", conf.MqttI.PubTopicPrefix)
	fmt.Printf("    Topic Read:           %s\n", conf.MqttI.TopicRead)
	fmt.Printf("    Broker:               %s\n", conf.MqttI.Broker)
	fmt.Printf("    User:                 %s\n", conf.MqttI.BrokerUser)
	fmt.Printf("    Pass:                 %s\n", conf.MqttI.BrokerPass)
	fmt.Printf("    Port:                 %s\n", conf.MqttI.BrokerPort)
	fmt.Printf("    Publish On Change?    %s\n", conf.MqttI.PubOnChange)
	fmt.Printf("    Topics File           %s\n", conf.MqttI.TopicsFile)
	fmt.Printf("    Topic Dump File:      %s\n", conf.MqttI.TopicDumpFile)
}
