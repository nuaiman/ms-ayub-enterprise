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
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);

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
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (name != '')
);

CREATE INDEX IF NOT EXISTS idx_godowns_name ON godowns(name);
CREATE INDEX IF NOT EXISTS idx_godowns_phone ON godowns(phone);

-- =====================================================
-- 6. CUSTOMERS
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
-- 7. LOTS
-- =====================================================
CREATE TABLE IF NOT EXISTS lots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    customer_id INTEGER NOT NULL,
    lot_number TEXT NOT NULL,
    product_name TEXT NOT NULL,
    weight_unit TEXT NOT NULL DEFAULT 'kg',
    quantity_unit TEXT NOT NULL DEFAULT 'units',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(customer_id, lot_number),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_lots_user_id ON lots(user_id);
CREATE INDEX IF NOT EXISTS idx_lots_customer_id ON lots(customer_id);
CREATE INDEX IF NOT EXISTS idx_lots_lot_number ON lots(lot_number);
CREATE INDEX IF NOT EXISTS idx_lots_product_name ON lots(product_name);

-- =====================================================
-- 8. STORES
-- =====================================================
CREATE TABLE IF NOT EXISTS stores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    lot_id INTEGER NOT NULL,
    godown_id INTEGER NOT NULL,
    weight REAL NOT NULL DEFAULT 0 CHECK (weight >= 0),
    quantity REAL NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    start_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active INTEGER NOT NULL DEFAULT 1,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE CASCADE,
    FOREIGN KEY (godown_id) REFERENCES godowns(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_stores_user_id ON stores(user_id);
CREATE INDEX IF NOT EXISTS idx_stores_lot_id ON stores(lot_id);
CREATE INDEX IF NOT EXISTS idx_stores_godown_id ON stores(godown_id);
CREATE INDEX IF NOT EXISTS idx_stores_start_date ON stores(start_date);
CREATE INDEX IF NOT EXISTS idx_stores_is_active ON stores(is_active);

-- =====================================================
-- 9. STORE ADJUSTMENTS
-- =====================================================
CREATE TABLE IF NOT EXISTS store_adjustments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    adjustment_type TEXT NOT NULL CHECK (
        adjustment_type IN ('delta', 'absolute')
    ),
    input_weight REAL NOT NULL DEFAULT 0,
    input_quantity REAL NOT NULL DEFAULT 0,
    weight_delta REAL NOT NULL DEFAULT 0,
    quantity_delta REAL NOT NULL DEFAULT 0,
    reason TEXT,
    notes TEXT,
    adjusted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_store_adjustments_store_id ON store_adjustments(store_id);
CREATE INDEX IF NOT EXISTS idx_store_adjustments_user_id ON store_adjustments(user_id);
CREATE INDEX IF NOT EXISTS idx_store_adjustments_adjusted_at ON store_adjustments(adjusted_at);
CREATE INDEX IF NOT EXISTS idx_store_adjustments_store_adjusted ON store_adjustments(store_id, adjusted_at DESC);
CREATE INDEX IF NOT EXISTS idx_store_adjustments_type ON store_adjustments(adjustment_type);

-- =====================================================
-- 10. STORE TRANSFERS
-- =====================================================
CREATE TABLE IF NOT EXISTS store_transfers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    from_godown_id INTEGER NOT NULL,
    to_godown_id INTEGER NOT NULL,
    weight_at_transfer REAL NOT NULL,
    quantity_at_transfer REAL NOT NULL,
    notes TEXT,
    transferred_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
    FOREIGN KEY (from_godown_id) REFERENCES godowns(id) ON DELETE RESTRICT,
    FOREIGN KEY (to_godown_id) REFERENCES godowns(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_store_transfers_store_id ON store_transfers(store_id);
CREATE INDEX IF NOT EXISTS idx_store_transfers_user_id ON store_transfers(user_id);
CREATE INDEX IF NOT EXISTS idx_store_transfers_transferred_at ON store_transfers(transferred_at);
CREATE INDEX IF NOT EXISTS idx_store_transfers_store_transferred ON store_transfers(store_id, transferred_at DESC);
CREATE INDEX IF NOT EXISTS idx_store_transfers_from_godown ON store_transfers(from_godown_id);
CREATE INDEX IF NOT EXISTS idx_store_transfers_to_godown ON store_transfers(to_godown_id);

-- =====================================================
-- 11. LOT TRANSFERS (ACCOUNT TRANSFERS)
-- =====================================================
CREATE TABLE IF NOT EXISTS lot_transfers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    lot_id INTEGER NOT NULL,
    from_customer_id INTEGER NOT NULL,
    to_customer_id INTEGER NOT NULL,
    notes TEXT,
    transferred_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE CASCADE,
    FOREIGN KEY (from_customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    FOREIGN KEY (to_customer_id) REFERENCES customers(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_lot_transfers_lot_id ON lot_transfers(lot_id);
CREATE INDEX IF NOT EXISTS idx_lot_transfers_user_id ON lot_transfers(user_id);
CREATE INDEX IF NOT EXISTS idx_lot_transfers_transferred_at ON lot_transfers(transferred_at);
CREATE INDEX IF NOT EXISTS idx_lot_transfers_lot_transferred ON lot_transfers(lot_id, transferred_at DESC);
CREATE INDEX IF NOT EXISTS idx_lot_transfers_from_customer ON lot_transfers(from_customer_id);
CREATE INDEX IF NOT EXISTS idx_lot_transfers_to_customer ON lot_transfers(to_customer_id);

-- =====================================================
-- 12. GODOWN BILLS (MONTHLY)
-- =====================================================
CREATE TABLE IF NOT EXISTS godown_bills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    godown_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    month_year TEXT NOT NULL,
    bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (
        bill_type IN ('weight', 'quantity', 'fixed')
    ),
    rate REAL NOT NULL DEFAULT 0 CHECK (rate >= 0),
    weight_at_billing REAL NOT NULL DEFAULT 0,
    quantity_at_billing REAL NOT NULL DEFAULT 0,
    weight_unit_at_billing TEXT,
    quantity_unit_at_billing TEXT,
    total_amount REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        month_year GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]' AND
        substr(month_year, 6, 2) BETWEEN '01' AND '12'
    ),

    UNIQUE(godown_id, store_id, month_year),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (godown_id) REFERENCES godowns(id) ON DELETE RESTRICT,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_godown_bills_user_id ON godown_bills(user_id);
CREATE INDEX IF NOT EXISTS idx_godown_bills_godown_id ON godown_bills(godown_id);
CREATE INDEX IF NOT EXISTS idx_godown_bills_store_id ON godown_bills(store_id);
CREATE INDEX IF NOT EXISTS idx_godown_bills_month_year ON godown_bills(month_year);
CREATE INDEX IF NOT EXISTS idx_godown_bills_godown_month ON godown_bills(godown_id, month_year);
CREATE INDEX IF NOT EXISTS idx_godown_bills_bill_type ON godown_bills(bill_type);
CREATE INDEX IF NOT EXISTS idx_godown_bills_created_at ON godown_bills(created_at);

-- =====================================================
-- 13. MAJHI BILLS (ONE-TIME)
-- Source is EITHER a store OR a delivery item.
-- =====================================================
CREATE TABLE IF NOT EXISTS majhi_bills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    majhi_id INTEGER NOT NULL,
    store_id INTEGER,
    delivery_item_id INTEGER,
    bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (
        bill_type IN ('weight', 'quantity', 'job')
    ),
    rate REAL NOT NULL DEFAULT 0 CHECK (rate >= 0),
    weight_at_billing REAL NOT NULL DEFAULT 0,
    quantity_at_billing REAL NOT NULL DEFAULT 0,
    weight_unit_at_billing TEXT,
    quantity_unit_at_billing TEXT,
    total_amount REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        (store_id IS NOT NULL AND delivery_item_id IS NULL) OR
        (store_id IS NULL AND delivery_item_id IS NOT NULL)
    ),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (majhi_id) REFERENCES majhis(id) ON DELETE RESTRICT,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
    FOREIGN KEY (delivery_item_id) REFERENCES delivery_items(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_majhi_bills_user_id ON majhi_bills(user_id);
CREATE INDEX IF NOT EXISTS idx_majhi_bills_majhi_id ON majhi_bills(majhi_id);
CREATE INDEX IF NOT EXISTS idx_majhi_bills_store_id ON majhi_bills(store_id);
CREATE INDEX IF NOT EXISTS idx_majhi_bills_delivery_item_id ON majhi_bills(delivery_item_id);
CREATE INDEX IF NOT EXISTS idx_majhi_bills_bill_type ON majhi_bills(bill_type);
CREATE INDEX IF NOT EXISTS idx_majhi_bills_created_at ON majhi_bills(created_at);

-- =====================================================
-- 14. CUSTOMER STORE BILLS (MONTHLY)
-- =====================================================
CREATE TABLE IF NOT EXISTS customer_store_bills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    customer_id INTEGER NOT NULL,
    month_year TEXT NOT NULL,
    bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (
        bill_type IN ('weight', 'quantity')
    ),
    rate REAL NOT NULL DEFAULT 0 CHECK (rate >= 0),
    weight_at_billing REAL NOT NULL DEFAULT 0,
    quantity_at_billing REAL NOT NULL DEFAULT 0,
    weight_unit_at_billing TEXT,
    quantity_unit_at_billing TEXT,
    total_amount REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        month_year GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]' AND
        substr(month_year, 6, 2) BETWEEN '01' AND '12'
    ),

    UNIQUE(customer_id, store_id, month_year),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_customer_store_bills_user_id ON customer_store_bills(user_id);
