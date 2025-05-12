package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Manejador de cierre de conexiones ordenado (Al apretar Ctrl+C)
	// Se crea un canal para recibir señales de interrupción
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//Armado de URL de conexión
	url := url.URL{
		// Se utiliza WebSocket simple (ws)
		Scheme: "ws",
		Host:   "echo.websocket.events",
		Path:   "/",
	}
	log.Printf("Conectando con %s", url.String())

	// Conexión al WebSocket
	conexion, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	// Si la conexión falla, se muestra un error
	if err != nil {
		log.Fatal("Error de conexión:", err)
	}
	defer conexion.Close()

	// Canal para recibir mensajes
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			// Lee mensajes del WebSocket
			_, mensaje, err := conexion.ReadMessage()
			if err != nil {
				log.Println("Error al leer mensaje:", err)
				return
			}
			log.Printf("Mensaje recibido: %s", mensaje)
		}
	}()

	// Envío de mensajes al WebSocket
	msj := "Hola, soy un cliente de WebSocket"
	// Se establece un mensaje de tipo texto plano y se lo convierte a bytes
	err = conexion.WriteMessage(websocket.TextMessage, []byte(msj))
	if err != nil {
		log.Println("Error al enviar mensaje:", err)
		return
	}
	log.Printf("Mensaje enviado: %s", msj)

	// Espera a que se reciba una señal de interrupción
	select {
	//Caso en el que la conexión se cierra
	case <-done:
	// Caso en el que no se enviaron mensajes por más de 5 segundos
	case <-time.After(5 * time.Second):
		log.Println("Tiempo de espera agotado")
	}

	log.Println("Cerrando conexión")
}
