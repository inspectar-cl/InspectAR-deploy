package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fmt.Println("🚀 Iniciando middleware MQTT...")

	// Configuración del broker
	broker := "emqx:1883"
	clientID := "go-subscriber"
	maxRetries := 10
	retryDelay := 5 * time.Second

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(retryDelay)
	opts.SetDefaultPublishHandler(messageHandler)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Conectado al broker MQTT")
		if token := c.Subscribe("sensors/#", 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println("❌ Error al suscribirse:", token.Error())
		} else {
			fmt.Println("📡 Suscrito al topic 'sensors/#'")
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("🔌 Conexión perdida:", err)
		fmt.Println("🔄 Intentando reconectar...")
	}

	opts.OnReconnecting = func(c mqtt.Client, opts *mqtt.ClientOptions) {
		fmt.Println("🔄 Reconectando al broker MQTT...")
	}

	// Crear cliente y conectar con reintentos
	client := mqtt.NewClient(opts)

	fmt.Printf("🔌 Intentando conectar a %s...\n", broker)
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("📡 Intento %d/%d de conexión...\n", i+1, maxRetries)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("❌ Error en intento %d: %v\n", i+1, token.Error())
			if i < maxRetries-1 {
				fmt.Printf("⏳ Esperando %v antes del siguiente intento...\n", retryDelay)
				time.Sleep(retryDelay)
			}
		} else {
			fmt.Println("✅ Conexión establecida exitosamente!")
			break
		}

		if i == maxRetries-1 {
			fmt.Printf("💥 No se pudo conectar después de %d intentos. Terminando...\n", maxRetries)
			return
		}
	}

	// Esperar señal de salida
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Cerrar conexión
	fmt.Println("🛑 Desconectando...")
	client.Disconnect(250)
}

// Manejador de mensajes
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📨 Mensaje recibido en %s: %s\n", msg.Topic(), string(msg.Payload()))
	fmt.Println("----------------------------------------")

	// Enviar el mensaje recibido a la API REST
	apiURL := "http://iot-service:8090/lectura/"

	// Parsear el JSON y asegurar que 'valor' sea float64 y timestamp termine en Z
	var data map[string]interface{}
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Println("❌ Error parseando JSON:", err)
		return
	}
	// Normalizar 'valor' a float64
	valor, ok := data["valor"]
	if !ok {
		fmt.Println("❌ El campo 'valor' no está presente en el mensaje")
		return
	}
	switch v := valor.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			data["valor"] = f
		} else {
			fmt.Println("❌ No se pudo convertir 'valor' string a float64:", err)
			return
		}
	case float64:
		// ya está bien
	case float32:
		data["valor"] = float64(v)
	case int:
		data["valor"] = float64(v)
	case int64:
		data["valor"] = float64(v)
	case int32:
		data["valor"] = float64(v)
	default:
		fmt.Printf("❌ Tipo inesperado para 'valor': %T\n", v)
		return
	}
	if ts, ok := data["timestamp"].(string); ok && len(ts) > 0 && ts[len(ts)-1] != 'Z' {
		data["timestamp"] = ts + "Z"
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("❌ Error serializando JSON:", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Error enviando a la API REST:", err)
		return
	}
	defer resp.Body.Close()

	// Leer y mostrar la respuesta de la API
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("➡️  Enviado a API REST, status: %s, respuesta: %s\n", resp.Status, string(respBody))
	fmt.Printf("➡️  JSON enviado: %s\n", string(jsonData))
}
