package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// ListSubnets list subnet entities.
// API /platform/12/network/groupnets/{Groupnet}/subnets
// If filter {Groupnet} is specified, call API directly
// If not, use /platform/10/network/groupnets to query all groupnets and then call subnet API to query all
func ListSubnets(ctx context.Context, client *client.Client, subnetFilter *models.SubnetFilterType) (*[]powerscale.V12GroupnetSubnetExtended, error) {
	var subnetList []powerscale.V12GroupnetSubnetExtended
	if subnetFilter != nil && !subnetFilter.GroupnetName.IsNull() {
		networkSubnets, _, err := client.PscaleOpenAPIClient.NetworkGroupnetsApi.ListNetworkGroupnetsv12GroupnetSubnets(ctx, subnetFilter.GroupnetName.ValueString()).Execute()
		if err != nil {
			return nil, err
		}
		totalSubnets, err := ResumeSubnets(ctx, client, networkSubnets, subnetFilter.GroupnetName.ValueString())
		if err != nil {
			return totalSubnets, err
		}
		subnetList = *totalSubnets
	} else {
		networkGroupnets, _, err := client.PscaleOpenAPIClient.NetworkApi.ListNetworkv10NetworkGroupnets(ctx).Execute()
		if err != nil {
			return nil, err
		}

		for _, groupnet := range networkGroupnets.Groupnets {
			networkSubnets, _, err := client.PscaleOpenAPIClient.NetworkGroupnetsApi.ListNetworkGroupnetsv12GroupnetSubnets(ctx, *groupnet.Name).Execute()
			if err != nil {
				return nil, err
			}
			totalSubnets, err := ResumeSubnets(ctx, client, networkSubnets, *groupnet.Name)
			if err != nil {
				return &subnetList, err
			}
			subnetList = append(subnetList, *totalSubnets...)
		}
	}

	// Filter subnets based on names
	if subnetFilter != nil && subnetFilter.Names != nil && len(subnetFilter.Names) != 0 {
		var filteredSubnets []powerscale.V12GroupnetSubnetExtended
		for _, subnet := range subnetList {
			for _, name := range subnetFilter.Names {
				if subnet.Name != nil && *subnet.Name == name.ValueString() {
					filteredSubnets = append(filteredSubnets, subnet)
				}
			}
		}
		return &filteredSubnets, nil
	}

	return &subnetList, nil
}

// ResumeSubnets continue returning results from previous call using the resume token
func ResumeSubnets(ctx context.Context, client *client.Client, subnets *powerscale.V12GroupnetSubnets, groupnet string) (*[]powerscale.V12GroupnetSubnetExtended, error) {
	totalSubnets := subnets.Subnets
	for subnets.Resume != nil {
		subnets, _, err := client.PscaleOpenAPIClient.NetworkGroupnetsApi.ListNetworkGroupnetsv12GroupnetSubnets(ctx, groupnet).Resume(*subnets.Resume).Execute()
		if err != nil {
			return &totalSubnets, err
		}
		totalSubnets = append(totalSubnets, subnets.Subnets...)
	}
	return &totalSubnets, nil
}
