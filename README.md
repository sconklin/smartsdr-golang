# ssdr-infobroker

[![Build Status](https://travis-ci.org/sconklin/ssdr-infobroker.svg?branch=master)](https://travis-ci.org/sconklin/sdr-infobroker)
[![Go Report Card](https://goreportcard.com/badge/github.com/sconklin/ssdr-infobroker)](https://goreportcard.com/report/github.com/sconklin/ssdr-infobroker)
[![GoDoc](https://godoc.org/github.com/sconklin/ssdr-infobroker?status.svg)](https://godoc.org/github.com/sconklin/ssdr-infobroker)
[![MIT License](http://img.shields.io/badge/License-GPLv3-blue.svg)](./LICENSE)


This is intended to be an experimental gateway between the Flex Radio API and an MQTT broker to publish information about what's happening in the radio.

It is forked from [smartsdr-golang](https://github.com/baobrien/smartsdr-golang)

By default, UDP is blocked by the Ubuntu Firewall, open the ports like this:
sudo ufw allow from 172.31.0.0/8 to any port 4992 proto udp

## TODO
[x] - discovery client

[x] - wfm->radio command interface (traffic_cop)

[x] - radio->wfm command interface

[ ] - VITA49 stream rx'er and tx'er (hal_*)

[ ] - actual sample processing pipeline
