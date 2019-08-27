package main

import (
	"strings"
)

type TopicPub struct {
	Topic string
	Value string
}

// parses a line full of "key=value" pairs separated by a separator
func simpleSplit(prefix string, sep string, tokens []string) []TopicPub {
	var topics []TopicPub
	for _, pair := range tokens {
		if strings.TrimSpace(pair) == "" {
			continue
		}
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
func parseRadio(prefix string, tokens []string) []TopicPub {
	var topics []TopicPub

	for _, pair := range tokens {
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
	respsegs := strings.Split(status, " ")
	var startidx = 1
	prefix := respsegs[0]

	switch respsegs[0] {
	case "radio":
		switch respsegs[1] {
		case "filter_sharpness":
			prefix = prefix + "/" + strings.Title(respsegs[1]) + "/" + strings.Title(respsegs[2])
			startidx = 3
		case "static_net_params", "oscillator":
			prefix = prefix + "/" + strings.Title(respsegs[1])
			startidx = 2
		default:
			startidx = 1
		}
		_ = parseRadio(prefix, respsegs[startidx:])
	case "transmit", "waveform", "atu", "interlock":
		_ = simpleSplit(respsegs[0], " ", respsegs[1:])
	case "xvtr":
		// next field is xvtr number
		prefix = prefix + "/" + strings.Title(respsegs[1])
		startidx = 2
		_ = simpleSplit(prefix, " ", respsegs[2:])
	case "amplifier":
	case "memories":
	case "slice":
		// next field is slice number
		prefix = prefix + "/" + strings.Title(respsegs[1])
		startidx = 2
		_ = simpleSplit(prefix, " ", respsegs[2:])
	case "foundation":
	case "gps":
	case "scu":
	case "tx":
	case "eq":
		switch respsegs[1] {
		case "rx", "rxsc", "tx", "txsc":
			prefix = prefix + "/" + strings.Title(respsegs[1])
			startidx = 2
			_ = simpleSplit(prefix, " ", respsegs[2:])
		default:
			log.Infof("Unexpected case for eq: %v", respsegs)

		}
	case "usb_cable":
	default:
		log.Infof("Unknown Status key: %s", respsegs[0])
	}
}
