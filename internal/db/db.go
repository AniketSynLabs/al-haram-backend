package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const schemaName = "al_haram"

var DB *sql.DB

func Connect(dsn string) error {
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err = DB.PingContext(context.Background()); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	log.Println("✅ Database connected")
	return nil
}

func Migrate() error {
	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	log.Println("✅ Schema migrated (" + schemaName + ")")
	return nil
}

const schema = `
-- ── Schema ───────────────────────────────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS al_haram;
SET search_path TO al_haram, public;

-- ── Tables ───────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS al_haram.site_settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS al_haram.packages (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    tier            TEXT NOT NULL,
    nights          INT  NOT NULL DEFAULT 15,
    price           TEXT NOT NULL DEFAULT 'On Request',
    image           TEXT NOT NULL DEFAULT '',
    flight_type     TEXT NOT NULL DEFAULT '',
    makkah_hotel    TEXT NOT NULL DEFAULT '',
    makkah_dist     TEXT NOT NULL DEFAULT '',
    madinah_hotel   TEXT NOT NULL DEFAULT '',
    madinah_dist    TEXT NOT NULL DEFAULT '',
    shuttle_makkah  TEXT NOT NULL DEFAULT '',
    shuttle_madinah TEXT NOT NULL DEFAULT '',
    occupancy       TEXT NOT NULL DEFAULT '',
    features        TEXT NOT NULL DEFAULT '[]',
    documents       TEXT NOT NULL DEFAULT '[]',
    rate_note       TEXT NOT NULL DEFAULT '',
    sort_order      INT  NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS al_haram.services (
    id           TEXT    PRIMARY KEY,
    title        TEXT    NOT NULL,
    description  TEXT    NOT NULL DEFAULT '',
    icon         TEXT    NOT NULL DEFAULT '',
    price        TEXT    NOT NULL DEFAULT '',
    requirements TEXT    NOT NULL DEFAULT '',
    sort_order   INT     NOT NULL DEFAULT 0,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS al_haram.umrah_schedule (
    id              TEXT PRIMARY KEY,
    season          TEXT NOT NULL,
    departure_date  DATE NOT NULL,
    return_date     DATE NOT NULL,
    duration_nights INT  NOT NULL,
    airline         TEXT NOT NULL DEFAULT '',
    route           TEXT NOT NULL DEFAULT '',
    package_tier    TEXT NOT NULL DEFAULT '',
    rate            TEXT NOT NULL DEFAULT 'On Request',
    seats_left      INT  NOT NULL DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'Available',
    notes           TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS al_haram.enquiries (
    id             SERIAL      PRIMARY KEY,
    name           TEXT        NOT NULL,
    phone          TEXT        NOT NULL,
    email          TEXT        NOT NULL DEFAULT '',
    package_id     TEXT        NOT NULL DEFAULT '',
    service_id     TEXT        NOT NULL DEFAULT '',
    message        TEXT        NOT NULL DEFAULT '',
    travellers     TEXT        NOT NULL DEFAULT '1',
    preferred_date TEXT        NOT NULL DEFAULT '',
    status         TEXT        NOT NULL DEFAULT 'New',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS al_haram.bank_details (
    id           INT  PRIMARY KEY DEFAULT 1,
    account_name TEXT NOT NULL DEFAULT '',
    bank_name    TEXT NOT NULL DEFAULT '',
    account_no   TEXT NOT NULL DEFAULT '',
    ifsc         TEXT NOT NULL DEFAULT '',
    branch       TEXT NOT NULL DEFAULT '',
    note         TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS al_haram.policies (
    id      TEXT PRIMARY KEY,
    title   TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT ''
);

-- ── Seed data ─────────────────────────────────────────────────────────────────
INSERT INTO al_haram.site_settings (key, value) VALUES
    ('brand_name',    'Faxman Travels'),
    ('tagline',       'Sacred Journeys Since 2004'),
    ('hero_title',    'Your Journey to Makkah Starts Here'),
    ('hero_subtitle', 'Economy to Super Deluxe Umrah packages with flights, visa, hotels & full Ziyarat. 22+ years of safe, affordable, hassle-free travel.'),
    ('about_title',   'Serving Pilgrims with Trust Since 2004'),
    ('about_body',    'Faxman Travels was founded in 2004 by Mr. Shaikh Iftekhar Ahmad with a clear vision — to provide reliable, trustworthy, and customer-friendly travel services to pilgrims from Fatehpur and beyond. Since 2018, the agency has grown under the management of Mr. Shaikh Zaid Ahmad, carrying the legacy forward with modern travel solutions and a commitment to innovation.'),
    ('whatsapp',      '+918418021570'),
    ('email',         'faxman.travels@gmail.com'),
    ('address',       'Basement of Dulhan Marriage Hall, Mahajari, Fatehpur (U.P.) — Pincode 212601'),
    ('map_url',       'https://maps.google.com/maps?q=Mahajari+Fatehpur+UP+212601&t=&z=14&ie=UTF8&iwloc=&output=embed')
ON CONFLICT (key) DO NOTHING;

INSERT INTO al_haram.bank_details (id, account_name, bank_name, account_no, ifsc, branch, note) VALUES (
    1,
    'Faxman Travels',
    'Bank of Baroda',
    'Contact us for details',
    'Contact us for details',
    'Fatehpur Branch, U.P.',
    'After payment, share the transaction screenshot on WhatsApp for confirmation.'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO al_haram.policies (id, title, content) VALUES
    ('terms',   'Terms & Conditions',
     'All bookings are subject to availability and confirmation. The customer is responsible for providing correct documents and valid information. The company is not responsible for delays caused by airlines, embassy decisions, immigration authorities, natural disasters, or any other circumstances beyond our control. Prices are based on current airline fares and hotel rates; any future increase in fares or government regulations will be borne by the pilgrim. Customers must verify all travel documents before departure. Passport validity of minimum 6 months is mandatory.'),
    ('refund',  'Refund Policy',
     'Refunds for air tickets, visas, hotels, and packages are subject to individual airline rules, embassy policies, hotel cancellation terms, and applicable service charges. Some services may be non-refundable once confirmed. Refund processing times vary depending on the provider. Visa fees are non-refundable once submitted to the embassy. Customers are strongly advised to read all cancellation and refund terms carefully before making payment.'),
    ('privacy', 'Privacy Policy',
     'Customer data including passport details, phone numbers, email addresses, travel documents, and payment information is kept strictly confidential. We do not share personal information with third parties without explicit permission, except where required for visa processing, airline booking, or legal compliance. Basic site usage data may be collected to improve service quality. You may request deletion of your personal data at any time by contacting us directly.')
ON CONFLICT (id) DO NOTHING;

-- ── Seed packages ─────────────────────────────────────────────────────────────
INSERT INTO al_haram.packages
    (id, title, tier, nights, price, image, flight_type,
     makkah_hotel, makkah_dist, madinah_hotel, madinah_dist,
     shuttle_makkah, shuttle_madinah, occupancy,
     features, documents, rate_note, sort_order)
VALUES
('pkg_economy', 'Economy Umrah Package', 'Economy', 15, 'On Request',
 'https://i.pinimg.com/1200x/24/7b/36/247b36c6b1475fe7405647c54c4a616d.jpg',
 'Via Flight (connecting)',
 'Budget Hotel', '1500–2000 m',
 'Standard Hotel', 'Within 1 km',
 'Free 24×7 shuttle', 'Shuttle for Namaz only', '5–6 persons per room',
 '["Round-trip air ticket (via flight)","Umrah visa","Hotel on sharing basis (5–6 persons)","Breakfast, lunch & dinner","Local transfer by AC bus","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water","Free 24×7 Haram shuttle (Makkah)"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates based on current airline & hotel fares. Any future increase will be borne by the pilgrim.', 1),

('pkg_semi_deluxe', 'Semi Deluxe Umrah Package', 'Semi Deluxe', 15, 'On Request',
 'https://i.pinimg.com/736x/92/1f/f9/921ff99d18598b360b2a5705d1f7b653.jpg',
 'Option A: Direct deluxe flight + economy hotels | Option B: Via flight + deluxe hotels',
 'Economy / Deluxe Hotel (as per option)', '500–700 m (Option B)',
 'Economy / Deluxe Hotel (as per option)', '200–250 m (Option B)',
 'Shuttle provided', 'Shuttle provided', '4–5 persons per room',
 '["Round-trip air ticket","Umrah visa","Hotel on sharing basis (4–5 persons)","Breakfast, lunch & dinner","Local transfer by AC bus","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water","Two options: direct deluxe flight OR deluxe hotels"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates based on current airline & hotel fares. Any future increase will be borne by the pilgrim.', 2),

('pkg_deluxe', 'Deluxe Umrah Package', 'Deluxe', 15, 'On Request',
 'https://i.pinimg.com/736x/1e/f9/ed/1ef9ed0de60fb7ae04073b10f0b05951.jpg',
 'Direct / Deluxe Flight',
 'Deluxe Hotel', '500–700 m',
 'Deluxe Hotel', '200–250 m',
 'AC bus shuttle', 'AC bus shuttle', '4 persons per room',
 '["Round-trip air ticket (direct flight)","Umrah visa","Deluxe hotel on sharing basis (4 persons)","Breakfast, lunch & dinner","Local transfer by AC bus","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates based on current airline & hotel fares. Any future increase will be borne by the pilgrim.', 3),

('pkg_super_deluxe', 'Super Deluxe Umrah Package', 'Super Deluxe', 15, 'On Request',
 'https://i.pinimg.com/736x/28/43/23/28432334b5bf415ebd5e64aebc892a41.jpg',
 'Premium / Business Class Flight',
 'Le Méridien Towers Makkah (5-Star)', '100–150 m',
 'Super Deluxe Hotel', '100–150 m',
 'Private AC transport', 'Private AC transport', '2–3 persons per room',
 '["Round-trip premium/business flight","Umrah visa","5-Star hotel on sharing basis (2–3 persons)","Breakfast, lunch & dinner","Private AC transport","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water","Le Méridien Towers Makkah or equivalent"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates based on current airline & hotel fares. Any future increase will be borne by the pilgrim.', 4)
ON CONFLICT (id) DO UPDATE SET
    title           = EXCLUDED.title,
    tier            = EXCLUDED.tier,
    nights          = EXCLUDED.nights,
    price           = EXCLUDED.price,
    image           = EXCLUDED.image,
    flight_type     = EXCLUDED.flight_type,
    makkah_hotel    = EXCLUDED.makkah_hotel,
    makkah_dist     = EXCLUDED.makkah_dist,
    madinah_hotel   = EXCLUDED.madinah_hotel,
    madinah_dist    = EXCLUDED.madinah_dist,
    shuttle_makkah  = EXCLUDED.shuttle_makkah,
    shuttle_madinah = EXCLUDED.shuttle_madinah,
    occupancy       = EXCLUDED.occupancy,
    features        = EXCLUDED.features,
    documents       = EXCLUDED.documents,
    rate_note       = EXCLUDED.rate_note,
    sort_order      = EXCLUDED.sort_order;

INSERT INTO al_haram.services
    (id, title, description, icon, price, requirements, sort_order)
VALUES
    ('svc_saudi_visa',   'Saudi Visa Stamping',          'Hassle-free Saudi visa documentation and stamping assistance.',                  '🇸🇦', '', '', 1),
    ('svc_kuwait_visa',  'Kuwait Visa Stamping',          'Complete support for Kuwait visa documentation and stamping.',                   '🇰🇼', '', '', 2),
    ('svc_tourist_visa', 'Tourist Visa Assistance',       'Guidance and documentation support for tourist visas worldwide.',                '🌍', '', '', 3),
    ('svc_air_ticket',   'Air Ticket Booking',            'Domestic & international air ticket booking at competitive rates.',              '✈️', '', '', 4),
    ('svc_holiday',      'Holiday Packages',              'Curated holiday packages for families and individuals.',                         '🏖️', '', '', 5),
    ('svc_passport',     'Passport Assistance',           'New passport, renewal, and Tatkal application assistance.',                      '📘', '', '', 6),
    ('svc_emigration',   'Emigration Services',           'Emigration clearance and documentation for overseas workers.',                   '🛂', '', '', 7),
    ('svc_dummy_ticket', 'Dummy Ticket & Hotel Voucher',  'Verified dummy air tickets and hotel booking vouchers for visa purposes.',       '🎫', '', '', 8),
    ('svc_insurance',    'Travel Insurance',              'Affordable travel insurance plans for international and domestic trips.',        '🛡️', '', '', 9)
ON CONFLICT (id) DO NOTHING;
`
