package service

import (
	"encoding/json"

	"al-haram/internal/db"
	"al-haram/internal/model"
)

func ListPackages() ([]model.Package, error) {
	rows, err := db.DB.Query(`
		SELECT id,title,tier,nights,price,image,flight_type,
		       makkah_hotel,makkah_dist,madinah_hotel,madinah_dist,
		       shuttle_makkah,shuttle_madinah,occupancy,
		       features,documents,rate_note,sort_order
		FROM packages ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pkgs []model.Package
	for rows.Next() {
		var p model.Package
		var featsJSON, docsJSON string
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Tier, &p.Nights, &p.Price, &p.Image,
			&p.FlightType, &p.MakkahHotel, &p.MakkahDist,
			&p.MadinahHotel, &p.MadinahDist,
			&p.ShuttleMakkah, &p.ShuttleMadinah, &p.Occupancy,
			&featsJSON, &docsJSON, &p.RateNote, &p.SortOrder,
		); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(featsJSON), &p.Features)
		json.Unmarshal([]byte(docsJSON), &p.Documents)
		pkgs = append(pkgs, p)
	}
	if pkgs == nil {
		pkgs = []model.Package{}
	}
	return pkgs, nil
}

func UpdatePackage(id string, p model.Package) error {
	featsJSON, _ := json.Marshal(p.Features)
	docsJSON, _ := json.Marshal(p.Documents)
	_, err := db.DB.Exec(`
		UPDATE packages SET
			title=$1, tier=$2, nights=$3, price=$4, image=$5,
			flight_type=$6, makkah_hotel=$7, makkah_dist=$8,
			madinah_hotel=$9, madinah_dist=$10,
			shuttle_makkah=$11, shuttle_madinah=$12, occupancy=$13,
			features=$14, documents=$15, rate_note=$16, sort_order=$17
		WHERE id=$18`,
		p.Title, p.Tier, p.Nights, p.Price, p.Image,
		p.FlightType, p.MakkahHotel, p.MakkahDist,
		p.MadinahHotel, p.MadinahDist,
		p.ShuttleMakkah, p.ShuttleMadinah, p.Occupancy,
		string(featsJSON), string(docsJSON), p.RateNote, p.SortOrder,
		id,
	)
	return err
}
