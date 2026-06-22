package service

import (
	"time"

	"al-haram/internal/db"
	"al-haram/internal/model"
)

func ListSchedule() ([]model.ScheduleEntry, error) {
	rows, err := db.DB.Query(`
		SELECT id,season,departure_date,return_date,duration_nights,
		       airline,route,package_tier,rate,seats_left,status,notes
		FROM umrah_schedule ORDER BY departure_date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.ScheduleEntry
	for rows.Next() {
		var s model.ScheduleEntry
		var dep, ret time.Time
		if err := rows.Scan(
			&s.ID, &s.Season, &dep, &ret, &s.DurationNights,
			&s.Airline, &s.Route, &s.PackageTier,
			&s.Rate, &s.SeatsLeft, &s.Status, &s.Notes,
		); err != nil {
			return nil, err
		}
		s.DepartureDate = dep.Format("2006-01-02")
		s.ReturnDate = ret.Format("2006-01-02")
		entries = append(entries, s)
	}
	if entries == nil {
		entries = []model.ScheduleEntry{}
	}
	return entries, nil
}

func CreateScheduleEntry(s model.ScheduleEntry) (model.ScheduleEntry, error) {
	if s.ID == "" {
		s.ID = "sch_" + time.Now().Format("20060102150405")
	}
	_, err := db.DB.Exec(`
		INSERT INTO umrah_schedule
			(id,season,departure_date,return_date,duration_nights,
			 airline,route,package_tier,rate,seats_left,status,notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		s.ID, s.Season, s.DepartureDate, s.ReturnDate, s.DurationNights,
		s.Airline, s.Route, s.PackageTier, s.Rate, s.SeatsLeft, s.Status, s.Notes,
	)
	return s, err
}

func UpdateScheduleEntry(id string, s model.ScheduleEntry) error {
	_, err := db.DB.Exec(`
		UPDATE umrah_schedule SET
			season=$1, departure_date=$2, return_date=$3, duration_nights=$4,
			airline=$5, route=$6, package_tier=$7, rate=$8,
			seats_left=$9, status=$10, notes=$11
		WHERE id=$12`,
		s.Season, s.DepartureDate, s.ReturnDate, s.DurationNights,
		s.Airline, s.Route, s.PackageTier, s.Rate,
		s.SeatsLeft, s.Status, s.Notes,
		id,
	)
	return err
}

func DeleteScheduleEntry(id string) error {
	_, err := db.DB.Exec(`DELETE FROM umrah_schedule WHERE id=$1`, id)
	return err
}
