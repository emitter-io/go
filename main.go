package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())

	var dat map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &dat); err != nil {
		panic(err)
	}

	fmt.Println(dat)
}

func main() {
	fmt.Println("Hello world!")

	o := mqtt.NewClientOptions()
	o.AddBroker("tcp://api.emitter.io:8080")
	o.SetClientID("go-client")
	o.SetKeepAlive(60 * time.Second)
	o.SetDefaultPublishHandler(f)
	c := mqtt.NewClient(o)

	sToken := c.Connect()
	if sToken.Wait() && sToken.Error() != nil {
		panic("Error on Client.Connect(): " + sToken.Error().Error())
	}

	c.Subscribe("z3D7-osAGTU2mQvCvuGXLcMvPXLGGGcy/cluster", 0, nil)

	// stop after 10 seconds
	time.Sleep(10 * time.Second)

}
