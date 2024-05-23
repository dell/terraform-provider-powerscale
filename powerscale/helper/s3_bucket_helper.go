/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// ListS3Buckets list S3 Bucket entities.
func ListS3Buckets(ctx context.Context, client *client.Client, bucketFilter *models.S3BucketDatasourceFilter) ([]powerscale.V12S3Bucket, error) {
	listS3BucketParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv12S3Buckets(ctx)
	if bucketFilter != nil {
		if !bucketFilter.Zone.IsNull() {
			listS3BucketParam = listS3BucketParam.Zone(bucketFilter.Zone.ValueString())
		}
		if !bucketFilter.Owner.IsNull() {
			listS3BucketParam = listS3BucketParam.Owner(bucketFilter.Owner.ValueString())
		}
	}
	S3BucketResponse, _, err := listS3BucketParam.Execute()
	if err != nil {
		return nil, err
	}
	totalS3Buckets := S3BucketResponse.Buckets
	for S3BucketResponse.Resume != nil {
		resumeS3BucketParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv12S3Buckets(ctx).Resume(*S3BucketResponse.Resume)
		S3BucketResponse, _, err = resumeS3BucketParam.Execute()
		if err != nil {
			return totalS3Buckets, err
		}
		totalS3Buckets = append(totalS3Buckets, S3BucketResponse.Buckets...)
	}
	return totalS3Buckets, nil
}

// CreateS3Bucket create s3 bucket.
func CreateS3Bucket(ctx context.Context, client *client.Client, bucket powerscale.V10S3Bucket, zone string) (*powerscale.CreateResponse, error) {
	if !bucket.HasAcl() {
		bucket.SetAcl(make([]powerscale.V10S3BucketAclItem, 0))
	}
	param := client.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv12S3Bucket(ctx).V12S3Bucket(bucket)
	if len(zone) > 0 {
		param = param.Zone(zone)
	}
	response, _, err := param.Execute()
	return response, err
}

// GetS3Bucket gets S3 Bucket.
func GetS3Bucket(ctx context.Context, client *client.Client, bucketID string, zone string) (*powerscale.V12S3BucketsExtended, error) {
	param := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv12S3Bucket(ctx, bucketID)
	if len(zone) > 0 {
		param = param.Zone(zone)
	}
	response, _, err := param.Execute()
	return response, err
}

// UpdateS3Bucket update s3 bucket.
func UpdateS3Bucket(ctx context.Context, client *client.Client, bucketID string, zone string, bucketToUpdate powerscale.V10S3BucketExtendedExtended) error {
	if !bucketToUpdate.HasAcl() {
		bucketToUpdate.SetAcl(make([]powerscale.V10S3BucketAclItem, 0))
	}
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv12S3Bucket(ctx, bucketID).V12S3Bucket(bucketToUpdate)
	if len(zone) > 0 {
		updateParam = updateParam.Zone(zone)
	}
	_, err := updateParam.Execute()
	return err
}

// DeleteS3Bucket delete s3 bucket.
func DeleteS3Bucket(ctx context.Context, client *client.Client, bucketID string, zone string) error {
	param := client.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv12S3Bucket(ctx, bucketID)
	if len(zone) > 0 {
		param = param.Zone(zone)
	}
	_, err := param.Execute()
	return err
}

// ValidateS3BucketUpdate validates if update params contain params only for creating.
func ValidateS3BucketUpdate(plan models.S3BucketResource, state models.S3BucketResource) error {
	if !plan.Zone.Equal(state.Zone) &&
		!((plan.Zone.ValueString() == "System" && len(state.Zone.ValueString()) == 0) ||
			(state.Zone.ValueString() == "System" && len(plan.Zone.ValueString()) == 0)) {
		return fmt.Errorf("do not update field Zone")
	}

	if !plan.Name.Equal(state.Name) {
		return fmt.Errorf("do not update field Name")
	}

	if !plan.Path.Equal(state.Path) {
		return fmt.Errorf("do not update field Path")
	}

	return nil
}
