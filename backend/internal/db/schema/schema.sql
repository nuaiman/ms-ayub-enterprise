-- =====================================================
-- PRAGMAS
-- =====================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA busy_timeout = 5000;

-- =====================================================
-- 1. USERS 
-- =====================================================
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    username TEXT UNIQUE,
    email TEXT UNIQUE,
    phone TEXT,
    address TEXT,
    id_type TEXT CHECK (
        id_type IN ('nid', 'passport', 'driving_license', 'birth_certificate', 'trade_license', 'other')
    ),
    id_number TEXT,
    image_url TEXT,
    password TEXT,
    role TEXT NOT NULL DEFAULT 'staff' CHECK (
        role IN ('admin', 'manager', 'accounts', 'staff')
    ),
    is_active INTEGER NOT NULL DEFAULT 1,
    monthly_salary REAL DEFAULT 0,
    refresh_token TEXT,
    refresh_token_expiry DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- =====================================================
-- 2. SALARIES
-- =====================================================
CREATE TABLE IF NOT EXISTS salaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    employee_id INTEGER NOT NULL,
    month_year TEXT NOT NULL,
    bonus REAL DEFAULT 0 CHECK (bonus >= 0),
    deductions REAL DEFAULT 0 CHECK (deductions >= 0),
    status TEXT NOT NULL DEFAULT 'draft' CHECK (
        status IN ('draft', 'paid', 'cancelled')
    ),
    payment_date DATETIME,
    payment_method TEXT CHECK (
        payment_method IN ('cash', 'bank_transfer', 'check', 'mobile_banking')
    ),
    reference_number TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CHECK (
        month_year GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]' AND
        substr(month_year, 6, 2) BETWEEN '01' AND '12'
    ),
    
    UNIQUE(employee_id, month_year),
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (employee_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_salaries_user_id ON salaries(user_id);
CREATE INDEX IF NOT EXISTS idx_salaries_employee_id ON salaries(employee_id);
CREATE INDEX IF NOT EXISTS idx_salaries_month_year ON salaries(month_year);
CREATE INDEX IF NOT EXISTS idx_salaries_employee_month ON salaries(employee_id, month_year);
CREATE INDEX IF NOT EXISTS idx_salaries_status ON salaries(status);

-- =====================================================
-- 3. BROKERS
-- =====================================================
CREATE TABLE IF NOT EXISTS brokers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    phone TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (name != '')
);

CREATE INDEX IF NOT EXISTS idx_brokers_phone ON brokers(phone);
CREATE INDEX IF NOT EXISTS idx_brokers_name ON brokers(name);

-- =====================================================
-- 4. MAJHIS
-- =====================================================
CREATE TABLE IF NOT EXISTS majhis (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    phone TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (name != '')
);

CREATE INDEX IF NOT EXISTS idx_majhis_phone ON majhis(phone);
CREATE INDEX IF NOT EXISTS idx_majhis_name ON majhis(name);

-- =====================================================
-- 5. GODOWNS 
-- =====================================================
CREATE TABLE IF NOT EXISTS godowns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    phone TEXT,
    notes TEXT,
    is_active INTEGER NOT NULL DEFAULT 1,
    monthly_rent REAL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (name != '')
);

CREATE INDEX IF NOT EXISTS idx_godowns_name ON godowns(name);
CREATE INDEX IF NOT EXISTS idx_godowns_phone ON godowns(phone);
CREATE INDEX IF NOT EXISTS idx_godowns_is_active ON godowns(is_active);

