package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/spf13/viper"
)

func createConfigFile(fileName string, fileData string) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
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

func readConfigFile(configFileType string, configFileName string, configFilePath string) Configurations {

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configFilePath)
	viper.AddConfigPath(".")

	var configuration Configurations

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}

func assumeRole(configuration Configurations) *sts.AssumeRoleOutput {
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
	output, err := client.AssumeRole(context.TODO(), input)
	if err != nil {
		log.Println("Got an error assuming the role:")
		log.Fatal(err)
	}

	return output
}

func Configure(fullFilePath string) {
	splitPath := strings.Split(fullFilePath, ".")
	configFileType := splitPath[len(splitPath)-1]
	configFileNameAndPath := strings.Split(splitPath[0], "/")
	configFileName := configFileNameAndPath[len(configFileNameAndPath)-1]
	ConfigFilePath := strings.Join(configFileNameAndPath[:len(configFileNameAndPath)-1], "/")

	log.Printf(fmt.Sprintf("config file type: %s \t config file path: %s \t config file name: %s", configFileType, ConfigFilePath, configFileName))

	configuration := readConfigFile(configFileType, configFileName, ConfigFilePath)
	result := assumeRole(configuration)

	profileName := fmt.Sprintln("[default]")
	AccKey := fmt.Sprintln("aws_access_key_id", "=", *result.Credentials.AccessKeyId)
	SecKey := fmt.Sprintln("aws_secret_access_key", "=", *result.Credentials.SecretAccessKey)
	DefaultReg := fmt.Sprintln("region", "=", configuration.DefaultRegion)
	SessTkn := fmt.Sprintln("aws_session_token", "=", *result.Credentials.SessionToken)

	credFileData := fmt.Sprintf("%s%s%s%s", profileName, AccKey, SecKey, SessTkn)
	createConfigFile("credentials", credFileData)
	log.Println("AWS credentials updated in ~/.aws/credentials")

	configFileData := fmt.Sprintf("%s%s", profileName, DefaultReg)
	createConfigFile("config", configFileData)
	log.Println("AWS CLI config updated in ~/.aws/config")
}
