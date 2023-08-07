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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"math/big"
	"reflect"
	"strings"
)

// CopyFieldsToNonNestedModel copy OpenAPI struct source to destination of struct with terraform types.
// use this function when model struct contains only types.List/Object.
func CopyFieldsToNonNestedModel(ctx context.Context, source, destination interface{}) error {
	tflog.Debug(ctx, "Copy fields", map[string]interface{}{
		"source":      source,
		"destination": destination,
	})
	var err error
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination)

	// Check if destination is a pointer to a struct
	if destinationValue.Kind() != reflect.Ptr || destinationValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination is not a pointer to a struct")
	}

	// if source is a pointer, use the Elem() method to get the value that the pointer points to
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("source is not a struct")
	}

	// Get the type of the destination struct
	//destinationType := destinationValue.Elem().Type()
	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldTag := getFieldJSONTag(sourceValue, i)

		tflog.Debug(ctx, "Converting source field", map[string]interface{}{
			"sourceFieldTag":  sourceFieldTag,
			"sourceFieldKind": sourceValue.Field(i).Kind().String(),
		})

		sourceField := sourceValue.Field(i)
		if sourceField.Kind() == reflect.Ptr {
			sourceField = sourceField.Elem()
		}
		destinationField := getFieldByTfTag(destinationValue.Elem(), sourceFieldTag)
		structType := reflect.TypeOf(source)
		// For zero value (nil), the object still need to pass type information into it
		if !sourceField.IsValid() {
			destinationField = getFieldByTfTag(destinationValue.Elem(), sourceFieldTag)
			mapType, err := getStructAttrTypeFromType(ctx, structType.Field(i).Type)
			if err != nil {
				return err
			}
			destinationField.Set(reflect.ValueOf(types.ObjectNull(mapType)))
			continue
		}
		if destinationField.IsValid() && destinationField.CanSet() {
			tflog.Debug(ctx, "debugging source field", map[string]interface{}{
				"sourceField Interface": sourceField.Interface(),
			})
			// Convert the source value to the type of the destination field dynamically
			var destinationFieldValue attr.Value

			switch sourceField.Kind() {
			case reflect.String:
				destinationFieldValue = types.StringValue(sourceField.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				destinationFieldValue = types.Int64Value(sourceField.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				destinationFieldValue = types.Int64Value(sourceField.Int())
			case reflect.Float32, reflect.Float64:
				//destinationFieldValue = types.Float64Value(sourceField.Float())
				destinationFieldValue = types.NumberValue(big.NewFloat(sourceField.Float()))
			case reflect.Bool:
				destinationFieldValue = types.BoolValue(sourceField.Bool())
			case reflect.Array, reflect.Slice:
				destinationFieldValue, err = getSliceAttrValue(ctx, sourceField.Interface())
				if err != nil {
					return err
				}
			case reflect.Struct:
				destinationFieldValue, err = getStructValue(ctx, sourceField.Interface())
				if err != nil {
					return err
				}
			default:
				tflog.Error(ctx, "unsupported source field type", map[string]interface{}{
					"sourceField": sourceField,
				})
				continue
			}
			if destinationField.Type() == reflect.TypeOf(destinationFieldValue) {
				destinationField.Set(reflect.ValueOf(destinationFieldValue))
			}
		}
	}

	return nil
}

func getStructAttrTypeFromType(ctx context.Context, structType reflect.Type) (map[string]attr.Type, error) {
	attrTypeMap := make(map[string]attr.Type)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	for fieldIndex := 0; fieldIndex < structType.NumField(); fieldIndex++ {
		structField := structType.Field(fieldIndex)
		tag := structField.Tag.Get("json")
		tag = strings.TrimSuffix(tag, ",omitempty")
		structFieldKind := structField.Type.Kind()
		if structField.Type.Kind() == reflect.Ptr {
			structFieldKind = structField.Type.Elem().Kind()
		}
		switch structFieldKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			attrTypeMap[tag] = types.Int64Type
		case reflect.String:
			attrTypeMap[tag] = types.StringType
		case reflect.Float32, reflect.Float64:
			attrTypeMap[tag] = types.NumberType
		case reflect.Bool:
			attrTypeMap[tag] = types.BoolType
		case reflect.Struct:
			if structField.Type == reflect.TypeOf(powerscale.NullableString{}) {
				attrTypeMap[tag] = types.StringType
			} else {
				structAttrType, err := getStructAttrTypeFromType(ctx, structField.Type)
				if err != nil {
					return nil, err
				}
				attrTypeMap[tag] = types.ObjectType{AttrTypes: structAttrType}
			}
		case reflect.Array, reflect.Slice:
			structAttrType, err := getSliceAttrTypeFromType(ctx, structField.Type)
			if err != nil {
				return nil, err
			}
			attrTypeMap[tag] = structAttrType
		}
	}
	return attrTypeMap, nil

}

