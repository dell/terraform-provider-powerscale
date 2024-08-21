package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"terraform-provider-powerscale/client"
)

// GetS3GlobalSetting gets S3 GlobalSetting.
func GetS3GlobalSetting(ctx context.Context, client *client.Client) (*powerscale.V10S3SettingsGlobal, error) {
	param := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv10S3SettingsGlobal(ctx)
	response, _, err := param.Execute()
	return response, err
}

// UpdateS3GlobalSetting update s3 GlobalSetting.
func UpdateS3GlobalSetting(ctx context.Context, client *client.Client, GlobalSettingToUpdate powerscale.V10S3SettingsGlobalSettings) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv10S3SettingsGlobal(ctx)
	updateParam = updateParam.V10S3SettingsGlobal(GlobalSettingToUpdate)
	_, err := updateParam.Execute()
	return err
}
