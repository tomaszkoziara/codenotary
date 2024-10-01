package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomaszkoziara/codenotarybe/accounting"
	"github.com/tomaszkoziara/codenotarybe/api"
	"github.com/tomaszkoziara/codenotarybe/config"
	"github.com/tomaszkoziara/codenotarybe/store/immudbvault"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	immudbVaultClient := immudbvault.New(cfg.Vault.Ledger, cfg.Vault.Collection, cfg.Vault.APIKey)
	accountingService := accounting.New(immudbVaultClient)
	router := api.CreateRouter(accountingService)

	server := &http.Server{
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		fmt.Printf("cannot start listener: %v\n", err)
	}
	go func() {
		fmt.Println("starting server")
		err := server.Serve(listener)
		if err != nil {
			fmt.Printf("server has stopped listening: %v\n", err)
		}
	}()

	fmt.Println("server started")
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)
	<-wait

	fmt.Println("closing server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("erro while shuting down server: %v\n", err)
	}
}
