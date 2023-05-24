package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Config struct {
	systemManager *ssm.SSM
	parameters    map[ConfigParameterName]string
}

type ConfigParameterDefinition struct {
	ParameterName ConfigParameterName
	ParameterType ConfigParameterType
}

type ConfigParameterType string

const (
	StandardParameter    ConfigParameterType = "Standard"
	SecretParameter      ConfigParameterType = "Secret"
	EnvironmentParameter ConfigParameterType = "Environment"
)

type ConfigParameterName string

func New(awsSession *session.Session) *Config {
	return &Config{
		systemManager: ssm.New(awsSession),
		parameters:    map[ConfigParameterName]string{},
	}
}

func (s *Config) retrieveStandardParameters() error {
	paramNames := []*string{}
	for _, definition := range paramDefinitions {
		if definition.ParameterType == StandardParameter {
			paramNames = append(paramNames, aws.String(string(definition.ParameterName)))
		}
	}
	if len(paramNames) > 0 {
		paramOutput, err := s.systemManager.GetParameters(&ssm.GetParametersInput{
			Names:          paramNames,
			WithDecryption: aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("error while retrieving standard SSM parameters => %v", err.Error())
		}
		for _, param := range paramOutput.Parameters {
			if _, exists := s.parameters[ConfigParameterName(*param.Name)]; exists {
				return fmt.Errorf("duplicate parameter name: %s", string(*param.Name))
			}
			s.parameters[ConfigParameterName(*param.Name)] = *param.Value
		}
	}
	return nil
}

func (s *Config) retrieveSecretParameters() error {
	paramNames := []*string{}
	for _, definition := range paramDefinitions {
		if definition.ParameterType == SecretParameter {
			paramNames = append(paramNames, aws.String(string(definition.ParameterName)))
		}
	}
	if len(paramNames) > 0 {
		paramOutput, err := s.systemManager.GetParameters(&ssm.GetParametersInput{
			Names:          paramNames,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("error while retrieving secret SSM parameters => %v", err.Error())
		}
		for _, param := range paramOutput.Parameters {
			if _, exists := s.parameters[ConfigParameterName(*param.Name)]; exists {
				return fmt.Errorf("duplicate parameter name: %s", string(*param.Name))
			}
			s.parameters[ConfigParameterName(*param.Name)] = *param.Value
		}
	}
	return nil
}

func (s *Config) retrieveEnvironmentParameters() error {
	for _, definition := range paramDefinitions {
		if definition.ParameterType == EnvironmentParameter {
			if _, exists := s.parameters[definition.ParameterName]; exists {
				return fmt.Errorf("duplicate parameter name: %s", string(definition.ParameterName))
			}
			s.parameters[definition.ParameterName] = os.Getenv(string(definition.ParameterName))
		}
	}
	return nil
}

func (s *Config) initParameters() error {
	if err := s.retrieveStandardParameters(); err != nil {
		return err
	}
	if err := s.retrieveSecretParameters(); err != nil {
		return err
	}
	if err := s.retrieveEnvironmentParameters(); err != nil {
		return err
	}
	return nil
}

func (s *Config) InitializeParameters() error {
	return s.initParameters()
}
