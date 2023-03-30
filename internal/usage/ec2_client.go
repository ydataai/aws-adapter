// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"context"

	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2Client defines a interface for ec2 client
type EC2Client interface {
	GetGPUInstances(context.Context, string) (GPU, error)
}

type ec2Client struct {
	logger logging.Logger
	ec2    *ec2.EC2
}

// NewEC2Client initializes ec2 client
func NewEC2Client(logger logging.Logger, ec2 *ec2.EC2) EC2Client {
	return ec2Client{
		logger: logger,
		ec2:    ec2,
	}
}

// GetGPUInstances fetches gpu instances
func (sq ec2Client) GetGPUInstances(ctx context.Context, gpuInstaceType string) (GPU, error) {
	sq.logger.Infof("Starting to fetch running %s GPU instances", gpuInstaceType)

	inputs := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []*string{aws.String(gpuInstaceType)},
			},
		},
	}

	gpuInstances, err := sq.ec2.DescribeInstancesWithContext(ctx, &inputs)
	if err != nil {
		sq.logger.Error(err)
		return 0, err
	}
	gpuInstancesNumber := GPU(len(gpuInstances.Reservations))

	sq.logger.Infof("Number of running %s instances: %d", gpuInstaceType, gpuInstancesNumber)

	return gpuInstancesNumber, nil
}
