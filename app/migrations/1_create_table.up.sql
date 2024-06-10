CREATE TABLE IF NOT EXISTS bills (
        id TEXT PRIMARY KEY,
        currency TEXT NOT NULL,
        is_closed BOOLEAN NOT NULL,
        created_at TIMESTAMP NOT NULL,
        closed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS line_items (
        id SERIAL PRIMARY KEY,
        bill_id TEXT NOT NULL,
        description TEXT NOT NULL,
        amount REAL NOT NULL,
        currency TEXT NOT NULL,
        FOREIGN KEY (bill_id) REFERENCES bills (id)
)