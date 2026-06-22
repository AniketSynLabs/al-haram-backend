package service

import (
	"time"

	"al-haram/internal/db"
	"al-haram/internal/model"
)

func CreateEnquiry(e model.Enquiry) (model.Enquiry, error) {
	var id int
	err := db.DB.QueryRow(`
		INSERT INTO enquiries
			(name,phone,email,package_id,service_id,message,travellers,preferred_date)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		e.Name, e.Phone, e.Email, e.PackageID, e.ServiceID,
		e.Message, e.Travellers, e.PreferredDate,
	).Scan(&id)
	if err != nil {
		return model.Enquiry{}, err
	}
	e.ID = id
	e.Status = "New"
	e.CreatedAt = time.Now().Format(time.RFC3339)
	return e, nil
}

func ListEnquiries(status string) ([]model.Enquiry, error) {
	query := `
		SELECT id,name,phone,email,package_id,service_id,
		       message,travellers,preferred_date,status,created_at
		FROM enquiries`
	args := []any{}
	if status != "" {
		query += ` WHERE status=$1`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enquiries []model.Enquiry
	for rows.Next() {
		var e model.Enquiry
		var createdAt time.Time
		if err := rows.Scan(
			&e.ID, &e.Name, &e.Phone, &e.Email,
			&e.PackageID, &e.ServiceID, &e.Message,
			&e.Travellers, &e.PreferredDate, &e.Status, &createdAt,
		); err != nil {
			return nil, err
		}
		e.CreatedAt = createdAt.Format(time.RFC3339)
		enquiries = append(enquiries, e)
	}
	if enquiries == nil {
		enquiries = []model.Enquiry{}
	}
	return enquiries, nil
}

func UpdateEnquiryStatus(id, status string) error {
	_, err := db.DB.Exec(`UPDATE enquiries SET status=$1 WHERE id=$2`, status, id)
	return err
}
