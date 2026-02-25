// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// expandInt64List converts a types.List of Int64 into a []int.
func expandInt64List(ctx context.Context, list types.List, diags *diag.Diagnostics) []int {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var ids []types.Int64
	d := list.ElementsAs(ctx, &ids, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}
	out := make([]int, len(ids))
	for i, v := range ids {
		out[i] = int(v.ValueInt64())
	}
	return out
}
