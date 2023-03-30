// Package main for metering executable
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/ydataai/aws-adapter/internal/configuration"
	"github.com/ydataai/aws-adapter/internal/metering"
	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"
)

var (
	errChan chan error
)

func main() {
	applicationConfiguration := configuration.Application{}
	serverConfiguration := server.HTTPServerConfiguration{}
	restControllerConfiguration := config.RESTControllerConfiguration{}
	loggerConfiguration := logging.LoggerConfiguration{}
	meteringConfiguration := metering.Configuration{}

	configs := []config.ConfigurationVariables{
		&applicationConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
		&loggerConfiguration,
		&meteringConfiguration,
	}
	if err := config.InitConfigurationVariables(configs); err != nil {
		fmt.Printf("Failed to initialize env configurations. Err: %v", err)
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	awsConfig := &aws.Config{}
	if len(applicationConfiguration.Region) > 0 {
		awsConfig.Region = aws.String(applicationConfiguration.Region)
	}
	sess := session.Must(session.NewSession(awsConfig))
	meteringClient := metering.NewClient(meteringConfiguration, marketplacemetering.New(sess))

	restController := metering.NewRESTController(logger, meteringClient, restControllerConfiguration)

	serverCtx := context.Background()
	httpServer := server.NewServer(logger, serverConfiguration)
	httpServer.AddHealthz()
	httpServer.AddReadyz(nil)
	restController.Boot(httpServer)
	httpServer.Run(serverCtx)

	for err := range errChan {
		logger.Error(err)
	}
}
