package main

import (
	"fmt"
	"os"
	"strings"
)

func autoParseResponse(tokens []string) []TopicPub {
	var topics []TopicPub
	prefix := ""
	idx := 0

	// for as long as we have fields with no "=" at the beginning, append them to the topic
	for _, token := range tokens {
		if strings.Contains(token, "=") {
			break
		}
		prefix = prefix + strings.TrimSpace(token) + "/"
		idx = idx + 1
	}

	for _, pair := range tokens[idx:] {
		if strings.TrimSpace(pair) == "" {
			continue
		}

		var value string

		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			if len(strings.TrimSpace(parts[1])) == 0 {
				value = ""
			} else {
				value = parts[1]
			}
		} else {
			log.Debugf("%d parts in parsing %v", len(parts), parts)
		}

		topic := prefix + parts[0]
		log.Debugf("%s = %s", topic, value)
		nt := TopicPub{topic, value}
		topics = append(topics, nt)
	}
	return topics
}

func parseGps(status string) []TopicPub {
	var topics []TopicPub
	prefix := ""

	parts := strings.SplitN(status, " ", 2)

	prefix = parts[0]
	tokens := strings.Split(parts[1], "#")

	for _, pair := range tokens {
		if strings.TrimSpace(pair) == "" {
			continue
		}

		var value string

		parts := strings.Split(pair, "=")

		if len(parts) == 2 {
			if len(strings.TrimSpace(parts[1])) == 0 {
				value = ""
			} else {
				value = parts[1]
			}
		} else {
			log.Debugf("XXXXXXXXXXXX %d parts in parsing %v", len(parts), parts)
		}

		topic := prefix + "/" + parts[0]
		log.Debugf("%s = %s", topic, value)
		nt := TopicPub{topic, value}
		topics = append(topics, nt)
	}
	return topics
}

func parseProfile(status string) []TopicPub {
	var topics []TopicPub
	prefix := ""

	parts := strings.SplitN(status, " ", 3)

	prefix = parts[0] + "/" + parts[1]

	payload := strings.Split(parts[2], "=")

	prefix = prefix + "/" + payload[0]

	// this will be either "list" or "current", in either case the following works

	value := payload[1]
	log.Debugf("%s = %s", prefix, value)
	nt := TopicPub{prefix, value}
	topics = append(topics, nt)
	return topics
}

func parseClient(status string) []TopicPub {
	// client 0xDE997CDB connected local_ptt=1 client_id=C97FB7E3-C2D7-4AB7-AC8B-417A3F9ECEF5 program=SmartSDR-Win station=DESKTOP-M3NJRL1
	// convert this into suitable mqtt topics
	// perhaps .../client/connected/id 0xDE997CDB
	//            /client/0xDE997CDB/local_ptt 1
	//            /client/0xDE997CDB/client_id C97FB7E3-C2D7-4AB7-AC8B-417A3F9ECEF5
	// ets
	var topics []TopicPub
	prefix := ""

	parts := strings.SplitN(status, " ", 3)

	prefix = parts[0] + "/" + parts[1]

	payload := strings.Split(parts[2], "=")

	prefix = prefix + "/" + payload[0]

	// this will be either "list" or "current", in either case the following works

	value := payload[1]
	log.Debugf("%s = %s", prefix, value)
	nt := TopicPub{prefix, value}
	topics = append(topics, nt)
	return topics
}

func processStatus(handle uint32, status string) {
	var toPub []TopicPub

	respsegs := strings.Split(status, " ")

	switch respsegs[0] {
	case "radio", "transmit", "waveform", "atu", "interlock", "xvtr", "slice", "eq", "usb_cable", "memory", "wan":
		toPub = autoParseResponse(respsegs)
	case "amplifier":
		log.Infof("Status: %s", status)
	case "memories":
		log.Infof("Status: %s", status)
	case "foundation":
		log.Infof("Status: %s", status)
	case "gps":
		toPub = parseGps(status)
	case "scu":
		log.Infof("Status: %s", status)
	case "profile":
		_ = parseProfile(status)
	case "client":
		// Client needs a special handler
		log.Infof("Status: %s", status)
		toPub = parseClient(status)
	default:
		log.Infof("Unknown Status key: %s", respsegs[0])
		log.Infof("Status: %s", status)
	}

	for _, thing := range toPub {
		err := MqttEnqueue(thing)
		if err != nil {
			// TODO something else
			fmt.Println("Enqueue Error: ", err)
			os.Exit(1)
		}
	}
}
