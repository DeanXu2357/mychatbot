/*
Copyright © 2023 poyu <dean.xu.2357@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DeanXu2357/mychatbot/handler/discord"
	"github.com/DeanXu2357/mychatbot/llm"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Long:  `start the server`,
	Run:   RunServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func RunServer(cmd *cobra.Command, args []string) {
	fmt.Println("RunServer called")

	ctx := cmd.Context()

	ollama, errO := llm.NewOllama(
		viper.GetString("ollama.url"),
		viper.GetString("ollama.model"),
		viper.GetString("ollama.system"),
	)
	if errO != nil {
		log.Panic(errO)
	}

	token := viper.GetString("discord.token")
	discordHandler, errD := discord.New(token, ollama)
	if errD != nil {
		log.Panic(errD)
	}
	defer discordHandler.Close()
	if err := discordHandler.Handle(); err != nil {
		log.Panic(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	//r.GET("/discord/talk", discordHandler.Interaction)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", viper.GetString("server.port")),
		Handler:           r,
		ReadHeaderTimeout: 60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("listen:", err)
		}
	}()

	notifyCTX, _ := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-notifyCTX.Done()

	if err := srv.Shutdown(cmd.Context()); err != nil {
		log.Panicln("server shutdown:", err)
	}

	fmt.Println("shutting down gracefully, press Ctrl+C again to force close")
}
