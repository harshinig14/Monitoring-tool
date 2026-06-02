import { Injectable } from '@angular/core';

export interface PollingConfig {
  frequency: string;
}

export interface ThresholdConfig {
  cpu: number;
  memory: number;
  disk: number;
  network: number;
}

export interface SmtpConfig {
  server: string;
  port: number;
  username: string;
  password?: string;
  fromAddress: string;
  primaryRecipient: string;
  alternateRecipient: string;
  subjectTemplate: string;
  bodyMacroText: string;
}

@Injectable({
  providedIn: 'root'
})
export class ConfigStoreService {
  private readonly KEYS = {
    POLLING: 'syspulse_polling_config',
    THRESHOLDS: 'syspulse_threshold_config',
    SMTP: 'syspulse_smtp_config'
  };

  constructor() {
    this.initDefaults();
  }

  private initDefaults(): void {
    if (!localStorage.getItem(this.KEYS.POLLING)) {
      const defaultPolling: PollingConfig = { frequency: '1 minute' };
      this.savePollingConfig(defaultPolling);
    }
    if (!localStorage.getItem(this.KEYS.THRESHOLDS)) {
      const defaultThresholds: ThresholdConfig = { cpu: 80, memory: 80, disk: 80, network: 80 };
      this.saveThresholds(defaultThresholds);
    }
    if (!localStorage.getItem(this.KEYS.SMTP)) {
      const defaultSmtp: SmtpConfig = {
        server: 'smtp.syspulse.internal',
        port: 587,
        username: 'alerts@syspulse.internal',
        password: '••••••••••••••••',
        fromAddress: 'syspulse-noreply@internal.net',
        primaryRecipient: 'harshini14.ganesh@gmail.com',
        alternateRecipient: 'admin-fallback@internal.net',
        subjectTemplate: 'Threshold Alert Notification',
        bodyMacroText: `Alert Notification\n\nMachine Name: {{machineName}}\nMetric: {{metricName}}\nCurrent Value: {{currentValue}}\nConfigured Threshold: {{threshold}}\nTimestamp: {{timestamp}}\n\nPlease investigate the system immediately.`
      };
      this.saveSmtpConfig(defaultSmtp);
    }
  }

  // 1. Polling Config
  getPollingConfig(): PollingConfig {
    const data = localStorage.getItem(this.KEYS.POLLING);
    return data ? JSON.parse(data) : { frequency: '1 minute' };
  }

  savePollingConfig(config: PollingConfig): void {
    localStorage.setItem(this.KEYS.POLLING, JSON.stringify(config));
  }

  // 2. Thresholds Config
  getThresholds(): ThresholdConfig {
    const data = localStorage.getItem(this.KEYS.THRESHOLDS);
    return data ? JSON.parse(data) : { cpu: 80, memory: 80, disk: 80, network: 80 };
  }

  saveThresholds(config: ThresholdConfig): void {
    localStorage.setItem(this.KEYS.THRESHOLDS, JSON.stringify(config));
  }

  // 3. SMTP Config
  getSmtpConfig(): SmtpConfig {
    const data = localStorage.getItem(this.KEYS.SMTP);
    return data ? JSON.parse(data) : {
      server: '',
      port: 587,
      username: '',
      fromAddress: '',
      primaryRecipient: '',
      alternateRecipient: '',
      subjectTemplate: '',
      bodyMacroText: ''
    };
  }

  saveSmtpConfig(config: SmtpConfig): void {
    localStorage.setItem(this.KEYS.SMTP, JSON.stringify(config));
  }
}
