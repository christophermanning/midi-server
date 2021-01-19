# midi-server

Examples of how to send MIDI message events via WebSocket. Tested in Ubuntu 20.04 LTS.

First, connect the MIDI device, use `amidi -l` to view the connected
MIDI devices, and update the MIDI deviceID in the program if necessary.

## Golang Server

Install go and portmidi, then run:

```
go run main.go
```

## websocketd

Install websocketd, then run:

```
websocketd -port=8080 ./midi.sh
```
