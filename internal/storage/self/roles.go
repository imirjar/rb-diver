package self

import (
	"context"
	"fmt"

	"github.com/imirjar/rb-diver/internal/models"
)

func (r *SelfStore) GetRoles(ctx context.Context) ([]models.Role, error) {

	// if reportID
	// 1) get from report_roles role id where report id is
	// 2) eqit query + WHERE id in this ids

	query := "SELECT id, name FROM roles"

	rows, err := r.dbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	for rows.Next() {

		var role models.Role
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, fmt.Errorf("STORAGE: scan error %s", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *SelfStore) GetRolesByReportID(ctx context.Context, repID string) ([]models.Role, error) {

	query := fmt.Sprintf(`SELECT ro.*
		FROM roles ro
		JOIN report_roles rr ON ro.id = rr.role_id
		WHERE rr.report_id = %s;
		`, repID)

	rows, err := r.dbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	for rows.Next() {

		var role models.Role
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, fmt.Errorf("STORAGE: scan error %s", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}
