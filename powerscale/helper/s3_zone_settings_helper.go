package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// GetS3ZoneSetting gets S3 ZoneSetting.
func GetS3ZoneSetting(ctx context.Context, client *client.Client, state models.S3ZoneSettingsResource) (*powerscale.V10S3SettingsZone, error) {
	param := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv10S3SettingsZone(ctx)
	if len(state.Zone.ValueString()) > 0 {
		param = param.Zone(state.Zone.ValueString())
	}
	response, _, err := param.Execute()
	return response, err
}

// UpdateS3ZoneSetting update s3 ZoneSetting.
func UpdateS3ZoneSetting(ctx context.Context, client *client.Client, state models.S3ZoneSettingsResource) error {
	var toUpdate powerscale.V10S3SettingsZoneSettings
	err := ReadFromState(ctx, &state, &toUpdate)
	if err != nil {
		return err
	}
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv10S3SettingsZone(ctx)
	if len(state.Zone.ValueString()) > 0 {
		updateParam = updateParam.Zone(state.Zone.ValueString())
	}
	updateParam = updateParam.V10S3SettingsZone(toUpdate)
	_, err = updateParam.Execute()
	return err
}

// SetZoneSetting updates the S3 Zone Settings.
func SetZoneSetting(ctx context.Context, client *client.Client, s3ZSPlan models.S3ZoneSettingsResource) (models.S3ZoneSettingsResource, error) {
	err := UpdateS3ZoneSetting(ctx, client, s3ZSPlan)
	if err != nil {
		return models.S3ZoneSettingsResource{}, err
	}
	ZoneSettings, err := GetS3ZoneSetting(ctx, client, s3ZSPlan)
	if err != nil {
		return models.S3ZoneSettingsResource{}, err
	}
	var state models.S3ZoneSettingsResource
	err = CopyFieldsToNonNestedModel(ctx, ZoneSettings.GetSettings(), &state)
	if err != nil {
		return models.S3ZoneSettingsResource{}, err
	}
	state.Zone = s3ZSPlan.Zone
	return state, nil
}

// GetZoneSetting reads the S3 Zone Settings.
func GetZoneSetting(ctx context.Context, client *client.Client, s3ZoneSettingState models.S3ZoneSettingsResource) error {
	ZoneSettings, err := GetS3ZoneSetting(ctx, client, s3ZoneSettingState)
	if err != nil {
		return err
	}
	err = CopyFieldsToNonNestedModel(ctx, ZoneSettings.GetSettings(), &s3ZoneSettingState)
	if err != nil {
		return err
	}
	return nil
}
