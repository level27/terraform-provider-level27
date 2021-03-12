package level27

import (
	"fmt"
)

// User gets a user from the API
func (c *Client) User(method string, organisationID interface{}, userID interface{}, data interface{}) User {
	var user User
	// var users []User
	// var invitations []User

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("organisations/%v/users", organisationID)
		c.invokeAPI("GET", endpoint, nil, &user)
	case "CREATE":
		endpoint := "systems"
		c.invokeAPI("POST", endpoint, data, &user)
	case "UPDATE":
		endpoint := fmt.Sprintf("systems/%v", organisationID)
		c.invokeAPI("PUT", endpoint, data, &user)
	case "DELETE":
		endpoint := fmt.Sprintf("systems/%v", organisationID)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return user
}
