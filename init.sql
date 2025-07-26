DROP TABLE IF EXISTS payments;

CREATE TYPE processor AS ENUM ('default', 'fallback');

CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  correlation_id VARCHAR(255) NOT NULL,
  amount DECIMAL NOT NULL,
  processed_at TIMESTAMP NOT NULL,
  processed_by processor NOT NULL
);

-- Performance indexes
CREATE INDEX idx_payments_processed_at ON payments(processed_at);
CREATE INDEX idx_payments_processed_by ON payments(processed_by);
CREATE INDEX idx_payments_processed_at_by ON payments(processed_at, processed_by);

