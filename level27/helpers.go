package level27

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func parseID(id string, diags *diag.Diagnostics) (int, bool) {
	val, err := strconv.Atoi(id)
	if err != nil {
		diags.AddError("Invalid entity ID", fmt.Sprintf("'%s' is not a valid numeric entity ID", id))
		return 0, false
	}

	return val, true
}

func appComponentParseID(id string) (int, int, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	id0, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	id1, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	return id0, id1, nil
}

func appComponentURLParseID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2:attribute3", id)
	}

	return parts[0], parts[1], parts[2], nil
}

func organisationOrDefault(p *provider, orgAttr types.String, diags *diag.Diagnostics) int {
	if orgAttr.Null || orgAttr.Value == "" {
		// Organisation defaults to the org of your login
		login, err := p.GetLoginInfo()
		if err != nil {
			diags.AddError("Error fetching login info for default organisation ID", "API request failed:\n"+err.Error())
			return 0
		}

		return login.User.Organisation.ID
	}

	orgID, _ := strconv.Atoi(orgAttr.Value)
	return orgID
}

func tfString(val string) types.String {
	return types.String{Value: val}
}

func tfStringP(val *string) types.String {
	if val == nil {
		return types.String{Null: true}
	}

	return types.String{Value: *val}
}

func tfStringId(id int) types.String {
	return types.String{Value: strconv.Itoa(id)}
}

func tfStringIdP(id *int) types.String {
	if id == nil {
		return types.String{Null: true}
	}

	return tfStringId(*id)
}

func tfDStringP(val types.String) *string {
	if val.Null {
		return nil
	}
	return &val.Value
}

func tfDStringId(val types.String) int {
	id, err := strconv.Atoi(val.Value)
	if err != nil {
		panic("Invalid ID string: " + err.Error())
	}

	return id
}

func formatTeams(teams types.List) string {
	autoTeams := ""
	for i, v := range teams.Elems {
		if i != 0 {
			autoTeams += ","
		}

		autoTeams += v.(types.String).Value
	}

	return autoTeams
}