-- =====================================================
-- 6. GODOWN RENTS
-- =====================================================
CREATE TABLE IF NOT EXISTS rents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    godown_id INTEGER NOT NULL,
    month_year TEXT NOT NULL,
    amount REAL NOT NULL DEFAULT 0 CHECK (amount >= 0),
    status TEXT NOT NULL DEFAULT 'draft' CHECK (
        status IN ('draft', 'paid', 'cancelled')
    ),
    payment_date DATETIME,
    payment_method TEXT CHECK (
        payment_method IN ('cash', 'bank_transfer', 'check', 'mobile_banking')
    ),
    reference_number TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CHECK (
        month_year GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]' AND
        substr(month_year, 6, 2) BETWEEN '01' AND '12'
    ),
    
    UNIQUE(godown_id, month_year),
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (godown_id) REFERENCES godowns(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_rents_user_id ON rents(user_id);
CREATE INDEX IF NOT EXISTS idx_rents_godown_id ON rents(godown_id);
CREATE INDEX IF NOT EXISTS idx_rents_month_year ON rents(month_year);
CREATE INDEX IF NOT EXISTS idx_rents_godown_month ON rents(godown_id, month_year);
CREATE INDEX IF NOT EXISTS idx_rents_status ON rents(status);

-- =====================================================
-- 7. CUSTOMERS
-- =====================================================
CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_name TEXT,
    contact_person TEXT,
    phone TEXT UNIQUE,
    email TEXT,
    address TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (company_name IS NOT NULL OR contact_person IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_company ON customers(company_name);
CREATE INDEX IF NOT EXISTS idx_customers_contact_person ON customers(contact_person);

-- =====================================================
-- 8. ITEMS
-- =====================================================
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    customer_id INTEGER,
    product_name TEXT,
    category TEXT,
    is_active INTEGER NOT NULL DEFAULT 1,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_items_user_id ON items(user_id);
CREATE INDEX IF NOT EXISTS idx_items_customer_id ON items(customer_id);
CREATE INDEX IF NOT EXISTS idx_items_product_name ON items(product_name);
CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);
CREATE INDEX IF NOT EXISTS idx_items_is_active ON items(is_active);

-- =====================================================
-- 9. LOTS
-- =====================================================
CREATE TABLE IF NOT EXISTS lots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id INTEGER NOT NULL,
    lot_number INTEGER NOT NULL,
    customer_charge_type TEXT NOT NULL DEFAULT 'quantity' CHECK (customer_charge_type IN ('weight', 'quantity')),
    majhi_bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (majhi_bill_type IN ('weight', 'quantity', 'job')),
    customer_storage_rate REAL NOT NULL DEFAULT 0 CHECK (customer_storage_rate >= 0),
    unload_rate REAL NOT NULL DEFAULT 0 CHECK (unload_rate >= 0),
    majhi_id INTEGER,
    majhi_cut REAL NOT NULL DEFAULT 0 CHECK (majhi_cut >= 0),
    is_active INTEGER NOT NULL DEFAULT 1,
    notes TEXT,
    image_url TEXT,
    customer_last_paid_through DATETIME,
    customer_last_paid_amount REAL NOT NULL DEFAULT 0 CHECK (customer_last_paid_amount >= 0),
    customer_paid_unload_amount REAL NOT NULL DEFAULT 0 CHECK (customer_paid_unload_amount >= 0),
    majhi_total_paid REAL NOT NULL DEFAULT 0 CHECK (majhi_total_paid >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    FOREIGN KEY (majhi_id) REFERENCES majhis(id) ON DELETE SET NULL,
    UNIQUE(item_id, lot_number)
);

CREATE INDEX IF NOT EXISTS idx_lots_item_id ON lots(item_id);
CREATE INDEX IF NOT EXISTS idx_lots_lot_number ON lots(item_id, lot_number);
CREATE INDEX IF NOT EXISTS idx_lots_majhi_id ON lots(majhi_id);
CREATE INDEX IF NOT EXISTS idx_lots_is_active ON lots(is_active);

-- =====================================================
-- 10. STORES
-- =====================================================
CREATE TABLE IF NOT EXISTS stores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lot_id INTEGER NOT NULL,
    godown_id INTEGER NOT NULL,
    store_bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (store_bill_type IN ('weight', 'quantity')),
    godown_cut REAL NOT NULL DEFAULT 0 CHECK (godown_cut >= 0),
    quantity REAL NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_unit TEXT NOT NULL DEFAULT 'units',
    weight REAL NOT NULL DEFAULT 0 CHECK (weight >= 0),
    weight_unit TEXT NOT NULL DEFAULT 'kg',
    is_active INTEGER NOT NULL DEFAULT 1,
    billing_start DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    billing_end DATETIME,
    last_paid_through DATETIME,
    last_paid_amount REAL NOT NULL DEFAULT 0 CHECK (last_paid_amount >= 0),
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE CASCADE,
    FOREIGN KEY (godown_id) REFERENCES godowns(id) ON DELETE RESTRICT,
    UNIQUE(lot_id, godown_id),
    CHECK (billing_end IS NULL OR billing_start <= billing_end)
);

