package main

import (
	"strings"
)

type TopicPub struct {
	Topic string
	Value string
}

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
			log.Infof("%d parts in parsing %v", len(parts), parts)
		}

		topic := strings.Title(prefix) + strings.Title(parts[0])
		log.Infof("%s = %s", topic, value)
		nt := TopicPub{topic, value}
		topics = append(topics, nt)
	}
	return topics
}

func parseGps(status string) []TopicPub {
	var topics []TopicPub
	prefix := ""

	parts := strings.SplitN(status, " ", 2)

	prefix = strings.Title(parts[0])
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
			log.Infof("XXXXXXXXXXXX %d parts in parsing %v", len(parts), parts)
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
	// log.Infof("Status: %s", status)
	respsegs := strings.Split(status, " ")

	switch respsegs[0] {
	case "radio", "transmit", "waveform", "atu", "interlock", "xvtr", "slice", "eq", "usb_cable":
		_ = autoParseResponse(respsegs)
	case "amplifier":
		log.Infof("Status: %s", status)
	case "memories":
		log.Infof("Status: %s", status)
	case "foundation":
		log.Infof("Status: %s", status)
	case "gps":
		_ = parseGps(status)
	case "scu":
		log.Infof("Status: %s", status)
	case "tx":
		log.Infof("Status: %s", status)
	default:
		log.Infof("Unknown Status key: %s", respsegs[0])
	}
}
