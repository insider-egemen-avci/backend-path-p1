package services

import (
	"fmt"

	"insider-egemen-avci/backend-path-p1/internal/models"
)

func CheckRole(user *models.User, requiredRoles ...string) error {
	if len(requiredRoles) == 0 {
		return nil
	}

	for _, role := range requiredRoles {
		if role == user.Role {
			return nil
		}
	}

	return fmt.Errorf("unauthorized")
}
