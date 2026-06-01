CREATE TABLE users
(
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100),
    machine_name VARCHAR(255),
    os_type VARCHAR(50),
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE trace_results
(
    trace_id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    cpu_usage FLOAT,
    memory_usage FLOAT,
    disk_usage FLOAT,
    network_usage FLOAT,
    trace_date TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(user_id)
);

CREATE INDEX idx_trace_user
ON trace_results(user_id);

CREATE INDEX idx_trace_date
ON trace_results(trace_date);

CREATE INDEX idx_users_machine
ON users(machine_name);
