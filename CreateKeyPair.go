package main

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createKeyPair(svc *ec2.EC2, keyPairName string) (string, error) {
	
	describeKeyPairsInput := &ec2.DescribeKeyPairsInput{
		KeyNames: []*string{aws.String(keyPairName)},
	}
	_, err := svc.DescribeKeyPairs(describeKeyPairsInput)

	if err == nil {
		return keyPairName + ".pem", nil
	}


	createKeyPairInput := &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyPairName),
	}

	createKeyPairOutput, err := svc.CreateKeyPair(createKeyPairInput)
	if err != nil {
		return "", err
	}

	privateKey := *createKeyPairOutput.KeyMaterial
	pemPath := keyPairName + ".pem"
	err = ioutil.WriteFile(pemPath, []byte(privateKey), 0600)
	if err != nil {
		return "", err
	}

	return pemPath, nil
}