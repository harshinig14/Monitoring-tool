-- SQL Script to set up Configuration, Alerts, SMTP and Alerts Log Tables
-- Run this in pgAdmin or query tool using a superuser role

CREATE TABLE IF NOT EXISTS configurations (
    id SERIAL PRIMARY KEY,
    polling_frequency INTEGER,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alert_thresholds (
    id SERIAL PRIMARY KEY,
    cpu_threshold FLOAT NOT NULL DEFAULT 80,
    memory_threshold FLOAT NOT NULL DEFAULT 80,
    disk_threshold FLOAT NOT NULL DEFAULT 80,
    network_threshold FLOAT NOT NULL DEFAULT 80,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS email_configuration (
    id SERIAL PRIMARY KEY,
    smtp_server TEXT,
    smtp_port INTEGER,
    username TEXT,
    password TEXT,
    from_email TEXT,
    primary_recipient TEXT,
    alternate_recipient TEXT,
    subject_template TEXT,
    body_template TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alerts (
    alert_id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    machine_name VARCHAR(255),
    metric_name VARCHAR(50),
    current_value FLOAT,
    threshold_value FLOAT,
    severity VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS email_logs (
    email_id BIGSERIAL PRIMARY KEY,
    recipient VARCHAR(255),
    subject TEXT,
    status VARCHAR(50),
    sent_at TIMESTAMP DEFAULT NOW()
);

-- Optional: Insert initial configurations & seed thresholds
INSERT INTO configurations (polling_frequency) 
SELECT 60 WHERE NOT EXISTS (SELECT 1 FROM configurations);

INSERT INTO alert_thresholds (cpu_threshold, memory_threshold, disk_threshold, network_threshold)
SELECT 80.0, 80.0, 80.0, 80.0 WHERE NOT EXISTS (SELECT 1 FROM alert_thresholds);

INSERT INTO email_configuration (smtp_server, smtp_port, username, password, from_email, primary_recipient, alternate_recipient, subject_template, body_template) 
SELECT 'smtp.syspulse.internal', 587, 'alerts@syspulse.internal', '••••••••••••••••', 'syspulse-noreply@internal.net', 'harshini14.ganesh@gmail.com', 'admin-fallback@internal.net', 'Threshold Alert Notification', 'Alert Notification

Machine Name: {{machineName}}
Metric: {{metricName}}
Current Value: {{currentValue}}
Configured Threshold: {{threshold}}
Timestamp: {{timestamp}}

Please investigate the system immediately.'
WHERE NOT EXISTS (SELECT 1 FROM email_configuration);
