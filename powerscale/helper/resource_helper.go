/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"math"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
		}

		// For zero value (nil), the object still need to pass type information into it
		if !sourceField.IsValid() {
			destinationField = getFieldByTfTag(destinationValue.Elem(), sourceFieldTag)
			if !destinationField.IsValid() {
				continue
			}
			fieldType := structType.Field(i).Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			switch fieldType.Kind() {
			case reflect.String:
				destinationField.Set(reflect.ValueOf(types.StringNull()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				destinationField.Set(reflect.ValueOf(types.Int64Null()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				destinationField.Set(reflect.ValueOf(types.Int64Null()))
			case reflect.Float32, reflect.Float64:
				//destinationFieldValue = types.Float64Value(sourceField.Float())
				destinationField.Set(reflect.ValueOf(types.NumberNull()))
			case reflect.Bool:
				destinationField.Set(reflect.ValueOf(types.BoolNull()))
			case reflect.Array, reflect.Slice:
				mapType, err := getSliceAttrTypeFromType(ctx, structType.Field(i).Type)
				if err != nil {
					return err
				}
				destinationField.Set(reflect.ValueOf(types.ListNull(mapType)))
			case reflect.Struct:
				mapType, err := getStructAttrTypeFromType(ctx, structType.Field(i).Type)
				if err != nil {
					return err
				}
				destinationField.Set(reflect.ValueOf(types.ObjectNull(mapType)))
			default:
				tflog.Error(ctx, "unsupported source field type", map[string]interface{}{
					"sourceField": sourceField,
				})
			}
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
			} else if destinationField.Type().String() == "models.CaseInsensitiveStringValue" {
				tflog.Debug(ctx, "setting destination field to case insensitive string", map[string]interface{}{
					"sourceField": sourceField.String(),
				})
				destinationField.Set(reflect.ValueOf(models.CaseInsensitiveStringValue{StringValue: types.StringValue(sourceField.String())}))
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
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	sliceType = sliceType.Elem()
	// if a list of pointer
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
			if elemFieldVal.IsValid() {
				valueMap[tag] = types.Int64Value(elemFieldVal.Int())
			} else {
				valueMap[tag] = types.Int64Null()
			}
		case reflect.Float32, reflect.Float64:
			if elemFieldVal.IsValid() {
				// Due to accuracy issue, keep the precision as 4
				floatVal, err := strconv.ParseFloat(fmt.Sprintf("%.4f", elemFieldVal.Float()), 64)
				if err != nil {
					return types.ObjectNull(nil), err
				}
				valueMap[tag] = types.NumberValue(big.NewFloat(floatVal))
			} else {
				valueMap[tag] = types.NumberNull()
			}
		case reflect.String:
			if elemFieldVal.IsValid() {
				valueMap[tag] = types.StringValue(elemFieldVal.String())
			} else {
				valueMap[tag] = types.StringNull()
			}
		case reflect.Bool:
			if elemFieldVal.IsValid() {
				valueMap[tag] = types.BoolValue(elemFieldVal.Bool())
			} else {
				valueMap[tag] = types.BoolNull()
			}
		case reflect.Struct:
			if elemFieldType == reflect.TypeOf(powerscale.NullableString{}) {
				nullableString, ok := elemFieldVal.Interface().(powerscale.NullableString)
				if !ok {
					return types.ObjectNull(nil), fmt.Errorf("NullableString failed")
				}
				if !nullableString.IsSet() {
					valueMap[tag] = types.StringNull()
				} else {
					valueMap[tag] = types.StringValue(*nullableString.Get())
				}
			} else {
				if !elemFieldVal.IsValid() {
					mapType, err := getStructAttrTypeFromType(ctx, elemFieldType)
					if err != nil {
						return types.ObjectNull(nil), err
					}
					valueMap[tag] = types.ObjectNull(mapType)
				} else {
					valueMap[tag], err = getStructValue(ctx, elemFieldVal.Interface())
				}
				if err != nil {
					return types.ObjectNull(nil), err
				}
			}
		case reflect.Array, reflect.Slice:
			if !elemFieldVal.IsValid() {
				valueMap[tag], err = getSliceAttrValue(ctx, reflect.Zero(reflect.TypeOf(elem.Field(fieldIndex).Interface()).Elem()).Interface())
			} else {
				valueMap[tag], err = getSliceAttrValue(ctx, elemFieldVal.Interface())
			}
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
			return types.ListNull(sliceAttrType), nil
		}
		returnListValue, _ := types.ListValue(sliceAttrType, values)
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
			switch val := sourceValue.Field(i).Interface().(type) {
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
					if intVal.ValueInt64() > math.MaxInt32 || intVal.ValueInt64() < math.MinInt32 {
						return fmt.Errorf("value %d is out of range for int32", intVal.ValueInt64())
					}
					destinationField.Set(reflect.ValueOf(int32(intVal.ValueInt64()))) // #nosec G115 - Validated, Error returned if value is out of range
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int32 {
					if intVal.ValueInt64() > math.MaxInt32 || intVal.ValueInt64() < math.MinInt32 {
						return fmt.Errorf("value %d is out of range for int32", intVal.ValueInt64())
					}
					val := int32(intVal.ValueInt64()) // #nosec G115 - Validated, Error returned if value is out of range
					destinationField.Set(reflect.ValueOf(&val))
				}
			case basetypes.Int32Value:
				intVal, ok := sourceValue.Field(i).Interface().(basetypes.Int32Value)
				if !ok || intVal.IsNull() || intVal.IsUnknown() {
					continue
				}
				if destinationField.Kind() == reflect.Int32 {
					destinationField.Set(reflect.ValueOf(intVal.ValueInt32()))
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int32 {
					destinationField.Set(reflect.ValueOf(intVal.ValueInt32Pointer()))
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
			case basetypes.NumberValue:
				floatVal, ok := sourceValue.Field(i).Interface().(basetypes.NumberValue)
				if !ok || floatVal.IsNull() || floatVal.IsUnknown() {
					continue
				}
				bigFloat := floatVal.ValueBigFloat()
				if destinationField.Kind() == reflect.Ptr {
					if destinationField.Type().Elem().Kind() == reflect.Float64 {
						bigFloatVal, _ := bigFloat.Float64()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 64)
						destinationField.Set(reflect.ValueOf(&floatVal))
					}
					if destinationField.Type().Elem().Kind() == reflect.Float32 {
						bigFloatVal, _ := bigFloat.Float32()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 32)
						float32Val := float32(floatVal)
						destinationField.Set(reflect.ValueOf(&float32Val))
					}
				} else {
					if destinationField.Kind() == reflect.Float64 {
						bigFloatVal, _ := bigFloat.Float64()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 64)
						destinationField.Set(reflect.ValueOf(floatVal))
					}
					if destinationField.Kind() == reflect.Float32 {
						bigFloatVal, _ := bigFloat.Float32()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 32)
						destinationField.Set(reflect.ValueOf(float32(floatVal)))
					}
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
				list, err := getFieldListVal(ctx, listVal, destinationField.Interface())
				if err != nil {
					return err
				}
				if reflect.TypeOf(destinationField.Interface()).Kind() == reflect.Ptr {
					destinationField.Set(reflect.New(destinationField.Type().Elem()))
					destinationField.Elem().Set(list)
				} else {
					destinationField.Set(list)
				}
			case basetypes.SetValue:
				if val.IsNull() || val.IsUnknown() {
					continue
				}
				list, err := getFieldListVal(ctx, val, destinationField.Interface())
				if err != nil {
					return err
				}
				destinationField.Set(list)
			}
		}
	}
	return nil
}

