package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/muka/peerjs-go/server"
	"github.com/spf13/viper"
)

func fail(err error, msg string) {
	if err == nil {
		return
	}
	if err != nil {
		fmt.Printf("%s: %s \n", msg, err)
		os.Exit(1)
	}
}

func main() {

	viper.AutomaticEnv()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("peer")
	viper.SetConfigName("peer")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			fail(err, "Failed to read config file")
		}
	}

	opts := server.NewOptions()
	if viper.IsSet("Host") {
		opts.Host = viper.GetString("Host")
	}
	if viper.IsSet("Port") {
		opts.Port = viper.GetInt("Port")
	}
	if viper.IsSet("LogLevel") {
		opts.LogLevel = viper.GetString("LogLevel")
	}
	if viper.IsSet("ExpireTimeout") {
		opts.ExpireTimeout = viper.GetInt64("ExpireTimeout")
	}
	if viper.IsSet("AliveTimeout") {
		opts.AliveTimeout = viper.GetInt64("AliveTimeout")
	}
	if viper.IsSet("Key") {
		opts.Key = viper.GetString("Key")
	}
	if viper.IsSet("Path") {
		opts.Path = viper.GetString("Path")
	}
	if viper.IsSet("ConcurrentLimit") {
		opts.ConcurrentLimit = viper.GetInt("ConcurrentLimit")
	}
	if viper.IsSet("AllowDiscovery") {
		opts.AllowDiscovery = viper.GetBool("AllowDiscovery")
	}
	if viper.IsSet("CleanupOutMsgs") {
		opts.CleanupOutMsgs = viper.GetInt("CleanupOutMsgs")
	}

	s := server.New(opts)
	defer s.Stop()
	err := s.Start()
	fail(err, "Failed to start server")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
