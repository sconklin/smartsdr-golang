package main

import (
	"strings"
)

type TopicPub struct {
	Topic string
	Value string
}

// parses a line full of "key=value" pairs separated by a separator
func simpleSplit(prefix string, sep string, line string) []TopicPub {
	var topics []TopicPub
	pairs := strings.Split(line, sep)
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			log.Infof("simpleSplit: Bad parse of: %s", pair)
		}
		topic := strings.Title(prefix) + "/" + strings.Title(parts[0])
		value := parts[1]
		log.Infof("%s = %s", topic, value)
		nt := TopicPub{topic, value}
		topics = append(topics, nt)
	}
	return topics
}

// parseRadio parses lines with the radio keyword
func parseRadio(prefix string, sep string, line string) []TopicPub {
	var topics []TopicPub
	tokens := strings.Split(line, sep)
	startidx := 0

	if tokens[0] == "filter_sharpness" {
		prefix = prefix + "/" + strings.Title(tokens[0]) + "/" + strings.Title(tokens[1])
		startidx = 2
	} else if tokens[0] == "static_net_params" {
		prefix = prefix + "/" + strings.Title(tokens[0])
		startidx = 1
	} else if tokens[0] == "oscillator" {
		prefix = prefix + "/" + strings.Title(tokens[0])
		startidx = 1
	}

	for _, pair := range tokens[startidx:] {
		var value string
		parts := strings.Split(pair, "=")

		if len(parts) == 2 {
			if len(strings.TrimSpace(parts[1])) == 0 {
				value = ""
			} else {
				value = parts[1]
			}
		} else {
			log.Infof("%d parts in parsing %v", len(parts), parts)
		}

		topic := strings.Title(prefix) + "/" + strings.Title(parts[0])
		log.Infof("%s = %s", topic, value)
		nt := TopicPub{topic, value}
		topics = append(topics, nt)
	}
	return topics
}

func processStatus(handle uint32, status string) {
	// This will move to an MQTT publisher
	log.Infof("Status: %s", status)
	respsegs := strings.SplitN(status, " ", 2)
	if len(respsegs) == 2 {
		switch respsegs[0] {
		case "radio":
			// list := simpleSplit(respsegs[0], " ", respsegs[1])
			_ = parseRadio(respsegs[0], " ", respsegs[1])
		case "transmit":
		case "waveform":
		case "xvtr":
		case "atu":
		case "amplifier":
		case "memories":
		case "slice":
		case "foundation":
		case "gps":
		case "scu":
		case "tx":
		case "eq":
		case "usb_cable":
		case "interlock":
		default:
			log.Infof("Unknown Status key: %s", respsegs[0])
		}
	} else {
		log.Infof("Unknown Status: %s", status)
	}

}
