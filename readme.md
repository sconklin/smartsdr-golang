# smartsdr-golang

smartsdr-golang is an attempt to re-implement a 'waveform' plugin for FlexRadio SmartSDR in the go language. 

By default, UDP is blocked by th eUbuntu Firewall, open the ports like this:
sudo ufw allow from 172.31.0.0/8 to any port 4992 proto udp


## TODO
[x] - discovery client

[x] - wfm->radio command interface (traffic_cop)

[x] - radio->wfm command interface

[ ] - VITA49 stream rx'er and tx'er (hal_*)

[ ] - actual sample processing pipeline
