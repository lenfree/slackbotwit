package awswit

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

func newELB() *elb.ELB {
	return elb.New(session.New(), &aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
}

func getELBList(e *elb.ELB) *elb.DescribeLoadBalancersOutput {
	var params *elb.DescribeLoadBalancersInput

	resp, err := e.DescribeLoadBalancers(params)

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

func getELBName(e *elb.ELB, n string) *elb.LoadBalancerDescription {
	var params *elb.DescribeLoadBalancersInput

	resp, err := e.DescribeLoadBalancers(params)

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
	for _, elb := range resp.LoadBalancerDescriptions {
		if *elb.LoadBalancerName == n {
			return elb
		}
	}
	return nil
}
