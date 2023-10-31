/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"math"
	"math/big"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SmartPoolSettings The interface of various StoragepoolSettingsSettings.
type SmartPoolSettings interface {
	GetAutomaticallyManageIoOptimization() string
	GetAutomaticallyManageProtection() string
}

// SettingsSetter The interface has the wrapper methods to set ManageIoOptimization and  ManageProtection for smartpool
// settings.
type SettingsSetter interface {
	SetManageIoOptimization(types.Bool)
	SetManageIoOptimizationApplyToFiles(types.Bool)
	SetManageProtection(types.Bool)
	SetManageProtectionApplyToFiles(types.Bool)
}

// GetSmartPoolSettings Get SmartPool settings based on Onefs version.
func GetSmartPoolSettings(ctx context.Context, powerscaleClient *client.Client) (any, error) {
	if powerscaleClient.OnefsVersion.IsGreaterThan("9.4.0") {
		settings, _, err := powerscaleClient.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv16StoragepoolSettings(ctx).Execute()
		return settings, err
	}
	settings, _, err := powerscaleClient.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv5StoragepoolSettings(ctx).Execute()
	return settings, err
}

// customizeFields build ManageIoOptimization and ManageProtection related values for Terraform based on PowerScale settings.
func customizeFields(sts SmartPoolSettings, setter SettingsSetter) {
	// Compute the following values that align with OneFS UI
	if sts.GetAutomaticallyManageIoOptimization() == "none" {
		setter.SetManageIoOptimization(types.BoolValue(false))
		setter.SetManageIoOptimizationApplyToFiles(types.BoolValue(false))
	} else if sts.GetAutomaticallyManageIoOptimization() == "files_at_default" {
		setter.SetManageIoOptimization(types.BoolValue(true))
		setter.SetManageIoOptimizationApplyToFiles(types.BoolValue(false))
	} else if sts.GetAutomaticallyManageIoOptimization() == "all" {
		setter.SetManageIoOptimization(types.BoolValue(true))
		setter.SetManageIoOptimizationApplyToFiles(types.BoolValue(true))
	}

	if sts.GetAutomaticallyManageProtection() == "none" {
		setter.SetManageProtection(types.BoolValue(false))
		setter.SetManageProtectionApplyToFiles(types.BoolValue(false))
	} else if sts.GetAutomaticallyManageProtection() == "files_at_default" {
		setter.SetManageProtection(types.BoolValue(true))
		setter.SetManageProtectionApplyToFiles(types.BoolValue(false))
	} else if sts.GetAutomaticallyManageProtection() == "all" {
		setter.SetManageProtection(types.BoolValue(true))
		setter.SetManageProtectionApplyToFiles(types.BoolValue(true))
	}
}

// UpdateSmartPoolSettingsDatasourceModel update SmartPool Settings datasource model.
func UpdateSmartPoolSettingsDatasourceModel(ctx context.Context, settings any, model *models.SmartPoolSettingsDataSource) (string, string) {
	var sts SmartPoolSettings

	switch v := settings.(type) {
	case *powerscale.V16StoragepoolSettings:
		s, ok := settings.(*powerscale.V16StoragepoolSettings)
		if !ok {
			return "Error parameter SmartPoolSettings.", fmt.Sprintf("Unexpected type %T", v)
		}

		sts, _ = s.GetSettingsOk()
		err := CopyFields(ctx, s.Settings, model)
		if err != nil {
			return "Failed to map SmartPool settings fields", err.Error()
		}

		// default_transfer_limit_state and default_transfer_limit_pct are mutually exclusive
		// make sure all fields are in known state
		if s.Settings.DefaultTransferLimitState == nil {
			model.DefaultTransferLimitState = types.StringNull()
		}
		if s.Settings.DefaultTransferLimitPct == nil {
			model.DefaultTransferLimitPct = types.NumberNull()
		}
	case *powerscale.V5StoragepoolSettings:
		s, ok := settings.(*powerscale.V5StoragepoolSettings)
		if !ok {
			return "Error pararmeter SmartPoolSettings.", fmt.Sprintf("Unexpected type %T", v)
		}

		sts, _ = s.GetSettingsOk()
		err := CopyFields(ctx, s.Settings, model)
		if err != nil {
			return "Failed to map SmartPool settings fields", err.Error()
		}

		// default_transfer_limit_state and default_transfer_limit_pct are available in 9.5
		model.DefaultTransferLimitState = types.StringNull()
		model.DefaultTransferLimitPct = types.NumberNull()
	default:
		tflog.Error(ctx, fmt.Sprintf("Unknown type %T", v))
		return "Failed to parse SmartPoolSettings.", fmt.Sprintf("Unknown type %T", v)
	}

	model.ID = types.StringValue("smartpools_settings")
	customizeFields(sts, model)
	return "", ""
}

