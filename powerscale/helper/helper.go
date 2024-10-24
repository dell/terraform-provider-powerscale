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
	"golang.org/x/net/html"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
				destinationFieldValue = types.Int64Value(sourceField.Int())
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

func copySliceToTargetField(ctx context.Context, fields interface{}) attr.Value {
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
	return nil
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

// GetDataSourceByValue is a helper function that gathers data based on all data gathered by the datasource.
//
// Parameters:
// - fields: The fields to filter the data.
// - allData: The data to be filtered.
//
// Returns:
// - []interface{}: The filtered data.
// - error: An error if any occurred.
func GetDataSourceByValue(ctx context.Context, fields interface{}, allData interface{}) ([]interface{}, error) {

	if isPointer(fields) || isPointer(allData) {
		return nil, fmt.Errorf("Pointers are not supported")
	}

	filteredArray := reflect.ValueOf(allData)
	fieldsArray := reflect.ValueOf(fields)
	var err error

	for j := 0; j < fieldsArray.NumField(); j++ {

		field := fieldsArray.Type().Field(j).Name
		fieldValue := fieldsArray.FieldByName(field)

		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			if fieldValue.IsNil() {
				continue
			}
		} else {
			if fieldValue.IsZero() {
				continue
			}
		}

		filteredArray, err = FilterByField(ctx, filteredArray, fieldValue, field)

		if err != nil {
			return nil, err
		}

	}

	allFilteredData := make([]interface{}, filteredArray.Len())
	for i := 0; i < filteredArray.Len(); i++ {
		allFilteredData[i] = filteredArray.Index(i).Interface()

	}

	return allFilteredData, nil

}

// FilterByField filters the array of data sources based on the provided field.
//
// Parameters:
// - dataSources: The array of data sources to filter.
// - fieldValue: The value to filter the data sources by.
// - field: The name of the field to filter by.
//
// Returns:
// - reflect.Value: The filtered array of data sources.
// - error: An error if any occurred.
func FilterByField(ctx context.Context, dataSources reflect.Value, fieldValue reflect.Value, field string) (reflect.Value, error) {
	filteredData := reflect.MakeSlice(dataSources.Type(), 0, dataSources.Len())

	for i := 0; i < dataSources.Len(); i++ {

		dataSource := dataSources.Index(i).Interface()

		dataSourceValue := reflect.ValueOf(dataSource)
		fieldValueInDataSource := dataSourceValue.FieldByName(field)
		tflog.Debug(ctx, "kind-tf : "+fieldValueInDataSource.Kind().String()+" field name: "+field)

		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			for n := 0; n < fieldValue.Len(); n++ {

				interFieldValue, err := CheckAndConvertValue(fieldValue.Index(n))
				if err != nil {
					return reflect.Zero(nil), err
				}
				tflog.Debug(ctx, "fieldValueInDatasource: "+fieldValueInDataSource.Elem().Interface().(string))
				if fieldValueInDataSource.Elem().Interface() == interFieldValue.Interface() || fieldValueInDataSource.Interface() == interFieldValue.Interface() {
					filteredData = reflect.Append(filteredData, reflect.ValueOf(dataSource))
				}
			}
		} else {
			interFieldValue, err := CheckAndConvertValue(fieldValue)
			if err != nil {
				return reflect.Zero(nil), err
			}
			if fieldValueInDataSource.Elem().Interface() == interFieldValue.Interface() || fieldValueInDataSource.Interface() == interFieldValue.Interface() {
				filteredData = reflect.Append(filteredData, reflect.ValueOf(dataSource))
			}
		}
	}

	return filteredData, nil
}

// CheckAndConvertValue converts a reflect.Value to an attr.Type.
//
// It takes in a reflect.Value and checks its type. If the type is a
// types.StringType, it converts the input to a string, trims the quotes,
// and returns the resulting string. If the type is a types.Int64Type, it
// converts the input to an int, and returns the resulting int. If the
// type is a types.BoolType, it converts the input to a bool, and returns
// the resulting bool. If the type is none of the above, it returns an
// error.
//
// Returns:
// - reflect.Value: The converted value.
// - error: An error if the input type is not recognized.
func CheckAndConvertValue(input reflect.Value) (reflect.Value, error) {
	var valueRef reflect.Value
	switch ConvertType(input.Type()) {
	case types.StringType:
		value := fmt.Sprintf("%v", input)
		value = strings.Trim(value, "\"")
		valueRef = reflect.ValueOf(value)

		return valueRef, nil
	case types.Int64Type:
		value, err := strconv.Atoi(fmt.Sprintf("%v", input))
		if err != nil {
			return valueRef, nil
		}
		valueRef = reflect.ValueOf(value)

		return valueRef, nil
	case types.BoolType:
		value, err := strconv.ParseBool(fmt.Sprintf("%v", input))
		if err != nil {
			return valueRef, nil
		}
		valueRef = reflect.ValueOf(value)

		return valueRef, nil
	}

	return valueRef, fmt.Errorf("Value cannot be converted: %v", input)
}

