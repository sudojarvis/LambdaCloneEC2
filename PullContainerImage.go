package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/client"
)

func pullContainerImage(containerImgRef string, functionName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// res, err := cli.ImagePull(ctx, containerImgRef, pullopt)
	// if err != nil {
	// 	fmt.Println("Error pulling image:", err)
	// 	return
	// }
	// fmt.Println(res)

	fileName := functionName + ".tar"
	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader, err := cli.ImageSave(ctx, []string{containerImgRef})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		return err
	}

	// // Get the container ID from the containerImgRef
	// resp, err := cli.ContainerInspect(ctx, containerImgRef)
	// if err != nil {
	// 	return err
	// }
	// containerID := resp.ID

	// // Copy the entire contents of the container's filesystem
	// archiveReader, _, err := cli.CopyFromContainer(ctx, containerID, "/")
	// if err != nil {
	// 	return err
	// }
	// defer archiveReader.Close()

	// // Copy the contents from the archiveReader to the output file
	// _, err = io.Copy(outputFile, archiveReader)
	// if err != nil {
	// 	return err
	// }

	return nil
}




















// func pullContainerImage(containerImgRef string, functionName string) {
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		fmt.Println("Error creating Docker client:", err)
// 		return
// 	}

// 	sess, err := session.NewSession(&aws.Config{
// 		Region:      aws.String("us-east-1"),
// 		Credentials: credentials.NewSharedCredentials("", "default"),
// 	})

// 	if err != nil {
// 		fmt.Println("Error creating session")
// 	}

	
	

// 	// Create a new ECR client using the provided EC2 session
// 	ecrSvc := ecr.New(sess)



// 	// Get the authorization token from ECR
// 	autResp, err := ecrSvc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
// 		RegistryIds: []*string{aws.String("536619883796")}, // Replace with your AWS account ID
// 	})
// 	if err != nil {
// 		fmt.Println("Error getting authorization token:", err)
// 		return
// 	}

// 	// Extract authorization token data
// 	autRespData := autResp.AuthorizationData[0]
// 	authToken := *autRespData.AuthorizationToken

// 	// Decode the authorization token
// 	decodedToken, err := base64.StdEncoding.DecodeString(authToken)
// 	if err != nil {
// 		fmt.Println("Error decoding token:", err)
// 		return
// 	}

// 	// Split the decoded token into username and password
// 	tokenParts := strings.SplitN(string(decodedToken), ":", 2)
// 	if len(tokenParts) != 2 {
// 		fmt.Println("Invalid token format")
// 		return
// 	}
// 	username := tokenParts[0]
// 	password := tokenParts[1]

// 	// Prepare authentication configuration
// 	authConfig := registry.AuthConfig{
// 		Username:      username,
// 		Password:      password,
// 		ServerAddress: "536619883796.dkr.ecr.us-east-1.amazonaws.com",
// 	}

// 	// Encode authentication configuration
// 	authConfigBytes, err := json.Marshal(authConfig)
// 	if err != nil {
// 		fmt.Println("Error marshalling authConfig:", err)
// 		return
// 	}
// 	encodedAuthConfig := base64.URLEncoding.EncodeToString(authConfigBytes)

// 	// Options for pulling the image
// 	pullOptions := types.ImagePullOptions{
// 		RegistryAuth: encodedAuthConfig,
// 	}

// 	// Pull the container image
// 	ctx := context.Background()
// 	res, err := cli.ImagePull(ctx, containerImgRef, pullOptions)
// 	if err != nil {
// 		fmt.Println("Error pulling image:", err)
// 		return
// 	}
// 	defer res.Close()

// 	// Read the response from the pull operation
// 	if _, err := io.Copy(ioutil.Discard, res); err != nil {
// 		fmt.Println("Error reading response:", err)
// 		return
// 	}

// 	// Create the container
// 	info, err := cli.ContainerCreate(ctx, &container.Config{
// 		Image: containerImgRef,
// 	}, nil, nil, nil, functionName)
// 	if err != nil {
// 		fmt.Println("Error creating container:", err)
// 		return
// 	}

// 	// Create a TAR file to save the container's filesystem
// 	filepath := functionName + ".tar"
// 	file, err := os.Create(filepath)
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Copy the entire contents of the container's filesystem to the TAR file
// 	tarStream, _, err := cli.CopyFromContainer(ctx, info.ID, "/home")
// 	if err != nil {
// 		fmt.Println("Error copying from container:", err)
// 		return
// 	}
// 	defer tarStream.Close()

// 	// Copy the contents from the tarStream to the output file
// 	if _, err := io.Copy(file, tarStream); err != nil {
// 		fmt.Println("Error copying tar stream:", err)
// 		return
// 	}

// 	fmt.Println("Container filesystem copied to", filepath)
// }
