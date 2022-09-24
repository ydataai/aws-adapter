package service_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/ydataai/aws-adapter/mock"
	"github.com/ydataai/aws-adapter/pkg/clients"
	"github.com/ydataai/aws-adapter/pkg/common"
	"github.com/ydataai/aws-adapter/pkg/service"
	"github.com/ydataai/go-core/pkg/common/logging"
)

func TestAvailableGPU(t *testing.T) {
	loggerConfiguration := logging.LoggerConfiguration{}
	if err := loggerConfiguration.LoadFromEnvVars(); err != nil {
		fmt.Println(fmt.Errorf("could not set logging configuration. Err: %v", err))
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	t.Run("failure response", func(t *testing.T) {
		errM := errors.New("mock error")

		tt := []struct {
			name          string
			ec2M          func(context.Context, *gomock.Controller) clients.EC2ClientInterface
			serviceQuotaM func(context.Context, *gomock.Controller) clients.ServiceQuotaClientInterface
			err           error
		}{
			{
				name: "failure on ec2 request",
				ec2M: func(ctx context.Context, ctrl *gomock.Controller) clients.EC2ClientInterface {
					ec2M := mock.NewMockEC2ClientInterface(ctrl)
					ec2M.EXPECT().
						GetGPUInstances(gomock.Any()).Return(common.GPU(0), errM)

					return ec2M
				},
				serviceQuotaM: func(ctx context.Context, ctrl *gomock.Controller) clients.ServiceQuotaClientInterface {
					serviceQuotaM := mock.NewMockServiceQuotaClientInterface(ctrl)

					return serviceQuotaM
				},
				err: errM,
			},
			{
				name: "failure on service quota request",
				ec2M: func(ctx context.Context, ctrl *gomock.Controller) clients.EC2ClientInterface {
					ec2M := mock.NewMockEC2ClientInterface(ctrl)
					ec2M.EXPECT().
						GetGPUInstances(gomock.Any()).Return(common.GPU(0), nil)

					return ec2M
				},
				serviceQuotaM: func(ctx context.Context, ctrl *gomock.Controller) clients.ServiceQuotaClientInterface {
					serviceQuotaM := mock.NewMockServiceQuotaClientInterface(ctrl)

					serviceQuotaM.EXPECT().
						GetAvailableGPUInstances(gomock.Any(), gomock.Any()).
						Return(common.GPU(0), errM)

					return serviceQuotaM
				},
				err: errM,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ctx := context.Background()

				restServiceConfiguration := service.RESTServiceConfiguration{}

				restService := service.NewRESTService(
					logger,
					tc.ec2M(ctx, ctrl),
					tc.serviceQuotaM(ctx, ctrl),
					restServiceConfiguration,
				)

				_, err := restService.AvailableGPU(ctx)
				if err == nil {
					t.Fatal("should return an error")
				}
			})
		}

	})

	t.Run("successful response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		restServiceConfiguration := service.RESTServiceConfiguration{}

		ec2M := mock.NewMockEC2ClientInterface(ctrl)
		ec2M.EXPECT().
			GetGPUInstances(gomock.Any()).Return(common.GPU(2), nil)

		serviceQuotaM := mock.NewMockServiceQuotaClientInterface(ctrl)
		serviceQuotaM.EXPECT().
			GetAvailableGPUInstances(gomock.Any(), gomock.Any()).Return(common.GPU(4), nil)

		restService := service.NewRESTService(
			logger,
			ec2M,
			serviceQuotaM,
			restServiceConfiguration,
		)

		gpu, err := restService.AvailableGPU(ctx)
		if err != nil {
			t.Fatal("should not return any error")
		}

		if diff := cmp.Diff(gpu, common.GPU(2)); diff != "" {
			t.Fatalf("should be 2, got %v", gpu)
			t.Fatal(diff)
		}
	})
}