// GenerateSchemaAttributes generates schema attributes based on a map of attribute names and respective types.
//
// The function takes a map of attributes, where each attribute is a map of attribute types and a boolean flag
// indicating whether the attribute is a set. The function iterates over the attributes and for each attribute,
// it generates a schema attribute using the SchemaAttributeGeneration function. The generated schema attributes
// are stored in a map, where the attribute name is the key.
//
// The function returns the generated schema attributes as a map, where the attribute name is the key and the
// corresponding schema attribute is the value.
func GenerateSchemaAttributes(attributes map[string]map[attr.Type]bool) map[string]schema.Attribute {
	schemaAttributes := make(map[string]schema.Attribute)
	for field, attrMap := range attributes {
		for attrType, isSet := range attrMap {
			schemaAttributes[field] = SchemaAttributeGeneration(field, attrType, isSet)
		}
	}
	tflog.Info(context.Background(), fmt.Sprintf("Generated Schema Attributes: %v", schemaAttributes))
	return schemaAttributes
}

// SchemaAttributeGeneration generates a schema attribute based on the type.
//
// It takes in a field name, attribute type, and a boolean flag indicating
// whether the attribute is a set. If the attribute is a set, it returns a
// schema.SetAttribute with the specified element type and validators. If the
// attribute is not a set, it returns a schema.Attribute of the specified type
// with the field name and description.
//
// Returns a schema.Attribute.
func SchemaAttributeGeneration(field string, attrType attr.Type, isSet bool) schema.Attribute {

	if isSet {
		return schema.SetAttribute{
			Description:         "List of " + field,
			MarkdownDescription: "List of " + field,
			ElementType:         attrType,
			Optional:            true,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		}
	}
	switch attrType {
	case types.StringType:
		return schema.StringAttribute{
			Description:         "Value for" + field,
			MarkdownDescription: "Value for " + field,
			Optional:            true,
		}
	case types.Int64Type:
		return schema.Int64Attribute{
			Description:         "Value for " + field,
			MarkdownDescription: "Value for " + field,
			Optional:            true,
		}
	case types.BoolType:
		return schema.BoolAttribute{
			Description:         "Value for " + field,
			MarkdownDescription: "Value for " + field,
			Optional:            true,
		}
	}
	return nil
}

// TypeToMap converts any param into the specified type using its TFSDK tag.
//
// The function takes an interface{} parameter `t` and returns a map[string]map[attr.Type]bool.
// The map contains the field names as keys and a nested map as values.
// The nested map contains the converted type of the field as the key and a boolean value.
// The boolean value is true if the field is a slice or array, and false otherwise.
//
// Parameters:
// - t: The interface{} parameter to be converted.
//
// Returns:
// - map[string]map[attr.Type]bool: The converted map.
func TypeToMap(t interface{}) map[string]map[attr.Type]bool {
	r := reflect.TypeOf(t)
	m := make(map[string]map[attr.Type]bool)

	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		convType := ConvertType(field.Type)
		mTwo := make(map[attr.Type]bool)
		if convType == nil {
			continue
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			mTwo[convType] = true
		} else {
			mTwo[convType] = false
		}
		m[field.Tag.Get("tfsdk")] = mTwo
	}

	return m
}

// ConvertType converts a reflect.Type to an attr.Type.
//
// It takes in a reflect.Type and checks its kind. If the type is a
// slice or array, it recursively calls itself with the element type.
// It then checks the name of the type and returns the corresponding
// attr.Type. If the type is none of the above, it returns nil.
//
// Parameters:
// - intialType: The reflect.Type to be converted.
//
// Returns:
// - attr.Type: The converted attr.Type.
func ConvertType(intialType reflect.Type) attr.Type {
	if intialType.Kind() == reflect.Slice || intialType.Kind() == reflect.Array {
		return ConvertType(intialType.Elem())
	}
	switch intialType.Name() {
	case "StringValue":
		return types.StringType
	case "Int64Value":
		return types.Int64Type
	case "BoolValue":
		return types.BoolType
	}

	return nil
}

// isPointer checks if the given value is a pointer.
//
// value: The value to check.
// Returns: A boolean indicating whether the value is a pointer.
func isPointer(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}
