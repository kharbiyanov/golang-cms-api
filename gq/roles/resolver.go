package main

import (
	"cms-api/models"
	"cms-api/utils"
)

func GetRoleList() (interface{}, error) {
	var roles []models.Role
	for _, roleName := range utils.Roles.GetAllRoles() {
		role := models.Role{
			Name: roleName,
		}
		for _, p := range utils.Roles.GetFilteredPolicy(0, roleName) {
			access := models.Access{
				Object: p[1],
				Action: p[2],
			}
			role.Access = append(role.Access, access)
		}
		roles = append(roles, role)
	}
	return roles, nil
}
