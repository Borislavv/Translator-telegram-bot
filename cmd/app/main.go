package main

import (
	"flag"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Config vars.
var (
	environmentMode   string
	configurationPath string
)

func main() {
	askFlags()

	// Creating an instance of Config at first
	config := loadConfig()

	// Creating an instance of Manager (contains repos and config)
	manager := manager.New(config)

	// Creating an instance of TelegramGatewayService
	gateway := service.NewTelegramGateway(manager)

	//Creating an instance of UserService
	userService := service.NewUserService(manager)

	// Creating an instance of TelegramBotService
	bot := service.NewTelegramBot(manager, gateway, userService)

	// Close connection with database in defer
	defer manager.Repository.Close()

	for {
		// Handle batch of UpdatedMessages
		bot.HandlingMessages()

		// Timeout before new request
		time.Sleep(1 * time.Second)
	}
}

// askFlags - getting args. from cli
func askFlags() {
	flag.StringVar(&environmentMode, "env-mode", config.ProdMode, "one of env. modes: prod|dev")
	flag.StringVar(&configurationPath, "config-path", config.DefaultConfigPath, "path to config file")
	flag.Parse()
}

// loadConfig - loading a config file to struct
func loadConfig() *config.Config {
	config := config.New()

	_, err := toml.DecodeFile(configurationPath, config.Repository)
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
	}

	_, err = toml.DecodeFile(configurationPath, config.Integration)
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
	}

	return config
}