// UpdateSmartPoolSettingsResourceModel update SmartPool Settings resource model.
func UpdateSmartPoolSettingsResourceModel(ctx context.Context, settings any, model *models.SmartPoolSettingsResource) (string, string) {
	var sts SmartPoolSettings

	switch v := settings.(type) {
	case *powerscale.V16StoragepoolSettings:
		s, ok := settings.(*powerscale.V16StoragepoolSettings)
		if !ok {
			return "Error pararmeter SmartPoolSettings.", fmt.Sprintf("Unexpected type %T", v)
		}

		sts, _ = s.GetSettingsOk()
		err := CopyFields(ctx, s.Settings, model)
		if err != nil {
			return "Failed to map SmartPool settings fields", err.Error()
		}

		// default_transfer_limit_state and default_transfer_limit_pct are mutually exclusive
		// make sure all fields are in known state
		if s.Settings.DefaultTransferLimitState == nil {
			model.DefaultTransferLimitState = types.StringNull()
		}
		if s.Settings.DefaultTransferLimitPct == nil {
			model.DefaultTransferLimitPct = types.NumberNull()
		}

		resolveSpilloverTarget(model, s.Settings.SpilloverTarget.Type, s.Settings.SpilloverTarget.Name)

	case *powerscale.V5StoragepoolSettings:
		s, ok := settings.(*powerscale.V5StoragepoolSettings)
		if !ok {
			return "Error pararmeter SmartPoolSettings.", fmt.Sprintf("Unexpected type %T", v)
		}

		sts, _ = s.GetSettingsOk()
		err := CopyFields(ctx, s.Settings, model)
		if err != nil {
			return "Failed to map SmartPool settings fields", err.Error()
		}
		// default_transfer_limit_state and default_transfer_limit_pct are available in 9.5
		model.DefaultTransferLimitState = types.StringNull()
		model.DefaultTransferLimitPct = types.NumberNull()
		resolveSpilloverTarget(model, s.Settings.SpilloverTarget.Type, s.Settings.SpilloverTarget.Name)
	default:
		tflog.Error(ctx, fmt.Sprintf("Unknown type %T", v))
		return "Failed to parse SmartPoolSettings.", fmt.Sprintf("Unknown type %T", v)
	}

	model.ID = types.StringValue("smartpools_settings")
	customizeFields(sts, model)
	return "", ""
}

func buildTargetName(model *models.SmartPoolSettingsResource) (*string, string) {
	var nameVal *string
	var typeVal string
	if !model.SpilloverTarget.IsUnknown() && !model.SpilloverTarget.IsNull() {
		if n, ok := model.SpilloverTarget.Attributes()["name"]; ok {
			stringVal, ok := n.(basetypes.StringValue)
			if ok && !stringVal.IsNull() && !stringVal.IsUnknown() {
				nameStr := stringVal.ValueString()
				nameVal = &nameStr
			}
		}
		if n, ok := model.SpilloverTarget.Attributes()["type"]; ok {
			stringVal, ok := n.(basetypes.StringValue)
			if ok && !stringVal.IsNull() && !stringVal.IsUnknown() {
				typeVal = stringVal.ValueString()
			}
		}
	}
	return nameVal, typeVal
}

