package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
	"github.com/oklog/ulid"
)

type covid struct {
	Nombre   string
	Apellido string
}

/*
	"encoding/json"
	"context"
*/
func main() {
	natsURL := os.Getenv("NATS_URL")
	clusterID := os.Getenv("NATS_CLUSTER_ID")
	topic := os.Getenv("NATS_TOPIC")
	addr := ":" + os.Getenv("PORT")

	err := startPublisher(natsURL, clusterID, addr, topic)
	if err != nil {
		log.Fatal(err)
	}
}

type msg_COVID struct {
	Name         string
	Location     string
	Age          int
	Infectedtype string
	State        string
	CAMINO       string
}

func startPublisher(natsURL, clusterID, addr, topic string) error {
	publisher, err := nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: clusterID,
			ClientID:  "publisher",
			StanOptions: []stan.Option{
				stan.NatsURL(natsURL),
			},
			Marshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return err
	}

	h := handler{topic, publisher}
	http.HandleFunc("/", h.Handle)

	log.Print("Listening on ", addr)
	return http.ListenAndServe(addr, nil)
}

type handler struct {
	topic     string
	publisher message.Publisher
}

func (h handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	err := h.publish(w, r)
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(500)
		return
	}
}

func (h handler) publish(w http.ResponseWriter, r *http.Request) error {

	w.Header().Set("Content-Type", "application/json")
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	//failOnError(err, "Parsing JSON")
	body["CAMINO"] = "NATS"
	data, err := json.Marshal(body)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	uuid := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
	msg := message.NewMessage(uuid.String(), payload)
	if err := h.publisher.Publish(h.topic, msg); err != nil {
		return err
	}
	/*
			u1 := covid{payload}
		    json_data, err := json.Marshal(u1)
		    if err != nil {

		        log.Fatal(err)
		    }
		    fmt.Println(string(json_data))
	*/

	_, err = fmt.Fprint(w, "Sent message: ", string(data), " x ", msg.UUID, "\n")
	if err != nil {
		return err
	}
	return nil
}
