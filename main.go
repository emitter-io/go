package main

import (
	"encoding/json"
	"fmt"
	"time"

	"./emitter"
)

var f = func(client emitter.Emitter, msg emitter.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())

	var dat map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &dat); err != nil {
		panic(err)
	}

	fmt.Println(dat["time"])
}

func main() {
	fmt.Println("Hello world!")

	o := emitter.NewClientOptions()
	o.AddBroker("tcp://api.emitter.io:8080")
	o.SetClientID("go-client")
	o.SetKeepAlive(60 * time.Second)
	o.SetOnMessageHandler(f)
	o.SetOnPresenceHandler(func(_ emitter.Emitter, p emitter.PresenceEvent) {
		fmt.Printf("Occupancy: %v\n", p.Occupancy)
	})

	c := emitter.NewClient(o)

	sToken := c.Connect()
	if sToken.Wait() && sToken.Error() != nil {
		panic("Error on Client.Connect(): " + sToken.Error().Error())
	}

	//c.Subscribe("z3D7-osAGTU2mQvCvuGXLcMvPXLGGGcy", "cluster")

	r := emitter.NewPresenceRequest()
	r.Key = "X4-nUeHjiAygHMdN8wst82S3c2KcCMn7"
	r.Channel = "presence-demo/1"
	x := c.Presence(r)
	if x.Wait() && x.Error() != nil {
		panic("Error on Client.Presence(): " + x.Error().Error())
	}

	// stop after 10 seconds
	time.Sleep(10 * time.Second)

}
