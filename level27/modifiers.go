package level27

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func DefaultValue(val attr.Value) tfsdk.AttributePlanModifier {
	if val.IsNull() {
		panic("DefaultValue does not work with null values.")
	}

	return DefaultValueModifier{
		DefaultValue: val,
	}
}

type DefaultValueModifier struct {
	DefaultValue attr.Value
}

func (DefaultValueModifier) Description(context.Context) string {
	return "Provides a default value for an attribute"
}

func (DefaultValueModifier) MarkdownDescription(context.Context) string {
	return "Provides a default value for an attribute"
}

func (m DefaultValueModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig.IsNull() && req.AttributePlan.IsUnknown() {
		resp.AttributePlan = m.DefaultValue
	}
}

// UseStateForUnknown returns a UseStateForUnknownModifier.
func UseStateForUnknown2() tfsdk.AttributePlanModifier {
	return UseStateForUnknownModifier2{}
}

