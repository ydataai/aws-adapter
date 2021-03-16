package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ydataai/aws-quota-provider/pkg/common"
)

type configuration struct {
	logLevel  logrus.Level
	awsRegion string
}

func (c *configuration) LoadEnvVars() error {
	logLevel, err := common.VariableFromEnvironment("LOG_LEVEL")
	if err != nil {
		return err
	}

	awsRegion, err := common.VariableFromEnvironment("REGION")
	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	c.logLevel = level
	c.awsRegion = awsRegion

	return nil
}
