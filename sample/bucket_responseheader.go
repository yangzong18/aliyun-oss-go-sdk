package sample

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// BucketResponseHeaderSample shows how to set, get and delete the bucket's response header.
func BucketResponseHeaderSample() {
	// New client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		HandleError(err)
	}

	// Create the bucket with default parameters
	err = client.CreateBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Set bucket's response header.
	responseHeader := oss.PutBucketResponseHeader{
		Rule:[]oss.BucketResponseHeaderRule{
			{
				Name: "rule1",
				Filters: []string{
					"Put*","GetObject",
				},
				HideHeaders: "Filters",
				ReplaceHeaders: oss.BucketResponseHeaders{
					Header: "Content-Type",
					Value: "type",
				},
			},
			{
				Name: "rule2",
				Filters: []string{
					"*",
				},
				HideHeaders: "Filters",
				AddHeaders: oss.BucketResponseHeaders{
					Header: "a",
					Value: "b",
				},
			},
		},
	}
	err = client.PutBucketResponseHeader(bucketName, responseHeader)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("Bucket Response Header Set Success!")

	// Get bucket's response header.
	header, err := client.GetBucketResponseHeader(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Printf("header:%#v\n",header)


	// Delete bucket's response header.
	err = client.DeleteBucketEncryption(bucketName)
	fmt.Println("Bucket Response Header Delete Success!")

	fmt.Println("BucketResponseHeaderSample completed")
}
