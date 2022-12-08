package main

import (
	"fmt"
	"log"

	"concurrency-demo/internal/command"
	"concurrency-demo/internal/provider"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rootCommand := &cobra.Command{
		Use:   "cli",
		Short: "makes api call to newsapi endpoint",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			app := provider.New()
			app.Provide()
		},
	}

	httpRequest := command.HTTPRequestsCommand()
	rootCommand.AddCommand(httpRequest)

	rootCommand.CompletionOptions.DisableDefaultCmd = true

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
	}

}
