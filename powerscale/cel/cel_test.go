package cel

import (
	"terraform-provider-powerscale/powerscale/models"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type testStruct1 struct {
	Name  string      `tfsdk:"name"`
	Val   types.Int32 `tfsdk:"val"`
	Story testStruct2 `tfsdk:"story"`
}

type testStruct2 struct {
	Name types.String `tfsdk:"name"`
	Val  int          `tfsdk:"val"`
}

func TestFilter(t *testing.T) {
	outs, err := filterCel(
		[]testStruct1{
			{
				Name: "hello",
				Val:  types.Int32Value(1),
				Story: testStruct2{
					Name: types.StringValue("hello int"),
					Val:  10,
				},
			},
			{
				Name: "world",
				Val:  types.Int32Value(2),
				Story: testStruct2{
					Name: types.StringValue("world int"),
					Val:  20,
				},
			},
		},
		`self.story.name == "hello int"`,
	)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(outs) != 1 {
		t.Errorf("expected 1, got %d", len(outs))
	}
}

func TestConvertList(t *testing.T) {
	inpModel := []models.FilePoolPolicyDetailModel{
		{
			Actions: []models.V1FilepoolDefaultPolicyAction{
				{
					ActionType: types.StringValue("set_requested_protection"),
				},
			},
		},
		{
			Actions: []models.V1FilepoolDefaultPolicyAction{
				{
					ActionType: types.StringValue("set_data_access_pattern"),
				},
			},
		},
		{
			Name: types.StringValue("policy_test"),
		},
	}

	inpStr := "has(self.actions[0].action_type)"

	outModels, err := FilterOptionalCel(inpModel, &inpStr)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(outModels) != 2 {
		t.Errorf("expected 2, got %d", len(outModels))
	}
}

func TestConvertListModel(t *testing.T) {
	inpModel := []models.FilePoolPolicyDetailModel{
		{
			FileMatchingPattern: &models.V1FilepoolPolicyFileMatchingPattern{
				OrCriteria: []models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem{
					{
						AndCriteria: []models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem{
							{
								Operator: types.StringValue(">"),
								Type:     types.StringValue("name"),
								Value:    types.StringValue("hello"),
							},
						},
					},
				},
			},
		},
		{
			FileMatchingPattern: &models.V1FilepoolPolicyFileMatchingPattern{
				OrCriteria: []models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem{
					{
						AndCriteria: []models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem{
							{
								Operator: types.StringValue(">"),
								Type:     types.StringValue("name"),
								Value:    types.StringValue("hello"),
							},
						},
					},
				},
			},
		},
		{
			FileMatchingPattern: &models.V1FilepoolPolicyFileMatchingPattern{
				OrCriteria: []models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem{
					{
						AndCriteria: []models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem{
							{
								Operator: types.StringValue("=="),
								Type:     types.StringValue("name"),
								Value:    types.StringValue("hello"),
							},
						},
					},
				},
			},
		},
	}

	inpStr := `
	has(self.file_matching_pattern) && 
	self.file_matching_pattern.or_criteria.exists(
		oc, oc.and_criteria.all(
			ac, ac.operator == ">" || ac.operator == "<"
		)
	)
	`

	outModels, err := FilterOptionalCel(inpModel, &inpStr)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(outModels) != 2 {
		t.Errorf("expected 2, got %d", len(outModels))
	}
}