func assignObjectToField(ctx context.Context, source basetypes.ObjectValue, destination interface{}) error {
	destElemType := reflect.TypeOf(destination).Elem()
	isPtr := false
	if destElemType.Kind() == reflect.Ptr {
		isPtr = true
		destElemType = destElemType.Elem()
	}
	targetObject := reflect.New(destElemType).Elem()
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
				if destinationField.Type() == reflect.TypeOf(powerscale.NullableString{}) {
					addr, ok := destinationField.Addr().Interface().(*powerscale.NullableString)
					if !ok {
						continue
					}
					addr.Set(&targetValue)
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
					if intVal.ValueInt64() > math.MaxInt32 || intVal.ValueInt64() < math.MinInt32 {
						return fmt.Errorf("value %d is out of range for int32", intVal.ValueInt64())
					}
					destinationField.Set(reflect.ValueOf(int32(intVal.ValueInt64()))) // #nosec G115 - Validated, Error returned if value is out of range
				}
				if destinationField.Kind() == reflect.Ptr && destinationField.Type().Elem().Kind() == reflect.Int32 {
					if intVal.ValueInt64() > math.MaxInt32 || intVal.ValueInt64() < math.MinInt32 {
						return fmt.Errorf("value %d is out of range for int32", intVal.ValueInt64())
					}
					val := int32(intVal.ValueInt64()) // #nosec G115 - Validated, Error returned if value is out of range
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
			case basetypes.NumberType{}:
				floatVal, ok := val.(basetypes.NumberValue)
				if !ok || floatVal.IsNull() || floatVal.IsUnknown() {
					continue
				}
				bigFloat := floatVal.ValueBigFloat()
				if destinationField.Kind() == reflect.Ptr {
					if destinationField.Type().Elem().Kind() == reflect.Float64 {
						bigFloatVal, _ := bigFloat.Float64()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 64)
						destinationField.Set(reflect.ValueOf(&floatVal))
					}
					if destinationField.Type().Elem().Kind() == reflect.Float32 {
						bigFloatVal, _ := bigFloat.Float32()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 32)
						float32Val := float32(floatVal)
						destinationField.Set(reflect.ValueOf(&float32Val))
					}
				} else {
					if destinationField.Kind() == reflect.Float64 {
						bigFloatVal, _ := bigFloat.Float64()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 64)
						destinationField.Set(reflect.ValueOf(floatVal))
					}
					if destinationField.Kind() == reflect.Float32 {
						bigFloatVal, _ := bigFloat.Float32()
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", bigFloatVal), 32)
						destinationField.Set(reflect.ValueOf(float32(floatVal)))
					}
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
					list, err := getFieldListVal(ctx, listVal, destinationField.Interface())
					if err != nil {
						return err
					}
					if reflect.TypeOf(destinationField.Interface()).Kind() == reflect.Ptr {
						destinationField.Set(reflect.New(destinationField.Type().Elem()))
						destinationField.Elem().Set(list)
					} else {
						destinationField.Set(list)
					}
				}
			}
		}
	}
	if isPtr {
		reflect.ValueOf(destination).Elem().Set(targetObject.Addr())
	} else {
		reflect.ValueOf(destination).Elem().Set(targetObject)
	}
	return nil
}

