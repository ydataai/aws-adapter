package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/aws-quota-provider/pkg/clients"
	"github.com/ydataai/aws-quota-provider/pkg/common"
	"github.com/ydataai/aws-quota-provider/pkg/controller"
	"github.com/ydataai/aws-quota-provider/pkg/server"
	"github.com/ydataai/aws-quota-provider/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/sirupsen/logrus"
)

func main() {
	restServiceConfiguration := service.RESTServiceConfiguration{}
	serverConfiguration := server.Configuration{}
	restControllerConfiguration := controller.RESTControllerConfiguration{}
	applicationConfiguration := configuration{}

	err := initConfigurationurationVariables([]common.ConfigurationVariables{
		&restServiceConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
		&applicationConfiguration,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("could not set configuration variables. Err: %v", err))
		os.Exit(1)
	}

	var log = logrus.New()
	log.SetLevel(applicationConfiguration.logLevel)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(applicationConfiguration.awsRegion),
	}))

	ec2Service := ec2.New(sess)
	serviceQuotaService := servicequotas.New(sess)

	ec2Client := clients.NewEC2Client(log, ec2Service)
	serviceQuotaClient := clients.NewServiceQuotaClient(log, serviceQuotaService)

	restService := service.NewRESTService(log, ec2Client, serviceQuotaClient, restServiceConfiguration)
	restController := controller.NewRESTController(log, restService, restControllerConfiguration)

	serverCtx := context.Background()

	s := server.NewServer(log, serverConfiguration)
	restController.Boot(s)

	s.Run(serverCtx)

	for err := range s.ErrCh {
		log.Error(err)
	}

}
func initConfigurationurationVariables(configurations []common.ConfigurationVariables) error {
	for _, configuration := range configurations {
		if err := configuration.LoadEnvVars(); err != nil {
			return err
		}
	}

	return nil
}
