package cel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type hola struct {
	Name  string      `tfsdk:"name"`
	Val   types.Int32 `tfsdk:"val"`
	Story hola2       `tfsdk:"story"`
}

type hola2 struct {
	Name types.String `tfsdk:"name"`
	Val  int          `tfsdk:"val"`
}

func filterHolas(inputs []hola) ([]hola, error) {
	return filterCel(inputs, `self.story.name == "hello int"`)
}

func TestFilter(t *testing.T) {
	outs, err := filterHolas([]hola{
		{
			Name: "hello",
			Val:  types.Int32Value(1),
			Story: hola2{
				Name: types.StringValue("hello int"),
				Val:  10,
			},
		},
		{
			Name: "world",
			Val:  types.Int32Value(2),
			Story: hola2{
				Name: types.StringValue("world int"),
				Val:  20,
			},
		},
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(outs) != 1 {
		t.Errorf("expected 1, got %d", len(outs))
	}
}
