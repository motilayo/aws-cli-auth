package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/spf13/viper"

	c "aws-cli-auth/config"
)

func execCmd(val string) {
	cmd := exec.Command("eval", val)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // standard output
	cmd.Stderr = &stderr // standard error
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var configuration c.Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configuration.DefaultRegion))
	if err != nil {
		panic(err)
	}

	client := sts.NewFromConfig(cfg)

	tokenCode, _ := stscreds.StdinTokenProvider()

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(configuration.RoleArn),
		RoleSessionName: aws.String(configuration.SessionName),
		SerialNumber:    aws.String(configuration.MFASerial),
		TokenCode:       aws.String(tokenCode),
	}
	result, err := client.AssumeRole(context.TODO(), input)
	if err != nil {
		fmt.Println("Got an error assuming the role:")
		fmt.Println(err)
		return
	}

	exportAccKey := fmt.Sprintf(`"AWS_ACCESS_KEY_ID=%s"`, *result.Credentials.AccessKeyId)
	exportSecKey := fmt.Sprintf(`"AWS_SECRET_ACCESS_KEY=%s"`, *result.Credentials.SecretAccessKey)
	exportDefaultReg := fmt.Sprintf(`"AWS_DEFAULT_REGION=%s"`, configuration.DefaultRegion)
	exportSessTkn := fmt.Sprintf(`"AWS_SESSION_TOKEN=%s"`, *result.Credentials.SessionToken)

	execCmd(exportAccKey)
	execCmd(exportSecKey)
	execCmd(exportDefaultReg)
	execCmd(exportSessTkn)
}
