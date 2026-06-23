package service

import (
	"sync"
	"time"

	"al-haram/internal/db"
	"al-haram/internal/model"
)

// BootstrapData is everything the public site needs on first load.
type BootstrapData struct {
	Settings    map[string]string  `json:"settings"`
	BankDetails model.BankDetails  `json:"bankDetails"`
	Packages    []model.Package    `json:"packages"`
	Services    []model.Service    `json:"services"`
	Schedule    []model.ScheduleEntry `json:"schedule"`
	Policies    []model.Policy     `json:"policies"`
	Gallery     []model.GalleryItem `json:"gallery"`
}

// cache holds the last successful bootstrap result so repeated page loads
// never hit the database.
var (
	cacheMu      sync.RWMutex
	cachedData   *BootstrapData
	cacheExpires time.Time
	cacheTTL     = 60 * time.Second
)

func InvalidateBootstrapCache() {
	cacheMu.Lock()
	cachedData = nil
	cacheMu.Unlock()
}

func GetBootstrap() (*BootstrapData, error) {
	// Return cached result if still fresh.
	cacheMu.RLock()
	if cachedData != nil && time.Now().Before(cacheExpires) {
		d := cachedData
		cacheMu.RUnlock()
		return d, nil
	}
	cacheMu.RUnlock()

	// One connection, one transaction — read everything.
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Settings
	srows, err := tx.Query(`SELECT key, value FROM site_settings`)
	if err != nil {
		return nil, err
	}
	settings := map[string]string{}
	for srows.Next() {
		var k, v string
		if err := srows.Scan(&k, &v); err != nil {
			srows.Close()
			return nil, err
		}
		settings[k] = v
	}
	srows.Close()

	// Bank details
	var bank model.BankDetails
	_ = tx.QueryRow(`SELECT account_name,bank_name,account_no,ifsc,branch,note FROM bank_details WHERE id=1`).
		Scan(&bank.AccountName, &bank.BankName, &bank.AccountNumber, &bank.IFSCCode, &bank.Branch, &bank.Note)

	// Packages
	prows, err := tx.Query(`SELECT ` + packageCols + ` FROM packages ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	var packages []model.Package
	for prows.Next() {
		p, err := scanPackage(prows)
		if err != nil {
			prows.Close()
			return nil, err
		}
		packages = append(packages, p)
	}
	prows.Close()
	if packages == nil {
		packages = []model.Package{}
	}

	// Services
	svrows, err := tx.Query(`SELECT id,title,description,icon,price,requirements,sort_order,is_active FROM services ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	var services []model.Service
	for svrows.Next() {
		var s model.Service
		if err := svrows.Scan(&s.ID, &s.Title, &s.Description, &s.Icon, &s.Price, &s.Requirements, &s.SortOrder, &s.IsActive); err != nil {
			svrows.Close()
			return nil, err
		}
		services = append(services, s)
	}
	svrows.Close()
	if services == nil {
		services = []model.Service{}
	}

	// Schedule
	scrows, err := tx.Query(`SELECT id,season,departure_date,return_date,duration_nights,airline,route,package_tier,rate,seats_left,status,notes FROM umrah_schedule ORDER BY departure_date`)
	if err != nil {
		return nil, err
	}
	var schedule []model.ScheduleEntry
	for scrows.Next() {
		var e model.ScheduleEntry
		if err := scrows.Scan(&e.ID, &e.Season, &e.DepartureDate, &e.ReturnDate, &e.DurationNights, &e.Airline, &e.Route, &e.PackageTier, &e.Rate, &e.SeatsLeft, &e.Status, &e.Notes); err != nil {
			scrows.Close()
			return nil, err
		}
		schedule = append(schedule, e)
	}
	scrows.Close()
	if schedule == nil {
		schedule = []model.ScheduleEntry{}
	}

	// Policies
	porows, err := tx.Query(`SELECT id,title,content FROM policies ORDER BY id`)
	if err != nil {
		return nil, err
	}
	var policies []model.Policy
	for porows.Next() {
		var p model.Policy
		if err := porows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			porows.Close()
			return nil, err
		}
		policies = append(policies, p)
	}
	porows.Close()
	if policies == nil {
		policies = []model.Policy{}
	}

	// Gallery
	grows, err := tx.Query(`SELECT id,url,type,caption,sort_order,created_at FROM gallery_items ORDER BY sort_order,created_at`)
	if err != nil {
		return nil, err
	}
	var gallery []model.GalleryItem
	for grows.Next() {
		var g model.GalleryItem
		if err := grows.Scan(&g.ID, &g.URL, &g.Type, &g.Caption, &g.SortOrder, &g.CreatedAt); err != nil {
			grows.Close()
			return nil, err
		}
		gallery = append(gallery, g)
	}
	grows.Close()
	if gallery == nil {
		gallery = []model.GalleryItem{}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := &BootstrapData{
		Settings:    settings,
		BankDetails: bank,
		Packages:    packages,
		Services:    services,
		Schedule:    schedule,
		Policies:    policies,
		Gallery:     gallery,
	}

	cacheMu.Lock()
	cachedData = result
	cacheExpires = time.Now().Add(cacheTTL)
	cacheMu.Unlock()

	return result, nil
}
