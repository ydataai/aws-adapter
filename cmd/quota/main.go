package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"

	"github.com/ydataai/aws-adapter/pkg/clients"
	"github.com/ydataai/aws-adapter/pkg/controller"
	"github.com/ydataai/aws-adapter/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

var (
	errChan chan error
)

func main() {
	restServiceConfiguration := service.RESTServiceConfiguration{}
	serverConfiguration := server.HTTPServerConfiguration{}
	restControllerConfiguration := config.RESTControllerConfiguration{}
	applicationConfiguration := Configuration{}
	loggerConfiguration := logging.LoggerConfiguration{}

	if err := config.InitConfigurationVariables([]config.ConfigurationVariables{
		&restServiceConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
		&applicationConfiguration,
	}); err != nil {
		fmt.Println(fmt.Errorf("could not set configuration variables. Err: %v", err))
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(applicationConfiguration.AWSRegion),
	}))

	ec2Service := ec2.New(sess)
	serviceQuotaService := servicequotas.New(sess)

	ec2Client := clients.NewEC2Client(logger, ec2Service)
	serviceQuotaClient := clients.NewServiceQuotaClient(logger, serviceQuotaService)

	restService := service.NewRESTService(logger, ec2Client, serviceQuotaClient, restServiceConfiguration)
	restController := controller.NewRESTController(logger, restService, restControllerConfiguration)

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
