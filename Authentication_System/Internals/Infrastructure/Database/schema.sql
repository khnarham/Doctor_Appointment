CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'patient',
    gender VARCHAR(10) NOT NULL DEFAULT '',
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email_lower ON users(LOWER(email));
CREATE INDEX idx_users_created_at ON users(created_at);

CREATE TABLE email_verifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_verifications_expires_at ON email_verifications(expires_at);

CREATE INDEX idx_email_verifications_user_id ON email_verifications(user_id);

CREATE TABLE patients (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NULL REFERENCES users(id) ON DELETE SET NULL,
    walk_in_name VARCHAR(255) NULL, 
    walk_in_phone VARCHAR(20) NOT NULL DEFAULT '',
    gender VARCHAR(10) NOT NULL,
    date_of_birth DATE NOT NULL,
    is_walk_in BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_patients_user_id ON patients(user_id);
CREATE INDEX idx_patients_walk_in_phone ON patients(walk_in_phone) WHERE walk_in_phone != '';