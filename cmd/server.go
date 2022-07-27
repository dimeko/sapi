package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dimeko/sapi/api"
	"github.com/dimeko/sapi/app"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var srvCmd = &cobra.Command{
	Use:   "server",
	Short: "Starting server",
	Run:   start,
}

func init() {
	rootCmd.AddCommand(srvCmd)
}

func start(command *cobra.Command, args []string) {
	StartServer()
}

func StartServer() {
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		panic("Cannot find .env file")
	}

	port := os.Getenv("APP_PORT")
	app := app.New()
	api := api.New(app)
	httpServer := &http.Server{
		Handler:      api.Router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Println("Starting server on port:", port)
		log.Fatal(httpServer.ListenAndServe())
	}()

	shutdown := make(chan os.Signal, 0)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("Shutting down server gracefully in 1 second.")
	time.Sleep(time.Second)
	defer cancel()

	log.Fatal(httpServer.Shutdown(ctx))
	os.Exit(0)
}
