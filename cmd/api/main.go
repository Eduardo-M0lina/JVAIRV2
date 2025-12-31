package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// Inicializar contenedor de dependencias
	configPath := filepath.Join(".", "configs")
	container, err := NewContainer(configPath)
	if err != nil {
		log.Fatalf("Error al inicializar el contenedor: %v", err)
	}
	defer func() {
		if err := container.Close(); err != nil {
			log.Printf("Error al cerrar el contenedor: %v", err)
		}
	}()

	log.Println("Conexión a la base de datos establecida correctamente")

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", container.Config.Server.Port),
		Handler:      container.Router,
		ReadTimeout:  container.Config.Server.ReadTimeout,
		WriteTimeout: container.Config.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para recibir errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("Servidor iniciado en %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para recibir señales del sistema operativo
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Bloquear hasta recibir una señal o un error
	select {
	case err := <-serverErrors:
		log.Fatalf("Error al iniciar el servidor: %v", err)

	case <-shutdown:
		log.Println("Apagando el servidor...")

		// Crear un contexto con timeout para el apagado
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Apagar el servidor
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Error al apagar el servidor: %v", err)
			err = server.Close()
			if err != nil {
				log.Printf("Error al cerrar el servidor: %v", err)
			}
		}

		log.Println("Servidor apagado correctamente")
	}
}
