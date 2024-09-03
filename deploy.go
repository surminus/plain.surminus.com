package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func deploy() error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)

	options := &ec2.RunInstancesInput{
		MaxCount:     1,
		MinCount:     1,
		ImageId:      "",
		InstanceType: "",
		KeyName:      "",
	}

	client.RunInstances(ctx, &ec2.RunInstancesInput{
		MaxCount: 1,
	})

	return nil
}
