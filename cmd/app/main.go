package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/telegram"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
)

// Config vars.
var (
	environmentMode   string
	configurationPath string
)

func main() {
	askFlags()

	fmt.Println("Initialization")

	// Init. channels for communication with gorutines
	// Structures are used, not pointers for communication through channels,
	// 	since collision occurs when getting target structure by pointer in goroutines
	messagesChannel := make(chan *model.UpdatedMessage, 128)
	notificationsChannel := make(chan *modelDB.NotificationQueue, 128)
	errorsChannel := make(chan string, 256)
	storeChannel := make(chan *model.UpdatedMessage, 128)

	// Creating an instance of Config at first and load it
	config := config.New().Load(configurationPath, environmentMode)

	// Creating an instance of Manager (contains repos and config)
	manager := manager.New(config)

	// Creating an instance of TelegramGatewayService
	telegramGateway := telegram.NewTelegramGateway(manager)

	// Creating an instance of UserService
	userService := service.NewUserService(manager)

	// Creating an instance of ChatService
	chatService := service.NewChatService(manager)

	// Creating an instance of TranslatorGateway
	translatorGateway := translator.NewTranslatorGateway(manager)

	// Creating an instance of TranslatorService
	translator := translator.NewTranslatorService(translatorGateway)

	// Creating an instace of TelegramService
	telegramService := telegram.NewTelegramService(
		manager,
		telegramGateway,
		userService,
		chatService,
		translator,
		notificationsChannel,
		messagesChannel,
		errorsChannel,
		storeChannel,
	)

	// Creating an instance of TelegramBotService
	bot := telegram.NewTelegramBot(telegramService, errorsChannel)

	// Close connection with database in defer
	defer manager.Repository.Close()

	fmt.Println("Handling messages ...")

	go func() {
		for {
			bot.ProcessErrors()
		}
	}()

	var m sync.Mutex

	bot.ProcessMessages(&m)

	for {
		runtime.Gosched()
		// bot.ProcessNotifications()

		// Timeout
		time.Sleep(1 * time.Second)
	}
}

// askFlags - getting args. from cli
func askFlags() {
	flag.StringVar(&environmentMode, "env-mode", config.ProdMode, "one of env. modes: prod|dev")
	flag.StringVar(&configurationPath, "config-path", config.DefaultConfigPath, "path to config file")
	flag.Parse()
}
