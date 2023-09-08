package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
)

type Instance struct {
	Name string
	Id   string
	Ip   string
}

var options []string

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	check(err)
	clientEC2 := ec2.NewFromConfig(cfg)
	res, err := clientEC2.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	check(err)
	var instances []Instance
	for _, reservation := range res.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name == "terminated" || instance.State.Name == "stopped" {
				continue
			}
			var i Instance
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					i.Name = *tag.Value
				}
			}
			i.Id = *instance.InstanceId
			i.Ip = *instance.PrivateIpAddress
			instances = append(instances, i)
			options = append(options, fmt.Sprintf("%-28s  %-24s%s", i.Name, i.Id, i.Ip))
		}
	}
	answers := struct{ Instance string }{}
	err = survey.Ask([]*survey.Question{
		{
			Name: "instance",
			Prompt: &survey.Select{
				Message: "Which instance would you like to connect to:",
				Options: options,
			},
		},
	}, &answers)
	check(err)

	var target string
	for _, i := range instances {
		if strings.Contains(answers.Instance, i.Id) {
			target = i.Id
		}
	}

	clientSSM := ssm.NewFromConfig(cfg)
	input := &ssm.StartSessionInput{
		Target: aws.String(target),
	}
	ses, err := clientSSM.StartSession(context.TODO(), input)
	check(err)
	fmt.Println("started session", string(*ses.SessionId))

	// I followed the example for the rest this
	// https://github.com/mmmorris1975/aws-runas/blob/master/cli/ssm_cmd.go
	var inJ, outJ []byte
	outJ, err = json.Marshal(ses)
	check(err)

	inJ, err = json.Marshal(input)
	check(err)

	ep, err := ssm.NewDefaultEndpointResolver().ResolveEndpoint(cfg.Region, ssm.EndpointResolverOptions{})
	check(err)

	client := ssm.NewFromConfig(cfg)

	c := exec.Command("session-manager-plugin", string(outJ), cfg.Region, "StartSession", "", string(inJ), ep.URL)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()

	term, err := client.TerminateSession(context.TODO(), &ssm.TerminateSessionInput{
		SessionId: ses.SessionId,
	})
	check(err)

	fmt.Println("terminated session", string(*term.SessionId))
}