CREATE INDEX IF NOT EXISTS idx_customer_store_bills_customer_id ON customer_store_bills(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_store_bills_store_id ON customer_store_bills(store_id);
CREATE INDEX IF NOT EXISTS idx_customer_store_bills_month_year ON customer_store_bills(month_year);
CREATE INDEX IF NOT EXISTS idx_customer_store_bills_customer_month ON customer_store_bills(customer_id, month_year);
CREATE INDEX IF NOT EXISTS idx_customer_store_bills_bill_type ON customer_store_bills(bill_type);
CREATE INDEX IF NOT EXISTS idx_customer_store_bills_created_at ON customer_store_bills(created_at);

-- =====================================================
-- 15. EXPENSES
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
-- 16. INCOMES
-- =====================================================
CREATE TABLE IF NOT EXISTS incomes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    amount REAL NOT NULL DEFAULT 0 CHECK (amount >= 0),
    income_date DATETIME NOT NULL,
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_incomes_date ON incomes(income_date);
CREATE INDEX IF NOT EXISTS idx_incomes_user_id ON incomes(user_id);

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
-- 18. BILL PAYMENTS
-- =====================================================
CREATE TABLE IF NOT EXISTS bill_payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    bill_type TEXT NOT NULL CHECK (
        bill_type IN ('godown', 'majhi', 'customer_store', 'customer_delivery', 'customer_additional')
    ),
    bill_id INTEGER NOT NULL,
    amount REAL NOT NULL CHECK (amount > 0),
    payment_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payment_method TEXT CHECK (
        payment_method IN ('cash', 'bank_transfer', 'check', 'mobile_banking')
    ),
    reference_number TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_bill_payments_user_id ON bill_payments(user_id);
CREATE INDEX IF NOT EXISTS idx_bill_payments_bill ON bill_payments(bill_type, bill_id);
CREATE INDEX IF NOT EXISTS idx_bill_payments_payment_date ON bill_payments(payment_date);
CREATE INDEX IF NOT EXISTS idx_bill_payments_created_at ON bill_payments(created_at);

-- =====================================================
-- 19. DELIVERIES (Master)
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
CREATE INDEX IF NOT EXISTS idx_deliveries_delivery_date ON deliveries(delivery_date);

-- =====================================================
-- 20. DELIVERY ITEMS (Detail)
-- =====================================================
CREATE TABLE IF NOT EXISTS delivery_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    delivery_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    majhi_id INTEGER NOT NULL,
    vehicle_number TEXT,
    driver_number TEXT,
    quantity REAL NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    weight REAL NOT NULL DEFAULT 0 CHECK (weight >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE RESTRICT,
    FOREIGN KEY (majhi_id) REFERENCES majhis(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_delivery_items_delivery_id ON delivery_items(delivery_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_store_id ON delivery_items(store_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_majhi_id ON delivery_items(majhi_id);

-- =====================================================
-- 21. CUSTOMER DELIVERY BILLS (ONE-TIME)
-- =====================================================
CREATE TABLE IF NOT EXISTS customer_delivery_bills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    customer_id INTEGER NOT NULL,
    delivery_item_id INTEGER NOT NULL,
    bill_type TEXT NOT NULL DEFAULT 'quantity' CHECK (
        bill_type IN ('weight', 'quantity')
    ),
    rate REAL NOT NULL DEFAULT 0 CHECK (rate >= 0),
    weight_at_billing REAL NOT NULL DEFAULT 0,
    quantity_at_billing REAL NOT NULL DEFAULT 0,
    weight_unit_at_billing TEXT,
    quantity_unit_at_billing TEXT,
    total_amount REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(delivery_item_id),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    FOREIGN KEY (delivery_item_id) REFERENCES delivery_items(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_customer_delivery_bills_user_id ON customer_delivery_bills(user_id);
CREATE INDEX IF NOT EXISTS idx_customer_delivery_bills_customer_id ON customer_delivery_bills(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_delivery_bills_delivery_item_id ON customer_delivery_bills(delivery_item_id);
CREATE INDEX IF NOT EXISTS idx_customer_delivery_bills_bill_type ON customer_delivery_bills(bill_type);
CREATE INDEX IF NOT EXISTS idx_customer_delivery_bills_created_at ON customer_delivery_bills(created_at);

-- =====================================================
-- 22. DAMAGES
-- =====================================================
CREATE TABLE IF NOT EXISTS damages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    store_id INTEGER NOT NULL,
    quantity REAL NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_unit TEXT NOT NULL DEFAULT 'units',
    weight REAL NOT NULL DEFAULT 0 CHECK (weight >= 0),
    weight_unit TEXT NOT NULL DEFAULT 'kg',
    damage_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason TEXT NOT NULL,
    amount REAL NOT NULL DEFAULT 0 CHECK (amount >= 0),
    notes TEXT,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (quantity > 0 OR weight > 0),
    CHECK (reason != ''),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_damages_user_id ON damages(user_id);
CREATE INDEX IF NOT EXISTS idx_damages_store_id ON damages(store_id);
CREATE INDEX IF NOT EXISTS idx_damages_damage_date ON damages(damage_date);
CREATE INDEX IF NOT EXISTS idx_damages_store_date ON damages(store_id, damage_date DESC);

-- =====================================================
-- 23. CUSTOMER ADDITIONAL BILLS (ONE-TIME)
-- =====================================================
CREATE TABLE IF NOT EXISTS customer_additional_bills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    customer_id INTEGER NOT NULL,
    amount REAL NOT NULL CHECK (amount >= 0),
    description TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (description != ''),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_customer_additional_bills_user_id ON customer_additional_bills(user_id);
CREATE INDEX IF NOT EXISTS idx_customer_additional_bills_customer_id ON customer_additional_bills(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_additional_bills_created_at ON customer_additional_bills(created_at);

-- =====================================================
-- 24. INVENTORY AUDIT
-- =====================================================
CREATE TABLE IF NOT EXISTS inventory_audit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    store_id INTEGER NOT NULL,
    delivery_item_id INTEGER,
    damage_id INTEGER,
    action TEXT NOT NULL CHECK (action IN (
        'delivery_out',
        'delivery_update',
        'delivery_return',
        'damage_out',
        'damage_update',
        'damage_return'
    )),
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
    FOREIGN KEY (damage_id) REFERENCES damages(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_audit_store ON inventory_audit(store_id);
CREATE INDEX IF NOT EXISTS idx_inventory_audit_delivery ON inventory_audit(delivery_item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_audit_damage ON inventory_audit(damage_id);
CREATE INDEX IF NOT EXISTS idx_inventory_audit_created_at ON inventory_audit(created_at);

-- =====================================================
-- 25. INVOICES
-- =====================================================
CREATE TABLE IF NOT EXISTS invoices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    entity_type TEXT NOT NULL CHECK (
        entity_type IN ('customer', 'majhi', 'godown', 'broker')
    ),
    entity_id INTEGER NOT NULL,
    subtotal REAL NOT NULL DEFAULT 0,
    discount_amount REAL NOT NULL DEFAULT 0,
    total REAL NOT NULL DEFAULT 0,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_invoices_user_id ON invoices(user_id);
CREATE INDEX IF NOT EXISTS idx_invoices_entity ON invoices(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_invoices_created_at ON invoices(created_at);

-- =====================================================
-- 26. INVOICE ITEMS
-- Records the individual bills an invoice was built from.
-- Polymorphic reference (bill_type + bill_id) - no FK, since the
-- target table varies. Enforced in app code.
-- =====================================================
CREATE TABLE IF NOT EXISTS invoice_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    invoice_id INTEGER NOT NULL,
    bill_type TEXT NOT NULL CHECK (
        bill_type IN ('godown', 'majhi', 'customer_store', 'customer_delivery', 'customer_additional')
    ),
    bill_id INTEGER NOT NULL,
    amount REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(invoice_id, bill_type, bill_id),

    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice_id ON invoice_items(invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_bill ON invoice_items(bill_type, bill_id);

-- =====================================================
-- 27. DISCOUNTS
-- =====================================================
CREATE TABLE IF NOT EXISTS discounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    invoice_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('flat', 'percent')),
    value REAL NOT NULL CHECK (value >= 0),
    computed_amount REAL NOT NULL DEFAULT 0,
    reason TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_discounts_invoice_id ON discounts(invoice_id);

-- =====================================================
-- INVENTORY MANAGEMENT TRIGGERS
-- =====================================================

-- ---------------------------------------------------------------------------
-- BEFORE INSERT: reject if the store does not have enough stock.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS check_stock_before_delivery_insert
BEFORE INSERT ON delivery_items
WHEN (SELECT quantity FROM stores WHERE id = NEW.store_id) < NEW.quantity
  OR (SELECT weight   FROM stores WHERE id = NEW.store_id) < NEW.weight
BEGIN
    SELECT RAISE(ABORT, 'Insufficient stock in store');
END;

-- ---------------------------------------------------------------------------
-- BEFORE UPDATE (same store).
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS check_stock_before_delivery_update_same_store
BEFORE UPDATE ON delivery_items
WHEN NEW.store_id = OLD.store_id
  AND (
      (SELECT quantity FROM stores WHERE id = NEW.store_id) + OLD.quantity < NEW.quantity
      OR
      (SELECT weight   FROM stores WHERE id = NEW.store_id) + OLD.weight   < NEW.weight
  )
BEGIN
    SELECT RAISE(ABORT, 'Insufficient stock in store for update');
END;

-- ---------------------------------------------------------------------------
-- BEFORE UPDATE (store changed).
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS check_stock_before_delivery_update_new_store
BEFORE UPDATE ON delivery_items
WHEN NEW.store_id != OLD.store_id
  AND (
      (SELECT quantity FROM stores WHERE id = NEW.store_id) < NEW.quantity
      OR
      (SELECT weight   FROM stores WHERE id = NEW.store_id) < NEW.weight
  )
BEGIN
    SELECT RAISE(ABORT, 'Insufficient stock in new store');
END;

-- ---------------------------------------------------------------------------
-- AFTER INSERT: audit + decrement the store.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS delivery_item_out
AFTER INSERT ON delivery_items
BEGIN
    INSERT INTO inventory_audit (
        store_id, delivery_item_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        NEW.store_id, NEW.id, 'delivery_out',
        -NEW.quantity, -NEW.weight,
        s.quantity, s.quantity - NEW.quantity,
        s.weight,   s.weight   - NEW.weight,
        (SELECT user_id FROM deliveries WHERE id = NEW.delivery_id)
    FROM stores s WHERE s.id = NEW.store_id;

    UPDATE stores
    SET quantity   = quantity - NEW.quantity,
        weight     = weight   - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- ---------------------------------------------------------------------------
-- AFTER UPDATE: audit + adjust the store(s).
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS delivery_item_update
AFTER UPDATE ON delivery_items
WHEN NEW.store_id = OLD.store_id
  AND (NEW.quantity != OLD.quantity OR NEW.weight != OLD.weight)
BEGIN
    INSERT INTO inventory_audit (
        store_id, delivery_item_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        NEW.store_id, NEW.id, 'delivery_update',
        -(NEW.quantity - OLD.quantity), -(NEW.weight - OLD.weight),
        s.quantity + (NEW.quantity - OLD.quantity),
        s.quantity,
        s.weight   + (NEW.weight - OLD.weight),
        s.weight,
        (SELECT user_id FROM deliveries WHERE id = NEW.delivery_id)
    FROM stores s WHERE s.id = NEW.store_id;

    UPDATE stores
    SET quantity   = quantity + OLD.quantity - NEW.quantity,
        weight     = weight   + OLD.weight   - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

CREATE TRIGGER IF NOT EXISTS delivery_item_store_move
AFTER UPDATE ON delivery_items
WHEN NEW.store_id != OLD.store_id
BEGIN
    INSERT INTO inventory_audit (
        store_id, delivery_item_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        OLD.store_id, OLD.id, 'delivery_return',
        OLD.quantity, OLD.weight,
        s.quantity, s.quantity + OLD.quantity,
        s.weight,   s.weight   + OLD.weight,
        (SELECT user_id FROM deliveries WHERE id = NEW.delivery_id)
    FROM stores s WHERE s.id = OLD.store_id;

    UPDATE stores
    SET quantity   = quantity + OLD.quantity,
        weight     = weight   + OLD.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.store_id;

    INSERT INTO inventory_audit (
        store_id, delivery_item_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        NEW.store_id, NEW.id, 'delivery_out',
        -NEW.quantity, -NEW.weight,
        s.quantity, s.quantity - NEW.quantity,
        s.weight,   s.weight   - NEW.weight,
        (SELECT user_id FROM deliveries WHERE id = NEW.delivery_id)
    FROM stores s WHERE s.id = NEW.store_id;

    UPDATE stores
    SET quantity   = quantity - NEW.quantity,
        weight     = weight   - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- ---------------------------------------------------------------------------
-- AFTER DELETE: audit + restore the store.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS delivery_item_return
AFTER DELETE ON delivery_items
BEGIN
    INSERT INTO inventory_audit (
        store_id, delivery_item_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        OLD.store_id, OLD.id, 'delivery_return',
        OLD.quantity, OLD.weight,
        s.quantity, s.quantity + OLD.quantity,
        s.weight,   s.weight   + OLD.weight,
        (SELECT user_id FROM deliveries WHERE id = OLD.delivery_id)
    FROM stores s WHERE s.id = OLD.store_id;

    UPDATE stores
    SET quantity   = quantity + OLD.quantity,
        weight     = weight   + OLD.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.store_id;
END;

-- ---------------------------------------------------------------------------
-- DAMAGES: BEFORE INSERT / UPDATE guards.
-- A damage always reduces store stock, so it can never exceed what's on hand.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS check_stock_before_damage_insert
BEFORE INSERT ON damages
WHEN (SELECT quantity FROM stores WHERE id = NEW.store_id) < NEW.quantity
  OR (SELECT weight   FROM stores WHERE id = NEW.store_id) < NEW.weight
BEGIN
    SELECT RAISE(ABORT, 'Insufficient stock in store for damage');
END;

CREATE TRIGGER IF NOT EXISTS check_stock_before_damage_update_same_store
BEFORE UPDATE ON damages
WHEN NEW.store_id = OLD.store_id
  AND (
      (SELECT quantity FROM stores WHERE id = NEW.store_id) + OLD.quantity < NEW.quantity
      OR
      (SELECT weight   FROM stores WHERE id = NEW.store_id) + OLD.weight   < NEW.weight
  )
BEGIN
    SELECT RAISE(ABORT, 'Insufficient stock in store for damage update');
END;

CREATE TRIGGER IF NOT EXISTS check_stock_before_damage_update_new_store
BEFORE UPDATE ON damages
WHEN NEW.store_id != OLD.store_id
  AND (
      (SELECT quantity FROM stores WHERE id = NEW.store_id) < NEW.quantity
      OR
      (SELECT weight   FROM stores WHERE id = NEW.store_id) < NEW.weight
  )
BEGIN
    SELECT RAISE(ABORT, 'Insufficient stock in new store for damage');
END;

-- ---------------------------------------------------------------------------
-- DAMAGES: AFTER INSERT. Audit + decrement.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS damage_out
AFTER INSERT ON damages
BEGIN
    INSERT INTO inventory_audit (
        store_id, damage_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        NEW.store_id, NEW.id, 'damage_out',
        -NEW.quantity, -NEW.weight,
        s.quantity, s.quantity - NEW.quantity,
        s.weight,   s.weight   - NEW.weight,
        NEW.user_id
    FROM stores s WHERE s.id = NEW.store_id;

    UPDATE stores
    SET quantity   = quantity - NEW.quantity,
        weight     = weight   - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- ---------------------------------------------------------------------------
-- DAMAGES: AFTER UPDATE, same store. Audit + adjust.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS damage_update
AFTER UPDATE ON damages
WHEN NEW.store_id = OLD.store_id
  AND (NEW.quantity != OLD.quantity OR NEW.weight != OLD.weight)
BEGIN
    INSERT INTO inventory_audit (
        store_id, damage_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        NEW.store_id, NEW.id, 'damage_update',
        -(NEW.quantity - OLD.quantity), -(NEW.weight - OLD.weight),
        s.quantity + (NEW.quantity - OLD.quantity),
        s.quantity,
        s.weight   + (NEW.weight - OLD.weight),
        s.weight,
        NEW.user_id
    FROM stores s WHERE s.id = NEW.store_id;

    UPDATE stores
    SET quantity   = quantity + OLD.quantity - NEW.quantity,
        weight     = weight   + OLD.weight   - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- ---------------------------------------------------------------------------
-- DAMAGES: AFTER UPDATE, store changed. Restore old, apply new.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS damage_store_move
AFTER UPDATE ON damages
WHEN NEW.store_id != OLD.store_id
BEGIN
    INSERT INTO inventory_audit (
        store_id, damage_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        OLD.store_id, OLD.id, 'damage_return',
        OLD.quantity, OLD.weight,
        s.quantity, s.quantity + OLD.quantity,
        s.weight,   s.weight   + OLD.weight,
        NEW.user_id
    FROM stores s WHERE s.id = OLD.store_id;

    UPDATE stores
    SET quantity   = quantity + OLD.quantity,
        weight     = weight   + OLD.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.store_id;

    INSERT INTO inventory_audit (
        store_id, damage_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        NEW.store_id, NEW.id, 'damage_out',
        -NEW.quantity, -NEW.weight,
        s.quantity, s.quantity - NEW.quantity,
        s.weight,   s.weight   - NEW.weight,
        NEW.user_id
    FROM stores s WHERE s.id = NEW.store_id;

    UPDATE stores
    SET quantity   = quantity - NEW.quantity,
        weight     = weight   - NEW.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.store_id;
END;

-- ---------------------------------------------------------------------------
-- DAMAGES: AFTER DELETE. Audit + restore.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS damage_return
AFTER DELETE ON damages
BEGIN
    INSERT INTO inventory_audit (
        store_id, damage_id, action,
        quantity_change, weight_change,
        previous_quantity, new_quantity,
        previous_weight, new_weight,
        created_by
    )
    SELECT
        OLD.store_id, OLD.id, 'damage_return',
        OLD.quantity, OLD.weight,
        s.quantity, s.quantity + OLD.quantity,
        s.weight,   s.weight   + OLD.weight,
        OLD.user_id
    FROM stores s WHERE s.id = OLD.store_id;

    UPDATE stores
    SET quantity   = quantity + OLD.quantity,
        weight     = weight   + OLD.weight,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.store_id;
END;

-- ---------------------------------------------------------------------------
-- Auto-deactivate a store when it empties.
-- ---------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS deactivate_store_when_empty
AFTER UPDATE ON stores
WHEN NEW.quantity = 0
  AND NEW.weight = 0
  AND NEW.is_active = 1
  AND (OLD.quantity > 0 OR OLD.weight > 0)
BEGIN
    UPDATE stores
    SET is_active = 0,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
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

CREATE TRIGGER IF NOT EXISTS update_customers_timestamp
AFTER UPDATE ON customers
BEGIN
    UPDATE customers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
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

CREATE TRIGGER IF NOT EXISTS update_store_adjustments_timestamp
AFTER UPDATE ON store_adjustments
BEGIN
    UPDATE store_adjustments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_store_transfers_timestamp
AFTER UPDATE ON store_transfers
BEGIN
    UPDATE store_transfers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_lot_transfers_timestamp
AFTER UPDATE ON lot_transfers
BEGIN
    UPDATE lot_transfers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_expenses_timestamp
AFTER UPDATE ON expenses
BEGIN
    UPDATE expenses SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_incomes_timestamp
AFTER UPDATE ON incomes
BEGIN
    UPDATE incomes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_majhi_bills_timestamp
AFTER UPDATE ON majhi_bills
BEGIN
    UPDATE majhi_bills SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_godown_bills_timestamp
AFTER UPDATE ON godown_bills
BEGIN
    UPDATE godown_bills SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_customer_store_bills_timestamp
AFTER UPDATE ON customer_store_bills
BEGIN
    UPDATE customer_store_bills SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
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

CREATE TRIGGER IF NOT EXISTS update_customer_delivery_bills_timestamp
AFTER UPDATE ON customer_delivery_bills
BEGIN
    UPDATE customer_delivery_bills SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_damages_timestamp
AFTER UPDATE ON damages
BEGIN
    UPDATE damages SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_customer_additional_bills_timestamp
AFTER UPDATE ON customer_additional_bills
BEGIN
    UPDATE customer_additional_bills SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_invoices_timestamp
AFTER UPDATE ON invoices
BEGIN
    UPDATE invoices SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_discounts_timestamp
AFTER UPDATE ON discounts
BEGIN
    UPDATE discounts SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;