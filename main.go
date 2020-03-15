package main

import "encoding/json"
import "github.com/ian-kent/go-log/log"
import "github.com/akamensky/argparse"
import "io/ioutil"
import "os"
import "os/signal"
import "syscall"
import "net/http"

type config struct {
	Serving     ServeConfig `json:"Server"`
	ProjectName string      `json:"ProjectName"`
}

type arguments struct {
	ConfigPath string
}

// Config struct used to store the settings from the config file
var Config config

// Args struct used to store user passed arguments
var Args *arguments

// Server the http server used to host the graphql endpoint
var Server *http.Server

func shutdown(signals chan os.Signal, shutdown chan bool) {
	signal := <-signals
	// TODO: Change this to shutdown gracefully
	Server.Close()
	log.Debug("Received signal %q", signal)
	log.Info("Release Queue Shutting Down")
	shutdown <- true
}

func parseArguments() error {
	parser := argparse.NewParser("Release Queue Server", "Records pending releases and provides stats on release frequency etc")
	configPath := parser.String("c", "config", &argparse.Options{Default: "config.json"})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Error("Error Reading arguments: %q", err)
		return err
	}
	Args = &arguments{
		ConfigPath: *configPath,
	}
	log.Debug("Config path set to %q", Args.ConfigPath)
	return nil
}

func loadConfig(config *config) error {

	data, err := ioutil.ReadFile(Args.ConfigPath)
	if err != nil {
		log.Error("Failed to read config file %q. Error was %q", Args.ConfigPath, err)
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	log.Info("Config loaded from %q", Args.ConfigPath)
	return nil
}

func main() {
	log.Info("Staring Release Queue Server")
	signals := make(chan os.Signal, 1)
	shutdownReady := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	err := parseArguments()
	if err != nil {
		os.Exit(1)
	}
	err = loadConfig(&Config)
	if err != nil {
		os.Exit(1)
	}
	Server = serveGraphQL(Config.Serving.Port)
	go shutdown(signals, shutdownReady)
	<-shutdownReady
	log.Info("Release Queue Server Shutdown")
}
