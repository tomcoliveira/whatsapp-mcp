package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	qrcode "github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	_ "modernc.org/sqlite"
)

func main() {
	// Setup do client WhatsMeow com armazenamento local em SQLite
	dbLog := log.New(os.Stdout, "DB ", log.LstdFlags)
	container, err := sqlstore.New("sqlite", "file:store.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatal(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatal("Nenhum dispositivo encontrado:", err)
	}
	client := whatsmeow.NewClient(deviceStore, nil)

	// Gerencia o canal de QR code
	if client.Store.ID == nil {
		// ainda não está logado — escuta o QR
		qrChan, _ := client.GetQRChannel(context.Background())

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					fmt.Println("QR recebido. Salvando em qr.png...")
					err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "./qr.png")
					if err != nil {
						log.Println("Erro ao gerar imagem do QR:", err)
					} else {
						log.Println("QR salvo com sucesso como qr.png")
					}
				}
			}
		}()
	}

	// Conecta ao WhatsApp
	err = client.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar ao WhatsApp:", err)
	}

	// Servidor HTTP simples para expor o QR
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	log.Println("Servidor iniciado em http://localhost:8080 (acesse /qr.png)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
