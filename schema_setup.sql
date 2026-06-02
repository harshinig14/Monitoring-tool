-- SQL Script to set up Configuration, Alerts and SMTP Tables
-- Run this in pgAdmin or query tool using a superuser role

CREATE TABLE IF NOT EXISTS configurations (
    id SERIAL PRIMARY KEY,
    polling_frequency INTEGER,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alert_thresholds (
    id SERIAL PRIMARY KEY,
    cpu_threshold FLOAT,
    memory_threshold FLOAT,
    disk_threshold FLOAT,
    network_threshold FLOAT
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
    body_template TEXT
);

-- Optional: Insert initial configurations
INSERT INTO configurations (polling_frequency) VALUES (60);
INSERT INTO alert_thresholds (cpu_threshold, memory_threshold, disk_threshold, network_threshold) VALUES (80.0, 80.0, 80.0, 80.0);
INSERT INTO email_configuration (smtp_server, smtp_port, username, password, from_email, primary_recipient, alternate_recipient, subject_template, body_template) 
VALUES ('smtp.syspulse.internal', 587, 'alerts@syspulse.internal', '••••••••••••••••', 'syspulse-noreply@internal.net', 'harshini14.ganesh@gmail.com', 'admin-fallback@internal.net', 'Threshold Alert Notification', 'Alert Notification

Machine Name: {{machineName}}
Metric: {{metricName}}
Current Value: {{currentValue}}
Configured Threshold: {{threshold}}
Timestamp: {{timestamp}}

Please investigate the system immediately.');
