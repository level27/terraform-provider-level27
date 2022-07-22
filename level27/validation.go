package level27

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func validateID() tfsdk.AttributeValidator {
	return idValidator{}
}

type idValidator struct{}

func (idValidator) Description(context.Context) string {
	return "Ensures values are well-formed numeric IDs for the API"
}

func (i idValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (idValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	// I'm assuming that terraform enforces the type is correct here.
	str := req.AttributeConfig.(types.String)

	_, err := strconv.Atoi(str.Value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.AttributePath, "Invalid entity ID", "Value must be a valid integer ID")
	}
}

func validateTwoID() tfsdk.AttributeValidator {
	return twoIdValidator{}
}

type twoIdValidator struct{}

func (twoIdValidator) Description(context.Context) string {
	return "Ensures values are well-formed combined IDs, in id1:id2 form"
}

func (i twoIdValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (twoIdValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	// I'm assuming that terraform enforces the type is correct here.
	str := req.AttributeConfig.(types.String)

	_, _, err := appComponentParseID(str.Value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.AttributePath, "Invalid entity ID", "Value must be a valid combined ID in id1:id2 form")
	}
}

