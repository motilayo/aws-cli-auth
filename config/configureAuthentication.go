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

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"path/filepath"
)

func createConfigFile(fileName string, fileData string) {
	color.Set(color.FgRed)
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	credentialsFile, err := os.Create(dirname + "/.aws/" + fileName)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = credentialsFile.WriteString(fileData)
	if err != nil {
		credentialsFile.Close()
		log.Fatalln(err)
	}

	err = credentialsFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func readConfigFile(configFileType string, configFileName string, configFilePath string) Configurations {
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configFilePath)

	log.Println(viper.ConfigFileUsed())
	var configuration Configurations

	err := viper.ReadInConfig()
	color.Set(color.FgRed)

	if err != nil {
		log.Fatalln(fmt.Sprintf("Error reading config file, %s", err))
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Unable to decode into struct, %v", err))
	}

	return configuration
}

func assumeRole(configuration Configurations) *sts.AssumeRoleOutput {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configuration.DefaultRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(configuration.User.AccKeyId, configuration.User.SecAccKey, "")))
	if err != nil {
		log.Fatalln(err)
	}

	client := sts.NewFromConfig(cfg)

	color.Set(color.FgHiBlue)
	tokenCode, _ := stscreds.StdinTokenProvider()
	fmt.Println()

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(configuration.RoleArn),
		RoleSessionName: aws.String(configuration.SessionName),
		SerialNumber:    aws.String(configuration.MFASerial),
		TokenCode:       aws.String(tokenCode),
	}
	output, err := client.AssumeRole(context.TODO(), input)
	color.Set(color.FgRed)
	if err != nil {
		log.Fatalln(err)
	}

	return output
}

func Configure(fullFilePath string) {
	fullFilePath, err := filepath.Abs(strings.Replace(fullFilePath, "~", os.Getenv("HOME"), 1))
	if err != nil {
		log.Fatalln(err)
	}
	filePath, fileName := filepath.Split(fullFilePath)

	fileType := filepath.Ext(fileName)
	fileType = fileType[1:]

	color.Set(color.FgYellow)
	log.Println(fmt.Sprintf("config file type: %s \t config file path: %s \t config file name: %s", fileType, filePath, fileName))
	fmt.Println()

	configuration := readConfigFile(fileType, fileName, filePath)
	result := assumeRole(configuration)

	profileName := fmt.Sprintln("[default]")
	AccKey := fmt.Sprintln("aws_access_key_id", "=", *result.Credentials.AccessKeyId)
	SecKey := fmt.Sprintln("aws_secret_access_key", "=", *result.Credentials.SecretAccessKey)
	DefaultReg := fmt.Sprintln("region", "=", configuration.DefaultRegion)
	SessTkn := fmt.Sprintln("aws_session_token", "=", *result.Credentials.SessionToken)

	credFileData := fmt.Sprintf("%s%s%s%s", profileName, AccKey, SecKey, SessTkn)
	createConfigFile("credentials", credFileData)
	color.Set(color.FgGreen)
	log.Println("AWS credentials updated in ~/.aws/credentials")
	fmt.Println()

	configFileData := fmt.Sprintf("%s%s", profileName, DefaultReg)
	createConfigFile("config", configFileData)
	color.Set(color.FgGreen)
	log.Println("AWS CLI config updated in ~/.aws/config")
	fmt.Println()
}
