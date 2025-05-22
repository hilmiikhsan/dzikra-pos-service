package main

import (
	"flag"
	"os"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/cmd"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
)

func main() {
	os.Args = initialize()

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	// seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Info().Msg("No command provided, defaulting to 'serve'")
		runServe(serverCmd, os.Args[1:])
		return
	}

	switch os.Args[1] {
	case "serve":
		runServe(serverCmd, os.Args[2:])
	case "server":
		cmd.RunServerHTTP(serverCmd, os.Args[2:])
	case "grpc":
		cmd.RunServeGRPC()
	case "seed":
		// cmd.RunSeed(seedCmd, os.Args[2:])
	default:
		log.Info().Msg("Invalid command, defaulting to 'serve'")
		if os.Args[1][0] == '-' {
			runServe(serverCmd, os.Args[1:])
		} else {
			runServe(serverCmd, os.Args[2:])
		}
	}
}

func runServe(serverCmd *flag.FlagSet, args []string) {
	go func() {
		log.Info().Msg("Starting gRPC server…")
		cmd.RunServeGRPC()
	}()

	log.Info().Msg("Starting HTTP server…")

	cmd.RunServerHTTP(serverCmd, args)
}

func initialize() []string {
	configPath := flag.String("config_path", "./", "path to config file")
	configFilename := flag.String("config_filename", ".env", "config file name")
	flag.Parse()

	var cfgFile string
	if *configPath == "./" {
		cfgFile = *configPath + *configFilename
	} else {
		cfgFile = *configPath + "/" + *configFilename
	}

	log.Info().Msgf("Initializing configuration with config: %s", cfgFile)

	config.Configuration(
		config.WithPath(*configPath),
		config.WithFilename(*configFilename),
	).Initialize()

	adapter.Adapters = &adapter.Adapter{}

	var newArgs []string
	for _, arg := range os.Args {
		if strings.Contains(arg, "config_path") || strings.Contains(arg, "config_filename") {
			continue
		}

		newArgs = append(newArgs, arg)
	}

	return newArgs
}
