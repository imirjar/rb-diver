package models

type ReportRoles struct {
	ReportID string
	RoleID   string
}

func (rr *ReportRoles) getRoles() []Role {
	return []Role{}
}

func (rr *ReportRoles) getReports() []Report {
	return []Report{}
}