func buildAutomaticallyManageValues(model *models.SmartPoolSettingsResource) (*string, *string) {
	var manageIoOptimization, manageProtection *string
	if !model.ManageIoOptimization.IsUnknown() && !model.ManageIoOptimization.IsNull() {
		ioOptimization := "all"
		if !model.ManageIoOptimization.ValueBool() {
			ioOptimization = "none"
		} else if model.ManageIoOptimizationApplyToFiles.IsUnknown() ||
			model.ManageIoOptimizationApplyToFiles.IsNull() ||
			!model.ManageIoOptimizationApplyToFiles.ValueBool() {
			ioOptimization = "files_at_default"
		}
		manageIoOptimization = &ioOptimization
	}

	if !model.ManageProtection.IsUnknown() && !model.ManageProtection.IsNull() {
		protection := "all"
		if !model.ManageProtection.ValueBool() {
			protection = "none"
		} else if model.ManageProtectionApplyToFiles.IsUnknown() ||
			model.ManageProtectionApplyToFiles.IsNull() ||
			!model.ManageProtectionApplyToFiles.ValueBool() {
			protection = "files_at_default"
		}
		manageProtection = &protection
	}
	return manageIoOptimization, manageProtection
}

// BigFloatToInt32 converts big.float value to int32.
func BigFloatToInt32(x *big.Float) (*int32, error) {
	if x == nil {
		return nil, nil
	}

	v64, accuracy := x.Int64()
	if accuracy == big.Below || accuracy == big.Above {
		return nil, fmt.Errorf("error converting value %v to int", x)
	}

	if v64 < math.MinInt32 || v64 > math.MaxInt32 {
		return nil, fmt.Errorf("error converting value %v to int", x)
	}
	ret := int32(v64)
	return &ret, nil
}

// UpdateSmartPoolSettings apply SmartPool Settings changes on PowerScale.
func UpdateSmartPoolSettings(ctx context.Context, client *client.Client, model *models.SmartPoolSettingsResource) error {
	if client.OnefsVersion.IsGreaterThan("9.4.0") {
		updateParam := client.PscaleOpenAPIClient.StoragepoolApi.UpdateStoragepoolv16StoragepoolSettings(ctx)
		settings := powerscale.V16StoragepoolSettingsExtended{}

		err := ReadFromState(ctx, model, &settings)
		if err != nil {
			return err
		}

		limitPct32, err := BigFloatToInt32(model.DefaultTransferLimitPct.ValueBigFloat())
		if err != nil {
			return err
		}
		settings.DefaultTransferLimitPct = limitPct32

		manageIoOptimization, manageProtection := buildAutomaticallyManageValues(model)
		settings.AutomaticallyManageIoOptimization = manageIoOptimization
		settings.AutomaticallyManageProtection = manageProtection

		if settings.SpilloverTarget != nil {
			nameVal, typeVal := buildTargetName(model)
			settings.SpilloverTarget.Type = typeVal
			settings.SpilloverTarget.NameOrId = nameVal
		}

		updateParam = updateParam.V16StoragepoolSettings(settings)
		_, err = updateParam.Execute()
		return err
	}

	// for PowerScale 9.4
	updateParam := client.PscaleOpenAPIClient.StoragepoolApi.UpdateStoragepoolv5StoragepoolSettings(ctx)
	settings := powerscale.V5StoragepoolSettingsExtended{}

	err := ReadFromState(ctx, model, &settings)
	if err != nil {
		return err
	}

	manageIoOptimization, manageProtection := buildAutomaticallyManageValues(model)
	settings.AutomaticallyManageIoOptimization = manageIoOptimization
	settings.AutomaticallyManageProtection = manageProtection

	if settings.SpilloverTarget != nil {
		nameVal, typeVal := buildTargetName(model)
		settings.SpilloverTarget.Type = typeVal
		settings.SpilloverTarget.NameOrId = nameVal
	}

	updateParam = updateParam.V5StoragepoolSettings(settings)
	_, err = updateParam.Execute()
	return err
}

// resolveSpilloverTarget assign values for Terraform SpilloverTarget model.
func resolveSpilloverTarget(model *models.SmartPoolSettingsResource, targetType, targetName string) {
	spilloverTargetType := map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
	}
	spilloverTargetMap := make(map[string]attr.Value)
	spilloverTargetMap["name"] = types.StringValue(targetName)
	// targetType may be tier, nodepool, anywhere, invalid from GET endpoint, but for PUT, only storagepool or anywhere is valid.
	spilloverTargetMap["type"] = types.StringValue(targetType)
	if targetType != "anywhere" && targetType != "invalid" {
		spilloverTargetMap["type"] = types.StringValue("storagepool")
	}

	spilloverTargetObject, _ := types.ObjectValue(spilloverTargetType, spilloverTargetMap)
	model.SpilloverTarget = spilloverTargetObject
}
