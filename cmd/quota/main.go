// Package main for quota executable
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/aws-adapter/internal/configuration"
	"github.com/ydataai/aws-adapter/internal/usage"
	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

var (
	errChan chan error
)

func main() {
	applicationConfiguration := configuration.Application{}
	restServiceConfiguration := usage.ServiceConfiguration{}
	serverConfiguration := server.HTTPServerConfiguration{}
	restControllerConfiguration := config.RESTControllerConfiguration{}
	loggerConfiguration := logging.LoggerConfiguration{}

	if err := config.InitConfigurationVariables([]config.ConfigurationVariables{
		&restServiceConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
		&applicationConfiguration,
		&loggerConfiguration,
	}); err != nil {
		fmt.Println(fmt.Errorf("could not set configuration variables. Err: %v", err))
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	awsConfig := &aws.Config{}
	if len(applicationConfiguration.Region) > 0 {
		awsConfig.Region = aws.String(applicationConfiguration.Region)
	}
	sess := session.Must(session.NewSession(awsConfig))

	ec2Service := ec2.New(sess)
	serviceQuotaService := servicequotas.New(sess)

	ec2Client := usage.NewEC2Client(logger, ec2Service)
	serviceQuotaClient := usage.NewQuotaClient(logger, serviceQuotaService)

	restService := usage.NewService(logger, ec2Client, serviceQuotaClient, restServiceConfiguration)
	restController := usage.NewRESTController(logger, restService, restControllerConfiguration)

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
