package service

import (
	"encoding/json"

	"al-haram/internal/db"
	"al-haram/internal/model"
)

const packageCols = `
	id, title, tier, nights, price, image, flight_type,
	makkah_hotel, makkah_dist, madinah_hotel, madinah_dist,
	shuttle_makkah, shuttle_madinah, occupancy,
	features, documents, rate_note, sort_order`

func scanPackage(rows interface {
	Scan(...any) error
}) (model.Package, error) {
	var p model.Package
	var feats, docs string
	err := rows.Scan(
		&p.ID, &p.Title, &p.Tier, &p.Nights, &p.Price, &p.Image,
		&p.FlightType, &p.MakkahHotel, &p.MakkahDist,
		&p.MadinahHotel, &p.MadinahDist,
		&p.ShuttleMakkah, &p.ShuttleMadinah, &p.Occupancy,
		&feats, &docs, &p.RateNote, &p.SortOrder,
	)
	if err != nil {
		return p, err
	}
	json.Unmarshal([]byte(feats), &p.Features)
	json.Unmarshal([]byte(docs), &p.Documents)
	return p, nil
}

func marshalPackage(p model.Package) (string, string) {
	feats, _ := json.Marshal(p.Features)
	docs, _ := json.Marshal(p.Documents)
	return string(feats), string(docs)
}

func ListPackages() ([]model.Package, error) {
	rows, err := db.DB.Query(`SELECT ` + packageCols + ` FROM packages ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pkgs []model.Package
	for rows.Next() {
		p, err := scanPackage(rows)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, p)
	}
	if pkgs == nil {
		pkgs = []model.Package{}
	}
	return pkgs, nil
}

func CreatePackage(p model.Package) (model.Package, error) {
	feats, docs := marshalPackage(p)
	_, err := db.DB.Exec(
		`INSERT INTO packages (`+packageCols[2:]+`) VALUES`+
			` ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		p.ID, p.Title, p.Tier, p.Nights, p.Price, p.Image,
		p.FlightType, p.MakkahHotel, p.MakkahDist,
		p.MadinahHotel, p.MadinahDist,
		p.ShuttleMakkah, p.ShuttleMadinah, p.Occupancy,
		feats, docs, p.RateNote, p.SortOrder,
	)
	return p, err
}

func UpdatePackage(id string, p model.Package) error {
	feats, docs := marshalPackage(p)
	_, err := db.DB.Exec(`
		UPDATE packages SET
			title=$1,        tier=$2,     nights=$3,         price=$4,
			image=$5,        flight_type=$6,
			makkah_hotel=$7, makkah_dist=$8,
			madinah_hotel=$9, madinah_dist=$10,
			shuttle_makkah=$11, shuttle_madinah=$12, occupancy=$13,
			features=$14,    documents=$15, rate_note=$16,   sort_order=$17
		WHERE id=$18`,
		p.Title, p.Tier, p.Nights, p.Price,
		p.Image, p.FlightType,
		p.MakkahHotel, p.MakkahDist,
		p.MadinahHotel, p.MadinahDist,
		p.ShuttleMakkah, p.ShuttleMadinah, p.Occupancy,
		feats, docs, p.RateNote, p.SortOrder,
		id,
	)
	return err
}

func DeletePackage(id string) error {
	_, err := db.DB.Exec(`DELETE FROM packages WHERE id=$1`, id)
	return err
}
