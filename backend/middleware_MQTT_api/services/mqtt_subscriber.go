package services

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SubscribeAndListen(broker string, port int, topic string) string {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go-subscriber")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Sprintf("Error de conexión: %v", token.Error())
	}
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Mensaje recibido en %s: %s\n", msg.Topic(), string(msg.Payload()))
	})
	select {} // Mantener la suscripción activa
	// Nunca se llega aquí, pero por compatibilidad:
	// return "Suscripción finalizada"
}
