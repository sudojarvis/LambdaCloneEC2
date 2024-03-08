package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/lambda"
)




func main() {
	
	sess, err := getSession()
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	lambdaSvc, err := getLambdaService(sess)
	if err != nil {
		fmt.Println("Error creating Lambda service:", err)
		return
	}

	config, err := lambdaSvc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String("hello"),  // hardcoded the function name
	})


	listFunctionOutput, err := lambdaSvc.ListFunctions(&lambda.ListFunctionsInput{})
	if err != nil {
		fmt.Println("Error listing functions:", err)
		return
	}

	fmt.Println("List of Lambda functions:")
	for _, function := range listFunctionOutput.Functions {
		fmt.Println("Name:", *function.FunctionName)
	}

	if err != nil {
		fmt.Println("Error getting function configuration:", err)
		return
	}

	fmt.Println(config)

	packageType := *config.Configuration.PackageType

	ec2Svc, err := getEC2Service(sess)

	if err != nil {
		fmt.Println("Error creating EC2 service:", err)
		return
	}

	// hardcoded the keyPairName

	
	keyPairName := "xyz"

	pemPath, err := createKeyPair(ec2Svc, keyPairName)



	if err != nil {
		fmt.Println("Error creating key pair:", err)
		return
	}

	// pemPath := keyPairName + ".pem"



	runParams := &ec2.RunInstancesInput{

		ImageId:      aws.String("ami-07d9b9ddc6cd8dd30"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      aws.String(keyPairName),
	}

	runResult, err := ec2Svc.RunInstances(runParams)	

	if err != nil {
		fmt.Println("Error running instance:", err)
		return
	}

	instanceID := *runResult.Instances[0].InstanceId
	fmt.Println("Created instance", instanceID)

	time.Sleep(40 * time.Second)

	err = modifyInstanceMetadataOptions(ec2Svc, instanceID)
	if err != nil {
		fmt.Println("Error modifying instance metadata options:", err)
		return
	}

	result, err := describeInstance(ec2Svc, instanceID)
	if err != nil {
		fmt.Println("Error describing instance:", err)
		return
	}

	publicDNS := *result.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicDnsName
	fmt.Println("Public DNS:", publicDNS)


	

	if packageType == "Zip" {

		fileName := *config.Configuration.FunctionName + ".zip"
		url := *config.Code.Location
		resp, err := http.Get(url)	
		if err != nil {
			fmt.Println("Error getting zip from URL:", err)
			return
		}

		defer resp.Body.Close()

		destFile, err := os.Create(fileName)

		if err != nil {
			fmt.Println("Error creating zip from URL:", err)
			return
		}

		defer destFile.Close()
		_, err = io.Copy(destFile, resp.Body)
		


		if err != nil {
			fmt.Println("Error copying zip from URL:", err)
		}
	
		
		/// hardcoded the keyPath, localFilePath, remoteUser, remoteHost, remoteFolderPath
		keyPath := pemPath
		localFilePath := fileName
		remoteUser := "ubuntu"
		remoteHost := publicDNS
		remoteFolderPath := ""


		// Create the command
		cmd := exec.Command("sh", "-c", fmt.Sprintf("chmod 400 %q && yes | scp -o 'StrictHostKeyChecking no' -i %q -r %q %s@%s:%s", keyPath, keyPath, localFilePath, remoteUser, remoteHost, remoteFolderPath))

		// Run the command
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}


		os.Remove( fileName )

		fmt.Println("EXecuted")
		// fmt.Println(" ")


	} else if packageType == "Image" {

		containerImgRef := *config.Code.ImageUri
		fmt.Println(containerImgRef)

		functionName := *config.Configuration.FunctionName

		// err := pullContainerImage(containerImgRef, functionName)
		pullContainerImage(containerImgRef,functionName)
		// if err != nil {
		// 	fmt.Println("Error pulling container image:", err)
		// 	return
		// }

		/// hardcoded the keyPath, localFilePath, remoteUser, remoteHost, remoteFolderPath
		keyPath := pemPath
		localFilePath :=  functionName + ".tar"
		remoteUser := "ubuntu"
		remoteHost := publicDNS
		remoteFolderPath := ""



		cmd := exec.Command("sh", "-c", fmt.Sprintf("chmod 400 %q && yes | scp -o 'StrictHostKeyChecking no' -i %q -r %q %s@%s:%s", keyPath, keyPath, localFilePath, remoteUser, remoteHost, remoteFolderPath))


		err = cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	


	

		if err := cmd.Wait(); err != nil {
			fmt.Println("Error waiting for command to finish:", err)
			return
		}
		

		os.Remove( functionName + ".tar" )

		fmt.Println("EXecuted")



	}

}





func getSession() (*session.Session, error) {
	awsConfigPath:= ".aws/config"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		SharedConfigFiles: []string{awsConfigPath}, 
	}))

	return sess, nil

}

// func getSession() (*session.Session, error) {
// 	return session.NewSession(&aws.Config{
// 		Region:      aws.String("us-east-1"),
// 		Credentials: credentials.NewSharedCredentials("", "default"),
		
// 	})
// }

func getLambdaService(sess *session.Session) (*lambda.Lambda, error) {
	return lambda.New(sess), nil
}


func getEC2Service(sess *session.Session) (*ec2.EC2, error) {
	return ec2.New(sess), nil
}





// func runInstance(svc *ec2.EC2, keyPairName string) (*ec2.RunInstancesOutput, error) {
// 	runParams := &ec2.RunInstancesInput{
// 		ImageId:      aws.String("ami-07d9b9ddc6cd8dd30"),
// 		InstanceType: aws.String("t2.micro"),
// 		MinCount:     aws.Int64(1),
// 		MaxCount:     aws.Int64(1),
// 		KeyName:      aws.String(keyPairName),
// 	}

// 	return svc.RunInstances(runParams)
// }


func describeInstance(svc *ec2.EC2, instanceID string) (*ec2.DescribeInstancesOutput, error) {
	return svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
}


func modifyInstanceMetadataOptions(svc *ec2.EC2, instanceID string) error {
	_, err := svc.ModifyInstanceMetadataOptions(&ec2.ModifyInstanceMetadataOptionsInput{
		DryRun:                  new(bool),
		HttpEndpoint:            aws.String("enabled"),
		HttpProtocolIpv6:        aws.String("disabled"),
		HttpPutResponseHopLimit: aws.Int64(1),
		InstanceId:              aws.String(instanceID),
	})
	return err
}
