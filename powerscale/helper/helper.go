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
	"bytes"
	"context"
	powerscale "dell/powerscale-go-client"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"terraform-provider-powerscale/client"

	"golang.org/x/net/html"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CopyFields copy the source of a struct to destination of struct with terraform types.
// Unsigned integers are not properly handled.
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
	// destinationType := destinationValue.Elem().Type()
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
		if !sourceField.IsValid() {
			tflog.Error(ctx, "source field is not valid", map[string]interface{}{
				"sourceFieldTag": sourceFieldTag,
				"sourceField":    sourceField,
			})
			continue
		}

		destinationField := getFieldByTfTag(destinationValue.Elem(), sourceFieldTag)
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
				if sourceField.Uint() > math.MaxInt64 {
					return fmt.Errorf("source field value is too large for int64")
				}
				destinationFieldValue = types.Int64Value(int64(sourceField.Uint())) // #nosec G115 --- validated, Error returned if value is too large for int64
			case reflect.Float32, reflect.Float64:
				// destinationFieldValue = types.Float64Value(sourceField.Float())
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
				} else if /* check if destination is types.Set */ _, ok := destinationField.Interface().(types.Set); ok {
					destinationFieldValue = copySliceToSetType(ctx, sourceField.Interface())

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

func getFieldJSONTag(sourceValue reflect.Value, i int) string {
	sourceFieldTag := sourceValue.Type().Field(i).Tag.Get("json")
	sourceFieldTag = strings.TrimSuffix(sourceFieldTag, ",omitempty")
	return sourceFieldTag
}

func getFieldByTfTag(destinationValue reflect.Value, tagValue string) reflect.Value {
	for j := 0; j < destinationValue.NumField(); j++ {
		field := destinationValue.Type().Field(j)
		if field.Tag.Get("tfsdk") == tagValue {
			return destinationValue.Field(j)
		}
	}
	return reflect.Value{}
}

func copySliceToSetType(ctx context.Context, fields any) types.Set {
	listVal := copySliceToTargetField(ctx, fields)
	if listVal.IsUnknown() {
		return types.SetUnknown(listVal.ElementType(ctx))
	}
	setValue, _ := types.SetValue(listVal.ElementType(ctx), listVal.Elements())
	return setValue
}

func copySliceToTargetField(ctx context.Context, fields interface{}) types.List {
	var objects []attr.Value
	attrTypeMap := make(map[string]attr.Type)

	// get the attrType for Object
	structElem := reflect.ValueOf(fields).Type().Elem()
	switch structElem.Kind() {
	case reflect.String:
		listValue, _ := types.ListValueFrom(ctx, types.StringType, fields)
		return listValue
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		listValue, _ := types.ListValueFrom(ctx, types.Int64Type, fields)
		return listValue
	case reflect.Float32, reflect.Float64:
		listValue, _ := types.ListValueFrom(ctx, types.Float64Type, fields)
		return listValue
	case reflect.Bool:
		listValue, _ := types.ListValueFrom(ctx, types.BoolType, fields)
		return listValue
	case reflect.Struct:
		for fieldIndex := 0; fieldIndex < structElem.NumField(); fieldIndex++ {
			field := structElem.Field(fieldIndex)
			tag := field.Tag.Get("json")
			tag = strings.TrimSuffix(tag, ",omitempty")
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			switch fieldType.Kind() {
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
				eleField := elem.Field(fieldIndex)
				eleFieldType := eleField.Type()
				if eleFieldType.Kind() == reflect.Ptr {
					eleFieldType = eleFieldType.Elem()
					eleField = eleField.Elem()
				}
				switch eleFieldType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					valueMap[tag] = types.Int64Value(eleField.Int())
				case reflect.String:
					valueMap[tag] = types.StringValue(eleField.String())
				case reflect.Float32, reflect.Float64:
					valueMap[tag] = types.NumberValue(big.NewFloat(eleField.Float()))
				}
			}
			object, _ := types.ObjectValue(attrTypeMap, valueMap)
			objects = append(objects, object)
		}
		listValue, _ := types.ListValue(types.ObjectType{AttrTypes: attrTypeMap}, objects)
		return listValue
	}
	return types.ListUnknown(types.StringType)
}

// ParseBody parses the error message from an openApi error response.
func ParseBody(body []byte) (string, error) {
	contentType := http.DetectContentType(body)

	if strings.Contains(contentType, "text/html") {
		return parseHTMLBody(body)
	}

	var parsedData map[string][]map[string]string
	err := json.Unmarshal(body, &parsedData)
	if err != nil {
		return "", err
	}

	var message string
	if errors, ok := parsedData["errors"]; ok {
		for _, e := range errors {
			message = message + e["message"] + " "
		}
		return message + " ", nil
	}
	return "", fmt.Errorf("no message field found in body")
}

// GetErrorString extracts the error message from an openApi error response.
func GetErrorString(err error, errStr string) string {
	err1, ok := err.(*powerscale.GenericOpenAPIError)
	message := ""
	msgStr := ""
	if ok {
		message, err := ParseBody(err1.Body())
		if err != nil {
			msgStr = errStr + err.Error()
		}
		errStr = errStr + message
	}
	if message == "" {
		msgStr = errStr + err.Error()
	}
	return msgStr
}

func parseHTMLBody(body []byte) (string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	var message string
	titleNode := findNode(doc, "title")
	if titleNode != nil {
		extractText(titleNode, &message)
	} else {
		return "", fmt.Errorf("no title element found in HTML body")
	}

	if message == "" {
		return "", fmt.Errorf("no message found in HTML title")
	}

	return message, nil
}

func findNode(n *html.Node, data string) *html.Node {
	if n.Type == html.ElementNode && n.Data == data {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if bodyNode := findNode(c, data); bodyNode != nil {
			return bodyNode
		}
	}

	return nil
}

func extractText(n *html.Node, message *string) {
	if n.Type == html.TextNode {
		*message += strings.TrimSpace(n.Data) + " "
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, message)
	}
}

// GetClusterVersion retrieves the cluster version.
func GetClusterVersion(ctx context.Context, client *client.Client) (string, error) {
	clusterVersion, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterVersion(ctx).Execute()
	if err != nil {
		return "", err
	}
	return clusterVersion.Nodes[0].Release, err
}

// For List set explicitly from plan
// This is to keep state in similar order to plan
// Lists returned from the array are not always in the same order as they appear in the plan
func ListCheck(list types.List, elementType attr.Type) types.List {
	if list.IsUnknown() {
		return types.ListNull(elementType)
	}
	return list
}

func DefaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
