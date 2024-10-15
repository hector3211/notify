-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT,
    last_name TEXT,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    role TEXT DEFAULT "user",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Invoices table with user_id as foreign key
CREATE TABLE IF NOT EXISTS invoices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    invoice TEXT,
    status TEXT DEFAULT "pending",
    install_date TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

