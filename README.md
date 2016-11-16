# Emitter Golang SDK
This repository contains Go/Golang client for [Emitter](https://emitter.io) (see also on [Emitter GitHub](https://github.com/emitter-io/emitter)). Emitter is an **open-source** real-time communication service for connecting online devices. Infrastructure and APIs for IoT, gaming, apps and real-time web. At its core, emitter.io is a distributed, scalable and fault-tolerant publish-subscribe messaging platform based on MQTT protocol and featuring message storage.

This library provides a nicer MQTT interface fine-tuned and extended with specific features provided by [Emitter](https://emitter.io). The code uses the [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang) for handling all the network communication and MQTT protocol, and is released under the same license (EPL v1).

## Installation and Build

This client, similarly to the Eclipse Paho client is designed to work with the standard Go tools, so installation is as easy as:

```
go get github.com/emitter-io/go
```

The client depends on Eclipse Paho MQTT Go Client and Google's websockets package, also easily installed with the command:

```
go get github.com/eclipse/paho.mqtt.golang
go get golang.org/x/net/websocket
go get github.com/satori/go.uuid
```