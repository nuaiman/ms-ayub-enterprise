-- =====================================================
-- SEED DATA - ALL TABLES EXCEPT USERS
-- =====================================================

-- =====================================================
-- 1. BROKERS
-- =====================================================
INSERT OR IGNORE INTO brokers (id, name, phone, notes) VALUES
(1, 'Broker Rahman', '01720000001', 'Experienced broker in Dhaka area'),
(2, 'Broker Khan', '01720000002', 'Specializes in Chittagong routes'),
(3, 'Broker Ahmed', '01720000003', 'Handles all major godowns'),
(4, 'Broker Hossain', '01720000004', 'New broker - verified'),
(5, 'Broker Karim', '01720000005', '10 years experience in transport'),
(6, 'Broker Rahim', '01720000006', 'Specializes in local deliveries');

-- =====================================================
-- 2. MAJHIS
-- =====================================================
INSERT OR IGNORE INTO majhis (id, name, phone, notes) VALUES
(1, 'Majhi Karim', '01730000001', 'Senior majhi - 15 years experience'),
(2, 'Majhi Rahim', '01730000002', 'Handles loading/unloading'),
(3, 'Majhi Jabbar', '01730000003', 'Specializes in heavy cargo'),
(4, 'Majhi Shamsu', '01730000004', 'Works in Chittagong region'),
(5, 'Majhi Kamal', '01730000005', 'Expert in rice and grains'),
(6, 'Majhi Selim', '01730000006', 'Handles construction materials');

-- =====================================================
-- 3. CUSTOMERS (with address)
-- =====================================================
INSERT OR IGNORE INTO customers (id, company_name, contact_person, phone, email, address, notes) VALUES
(1, 'ABC Trading Ltd', 'Mr. Hasan', '01750000001', 'hasan@abctrading.com', 'House 12, Road 5, Dhanmondi, Dhaka', 'Regular customer - 5 years'),
(2, 'XYZ Industries', 'Ms. Fatima', '01750000002', 'fatima@xyzind.com', 'Plot 7, Sector 3, Uttara, Dhaka', 'Premium customer - bulk orders'),
(3, 'RFL Group', 'Mr. Kamal', '01750000003', 'kamal@rflgroup.com', 'RFL Plaza, 8th Floor, Motijheel, Dhaka', 'Large corporate client'),
(4, 'Pran Foods', 'Mr. Rana', '01750000004', 'rana@pranfoods.com', 'Pran Tower, 25/1, Gulshan Avenue, Dhaka', 'Food industry client'),
(5, 'Beximco Ltd', 'Ms. Sonia', '01750000005', 'sonia@beximco.com', 'Beximco Industrial Park, Gazipur', 'Textile sector client'),
(6, 'Square Group', 'Mr. Alam', '01750000006', 'alam@squaregroup.com', 'Square Centre, 48, Mohakhali, Dhaka', 'Pharmaceutical client'),
(7, 'City Mart', 'Mr. Sajib', '01750000007', 'sajib@citymart.com', 'House 15, Road 12, Banani, Dhaka', 'Retail chain'),
(8, 'Green Agro', 'Mrs. Lima', '01750000008', 'lima@greenagro.com', 'Village: Polashbari, Upazila: Sreepur, Gazipur', 'Agricultural products');

-- =====================================================
-- 4. GODOWNS
-- =====================================================
INSERT OR IGNORE INTO godowns (id, name, phone, notes, is_active, monthly_rent) VALUES
(1, 'Godown A - Dhaka', '01740000001', 'Main godown in Dhaka - 5000 sq ft', 1, 15000),
(2, 'Godown B - Chittagong', '01740000002', 'Godown in Chittagong port area', 1, 12000),
(3, 'Godown C - Sylhet', '01740000003', 'Godown in Sylhet - 3000 sq ft', 1, 8000),
(4, 'Godown D - Khulna', '01740000004', 'Godown in Khulna - 2000 sq ft', 1, 6000),
(5, 'Godown E - Rajshahi', '01740000005', 'Godown in Rajshahi - 2500 sq ft', 1, 7000),
(6, 'Godown F - Barishal', '01740000006', 'Godown in Barishal - 1500 sq ft', 1, 5000);

-- =====================================================
-- 5. ITEMS (user_id = 1, bootstrap admin)
-- =====================================================
INSERT OR IGNORE INTO items (id, user_id, customer_id, product_name, category, is_active, notes) VALUES
(1, 1, 1, 'Steel Rod 12mm', 'Construction', 1, 'High-tensile steel for RCC work'),
(2, 1, 4, 'Portland Cement 50kg', 'Building Materials', 1, 'Grade 43, standard quality');

-- =====================================================
-- 6. LOTS (UPDATED with customer_paid_unload_amount)
-- =====================================================
INSERT OR IGNORE INTO lots (
    id, item_id, lot_number, customer_charge_type, majhi_bill_type, 
    customer_storage_rate, unload_rate, majhi_id, majhi_cut, 
    is_active, notes,
    customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid
) VALUES
(1, 1, 1, 'quantity', 'quantity', 5.00, 2.50, 1, 1.00, 1, 'First batch of 12mm rods', NULL, 0, 0, 0),
(2, 2, 1, 'weight', 'weight', 10.00, 3.00, 2, 0.50, 1, 'First batch of cement', NULL, 0, 0, 0);