func getSliceAttrTypeFromType(ctx context.Context, sliceType reflect.Type) (attr.Type, error) {
	sliceType = sliceType.Elem()
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	switch sliceType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.ListType{ElemType: types.Int64Type}, nil
	case reflect.String:
		return types.ListType{ElemType: types.StringType}, nil
	case reflect.Float32, reflect.Float64:
		return types.ListType{ElemType: types.NumberType}, nil
	case reflect.Bool:
		return types.ListType{ElemType: types.BoolType}, nil
	case reflect.Struct:
		structAttrType, err := getStructAttrTypeFromType(ctx, sliceType)
		if err != nil {
			return nil, err
		}
		return types.ListType{ElemType: types.ObjectType{AttrTypes: structAttrType}}, nil
	case reflect.Array, reflect.Slice:
		sliceAttrType, err := getSliceAttrTypeFromType(ctx, sliceType)
		if err != nil {
			return nil, err
		}
		return types.ListType{ElemType: sliceAttrType}, nil
	default:
		return nil, fmt.Errorf("unknown type")
	}
}

func getStructValue(ctx context.Context, structObj interface{}) (basetypes.ObjectValue, error) {
	elem := reflect.ValueOf(structObj)
	attrType, err := getStructAttrTypeFromType(ctx, reflect.TypeOf(structObj))
	if err != nil {
		return types.ObjectNull(nil), err
	}
	valueMap := make(map[string]attr.Value)
	// iterate the listObject
	for fieldIndex := 0; fieldIndex < elem.NumField(); fieldIndex++ {
		tag := elem.Type().Field(fieldIndex).Tag.Get("json")
		tag = strings.TrimSuffix(tag, ",omitempty")
		elemFieldVal := elem.Field(fieldIndex)
		elemFieldType := elemFieldVal.Type()
		if elemFieldType.Kind() == reflect.Ptr {
			elemFieldVal = elemFieldVal.Elem()
			elemFieldType = elemFieldType.Elem()
		}
		switch elemFieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valueMap[tag] = types.Int64Value(elemFieldVal.Int())
		case reflect.String:
			valueMap[tag] = types.StringValue(elemFieldVal.String())
		case reflect.Bool:
			valueMap[tag] = types.BoolValue(elemFieldVal.Bool())
		case reflect.Struct:
			if elemFieldType == reflect.TypeOf(powerscale.NullableString{}) {
				nullableString, ok := elemFieldVal.Interface().(powerscale.NullableString)
				if !ok {
					return types.ObjectNull(nil), fmt.Errorf("NullableString failed")
				}
				valueMap[tag] = types.StringValue(*nullableString.Get())
			} else {
				valueMap[tag], err = getStructValue(ctx, elemFieldVal.Interface())
				if err != nil {
					return types.ObjectNull(nil), err
				}
			}
		case reflect.Array, reflect.Slice:
			valueMap[tag], err = getSliceAttrValue(ctx, elemFieldVal.Interface())
			if err != nil {
				return types.ObjectNull(nil), err
			}
		}
	}
	object, _ := types.ObjectValue(attrType, valueMap)
	return object, nil
}

func getSliceAttrValue(ctx context.Context, sliceObject interface{}) (attr.Value, error) {
	sliceValue := reflect.ValueOf(sliceObject)
	switch sliceValue.Type().Elem().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		listValue, _ := types.ListValueFrom(ctx, types.Int64Type, sliceObject)
		return listValue, nil
	case reflect.String:
		listValue, _ := types.ListValueFrom(ctx, types.StringType, sliceObject)
		return listValue, nil
	case reflect.Float32, reflect.Float64:
		listValue, _ := types.ListValueFrom(ctx, types.NumberType, sliceObject)
		return listValue, nil
	case reflect.Bool:
		listValue, _ := types.ListValueFrom(ctx, types.BoolType, sliceObject)
		return listValue, nil
	case reflect.Struct:
		var values []attr.Value
		sliceElemType, err := getStructAttrTypeFromType(ctx, sliceValue.Type().Elem())
		if err != nil {
			return nil, err
		}
		for index := 0; index < sliceValue.Len(); index++ {
			sliceElemValue, err := getStructValue(ctx, sliceValue.Index(index).Interface())
			if err != nil {
				return nil, err
			}
			values = append(values, sliceElemValue)
		}
		if len(values) == 0 {
			return types.ListNull(types.ObjectType{AttrTypes: sliceElemType}), nil
		}
		returnListValue, _ := types.ListValue(types.ObjectType{AttrTypes: sliceElemType}, values)
		return returnListValue, nil
	case reflect.Array, reflect.Slice:
		var values []attr.Value
		sliceAttrType, err := getSliceAttrTypeFromType(ctx, sliceValue.Type().Elem())
		if err != nil {
			return nil, err
		}
		for index := 0; index < sliceValue.Len(); index++ {
			sliceElemValue, err := getSliceAttrValue(ctx, sliceValue.Index(index).Interface())
			if err != nil {
				return nil, err
			}
			values = append(values, sliceElemValue)
		}
		if len(values) == 0 {
			return types.ListNull(types.ListType{ElemType: sliceAttrType}), nil
		}
		returnListValue, _ := types.ListValue(types.ListType{ElemType: sliceAttrType}, values)
		return returnListValue, nil
	default:
		return nil, fmt.Errorf("unknown type")
	}
}

