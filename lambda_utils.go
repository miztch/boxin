package main

import "os"

// Returns this function is running on AWS Lambda or not
func isRunningOnLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