-- =====================================================
-- 7. STORES
-- =====================================================
INSERT OR IGNORE INTO stores (id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit, weight, weight_unit, is_active, billing_start, billing_end, last_paid_through, last_paid_amount, notes) VALUES
(1, 1, 1, 'quantity', 2.00, 100, 'rods', 500, 'kg', 1, '2026-09-01 00:00:00', NULL, NULL, 0, 'Initial stock of steel rods'),
(2, 1, 2, 'quantity', 1.50, 50, 'rods', 250, 'kg', 1, '2026-09-01 00:00:00', NULL, NULL, 0, 'Steel rods in Chittagong godown'),
(3, 2, 1, 'weight', 3.00, 0, 'bags', 2000, 'kg', 1, '2026-09-01 00:00:00', NULL, NULL, 0, 'Cement stock in Dhaka');

-- =====================================================
-- 8. DAMAGES
-- =====================================================
INSERT OR IGNORE INTO damages (id, store_id, user_id, quantity, quantity_unit, weight, weight_unit, damage_date, reason, amount, notes) VALUES
(1, 1, 1, 5, 'rods', 25, 'kg', '2026-09-05 00:00:00', 'Cracked during handling', 150.00, 'Minor damage, scrapped');

-- =====================================================
-- 9. DELIVERIES
-- =====================================================
INSERT OR IGNORE INTO deliveries (id, user_id, customer_id, delivery_date, receiver_name, receiver_phone, from_location, to_location, notes) VALUES
(1, 1, 1, '2026-09-06 10:00:00', 'Mr. Hasan', '01750000001', 'Godown A - Dhaka', 'Construction Site – Mirpur', 'First delivery of steel rods');

INSERT OR IGNORE INTO delivery_items (
    id, delivery_id, store_id, majhi_id, item_id, lot_id, 
    vehicle_number, driver_number, quantity, quantity_unit, 
    weight, weight_unit, loading_rate, majhi_cut, notes,
    customer_charge_type, customer_paid_unload_amount,
    majhi_bill_type, majhi_total_paid
) VALUES (
    1, 1, 1, 1, 1, 1, 
    'DHK-1234', '01730000001', 20, 'rods', 
    100, 'kg', 2.00, 1.00, 'First load',
    'quantity', 0,
    'quantity', 0
);

-- =====================================================
-- 11. TRANSPORTS
-- =====================================================
INSERT OR IGNORE INTO transports (id, user_id, customer_id, from_location, to_location, vehicle_quantity, delivery_type, notes, transport_date, office_commission_amount, customer_total_paid) VALUES
(1, 1, 1, 'Godown A - Dhaka', 'Chittagong Port', 1, 'local', 'Steel rods transport to Chittagong', '2026-09-07 08:00:00', 500.00, 0);

-- =====================================================
-- 12. VEHICLES
-- =====================================================
INSERT OR IGNORE INTO vehicles (id, user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone, joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost, demarage_amount, demarage_reason, broker_total_paid) VALUES
(1, 1, 1, 'DHAKA-1234', 1, 'Mr. Karim', '01720000001', 2000.00, 1500.00, 4000.00, 300.00, 500.00, 0.00, NULL, 0);

-- =====================================================
-- 13. EXPENSES
-- =====================================================
INSERT OR IGNORE INTO expenses (id, user_id, title, amount, expense_date, notes) VALUES
(1, 1, 'Electricity Bill - September', 2500.00, '2026-09-01 00:00:00', 'Monthly electricity bill'),
(2, 1, 'Office Supplies', 1200.00, '2026-09-05 00:00:00', 'Stationery and printer ink');

-- =====================================================
-- 14. RENTS
-- =====================================================
INSERT OR IGNORE INTO rents (id, user_id, godown_id, month_year, amount, status, payment_date, payment_method, reference_number, notes) VALUES
(1, 1, 1, '2026-09', 15000, 'paid', '2026-09-01 00:00:00', 'cash', 'RENT-001', 'September rent paid'),
(2, 1, 2, '2026-09', 12000, 'draft', NULL, NULL, NULL, 'Pending rent for Chittagong'),
(3, 1, 1, '2026-08', 15000, 'paid', '2026-08-01 00:00:00', 'bank_transfer', 'BT-2026-08', 'August rent paid via bank'),
(4, 1, 3, '2026-09', 8000, 'paid', '2026-09-02 00:00:00', 'mobile_banking', 'MB-123', 'Sylhet rent paid');

-- =====================================================
-- 19. LOGS
-- =====================================================
INSERT OR IGNORE INTO logs (id, user_id, action, description, entity_type, entity_id, ip_address, user_agent, created_at) VALUES
(1, 1, 'seed', 'Initial seed data load', 'system', 0, '127.0.0.1', 'seed-script', CURRENT_TIMESTAMP);

-- =====================================================
-- 20. INVENTORY AUDIT
-- =====================================================
-- No inventory audit entries inserted.