// ReadFromState read from model to openapi struct, model should not contain nested struct.
func ReadFromState(ctx context.Context, source, destination interface{}) error {
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination)
	if destinationValue.Kind() != reflect.Ptr || destinationValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination is not a pointer to a struct")
	}
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}
	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("source is not a struct")
	}
	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldTag := sourceValue.Type().Field(i).Tag.Get("tfsdk")
		destinationField, err := getFieldByJSONTag(destinationValue.Elem().Addr().Interface(), sourceFieldTag)
		if err != nil {
			// Not found, skip the field
			continue
		}
		if destinationField.IsValid() && destinationField.CanSet() {
			switch sourceValue.Field(i).Interface().(type) {
			case basetypes.StringValue:
				stringVal, ok := sourceValue.Field(i).Interface().(basetypes.StringValue)
				if !ok || stringVal.IsNull() || stringVal.IsUnknown() {
					continue
				}
				targetValue := stringVal.ValueString()
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.String {
					destinationField.Set(reflect.ValueOf(&targetValue))
				}
				if destinationField.Type().Kind() == reflect.String {
					destinationField.Set(reflect.ValueOf(targetValue))
				}
				if destinationField.Type() == reflect.TypeOf(powerscale.NullableString{}) {
					addr, ok := destinationField.Addr().Interface().(*powerscale.NullableString)
					if !ok {
						continue
					}
					addr.Set(&targetValue)
				}
			case basetypes.Int64Value:
				intVal, ok := sourceValue.Field(i).Interface().(basetypes.Int64Value)
				if !ok || intVal.IsNull() || intVal.IsUnknown() {
					continue
				}
				if destinationField.Kind() == reflect.Int64 {
					destinationField.Set(reflect.ValueOf(intVal.ValueInt64()))
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int64 {
					destinationField.Set(reflect.ValueOf(intVal.ValueInt64Pointer()))
				}
				if destinationField.Kind() == reflect.Int32 {
					destinationField.Set(reflect.ValueOf(int32(intVal.ValueInt64())))
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int32 {
					val := int32(intVal.ValueInt64())
					destinationField.Set(reflect.ValueOf(&val))
				}
			case basetypes.BoolValue:
				boolVal, ok := sourceValue.Field(i).Interface().(basetypes.BoolValue)
				if !ok || boolVal.IsNull() || boolVal.IsUnknown() {
					continue
				}
				if destinationField.Kind() == reflect.Ptr {
					destinationField.Set(reflect.ValueOf(boolVal.ValueBoolPointer()))
				} else {
					destinationField.Set(reflect.ValueOf(boolVal.ValueBool()))
				}
			case basetypes.ObjectValue:
				objVal, ok := sourceValue.Field(i).Interface().(basetypes.ObjectValue)
				if !ok || objVal.IsNull() || objVal.IsUnknown() {
					continue
				}
				err := assignObjectToField(ctx, objVal, destinationField.Addr().Interface())
				if err != nil {
					return err
				}
			case basetypes.ListValue:
				listVal, ok := sourceValue.Field(i).Interface().(basetypes.ListValue)
				if !ok || listVal.IsNull() || listVal.IsUnknown() {
					continue
				}
				err := assignListToField(ctx, listVal, destinationField.Addr().Interface())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func assignObjectToField(ctx context.Context, source basetypes.ObjectValue, destination interface{}) error {
	destElemVal := reflect.ValueOf(destination).Elem()
	destElemType := destElemVal.Type()
	targetObject := reflect.New(destElemType).Elem()
	// if target is pointer to a pointer
	if destElemVal.Kind() == reflect.Ptr {
		destElemVal = reflect.ValueOf(destination).Elem().Elem()
		destElemType = destElemVal.Type()
		targetObject = reflect.New(destElemType).Elem()
	}
	attrMap := source.Attributes()
	for key, val := range attrMap {
		destinationField, err := getFieldByJSONTag(targetObject.Addr().Interface(), key)
		if err != nil {
			// skip current field
			continue
		}
		if destinationField.IsValid() && destinationField.CanSet() {
			switch val.Type(ctx) {
			case basetypes.StringType{}:
				stringVal, ok := val.(basetypes.StringValue)
				if !ok || stringVal.IsNull() || stringVal.IsUnknown() {
					continue
				}
				targetValue := stringVal.ValueString()
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.String {
					destinationField.Set(reflect.ValueOf(&targetValue))
				}
				if destinationField.Type().Kind() == reflect.String {
					destinationField.Set(reflect.ValueOf(targetValue))
				}
			case basetypes.Int64Type{}:
				intVal, ok := val.(basetypes.Int64Value)
				if !ok || intVal.IsNull() || intVal.IsUnknown() {
					continue
				}
				if destinationField.Kind() == reflect.Int64 {
					destinationField.Set(reflect.ValueOf(intVal.ValueInt64()))
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int64 {
					destinationField.Set(reflect.ValueOf(intVal.ValueInt64Pointer()))
				}
				if destinationField.Kind() == reflect.Int32 {
					destinationField.Set(reflect.ValueOf(int32(intVal.ValueInt64())))
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int32 {
					val := int32(intVal.ValueInt64())
					destinationField.Set(reflect.ValueOf(&val))
				}
			case basetypes.BoolType{}:
				boolVal, ok := val.(basetypes.BoolValue)
				if !ok || boolVal.IsNull() || boolVal.IsUnknown() {
					continue
				}
				if destinationField.Kind() == reflect.Ptr {
					destinationField.Set(reflect.ValueOf(boolVal.ValueBoolPointer()))
				} else {
					destinationField.Set(reflect.ValueOf(boolVal.ValueBool()))
				}
			default:
				typeString := val.Type(ctx).String()
				if strings.HasPrefix(typeString, "types.ObjectType") {
					objVal, ok := val.(basetypes.ObjectValue)
					if !ok || objVal.IsNull() || objVal.IsUnknown() {
						continue
					}
					err := assignObjectToField(ctx, objVal, destinationField.Addr().Interface())
					if err != nil {
						return err
					}
				} else if strings.HasPrefix(typeString, "types.ListType") {
					listVal, ok := val.(basetypes.ListValue)
					if !ok || listVal.IsNull() || listVal.IsUnknown() {
						continue
					}
					err := assignListToField(ctx, listVal, destinationField.Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		}
	}
	destElemVal.Set(targetObject)
	return nil
}

func assignListToField(ctx context.Context, source basetypes.ListValue, destination interface{}) error {
	destVal := reflect.ValueOf(destination).Elem()
	// type of element of slice
	destType := destVal.Type()
	// if target is pointer to a pointer
	if destVal.Kind() == reflect.Ptr {
		destVal = destVal.Elem()
		destType = destVal.Type()
	}
	listLen := len(source.Elements())

	listElemType := source.ElementType(ctx)
	switch listElemType {
	case basetypes.StringType{}:
		tfsdk.ValueAs(ctx, source, destination)
	case basetypes.Int64Type{}:
		tfsdk.ValueAs(ctx, source, destination)
	case basetypes.BoolType{}:
		tfsdk.ValueAs(ctx, source, destination)
	default:
		targetList := reflect.MakeSlice(destType, listLen, listLen)
		typeString := listElemType.String()
		for i, listElem := range source.Elements() {
			if strings.HasPrefix(typeString, "types.ListType") {
				listVal, ok := listElem.(basetypes.ListValue)
				if !ok || listVal.IsNull() || listVal.IsUnknown() {
					continue
				}
				err := assignListToField(ctx, listVal, targetList.Index(i).Addr().Interface())
				if err != nil {
					return err
				}
			} else if strings.HasPrefix(typeString, "types.ObjectType") {
				objVal, ok := listElem.(basetypes.ObjectValue)
				if !ok || objVal.IsNull() || objVal.IsUnknown() {
					continue
				}
				err := assignObjectToField(ctx, objVal, targetList.Index(i).Addr().Interface())
				if err != nil {
					return err
				}
			}
		}
		destVal.Set(targetList)
	}
	return nil
}

// getFieldByJSONTag get field by tag, input destination is pointer.
func getFieldByJSONTag(destination interface{}, tag string) (reflect.Value, error) {
	destElemVal := reflect.ValueOf(destination).Elem()
	destElemType := destElemVal.Type()

	for i := 0; i < destElemType.NumField(); i++ {
		field := destElemType.Field(i)
		jsonTag := field.Tag.Get("json")
		if strings.Contains(jsonTag, ",") {
			jsonTag = strings.TrimSuffix(jsonTag, ",omitempty")
		}
		if jsonTag == tag {
			return destElemVal.Field(i), nil
		}
	}

	return reflect.Value{}, fmt.Errorf("field with tag %s not found in destination", tag)
}
