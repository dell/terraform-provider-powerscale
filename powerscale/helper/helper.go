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
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CopyFields copy the source of a struct to destination of struct with terraform types.
func CopyFields(ctx context.Context, source, destination interface{}) error {
	tflog.Debug(ctx, "Copy fields", map[string]interface{}{
		"source":      source,
		"destination": destination,
	})
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
		sourceFieldName := sourceValue.Type().Field(i).Name

		tflog.Debug(ctx, "Converting source field", map[string]interface{}{
			"sourceFieldName": sourceFieldName,
			"sourceFieldKind": sourceValue.Field(i).Kind().String(),
		})

		sourceField := sourceValue.Field(i)
		if sourceField.Kind() == reflect.Ptr {
			sourceField = sourceField.Elem()
		}
		if !sourceField.IsValid() {
			tflog.Error(ctx, "source field is not valid", map[string]interface{}{
				"sourceFieldName": sourceFieldName,
				"sourceField":     sourceField,
			})
			continue
		}

		destinationField := destinationValue.Elem().FieldByName(sourceFieldName)
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
				if destinationField.Type().Kind() == reflect.Slice {
					arr := reflect.ValueOf(sourceField.Interface())
					slice := reflect.MakeSlice(destinationField.Type(), arr.Len(), arr.Cap())
					for index := 0; index < arr.Len(); index++ {
						value := arr.Index(index)
						v := slice.Index(index)
						switch v.Kind() {
						case reflect.Ptr:
							newDes := reflect.New(v.Type().Elem()).Interface()
							err := CopyFields(ctx, value.Interface(), newDes)
							if err != nil {
								return err
							}
							slice.Index(index).Set(reflect.ValueOf(newDes))
						case reflect.Struct:
							newDes := reflect.New(v.Type()).Interface()
							err := CopyFields(ctx, value.Interface(), newDes)
							if err != nil {
								return err
							}
							slice.Index(index).Set(reflect.ValueOf(newDes).Elem())
						}
					}
					destinationField.Set(slice)
				} else {
					destinationFieldValue = copySliceToTargetField(ctx, sourceField.Interface())
				}
			case reflect.Struct:
				// placeholder for improvement, need to consider both go struct and types.Object
				switch destinationField.Kind() {
				case reflect.Ptr:
					newDes := reflect.New(destinationField.Type().Elem()).Interface()
					err := CopyFields(ctx, sourceField.Interface(), newDes)
					if err != nil {
						return err
					}
					destinationField.Set(reflect.ValueOf(newDes))
				case reflect.Struct:
					newDes := reflect.New(destinationField.Type()).Interface()
					err := CopyFields(ctx, sourceField.Interface(), newDes)
					if err != nil {
						return err
					}
					destinationField.Set(reflect.ValueOf(newDes).Elem())
				}
				continue

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

func copySliceToTargetField(ctx context.Context, fields interface{}) attr.Value {
	var objects []attr.Value
	attrTypeMap := make(map[string]attr.Type)

	// get the attrType for Object
	structElem := reflect.ValueOf(fields).Type().Elem()
	if structElem.Kind() == reflect.String {
		listValue, _ := types.ListValueFrom(ctx, types.StringType, fields)
		return listValue
	} else if structElem.Kind() == reflect.Struct {
		for fieldIndex := 0; fieldIndex < structElem.NumField(); fieldIndex++ {
			field := structElem.Field(fieldIndex)
			tag := field.Tag.Get("json")
			tag = strings.TrimSuffix(tag, ",omitempty")
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				attrTypeMap[tag] = types.Int64Type
			case reflect.String:
				attrTypeMap[tag] = types.StringType
			case reflect.Float32, reflect.Float64:
				attrTypeMap[tag] = types.NumberType
			}
		}
		// iterate the slice
		arr := reflect.ValueOf(fields)
		for index := 0; index < arr.Len(); index++ {
			valueMap := make(map[string]attr.Value)
			// iterate the fields
			elem := arr.Index(index)
			for fieldIndex := 0; fieldIndex < elem.NumField(); fieldIndex++ {
				tag := elem.Type().Field(fieldIndex).Tag.Get("json")
				tag = strings.TrimSuffix(tag, ",omitempty")
				switch elem.Field(fieldIndex).Type().Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					valueMap[tag] = types.Int64Value(elem.Field(fieldIndex).Int())
				case reflect.String:
					valueMap[tag] = types.StringValue(elem.Field(fieldIndex).String())
				case reflect.Float32, reflect.Float64:
					valueMap[tag] = types.NumberValue(big.NewFloat(elem.Field(fieldIndex).Float()))
				}
			}
			object, _ := types.ObjectValue(attrTypeMap, valueMap)
			objects = append(objects, object)
		}
		listValue, _ := types.ListValue(types.ObjectType{AttrTypes: attrTypeMap}, objects)
		return listValue
	}
	return nil
}

// ParseBody parses the error message from an openApi error response
func ParseBody(body []byte) (string, error) {
	var parsedData map[string]string
	err := json.Unmarshal(body, &parsedData)
	if err != nil {
		return "", err
	}
	message, ok := parsedData["message"]
	if !ok {
		return "", fmt.Errorf("no message field found in body")
	}
	return message, nil
}

// GetErrorString extracts the error message from an openApi error response
func GetErrorString(err error, errStr string) string {
	err1, ok := err.(*powerscale.GenericOpenAPIError)
	message := ""
	msgStr := ""
	if ok {
		message, _ := ParseBody(err1.Body())
		errStr = errStr + message
	}
	if message == "" {
		msgStr = errStr + err.Error()
	}
	return msgStr
}
