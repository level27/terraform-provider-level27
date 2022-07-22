package level27

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func parseID(id string, diags *diag.Diagnostics) (int, bool) {
	val, err := strconv.Atoi(id)
	if err != nil {
		diags.AddError("Invalid entity ID", fmt.Sprintf("'%s' is not a valid numeric entity ID", id))
		return 0, false
	}

	return val, true
}

func appComponentParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	return parts[0], parts[1], nil
}

func appComponentURLParseID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2:attribute3", id)
	}

	return parts[0], parts[1], parts[2], nil
}
