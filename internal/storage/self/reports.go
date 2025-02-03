package self

import (
	"context"
	"fmt"

	"github.com/imirjar/rb-diver/internal/models"
	_ "modernc.org/sqlite"
)

func (r *SelfStore) GetReport(ctx context.Context, id string) (models.Report, error) {

	var report models.Report
	// var role sql.NullInt64

	row := r.dbConn.QueryRowContext(ctx, "SELECT id, name, description, query FROM reports WHERE id=$1;", id)
	err := row.Scan(&report.Id, &report.Name, &report.Description, &report.Query)
	if err != nil {
		return models.Report{}, err
	}

	return report, nil
}

func (r *SelfStore) GetReports(ctx context.Context) ([]models.Report, error) {

	query := "SELECT id, name, description FROM reports"

	rows, err := r.dbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var reports []models.Report
	for rows.Next() {

		var rep models.Report
		err = rows.Scan(&rep.Id, &rep.Name, &rep.Description)
		if err != nil {
			return nil, fmt.Errorf("STORAGE: scan error %s", err)
		}
		reports = append(reports, rep)
	}

	return reports, nil
}

func (r *SelfStore) GetReportsByRole(ctx context.Context, roles []string) ([]models.Report, error) {

	query := `SELECT r.id, r.name, r.description
		FROM reports r
		JOIN report_roles rr ON r.id = rr.report_id
		WHERE rr.role_name = "admin";
		`

	rows, err := r.dbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var reports []models.Report
	for rows.Next() {

		var rep models.Report
		err = rows.Scan(&rep.Id, &rep.Name, &rep.Description)
		if err != nil {
			return nil, fmt.Errorf("STORAGE: scan error %s", err)
		}
		reports = append(reports, rep)
	}

	return reports, nil
}