CREATE INDEX IF NOT EXISTS idx_stores_lot_id ON stores(lot_id);
CREATE INDEX IF NOT EXISTS idx_stores_godown_id ON stores(godown_id);
CREATE INDEX IF NOT EXISTS idx_stores_lot_godown ON stores(lot_id, godown_id);
CREATE INDEX IF NOT EXISTS idx_stores_billing_start ON stores(billing_start);
CREATE INDEX IF NOT EXISTS idx_stores_billing_end ON stores(billing_end);
CREATE INDEX IF NOT EXISTS idx_stores_is_active ON stores(is_active);

-- =====================================================
-- 11. DAMAGES
-- =====================================================
CREATE TABLE IF NOT EXISTS damages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    store_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    quantity REAL NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_unit TEXT NOT NULL DEFAULT 'units',
    weight REAL NOT NULL DEFAULT 0 CHECK (weight >= 0),
    weight_unit TEXT NOT NULL DEFAULT 'kg',
    damage_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason TEXT NOT NULL,
    amount REAL DEFAULT 0 CHECK (amount >= 0),
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CHECK (quantity > 0 OR weight > 0)
);

CREATE INDEX IF NOT EXISTS idx_damages_store_id ON damages(store_id);
CREATE INDEX IF NOT EXISTS idx_damages_user_id ON damages(user_id);
CREATE INDEX IF NOT EXISTS idx_damages_damage_date ON damages(damage_date);
CREATE INDEX IF NOT EXISTS idx_damages_created_at ON damages(created_at);

