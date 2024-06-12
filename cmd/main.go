package main

import (
	"errors"
	"flag"
	"os"

	"github.com/bartosian/notibot/internal/core/controllers"
	"github.com/bartosian/notibot/internal/core/ports"
	"github.com/bartosian/notibot/internal/gateways/cherrygw"
	"github.com/bartosian/notibot/internal/gateways/discordgw"
	"github.com/bartosian/notibot/internal/gateways/twiliogw"
	"github.com/bartosian/notibot/pkg/l0g"
	"github.com/joho/godotenv"
)

var (
	requiredTwilioVars  = []string{"TWILIO_ACCOUNT_SID", "TWILIO_AUTH_TOKEN", "TWILIO_PHONE_NUMBER", "CLIENT_PHONE_NUMBER"}
	requiredDiscordVars = []string{"DISCORD_BOT_TOKEN", "DISCORD_CHANNEL"}
	requiredCherryVars  = []string{"CHERRY_AUTH_TOKEN", "CHERRY_SERVER_ID"}
	errEnvVarNotFound   = errors.New("environment variable not found")
)

type flags struct {
	envFilePath    string
	monitorDiscord bool
}

func main() {
	logger := l0g.NewLogger()
	flagSet := parseFlags()

	if flagSet.envFilePath != "" {
		loadEnvFile(flagSet.envFilePath, logger)
	}

	if !flagSet.monitorDiscord {
		logger.Info("all monitors disabled - dismissing", nil)
		os.Exit(0)
	}

	checkRequiredEnvVars(requiredDiscordVars, logger)
	checkRequiredEnvVars(requiredCherryVars, logger)
	checkRequiredEnvVars(requiredTwilioVars, logger)

	notifierGateway := twiliogw.NewTwilioGateway(logger)
	discordGateway, err := discordgw.NewDiscordGateway(logger)
	exitIfError(err, logger, "error creating discord gateway")

	cherryGateway, err := cherrygw.NewCherryGateway(logger)
	exitIfError(err, logger, "error creating cherry gateway")

	notifierController := controllers.NewNotifierController(notifierGateway, discordGateway, cherryGateway, logger)

	startMonitoring(notifierController)
}

func parseFlags() *flags {
	envFilePath := flag.String("env-file", "../.envrc", "path to the environment variables file")
	monitorDiscord := flag.Bool("discord", false, "receive notifications from discord channels")

	flag.Parse()

	return &flags{
		envFilePath:    *envFilePath,
		monitorDiscord: *monitorDiscord,
	}
}

func loadEnvFile(filePath string, logger l0g.Logger) {
	if err := godotenv.Load(filePath); err != nil {
		logger.Error("error loading .env file:", err, map[string]interface{}{"path": filePath})
		os.Exit(1)
	}
}

func checkRequiredEnvVars(envVars []string, logger l0g.Logger) {
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			logger.Error("missing environment variable", errEnvVarNotFound, map[string]interface{}{"var": envVar})
			os.Exit(1)
		}
	}
}

func exitIfError(err error, logger l0g.Logger, message string) {
	if err != nil {
		logger.Error(message, err, nil)
		os.Exit(1)
	}
}

func startMonitoring(notifierController ports.NotifierController) {
	if err := notifierController.MonitorDiscord(); err != nil {
		os.Exit(1)
	}
}
