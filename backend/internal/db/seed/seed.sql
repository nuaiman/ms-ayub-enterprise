-- =====================================================
-- SEED DATA
-- =====================================================
-- Idempotent: uses fixed IDs + INSERT OR IGNORE.
-- Safe to run on every startup.
-- =====================================================

PRAGMA foreign_keys = ON;

-- =====================================================
-- 1. CUSTOMER
-- =====================================================
INSERT OR IGNORE INTO customers (
    id, company_name, contact_person, phone, email, address, notes
) VALUES (
    1,
    'Acme Traders',
    'Rahim Uddin',
    '01700000001',
    'acme@example.com',
    '12 Mirpur Road, Dhaka',
    'Seed customer'
);

-- =====================================================
-- 2. MAJHI
-- =====================================================
INSERT OR IGNORE INTO majhis (
    id, name, phone, notes
) VALUES (
    1,
    'Karim Majhi',
    '01700000002',
    'Seed majhi'
);

-- =====================================================
-- 3. BROKER
-- =====================================================
INSERT OR IGNORE INTO brokers (
    id, name, phone, notes
) VALUES (
    1,
    'Hossain Broker',
    '01700000003',
    'Seed broker'
);

-- =====================================================
-- 4. GODOWN
-- =====================================================
INSERT OR IGNORE INTO godowns (
    id, name, phone, notes
) VALUES (
    1,
    'Main Godown',
    '01700000004',
    'Seed godown'
);