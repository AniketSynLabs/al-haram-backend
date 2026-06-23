package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

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

	// Cap the pool well below Railway's per-service connection limit.
	// Railway Hobby Postgres allows ~25 total connections; leave headroom
	// for the Railway internal admin connections and future replicas.
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(30 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)

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

CREATE TABLE IF NOT EXISTS al_haram.gallery_items (
    id         TEXT PRIMARY KEY,
    url        TEXT NOT NULL,
    type       TEXT NOT NULL DEFAULT 'image',
    caption    TEXT NOT NULL DEFAULT '',
    sort_order INT  NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ── Seed data ─────────────────────────────────────────────────────────────────
INSERT INTO al_haram.site_settings (key, value) VALUES
    ('brand_name',    'Faxman Travels'),
    ('tagline',       'Make Your Next Trip Awesome'),
    ('hero_title',    'Your Journey to Makkah Starts Here'),
    ('hero_subtitle', 'Economy to Super Deluxe Umrah packages with flights, visa, hotels & full Ziyarat. Trusted by pilgrims since 2004.'),
    ('about_title',   'Serving Pilgrims with Trust Since 2004'),
    ('about_body',    'Faxman Travels was founded in 2004 by Mr. Shaikh Iftekhar Ahmad with a clear vision — to provide reliable, trustworthy, and customer-friendly travel services to pilgrims from Fatehpur and beyond. Since 2018, the agency has grown under the management of Mr. Shaikh Zaid Siddiqui, carrying the legacy forward with modern travel solutions and a commitment to innovation.'),
    ('whatsapp',      '+918418021570'),
    ('email',         'faxman.travels@gmail.com'),
    ('address',       'Basement of Dulhan Marriage Hall, Mahajari, Fatehpur (U.P.) — Pincode 212601'),
    ('map_url',       'https://maps.google.com/maps?q=Mahajari+Fatehpur+UP+212601&t=&z=14&ie=UTF8&iwloc=&output=embed'),
    ('hero_image_url',    'https://i.pinimg.com/736x/f1/4d/bb/f14dbb8116639f8b2f2e75f5c02737d1.jpg'),
    ('about_image_1_url', 'https://i.pinimg.com/736x/6c/6c/e3/6c6ce35320c8f2454f87e0ac83589832.jpg'),
    ('about_image_2_url', 'https://i.pinimg.com/736x/69/00/95/6900952b20e2d8e77c63779a9aec0747.jpg'),
    ('about_image_3_url', 'https://i.pinimg.com/1200x/08/47/31/0847310bb2c5f810476817e0858a3867.jpg'),
    ('gallery_title',     'A Glimpse of the Holy Lands'),
    ('gallery_subtitle',  'Beautiful moments from Makkah and Madinah.'),
    ('stats_years',       '22+'),
    ('stats_packages',    '4'),
    ('stats_services',    '9+'),
    ('stats_pilgrims',    '500+'),
    ('founder_1_name',    'Mr. Shaikh Iftekhar Ahmad'),
    ('founder_1_role',    'Founder & Owner'),
    ('founder_1_bio',     'Established Faxman Travels in 2004 with a vision to provide reliable, trustworthy, and customer-friendly travel services. For over 22 years, his leadership has been the cornerstone of trust and excellence.'),
    ('founder_1_image_url', ''),
    ('founder_2_name',    'Mr. Shaikh Zaid Ahmad'),
    ('founder_2_role',    'Co-Owner & Managing Director'),
    ('founder_2_bio',     'Joined the business in 2018, actively carrying forward the vision and legacy built by his father. Continuously working towards expanding and modernising the business under his father''s valuable guidance.'),
    ('founder_2_image_url', '')
ON CONFLICT (key) DO NOTHING;

INSERT INTO al_haram.bank_details (id, account_name, bank_name, account_no, ifsc, branch, note) VALUES (
    1,
    'Shaikh Zaid Ahmad',
    'HDFC Bank',
    '50200032599285',
    'HDFC0001895',
    'Fatehpur Branch, U.P.',
    'After payment, share the transaction screenshot on WhatsApp for confirmation.'
) ON CONFLICT (id) DO UPDATE SET
    account_name = EXCLUDED.account_name,
    bank_name    = EXCLUDED.bank_name,
    account_no   = EXCLUDED.account_no,
    ifsc         = EXCLUDED.ifsc,
    branch       = EXCLUDED.branch,
    note         = EXCLUDED.note;

INSERT INTO al_haram.policies (id, title, content) VALUES
    ('terms',   'Terms & Conditions',
     'All bookings are subject to availability and confirmation. The customer is responsible for providing correct documents and valid information at the time of booking. The company is not responsible for any delays or cancellations caused by airlines, embassy decisions, immigration authorities, natural disasters, or any other circumstances beyond our control. Prices are based on current airline fares and hotel rates; any future increase in fares or government regulations will be borne by the pilgrim. Customers must verify all travel documents before departure. Passport validity of minimum 6 months is mandatory at the time of travel.'),
    ('refund',  'Refund Policy',
     'Refunds for air tickets, visas, hotels, and packages are subject to individual airline rules, embassy policies, hotel cancellation terms, and applicable service charges. Some services may be non-refundable once confirmed. Visa fees are non-refundable once submitted to the embassy. Refund processing times vary depending on the provider. Customers are strongly advised to read all cancellation and refund terms carefully before making any payment. Faxman Travels will assist in processing refund requests but cannot guarantee timelines set by third-party providers.'),
    ('privacy', 'Privacy Policy',
     'Customer data including passport details, phone numbers, email addresses, travel documents, and payment information is kept strictly confidential. We do not share personal information with third parties without explicit permission, except where required for visa processing, airline booking, or legal compliance. Basic site usage data may be collected to improve service quality. You may request deletion of your personal data at any time by contacting us directly.')
ON CONFLICT (id) DO UPDATE SET
    title   = EXCLUDED.title,
    content = EXCLUDED.content;

-- ── Seed packages ─────────────────────────────────────────────────────────────
INSERT INTO al_haram.packages
    (id, title, tier, nights, price, image, flight_type,
     makkah_hotel, makkah_dist, madinah_hotel, madinah_dist,
     shuttle_makkah, shuttle_madinah, occupancy,
     features, documents, rate_note, sort_order)
VALUES
('pkg_economy', 'Economy Umrah Package 2026', 'Economy', 16, 'From ₹92,000',
 'https://i.pinimg.com/1200x/24/7b/36/247b36c6b1475fe7405647c54c4a616d.jpg',
 'Oman Air (LKO-MCT-JED via Muscat)',
 'Nokhba Al-Khair or Similar', '1500–2000 m from Haram',
 'Markaziya or Similar', 'Within 1 km from Masjid an-Nabawi',
 'Free 24×7 Haram shuttle (150 m walk from drop point)', 'Shuttle for Namaz timings only', '5–6 persons per room',
 '["Round-trip air ticket (Oman Air)","Umrah visa","Hotel on sharing basis (5–6 persons per room)","Breakfast, lunch & dinner","Local transfer by AC bus","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates are based on current airline fares and hotel rates. Any future increase in airline or hotel rates will be borne by the pilgrim.', 1),

('pkg_semi_deluxe', 'Semi Deluxe Umrah Package 2026', 'Semi Deluxe', 15, 'From ₹1,07,000',
 'https://i.pinimg.com/736x/92/1f/f9/921ff99d18598b360b2a5705d1f7b653.jpg',
 'Oman Air / Air India (via Lucknow) — Direct or Via Flight',
 'Manar Al Azhar or Similar', '500–700 m from Haram',
 'Markaziya or Similar', '200–250 m from Masjid an-Nabawi',
 'Free 24×7 Haram shuttle', 'Walking distance / shuttle', '4–5 persons per room',
 '["Round-trip air ticket (Oman Air / Air India)","Umrah visa","Deluxe hotel on sharing basis (4–5 persons per room)","Breakfast, lunch & dinner","Local transfer by AC bus","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water","Type 1: Direct/deluxe flight with economy hotels","Type 2: Economy/via flight with deluxe hotels (500–700 m Makkah, 200–250 m Madinah)"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates are based on current airline fares and hotel rates. Any future increase in airline or hotel rates will be borne by the pilgrim.', 2),

('pkg_deluxe', 'Deluxe Umrah Package', 'Deluxe', 15, 'On Request',
 'https://i.pinimg.com/736x/1e/f9/ed/1ef9ed0de60fb7ae04073b10f0b05951.jpg',
 'Direct / Deluxe Flight',
 'Deluxe Hotel', '500–700 m from Haram',
 'Deluxe Hotel', '200–250 m from Masjid an-Nabawi',
 'Free 24×7 Haram shuttle', 'Shuttle provided', '3–4 persons per room',
 '["Round-trip air ticket (direct/deluxe flight)","Umrah visa","Deluxe hotel on sharing basis (3–4 persons per room)","Breakfast, lunch & dinner","Local transfer by AC bus","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates are based on current airline fares and hotel rates. Any future increase in airline or hotel rates will be borne by the pilgrim.', 3),

('pkg_super_deluxe', 'Super Deluxe Umrah Package', 'Super Deluxe', 15, 'On Request',
 'https://i.pinimg.com/736x/28/43/23/28432334b5bf415ebd5e64aebc892a41.jpg',
 'Premium Direct Flight',
 'Le Méridien Towers Makkah (5-Star)', '100–150 m from Haram',
 'Super Deluxe Hotel', '100–150 m from Masjid an-Nabawi',
 'Private AC transport to Haram', 'Private AC transport to Masjid an-Nabawi', '2–3 persons per room',
 '["Round-trip premium direct flight","Umrah visa","5-Star hotel on sharing basis (2–3 persons per room)","Breakfast, lunch & dinner","Private AC transport","Complete Ziyarat in Makkah & Madinah","Laundry service","Travel insurance","5 Ltr Zam Zam water","Le Méridien Towers Makkah or equivalent 5-Star"]',
 '["Valid passport (min. 6 months validity)","2 passport-size photos (white background)","PAN card"]',
 'Rates are based on current airline fares and hotel rates. Any future increase in airline or hotel rates will be borne by the pilgrim.', 4)
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

-- ── Seed schedule (Umrah 2026 / 1447-1448 Hijri) ─────────────────────────────
INSERT INTO al_haram.umrah_schedule
    (id, season, departure_date, return_date, duration_nights, airline, route, package_tier, rate, seats_left, status, notes)
VALUES
    ('sch_2026_01', '2026 / 1447', '2026-07-07', '2026-07-22', 16, 'Oman Air', 'LKO-MCT-JED / JED-MCT-LKO', 'Economy', '₹92,000', 0, 'Available', ''),
    ('sch_2026_02', '2026 / 1447', '2026-07-23', '2026-08-07', 16, 'Oman Air', 'LKO-MCT-JED / MED-MCT-LKO', 'Economy', '₹92,000', 0, 'Available', ''),
    ('sch_2026_03', '2026 / 1447', '2026-08-07', '2026-08-22', 16, 'Oman Air', 'LKO-MCT-JED / MED-MCT-LKO', 'Economy', '₹92,000', 0, 'Available', ''),
    ('sch_2026_04', '2026 / 1447', '2026-09-09', '2026-09-23', 15, 'Oman Air', 'LKO-MCT-JED / JED-MCT-LKO', 'Economy', '₹90,000', 0, 'Available', ''),
    ('sch_2026_05', '2026 / 1447', '2026-09-14', '2026-09-30', 17, 'Oman Air', 'LKO-MCT-JED / MED-MCT-LKO', 'Economy', '₹93,000', 0, 'Available', ''),
    ('sch_2026_06', '2026 / 1447', '2026-10-06', '2026-10-21', 16, 'Oman Air', 'LKO-MCT-JED / MED-MCT-LKO', 'Economy', '₹92,000', 0, 'Available', ''),
    ('sch_2026_07', '2026 / 1447', '2026-10-15', '2026-10-30', 16, 'Oman Air', 'LKO-MCT-JED / MED-MCT-LKO', 'Economy', '₹92,000', 0, 'Available', ''),
    ('sch_2026_s01', '2026 / 1448', '2026-06-27', '2026-07-11', 15, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 1 — 0845-1550 / 1800-0805'),
    ('sch_2026_s02', '2026 / 1448', '2026-06-30', '2026-07-15', 16, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Sold',      'Group 2 — Sold Out'),
    ('sch_2026_s03', '2026 / 1448', '2026-07-08', '2026-07-23', 16, 'Air India', 'LKO-DEL-JED', 'Semi Deluxe', '₹1,03,000', 0, 'Available', 'Group 3 — 1410-2140 / 2240-1335'),
    ('sch_2026_s04', '2026 / 1448', '2026-07-13', '2026-07-29', 17, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 4 — 0845-1710 / 1800-0805'),
    ('sch_2026_s05', '2026 / 1448', '2026-07-28', '2026-08-12', 16, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 5 — 0845-1710 / 1800-0805'),
    ('sch_2026_s06', '2026 / 1448', '2026-08-04', '2026-08-21', 18, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 6 — 0845-1710 / 1800-0805'),
    ('sch_2026_s07', '2026 / 1448', '2026-08-12', '2026-08-27', 16, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,10,000', 0, 'Available', 'Group 7 — 0845-1710 / 1800-0805'),
    ('sch_2026_s08', '2026 / 1448', '2026-08-24', '2026-09-09', 17, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 8 — 0845-1710 / 1800-0805'),
    ('sch_2026_s09', '2026 / 1448', '2026-09-01', '2026-09-16', 16, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 9 — 0845-1710 / 1800-0805'),
    ('sch_2026_s10', '2026 / 1448', '2026-09-09', '2026-09-26', 18, 'Oman Air', 'LKO-MCT-JED', 'Semi Deluxe', '₹1,07,000', 0, 'Available', 'Group 10 — 0845-1710 / 1800-0805')
ON CONFLICT (id) DO UPDATE SET
    season          = EXCLUDED.season,
    departure_date  = EXCLUDED.departure_date,
    return_date     = EXCLUDED.return_date,
    duration_nights = EXCLUDED.duration_nights,
    airline         = EXCLUDED.airline,
    route           = EXCLUDED.route,
    package_tier    = EXCLUDED.package_tier,
    rate            = EXCLUDED.rate,
    status          = EXCLUDED.status,
    notes           = EXCLUDED.notes;
`
