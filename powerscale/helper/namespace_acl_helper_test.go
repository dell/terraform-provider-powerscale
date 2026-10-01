/*
Copyright (c) 2026 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// buildACLEntry is a test helper that builds a types.Object representing an ACL entry.
func buildACLEntry(accesstype, op string, accessrights, inheritFlags []string, trusteeID, trusteeName, trusteeType string, unknownFields map[string]bool) basetypes.ObjectValue {
	attrs := make(map[string]attr.Value)

	// accesstype
	if unknownFields["accesstype"] {
		attrs["accesstype"] = types.StringUnknown()
	} else {
		attrs["accesstype"] = types.StringValue(accesstype)
	}

	// op
	if unknownFields["op"] {
		attrs["op"] = types.StringUnknown()
	} else if op == "" && !unknownFields["op_null"] {
		attrs["op"] = types.StringNull()
	} else if unknownFields["op_null"] {
		attrs["op"] = types.StringNull()
	} else {
		attrs["op"] = types.StringValue(op)
	}

	// accessrights
	if unknownFields["accessrights"] {
		attrs["accessrights"] = types.ListUnknown(types.StringType)
	} else {
		elems := make([]attr.Value, len(accessrights))
		for i, v := range accessrights {
			elems[i] = types.StringValue(v)
		}
		attrs["accessrights"], _ = types.ListValue(types.StringType, elems)
	}

	// inherit_flags
	if unknownFields["inherit_flags"] {
		attrs["inherit_flags"] = types.ListUnknown(types.StringType)
	} else {
		elems := make([]attr.Value, len(inheritFlags))
		for i, v := range inheritFlags {
			elems[i] = types.StringValue(v)
		}
		attrs["inherit_flags"], _ = types.ListValue(types.StringType, elems)
	}

	// trustee
	trusteeAttrs := make(map[string]attr.Value)
	if unknownFields["trustee.id"] {
		trusteeAttrs["id"] = types.StringUnknown()
	} else if trusteeID == "" {
		trusteeAttrs["id"] = types.StringNull()
	} else {
		trusteeAttrs["id"] = types.StringValue(trusteeID)
	}
	if unknownFields["trustee.name"] {
		trusteeAttrs["name"] = types.StringUnknown()
	} else if trusteeName == "" {
		trusteeAttrs["name"] = types.StringNull()
	} else {
		trusteeAttrs["name"] = types.StringValue(trusteeName)
	}
	if unknownFields["trustee.type"] {
		trusteeAttrs["type"] = types.StringUnknown()
	} else if trusteeType == "" {
		trusteeAttrs["type"] = types.StringNull()
	} else {
		trusteeAttrs["type"] = types.StringValue(trusteeType)
	}
	trusteeObj, _ := types.ObjectValue(aclTrusteeAttrTypes, trusteeAttrs)
	attrs["trustee"] = trusteeObj

	obj, _ := types.ObjectValue(aclEntryAttrTypes, attrs)
	return obj
}

func buildACLList(entries ...basetypes.ObjectValue) types.List {
	elems := make([]attr.Value, len(entries))
	for i, e := range entries {
		elems[i] = e
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: aclEntryAttrTypes}, elems)
	return list
}

// TestMergeCustomACLWithServerACL_PreservesUserAccessRights verifies that user-specified
// accessrights values are preserved even when the server returns canonical forms.
func TestMergeCustomACLWithServerACL_PreservesUserAccessRights(t *testing.T) {
	// User config: uses alias access rights "list", "traverse"
	userEntry := buildACLEntry(
		"allow", "",
		[]string{"list", "traverse", "dir_gen_write"},
		[]string{"container_inherit"},
		"UID:10", "", "",
		map[string]bool{"op": true, "trustee.name": true, "trustee.type": true},
	)

	// Server response: canonical forms "dir_gen_read", "dir_gen_execute" (OneFS rewrites aliases)
	serverEntry := buildACLEntry(
		"allow", "replace",
		[]string{"dir_gen_read", "dir_gen_execute", "dir_gen_write"},
		[]string{"container_inherit"},
		"UID:10", "root", "user",
		map[string]bool{},
	)

	userACL := buildACLList(userEntry)
	serverACL := buildACLList(serverEntry)

	merged, err := MergeCustomACLWithServerACL(userACL, serverACL)
	assert.NoError(t, err)

	mergedElements := merged.Elements()
	assert.Len(t, mergedElements, 1)

	mergedObj := mergedElements[0].(basetypes.ObjectValue)
	attrs := mergedObj.Attributes()

	// accessrights should be preserved from user config (NOT canonicalized)
	arList := attrs["accessrights"].(basetypes.ListValue)
	arElements := arList.Elements()
	assert.Len(t, arElements, 3)
	assert.Equal(t, "list", arElements[0].(basetypes.StringValue).ValueString())
	assert.Equal(t, "traverse", arElements[1].(basetypes.StringValue).ValueString())
	assert.Equal(t, "dir_gen_write", arElements[2].(basetypes.StringValue).ValueString())

	// inherit_flags should be preserved from user config
	ifList := attrs["inherit_flags"].(basetypes.ListValue)
	ifElements := ifList.Elements()
	assert.Len(t, ifElements, 1)
	assert.Equal(t, "container_inherit", ifElements[0].(basetypes.StringValue).ValueString())

	// op should be filled from server (was unknown in user config)
	opVal := attrs["op"].(basetypes.StringValue)
	assert.Equal(t, "replace", opVal.ValueString())

	// trustee.id should be preserved from user config
	trusteeObj := attrs["trustee"].(basetypes.ObjectValue)
	trusteeAttrs := trusteeObj.Attributes()
	assert.Equal(t, "UID:10", trusteeAttrs["id"].(basetypes.StringValue).ValueString())

	// trustee.name and trustee.type should be filled from server
	assert.Equal(t, "root", trusteeAttrs["name"].(basetypes.StringValue).ValueString())
	assert.Equal(t, "user", trusteeAttrs["type"].(basetypes.StringValue).ValueString())
}

// TestMergeCustomACLWithServerACL_PreservesInheritFlagsOrder verifies that
// user-specified inherit_flags ordering is preserved.
func TestMergeCustomACLWithServerACL_PreservesInheritFlagsOrder(t *testing.T) {
	// User specifies: container_inherit, object_inherit
	userEntry := buildACLEntry(
		"allow", "",
		[]string{"dir_gen_read"},
		[]string{"container_inherit", "object_inherit"},
		"UID:10", "", "",
		map[string]bool{"op": true, "trustee.name": true, "trustee.type": true},
	)

	// Server returns: object_inherit, container_inherit (different order)
	serverEntry := buildACLEntry(
		"allow", "replace",
		[]string{"dir_gen_read"},
		[]string{"object_inherit", "container_inherit"},
		"UID:10", "root", "user",
		map[string]bool{},
	)

	userACL := buildACLList(userEntry)
	serverACL := buildACLList(serverEntry)

	merged, err := MergeCustomACLWithServerACL(userACL, serverACL)
	assert.NoError(t, err)

	mergedObj := merged.Elements()[0].(basetypes.ObjectValue)
	ifList := mergedObj.Attributes()["inherit_flags"].(basetypes.ListValue)
	ifElements := ifList.Elements()

	// Should preserve user's order, not server's
	assert.Len(t, ifElements, 2)
	assert.Equal(t, "container_inherit", ifElements[0].(basetypes.StringValue).ValueString())
	assert.Equal(t, "object_inherit", ifElements[1].(basetypes.StringValue).ValueString())
}

// TestMergeCustomACLWithServerACL_NullUserACL verifies that when user ACL is null,
// the server ACL is returned as-is.
func TestMergeCustomACLWithServerACL_NullUserACL(t *testing.T) {
	serverEntry := buildACLEntry(
		"allow", "replace",
		[]string{"dir_gen_read"},
		[]string{"container_inherit"},
		"UID:10", "root", "user",
		map[string]bool{},
	)
	serverACL := buildACLList(serverEntry)
	nullACL := types.ListNull(types.ObjectType{AttrTypes: aclEntryAttrTypes})

	merged, err := MergeCustomACLWithServerACL(nullACL, serverACL)
	assert.NoError(t, err)
	assert.True(t, merged.Equal(serverACL))
}

// TestMergeCustomACLWithServerACL_UnknownUserACL verifies that when user ACL is unknown,
// the server ACL is returned as-is.
func TestMergeCustomACLWithServerACL_UnknownUserACL(t *testing.T) {
	serverEntry := buildACLEntry(
		"allow", "replace",
		[]string{"dir_gen_read"},
		[]string{"container_inherit"},
		"UID:10", "root", "user",
		map[string]bool{},
	)
	serverACL := buildACLList(serverEntry)
	unknownACL := types.ListUnknown(types.ObjectType{AttrTypes: aclEntryAttrTypes})

	merged, err := MergeCustomACLWithServerACL(unknownACL, serverACL)
	assert.NoError(t, err)
	assert.True(t, merged.Equal(serverACL))
}

// TestMergeCustomACLWithServerACL_MultipleACEs verifies correct merging with multiple ACE entries.
func TestMergeCustomACLWithServerACL_MultipleACEs(t *testing.T) {
	userEntry1 := buildACLEntry(
		"allow", "",
		[]string{"list", "traverse"},
		[]string{"container_inherit"},
		"UID:10", "", "",
		map[string]bool{"op": true, "trustee.name": true, "trustee.type": true},
	)
	userEntry2 := buildACLEntry(
		"deny", "",
		[]string{"dir_gen_write"},
		[]string{},
		"", "Everyone", "wellknown",
		map[string]bool{"op": true, "trustee.id": true},
	)

	serverEntry1 := buildACLEntry(
		"allow", "replace",
		[]string{"dir_gen_read", "dir_gen_execute"},
		[]string{"container_inherit"},
		"UID:10", "root", "user",
		map[string]bool{},
	)
	serverEntry2 := buildACLEntry(
		"deny", "replace",
		[]string{"dir_gen_write"},
		[]string{},
		"SID:S-1-1-0", "Everyone", "wellknown",
		map[string]bool{},
	)

	userACL := buildACLList(userEntry1, userEntry2)
	serverACL := buildACLList(serverEntry1, serverEntry2)

	merged, err := MergeCustomACLWithServerACL(userACL, serverACL)
	assert.NoError(t, err)

	mergedElements := merged.Elements()
	assert.Len(t, mergedElements, 2)

	// First ACE: accessrights preserved from user ("list", "traverse")
	ace1 := mergedElements[0].(basetypes.ObjectValue).Attributes()
	ar1 := ace1["accessrights"].(basetypes.ListValue).Elements()
	assert.Len(t, ar1, 2)
	assert.Equal(t, "list", ar1[0].(basetypes.StringValue).ValueString())
	assert.Equal(t, "traverse", ar1[1].(basetypes.StringValue).ValueString())

	// First ACE: trustee.name filled from server
	trustee1 := ace1["trustee"].(basetypes.ObjectValue).Attributes()
	assert.Equal(t, "UID:10", trustee1["id"].(basetypes.StringValue).ValueString())
	assert.Equal(t, "root", trustee1["name"].(basetypes.StringValue).ValueString())

	// Second ACE: trustee.id filled from server (was unknown in user config)
	ace2 := mergedElements[1].(basetypes.ObjectValue).Attributes()
	trustee2 := ace2["trustee"].(basetypes.ObjectValue).Attributes()
	assert.Equal(t, "SID:S-1-1-0", trustee2["id"].(basetypes.StringValue).ValueString())
	assert.Equal(t, "Everyone", trustee2["name"].(basetypes.StringValue).ValueString())
	assert.Equal(t, "wellknown", trustee2["type"].(basetypes.StringValue).ValueString())
}

// TestMergeCustomACLWithServerACL_AllFieldsSpecified verifies that when the user
// specifies all fields, none are overwritten by the server response.
func TestMergeCustomACLWithServerACL_AllFieldsSpecified(t *testing.T) {
	userEntry := buildACLEntry(
		"allow", "replace",
		[]string{"list", "traverse"},
		[]string{"container_inherit", "object_inherit"},
		"UID:10", "myuser", "user",
		map[string]bool{},
	)
	serverEntry := buildACLEntry(
		"allow", "replace",
		[]string{"dir_gen_read", "dir_gen_execute"},
		[]string{"object_inherit", "container_inherit"},
		"UID:10", "root", "user",
		map[string]bool{},
	)

	userACL := buildACLList(userEntry)
	serverACL := buildACLList(serverEntry)

	merged, err := MergeCustomACLWithServerACL(userACL, serverACL)
	assert.NoError(t, err)

	mergedObj := merged.Elements()[0].(basetypes.ObjectValue)
	attrs := mergedObj.Attributes()

	// All values should come from user, not server
	ar := attrs["accessrights"].(basetypes.ListValue).Elements()
	assert.Equal(t, "list", ar[0].(basetypes.StringValue).ValueString())
	assert.Equal(t, "traverse", ar[1].(basetypes.StringValue).ValueString())

	trustee := attrs["trustee"].(basetypes.ObjectValue).Attributes()
	assert.Equal(t, "myuser", trustee["name"].(basetypes.StringValue).ValueString())
}
