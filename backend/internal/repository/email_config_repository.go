package repository

import (
	"database/sql"
)

type EmailConfig struct {
	SMTPServer         string `json:"smtp_server"`
	SMTPPort           int    `json:"smtp_port"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	FromEmail          string `json:"from_email"`
	PrimaryRecipient   string `json:"primary_recipient"`
	AlternateRecipient string `json:"alternate_recipient"`
	SubjectTemplate    string `json:"subject_template"`
	BodyTemplate       string `json:"body_template"`
}

type EmailConfigRepository struct {
	DB *sql.DB
}

func NewEmailConfigRepository(db *sql.DB) *EmailConfigRepository {
	return &EmailConfigRepository{DB: db}
}

func (r *EmailConfigRepository) GetConfig() (*EmailConfig, error) {
	var c EmailConfig
	query := `SELECT smtp_server, smtp_port, username, password, from_email, primary_recipient, alternate_recipient, subject_template, body_template FROM email_configuration ORDER BY id DESC LIMIT 1`
	err := r.DB.QueryRow(query).Scan(&c.SMTPServer, &c.SMTPPort, &c.Username, &c.Password, &c.FromEmail, &c.PrimaryRecipient, &c.AlternateRecipient, &c.SubjectTemplate, &c.BodyTemplate)
	if err != nil {
		return &EmailConfig{
			SMTPServer:         "smtp.syspulse.internal",
			SMTPPort:           587,
			Username:           "alerts@syspulse.internal",
			Password:           "••••••••••••••••",
			FromEmail:          "syspulse-noreply@internal.net",
			PrimaryRecipient:   "harshini14.ganesh@gmail.com",
			AlternateRecipient: "admin-fallback@internal.net",
			SubjectTemplate:    "Threshold Alert Notification",
			BodyTemplate:       "Alert Notification\n\nMachine Name: {{machineName}}\nMetric: {{metricName}}\nCurrent Value: {{currentValue}}\nConfigured Threshold: {{threshold}}\nTimestamp: {{timestamp}}\n\nPlease investigate the system immediately.",
		}, nil
	}
	return &c, nil
}

func (r *EmailConfigRepository) SaveConfig(c *EmailConfig) error {
	query := `INSERT INTO email_configuration (smtp_server, smtp_port, username, password, from_email, primary_recipient, alternate_recipient, subject_template, body_template) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.DB.Exec(query, c.SMTPServer, c.SMTPPort, c.Username, c.Password, c.FromEmail, c.PrimaryRecipient, c.AlternateRecipient, c.SubjectTemplate, c.BodyTemplate)
	return err
}
