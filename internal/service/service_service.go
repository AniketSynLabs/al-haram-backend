package service

import (
	"al-haram/internal/db"
	"al-haram/internal/model"
)

func ListServices() ([]model.Service, error) {
	rows, err := db.DB.Query(`
		SELECT id,title,description,icon,price,requirements,sort_order,is_active
		FROM services ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var svcs []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(
			&s.ID, &s.Title, &s.Description, &s.Icon,
			&s.Price, &s.Requirements, &s.SortOrder, &s.IsActive,
		); err != nil {
			return nil, err
		}
		svcs = append(svcs, s)
	}
	if svcs == nil {
		svcs = []model.Service{}
	}
	return svcs, nil
}

func UpdateService(id string, s model.Service) error {
	_, err := db.DB.Exec(`
		UPDATE services SET
			title=$1, description=$2, icon=$3, price=$4,
			requirements=$5, sort_order=$6, is_active=$7
		WHERE id=$8`,
		s.Title, s.Description, s.Icon, s.Price,
		s.Requirements, s.SortOrder, s.IsActive,
		id,
	)
	return err
}
