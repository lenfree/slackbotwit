package awswit

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type awsEC2Filter struct {
	Tag   string
	Value string
}

func newEC2() *ec2.EC2 {
	return ec2.New(session.New(), &aws.Config{
		Region: aws.String("us-west-2"),
	})
}

func getEC2List(e *ec2.EC2, f awsEC2Filter) *ec2.DescribeInstancesOutput {
	var params *ec2.DescribeInstancesInput
	if f.Tag != "" {
		params = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(f.Tag),
					Values: []*string{aws.String(f.Value)},
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running"), aws.String("pending")},
				},
			},
		}
	} else {
		params = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running"), aws.String("pending")},
				},
			},
		}
	}

	resp, err := e.DescribeInstances(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Printf("error %s %s %s", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				fmt.Printf("%s %s %d %s",
					reqErr.Code(),
					reqErr.Message(),
					reqErr.StatusCode(),
					reqErr.RequestID(),
				)
			}
		}
	}
	return resp
}
