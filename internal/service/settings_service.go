package service

import (
	"al-haram/internal/db"
	"al-haram/internal/model"
)

func GetSettings() (map[string]string, error) {
	rows, err := db.DB.Query(`SELECT key, value FROM site_settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		settings[k] = v
	}
	return settings, nil
}

func UpdateSettings(settings map[string]string) error {
	for k, v := range settings {
		_, err := db.DB.Exec(`
			INSERT INTO site_settings (key,value) VALUES ($1,$2)
			ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value`, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBankDetails() (model.BankDetails, error) {
	var b model.BankDetails
	err := db.DB.QueryRow(`
		SELECT account_name,bank_name,account_no,ifsc,branch,note
		FROM bank_details WHERE id=1`).
		Scan(&b.AccountName, &b.BankName, &b.AccountNumber, &b.IFSCCode, &b.Branch, &b.Note)
	return b, err
}

func UpdateBankDetails(b model.BankDetails) error {
	_, err := db.DB.Exec(`
		UPDATE bank_details SET
			account_name=$1, bank_name=$2, account_no=$3,
			ifsc=$4, branch=$5, note=$6
		WHERE id=1`,
		b.AccountName, b.BankName, b.AccountNumber, b.IFSCCode, b.Branch, b.Note,
	)
	return err
}

func ListPolicies() ([]model.Policy, error) {
	rows, err := db.DB.Query(`SELECT id,title,content FROM policies ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []model.Policy
	for rows.Next() {
		var p model.Policy
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	if policies == nil {
		policies = []model.Policy{}
	}
	return policies, nil
}

func UpdatePolicy(id, title, content string) error {
	_, err := db.DB.Exec(`UPDATE policies SET title=$1,content=$2 WHERE id=$3`, title, content, id)
	return err
}
