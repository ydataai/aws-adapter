// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2Client defines a interface for ec2 client
type EC2Client interface {
	GetGPUInstances(string) (GPU, error)
}

type ec2Client struct {
	log logging.Logger
	ec2 *ec2.EC2
}

// NewEC2Client initializes ec2 client
func NewEC2Client(log logging.Logger, ec2 *ec2.EC2) EC2Client {
	return ec2Client{
		log: log,
		ec2: ec2,
	}
}

// GetGPUInstances fetches gpu instances
func (sq ec2Client) GetGPUInstances(gpuInstaceType string) (GPU, error) {
	sq.log.Info("Starting to featch running GPU instances")

	inputs := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []*string{aws.String(gpuInstaceType)},
			},
		},
	}

	gpuInstances, err := sq.ec2.DescribeInstances(&inputs)
	if err != nil {
		sq.log.Error(err)
		return 0, err
	}

	gpuInstancesNumber := GPU(len(gpuInstances.Reservations))

	return gpuInstancesNumber, nil
}
