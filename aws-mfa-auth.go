package main

import (
	// "bytes"
	"context"
	"fmt"
	"log"
	"os"

	// "os"
	// "os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/spf13/viper"

	c "config"
)

func createFile(fileName string, fileData string) {
	dirname, err := os.UserHomeDir()
    if err != nil {
        log.Fatal( err )
    }
	credentialsFile, err := os.Create(dirname + "/.aws/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = credentialsFile.WriteString(fileData)
	if err != nil {
		credentialsFile.Close()
		log.Fatal(err)
	}

	err = credentialsFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.SetConfigType("yml")
	var configuration c.Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configuration.DefaultRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(configuration.User.AccKeyId, configuration.User.SecAccKey, "")))
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
		log.Println("Got an error assuming the role:")
		log.Fatal(err)
		return
	}

	profileName := fmt.Sprintln("[default]")
	AccKey := fmt.Sprintln("aws_access_key_id", "=", *result.Credentials.AccessKeyId)
	SecKey := fmt.Sprintln("aws_secret_access_key", "=", *result.Credentials.SecretAccessKey)
	DefaultReg := fmt.Sprintln("region", "=", configuration.DefaultRegion)
	SessTkn := fmt.Sprintln("aws_session_token", "=", *result.Credentials.SessionToken)

	credFileData := fmt.Sprintf("%s%s%s%s", profileName, AccKey, SecKey, SessTkn)
	createFile("credentials", credFileData)

	configFileData := fmt.Sprintf("%s%s", profileName, DefaultReg)
	createFile("config", configFileData)

	log.Println("Temp credentials created")
}
