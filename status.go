package main

import (
	"strings"
)

type TopicPub struct {
	Topic string
	Value string
}

// parseResponse parses lines with the radio keyword
func parseResponse(prefixTokens int, tokens []string) []TopicPub {
	var topics []TopicPub

	prefix := ""

	for _, text := range tokens[0:prefixTokens] {
		prefix = prefix + "/" + text
	}

	for _, pair := range tokens[prefixTokens+1:] {
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
	pt := 0

	switch respsegs[0] {
	case "radio":
		switch respsegs[1] {
		case "filter_sharpness":
			pt = 3
		case "static_net_params", "oscillator":
			pt = 2
		default:
			pt = 1
		}
		_ = parseResponse(pt, respsegs)
	case "transmit", "waveform", "atu", "interlock":
		_ = parseResponse(1, respsegs)
	case "xvtr", "slice", "eq", "usb_cable":
		// next field is id of the thing
		_ = parseResponse(2, respsegs)
	case "amplifier":
	case "memories":
	case "foundation":
	case "gps":
	case "scu":
	case "tx":
	default:
		log.Infof("Unknown Status key: %s", respsegs[0])
	}
}
