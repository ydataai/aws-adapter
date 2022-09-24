package clients

import (
	"github.com/ydataai/aws-adapter/pkg/common"
	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2ClientInterface defines a interface for ec2 client
type EC2ClientInterface interface {
	GetGPUInstances(string) (common.GPU, error)
}

// EC2Client is the ec2 client
type EC2Client struct {
	log logging.Logger
	ec2 *ec2.EC2
}

// NewEC2Client initializes ec2 client
func NewEC2Client(log logging.Logger, ec2 *ec2.EC2) EC2Client {
	return EC2Client{
		log: log,
		ec2: ec2,
	}
}

// GetGPUInstances fetches gpu instances
func (sq EC2Client) GetGPUInstances(gpuInstaceType string) (common.GPU, error) {
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

	gpuInstancesNumber := common.GPU(len(gpuInstances.Reservations))

	return gpuInstancesNumber, nil
}