-- =====================================================
-- 12. DELIVERIES (Master)
-- =====================================================
CREATE TABLE IF NOT EXISTS deliveries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    customer_id INTEGER,
    delivery_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    receiver_name TEXT,
    receiver_phone TEXT,
    from_location TEXT,
    to_location TEXT,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_deliveries_user_id ON deliveries(user_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_customer_id ON deliveries(customer_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_date ON deliveries(delivery_date);
CREATE INDEX IF NOT EXISTS idx_deliveries_receiver_name ON deliveries(receiver_name);
CREATE INDEX IF NOT EXISTS idx_deliveries_receiver_phone ON deliveries(receiver_phone);

-- =====================================================
-- 13. DELIVERY ITEMS (Detail)
-- =====================================================
CREATE TABLE IF NOT EXISTS delivery_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    delivery_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    majhi_id INTEGER,
    item_id INTEGER NOT NULL,
    lot_id INTEGER NOT NULL,
    vehicle_number TEXT,
    driver_number TEXT,
    quantity REAL NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_unit TEXT NOT NULL DEFAULT 'units',
    weight REAL NOT NULL DEFAULT 0 CHECK (weight >= 0),
    weight_unit TEXT NOT NULL DEFAULT 'kg',
    loading_rate REAL NOT NULL DEFAULT 0 CHECK (loading_rate >= 0),
    majhi_cut REAL NOT NULL DEFAULT 0 CHECK (majhi_cut >= 0),
    notes TEXT,
    customer_charge_type TEXT NOT NULL DEFAULT 'quantity' CHECK (customer_charge_type IN ('weight', 'quantity')),
    customer_paid_unload_amount REAL NOT NULL DEFAULT 0 CHECK (customer_paid_unload_amount >= 0),
    majhi_bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (majhi_bill_type IN ('weight', 'quantity', 'job')),
    majhi_total_paid REAL NOT NULL DEFAULT 0 CHECK (majhi_total_paid >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE RESTRICT,
    FOREIGN KEY (majhi_id) REFERENCES majhis(id) ON DELETE SET NULL,
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT,
    FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_delivery_items_delivery_id ON delivery_items(delivery_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_store_id ON delivery_items(store_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_majhi_id ON delivery_items(majhi_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_item_id ON delivery_items(item_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_lot_id ON delivery_items(lot_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_vehicle_number ON delivery_items(vehicle_number);
CREATE INDEX IF NOT EXISTS idx_delivery_items_driver_number ON delivery_items(driver_number);

-- =====================================================
-- 14. TRANSPORTS
-- =====================================================
CREATE TABLE IF NOT EXISTS transports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    customer_id INTEGER,
    from_location TEXT NOT NULL,
    to_location TEXT,
    vehicle_quantity REAL NOT NULL DEFAULT 0,
    delivery_type TEXT CHECK (
        delivery_type IN ('local', 'district')
    ),
    notes TEXT,
    transport_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    office_commission_amount REAL DEFAULT 0,
    image_url TEXT,
    customer_total_paid REAL NOT NULL DEFAULT 0 CHECK (customer_total_paid >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (vehicle_quantity >= 0),
    CHECK (from_location != ''),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_transports_user_id ON transports(user_id);
CREATE INDEX IF NOT EXISTS idx_transports_customer_id ON transports(customer_id);
CREATE INDEX IF NOT EXISTS idx_transports_from_location ON transports(from_location);
CREATE INDEX IF NOT EXISTS idx_transports_to_location ON transports(to_location);
CREATE INDEX IF NOT EXISTS idx_transports_transport_date ON transports(transport_date);
CREATE INDEX IF NOT EXISTS idx_transports_delivery_type ON transports(delivery_type);
CREATE INDEX IF NOT EXISTS idx_transports_vehicle_quantity ON transports(vehicle_quantity);

-- =====================================================
-- 15. VEHICLES
-- =====================================================
CREATE TABLE IF NOT EXISTS vehicles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    transport_id INTEGER NOT NULL,
    vehicle_number TEXT NOT NULL,
    broker_id INTEGER,
    driver_name TEXT,
    driver_phone TEXT,
    joma_cost REAL DEFAULT 0,
    vehicle_cost REAL DEFAULT 0,
    customer_charge REAL DEFAULT 0,
    other_cost REAL DEFAULT 0,
    labour_cost REAL DEFAULT 0,
    demarage_amount REAL DEFAULT 0,
    demarage_reason TEXT,
    broker_total_paid REAL NOT NULL DEFAULT 0 CHECK (broker_total_paid >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (vehicle_number != ''),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (transport_id) REFERENCES transports(id) ON DELETE CASCADE,
    FOREIGN KEY (broker_id) REFERENCES brokers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_vehicles_user_id ON vehicles(user_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_transport_id ON vehicles(transport_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_broker_id ON vehicles(broker_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_vehicle_number ON vehicles(vehicle_number);
CREATE INDEX IF NOT EXISTS idx_vehicles_created_at ON vehicles(created_at);
CREATE INDEX IF NOT EXISTS idx_vehicles_driver_name ON vehicles(driver_name);
CREATE INDEX IF NOT EXISTS idx_vehicles_driver_phone ON vehicles(driver_phone);

-- =====================================================
-- 16. EXPENSES
-- =====================================================
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    amount REAL NOT NULL DEFAULT 0 CHECK (amount >= 0),
    expense_date DATETIME NOT NULL,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);

-- =====================================================
-- 17. LOGS
-- =====================================================
CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    action TEXT NOT NULL,
    description TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER NOT NULL,
    old_data TEXT,
    new_data TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_entity_type ON logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_logs_entity_id ON logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_logs_action ON logs(action);
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);
CREATE INDEX IF NOT EXISTS idx_logs_entity_created ON logs(entity_type, entity_id, created_at DESC);

-- =====================================================
-- INVENTORY MANAGEMENT TRIGGERS
-- =====================================================

-- TRIGGER: Check stock BEFORE delivery insert
CREATE TRIGGER IF NOT EXISTS check_stock_before_delivery_insert
BEFORE INSERT ON delivery_items
BEGIN
    SELECT CASE 
        WHEN (
            SELECT quantity FROM stores WHERE id = NEW.store_id
        ) < NEW.quantity THEN 
            RAISE(ABORT, 'Insufficient quantity in store')
        WHEN (
            SELECT weight FROM stores WHERE id = NEW.store_id
        ) < NEW.weight THEN 
            RAISE(ABORT, 'Insufficient weight in store')
    END;
END;

-- TRIGGER: Check stock BEFORE delivery update
CREATE TRIGGER IF NOT EXISTS check_stock_before_delivery_update
BEFORE UPDATE ON delivery_items
BEGIN
    SELECT CASE 
        WHEN (
            (SELECT quantity FROM stores WHERE id = NEW.store_id) + OLD.quantity
        ) < NEW.quantity THEN 
            RAISE(ABORT, 'Insufficient quantity in store for update')
        WHEN (
            (SELECT weight FROM stores WHERE id = NEW.store_id) + OLD.weight
        ) < NEW.weight THEN 
            RAISE(ABORT, 'Insufficient weight in store for update')
    END;
END;

-- TRIGGER: Update store on delivery insert
CREATE TRIGGER IF NOT EXISTS update_store_on_delivery_insert
AFTER INSERT ON delivery_items
BEGIN
    UPDATE stores 
    SET 
        quantity = quantity - NEW.quantity,
        weight = weight - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- TRIGGER: Update store on delivery update
CREATE TRIGGER IF NOT EXISTS update_store_on_delivery_update
AFTER UPDATE ON delivery_items
BEGIN
    UPDATE stores 
    SET 
        quantity = quantity + OLD.quantity - NEW.quantity,
        weight = weight + OLD.weight - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- TRIGGER: Restore store on delivery delete
CREATE TRIGGER IF NOT EXISTS restore_store_on_delivery_delete
AFTER DELETE ON delivery_items
BEGIN
    UPDATE stores 
    SET 
        quantity = quantity + OLD.quantity,
        weight = weight + OLD.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.store_id;
END;

-- TRIGGER: Prevent negative inventory
CREATE TRIGGER IF NOT EXISTS prevent_negative_inventory
AFTER UPDATE ON stores
BEGIN
    SELECT CASE 
        WHEN NEW.quantity < 0 THEN 
            RAISE(ABORT, 'Cannot have negative quantity in store')
        WHEN NEW.weight < 0 THEN 
            RAISE(ABORT, 'Cannot have negative weight in store')
    END;
END;

-- =====================================================
-- INVENTORY AUDIT TABLE & TRIGGERS
-- =====================================================
CREATE TABLE IF NOT EXISTS inventory_audit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    store_id INTEGER NOT NULL,
    delivery_item_id INTEGER,
    action TEXT NOT NULL CHECK (action IN ('delivery_out', 'delivery_update', 'delivery_return')),
    quantity_change REAL NOT NULL,
    weight_change REAL NOT NULL,
    previous_quantity REAL NOT NULL,
    new_quantity REAL NOT NULL,
    previous_weight REAL NOT NULL,
    new_weight REAL NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (delivery_item_id) REFERENCES delivery_items(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_audit_store ON inventory_audit(store_id);
CREATE INDEX IF NOT EXISTS idx_inventory_audit_delivery ON inventory_audit(delivery_item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_audit_created_at ON inventory_audit(created_at);

-- TRIGGER: Log inventory changes on delivery insert
CREATE TRIGGER IF NOT EXISTS log_inventory_on_delivery_insert
AFTER INSERT ON delivery_items
BEGIN
    INSERT INTO inventory_audit (
        store_id, 
        delivery_item_id, 
        action, 
        quantity_change, 
        weight_change,
        previous_quantity,
        new_quantity,
        previous_weight,
        new_weight,
        created_by
    )
    SELECT 
        NEW.store_id,
        NEW.id,
        'delivery_out',
        -NEW.quantity,
        -NEW.weight,
        (SELECT quantity FROM stores WHERE id = NEW.store_id) + NEW.quantity,
        (SELECT quantity FROM stores WHERE id = NEW.store_id),
        (SELECT weight FROM stores WHERE id = NEW.store_id) + NEW.weight,
        (SELECT weight FROM stores WHERE id = NEW.store_id),
        (SELECT user_id FROM deliveries WHERE id = NEW.delivery_id)
    FROM stores 
    WHERE id = NEW.store_id;
END;

-- TRIGGER: Log inventory changes on delivery update
CREATE TRIGGER IF NOT EXISTS log_inventory_on_delivery_update
AFTER UPDATE ON delivery_items
BEGIN
    INSERT INTO inventory_audit (
        store_id, 
        delivery_item_id, 
        action, 
        quantity_change, 
        weight_change,
        previous_quantity,
        new_quantity,
        previous_weight,
        new_weight,
        created_by
    )
    SELECT 
        NEW.store_id,
        NEW.id,
        'delivery_update',
        -NEW.quantity,
        -NEW.weight,
        (SELECT quantity FROM stores WHERE id = NEW.store_id) + NEW.quantity,
        (SELECT quantity FROM stores WHERE id = NEW.store_id),
        (SELECT weight FROM stores WHERE id = NEW.store_id) + NEW.weight,
        (SELECT weight FROM stores WHERE id = NEW.store_id),
        (SELECT user_id FROM deliveries WHERE id = NEW.delivery_id)
    FROM stores 
    WHERE id = NEW.store_id;
END;

-- TRIGGER: Log inventory changes on delivery delete
CREATE TRIGGER IF NOT EXISTS log_inventory_on_delivery_delete
AFTER DELETE ON delivery_items
BEGIN
    INSERT INTO inventory_audit (
        store_id, 
        delivery_item_id, 
        action, 
        quantity_change, 
        weight_change,
        previous_quantity,
        new_quantity,
        previous_weight,
        new_weight,
        created_by
    )
    SELECT 
        OLD.store_id,
        OLD.id,
        'delivery_return',
        OLD.quantity,
        OLD.weight,
        (SELECT quantity FROM stores WHERE id = OLD.store_id) - OLD.quantity,
        (SELECT quantity FROM stores WHERE id = OLD.store_id),
        (SELECT weight FROM stores WHERE id = OLD.store_id) - OLD.weight,
        (SELECT weight FROM stores WHERE id = OLD.store_id),
        (SELECT user_id FROM deliveries WHERE id = OLD.delivery_id)
    FROM stores 
    WHERE id = OLD.store_id;
END;

-- =====================================================
-- TIMESTAMP UPDATE TRIGGERS
-- =====================================================

CREATE TRIGGER IF NOT EXISTS update_users_timestamp
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_salaries_timestamp
AFTER UPDATE ON salaries
BEGIN
    UPDATE salaries SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_brokers_timestamp
AFTER UPDATE ON brokers
BEGIN
    UPDATE brokers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_majhis_timestamp
AFTER UPDATE ON majhis
BEGIN
    UPDATE majhis SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_godowns_timestamp
AFTER UPDATE ON godowns
BEGIN
    UPDATE godowns SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_rents_timestamp
AFTER UPDATE ON rents
BEGIN
    UPDATE rents SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_customers_timestamp
AFTER UPDATE ON customers
BEGIN
    UPDATE customers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_items_timestamp
AFTER UPDATE ON items
BEGIN
    UPDATE items SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_lots_timestamp
AFTER UPDATE ON lots
BEGIN
    UPDATE lots SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_stores_timestamp
AFTER UPDATE ON stores
BEGIN
    UPDATE stores SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_damages_timestamp
AFTER UPDATE ON damages
BEGIN
    UPDATE damages SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_deliveries_timestamp
AFTER UPDATE ON deliveries
BEGIN
    UPDATE deliveries SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_delivery_items_timestamp
AFTER UPDATE ON delivery_items
BEGIN
    UPDATE delivery_items SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_transports_timestamp
AFTER UPDATE ON transports
BEGIN
    UPDATE transports SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_vehicles_timestamp
AFTER UPDATE ON vehicles
BEGIN
    UPDATE vehicles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_expenses_timestamp
AFTER UPDATE ON expenses
BEGIN
    UPDATE expenses SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;