type listOrsetValue interface {
	basetypes.ListValue | basetypes.SetValue
	// the functions need to be explicitly defined till today due to GoLang issue
	// https://github.com/golang/go/issues/51183
	ElementType(context.Context) attr.Type
	Elements() []attr.Value
}

func getFieldListVal[T listOrsetValue](ctx context.Context, source T, destination interface{}) (reflect.Value, error) {
	destType := reflect.TypeOf(destination)
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	listLen := len(source.Elements())
	targetList := reflect.MakeSlice(destType, listLen, listLen)
	listElemType := source.ElementType(ctx)
	for i, listElem := range source.Elements() {
		switch listElemType {
		case basetypes.StringType{}:
			strVal, ok := listElem.(basetypes.StringValue)
			if !ok || strVal.IsNull() || strVal.IsUnknown() {
				continue
			}
			if destType.Elem().Kind() == reflect.Ptr {
				targetList.Index(i).Elem().Set(reflect.ValueOf(strVal.ValueStringPointer()))
			} else {
				targetList.Index(i).Set(reflect.ValueOf(strVal.ValueString()))
			}
		case basetypes.Int64Type{}:
			strVal, ok := listElem.(basetypes.Int64Value)
			if !ok || strVal.IsNull() || strVal.IsUnknown() {
				continue
			}
			if destType.Elem().Kind() == reflect.Ptr {
				targetList.Index(i).Elem().Set(reflect.ValueOf(strVal.ValueInt64Pointer()))
			} else {
				targetList.Index(i).Set(reflect.ValueOf(strVal.ValueInt64()))
			}
		case basetypes.BoolType{}:
			strVal, ok := listElem.(basetypes.BoolValue)
			if !ok || strVal.IsNull() || strVal.IsUnknown() {
				continue
			}
			if destType.Elem().Kind() == reflect.Ptr {
				targetList.Index(i).Elem().Set(reflect.ValueOf(strVal.ValueBoolPointer()))
			} else {
				targetList.Index(i).Set(reflect.ValueOf(strVal.ValueBool()))
			}
		default:
			typeString := listElemType.String()
			if strings.HasPrefix(typeString, "types.ListType") {
				listVal, ok := listElem.(basetypes.ListValue)
				if !ok || listVal.IsNull() || listVal.IsUnknown() {
					continue
				}
				val, err := getFieldListVal(ctx, listVal, targetList.Index(i).Interface())
				if err != nil {
					return targetList, err
				}
				if reflect.TypeOf(targetList.Index(i).Interface()).Kind() == reflect.Ptr {
					targetList.Index(i).Set(reflect.New(targetList.Index(i).Type().Elem()))
					targetList.Index(i).Elem().Set(val)
				} else {
					targetList.Index(i).Set(val)
				}
			} else if strings.HasPrefix(typeString, "types.ObjectType") {
				objVal, ok := listElem.(basetypes.ObjectValue)
				if !ok || objVal.IsNull() || objVal.IsUnknown() {
					continue
				}
				err := assignObjectToField(ctx, objVal, targetList.Index(i).Addr().Interface())
				if err != nil {
					return targetList, err
				}
			}
		}
	}
	return targetList, nil
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

// GetElementsChanges Returns element list changes between plan and state.
func GetElementsChanges(stateElements, planElements []attr.Value) (toAdd, toRemove []attr.Value) {

	// Note: the attr.Value's string value, attr.Value.String(), will be like "\"value\"",
	// please firstly use strings.Trim(attr.Value.String(), "\"").

	var duplicatedMembers []attr.Value
	for _, i := range stateElements {
		for _, j := range planElements {
			if i.Equal(j) {
				duplicatedMembers = append(duplicatedMembers, i)
			}
		}
	}

	for _, i := range stateElements {
		duplicated := false
		for _, member := range duplicatedMembers {
			if member.Equal(i) {
				duplicated = true
				break
			}
		}
		if duplicated {
			continue
		}
		toRemove = append(toRemove, i)
	}

	for _, i := range planElements {
		duplicated := false
		for _, member := range duplicatedMembers {
			if member.Equal(i) {
				duplicated = true
				break
			}
		}
		if duplicated {
			continue
		}
		toAdd = append(toAdd, i)
	}
	return
}

// IsListValueEquals checks if two terraform list are equal.
// It returns true if the lists have the same length and contain the same elements, otherwise it returns false.
func IsListValueEquals(a types.List, b types.List) bool {
	elementsA := a.Elements()
	elementsB := b.Elements()
	if len(elementsA) != len(elementsB) {
		return false
	}

	sort.Slice(elementsA, func(i, j int) bool {
		return elementsA[i].String() < elementsA[j].String()
	})
	sort.Slice(elementsB, func(i, j int) bool {
		return elementsB[i].String() < elementsB[j].String()
	})

	for i, v := range elementsA {
		if v.String() != elementsB[i].String() {
			return false
		}
	}
	return true
}
