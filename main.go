package main

import (
        "sync"
        "encoding/json"
        "fmt"
        "os"
        "gopkg.in/ini.v1"
        "net/http"
        "bytes"
        "io/ioutil"

        mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Request struct {
    LUX int               `json:"Lux"`
}

func main() {
        cfg, err := ini.Load("/data/msg.ini")
        if err != nil {
                conf, err = ini.Load("msg.ini")
                if err != nil {
                        fmt.Printf("Fail to read file: %v", err)
                        os.Exit(1)
                }
        }

        wt_token := "Bearer " + cfg.Section("token").Key("wt_bot_token").String()
        mqtt_broker := cfg.Section("server").Key("mqtt_broker").String()
        mqtt_port := cfg.Section("server").Key("mqtt_port").String()
        webex_msg := cfg.Section("message").Key("wt_msg").String()
        webex_to_email := cfg.Section("message").Key("person_to_message").String()
        const TOPIC = "/sense"

        jsonData := map[string]string{"toPersonEmail": webex_to_email, "text": webex_msg}
        jsonValue, _ := json.Marshal(jsonData)

        full_broker := "tcp://" + mqtt_broker + ":" + mqtt_port

        opts := mqtt.NewClientOptions().AddBroker(full_broker)

        client := mqtt.NewClient(opts)
        if token := client.Connect(); token.Wait() && token.Error() != nil {
                fmt.Println("FATAL ERROR")
        }

        var wg sync.WaitGroup
        wg.Add(1)

        if token := client.Subscribe(TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {
                data := Request{}
                s := string(msg.Payload())
                json.Unmarshal([]byte(s), &data)
                if data.LUX <= 10 {
                        request, _ := http.NewRequest("POST", "https://api.ciscospark.com/v1/messages", bytes.NewBuffer(jsonValue))
                        request.Header.Set("Content-Type", "application/json")
                        request.Header.Set("authorization", wt_token)
                        client := &http.Client{}
                        response, err := client.Do(request)
                        if err != nil {
                            fmt.Printf("The HTTP request failed with error %s\n", err)
                        } else {
                            data, _ := ioutil.ReadAll(response.Body)
                            fmt.Println(string(data))
                        }
                }
                wg.Add(1)
        }); token.Wait() && token.Error() != nil {
                fmt.Println("FATAL ERROR")
        }

        wg.Wait()
}
