package cel

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func convertGoToCelMap(in any) any {
	// use reflection to convert go type to cel type
	val := reflect.ValueOf(in)
	switch val.Kind() {
	case reflect.Ptr:
		// if pointer and not nil, deref the pointer and recurse
		if val.IsNil() {
			return nil
		}
		return convertGoToCelMap(val.Elem().Interface())
	case reflect.Struct:
		// if struct
		switch inV := in.(type) {
		case attr.Value:
			// if attr.Value, convert using terraform conversion rules
			if inV.IsNull() {
				return nil
			}
			return convertTerraformValueToCelType(inV)
		}
		return convertStructToCelMap(in)
	case reflect.Map:
		// if map, recurse for each value
		out := make(map[string]any, val.Len())
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := iter.Value()
			out[k] = convertGoToCelMap(v.Interface())
		}
		return out
	case reflect.Slice, reflect.Array:
		// if slice, recurse for each value
		out := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			out[i] = convertGoToCelMap(val.Index(i).Interface())
		}
		return out
	}

	return in
}

func convertTerraformValueToCelType(in attr.Value) any {
	switch inV := in.(type) {
	case types.String:
		return inV.ValueString()
	case types.Bool:
		return inV.ValueBool()
	case types.Int64:
		return inV.ValueInt64()
	case types.Int32:
		return inV.ValueInt32()
	case types.Float64:
		return inV.ValueFloat64()
	case types.Number:
		vbf := inV.ValueBigFloat()
		vf, _ := vbf.Float64()
		return vf
	case types.List:
		out := make([]any, 0)
		for _, v := range inV.Elements() {
			out = append(out, convertTerraformValueToCelType(v))
		}
		return out
	case types.Map:
		out := make(map[string]any)
		for k, v := range inV.Elements() {
			out[k] = convertTerraformValueToCelType(v)
		}
		return out
	case types.Object:
		out := make(map[string]any)
		for k, v := range inV.Attributes() {
			out[k] = convertTerraformValueToCelType(v)
		}
		return out
	default:
		panic("Found unknown Terraform Type " + in.Type(context.Background()).String())
	}
}

func convertStructToCelMap(in any) map[string]any {
	// use reflection to convert go struct to cel map
	val := reflect.ValueOf(in)
	itype := val.Type()
	ret := make(map[string]any)
	for i := 0; i < val.NumField(); i++ {
		value := convertGoToCelMap(val.Field(i).Interface())
		if value == nil {
			continue
		}
		// get field tfsdk tag
		field := itype.Field(i)
		tfTag := field.Tag.Get("tfsdk")
		ret[tfTag] = value
	}
	return ret
}

func filterCel[T any](inputs []T, filter string) ([]T, error) {
	env, err := cel.NewEnv( //cel.Variable("name", cel.StringType),
		cel.Variable("self", cel.MapType(cel.StringType, cel.DynType)))
	// Check err for environment setup errors.
	if err != nil {
		log.Fatalln(err)
	}
	ast, iss := env.Compile(filter)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	// cel.Types()
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalln(err)
	}

	ret := make([]T, 0)
	var ferr error
	for _, input := range inputs {
		mval := convertGoToCelMap(input)
		jval, _ := json.MarshalIndent(mval, "", "  ")
		log.Println("Input:", string(jval))
		out, _, err := prg.Eval(map[string]any{
			"self": mval,
		})
		if err != nil {
			ferr = errors.Join(ferr, err)
			continue
		}
		if out.Value().(bool) {
			ret = append(ret, input)
		}
	}

	if ferr != nil {
		return nil, ferr
	}

	return ret, nil
}

func FilterOptionalCel[T any](inputs []T, filter *string) ([]T, error) {
	if filter == nil {
		return inputs, nil
	}
	return filterCel(inputs, *filter)
}
