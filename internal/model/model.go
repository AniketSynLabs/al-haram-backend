package model

// Package represents an Umrah package tier.
type Package struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Tier           string   `json:"tier"`
	Nights         int      `json:"nights"`
	Price          string   `json:"price"`
	Image          string   `json:"image"`
	FlightType     string   `json:"flightType"`
	MakkahHotel    string   `json:"makkahHotel"`
	MakkahDist     string   `json:"makkahDistance"`
	MadinahHotel   string   `json:"madinahHotel"`
	MadinahDist    string   `json:"madinahDistance"`
	ShuttleMakkah  string   `json:"shuttleMakkah"`
	ShuttleMadinah string   `json:"shuttleMadinah"`
	Occupancy      string   `json:"occupancy"`
	Features       []string `json:"features"`
	Documents      []string `json:"documents"`
	RateNote       string   `json:"rateNote"`
	SortOrder      int      `json:"sortOrder"`
}

// Service represents a travel service offered.
type Service struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	Price        string `json:"price"`
	Requirements string `json:"requirements"`
	SortOrder    int    `json:"sortOrder"`
	IsActive     bool   `json:"isActive"`
}

// ScheduleEntry represents a single Umrah departure slot.
type ScheduleEntry struct {
	ID             string `json:"id"`
	Season         string `json:"season"`
	DepartureDate  string `json:"departureDate"`
	ReturnDate     string `json:"returnDate"`
	DurationNights int    `json:"durationNights"`
	Airline        string `json:"airline"`
	Route          string `json:"route"`
	PackageTier    string `json:"packageTier"`
	Rate           string `json:"rate"`
	SeatsLeft      int    `json:"seatsLeft"`
	Status         string `json:"status"`
	Notes          string `json:"notes"`
}

// Enquiry is a customer lead captured from the website.
type Enquiry struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	PackageID     string `json:"packageId"`
	ServiceID     string `json:"serviceId"`
	Message       string `json:"message"`
	Travellers    string `json:"travellers"`
	PreferredDate string `json:"preferredDate"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
}

// BankDetails holds the agency bank account information.
type BankDetails struct {
	AccountName   string `json:"accountName"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	IFSCCode      string `json:"ifscCode"`
	Branch        string `json:"branch"`
	Note          string `json:"note"`
}

// Policy holds a single policy document (terms / refund / privacy).
type Policy struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// GalleryItem is a single image or video in the public gallery.
type GalleryItem struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
	CreatedAt string `json:"createdAt"`
}
