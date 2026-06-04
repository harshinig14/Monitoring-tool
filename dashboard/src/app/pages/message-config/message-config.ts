import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-message-config',
  templateUrl: './message-config.html',
  styleUrl: './message-config.scss',
  standalone: false,
})
export class MessageConfig implements OnInit {
  // SMTP Fields
  public server: string = '';
  public port: number = 587;
  public username: string = '';
  public passwordHidden: string = '••••••••••••••••';
  public showPasswordState: boolean = false;
  public fromAddress: string = '';
  public primaryRecipient: string = '';
  public alternateRecipient: string = '';
  
  // Template Fields
  public subjectTemplate: string = '';
  public bodyMacroText: string = '';

  // UI state
  public saveSuccess: boolean = false;
  public testLoading: boolean = false;
  public testSuccess: boolean = false;
  public testError: string = '';
  
  // Email logs audit data
  public emailLogs: any[] = [];

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.apiService.getEmailConfig().subscribe({
      next: (config) => {
        if (config) {
          this.server = config.smtp_server;
          this.port = config.smtp_port;
          this.username = config.username;
          this.passwordHidden = config.password || '••••••••••••••••';
          this.fromAddress = config.from_email;
          this.primaryRecipient = config.primary_recipient;
          this.alternateRecipient = config.alternate_recipient;
          this.subjectTemplate = config.subject_template;
          this.bodyMacroText = config.body_template;
        }
      },
      error: (err) => {
        console.warn('Failed to load email config from DB, using fallback defaults', err);
      }
    });

    this.loadEmailLogs();
  }

  public togglePassword(): void {
    this.showPasswordState = !this.showPasswordState;
  }

  public loadEmailLogs(): void {
    this.apiService.getEmailLogs().subscribe({
      next: (logs) => {
        this.emailLogs = logs || [];
      },
      error: (err) => {
        console.warn('Failed to load email logs:', err);
      }
    });
  }

  public saveSmtpSettings(): void {
    const payload = {
      smtp_server: this.server,
      smtp_port: this.port,
      username: this.username,
      password: this.passwordHidden,
      from_email: this.fromAddress,
      primary_recipient: this.primaryRecipient,
      alternate_recipient: this.alternateRecipient,
      subject_template: this.subjectTemplate,
      body_template: this.bodyMacroText
    };

    this.apiService.updateEmailConfig(payload).subscribe({
      next: (res) => {
        if (res && res.success) {
          this.saveSuccess = true;
          this.loadEmailLogs();
          setTimeout(() => {
            this.saveSuccess = false;
          }, 3000);
        }
      },
      error: (err) => {
        console.error('Failed to save email config to DB:', err);
      }
    });
  }

  public sendTestEmail(): void {
    this.testLoading = true;
    this.testSuccess = false;
    this.testError = '';

    this.apiService.sendTestEmail().subscribe({
      next: (res) => {
        this.testLoading = false;
        if (res && res.success) {
          this.testSuccess = true;
          this.loadEmailLogs();
          setTimeout(() => {
            this.testSuccess = false;
          }, 4000);
        } else {
          this.testError = 'Server returned failure.';
          this.loadEmailLogs();
        }
      },
      error: (err) => {
        this.testLoading = false;
        this.testError = err.error?.error || err.message || 'SMTP transmission error';
        this.loadEmailLogs();
      }
    });
  }

  // Parses subject macro variables
  public get parsedSubjectPreview(): string {
    return this.parseMacros(this.subjectTemplate);
  }

  // Parses body macro variables
  public get parsedBodyPreview(): string {
    return this.parseMacros(this.bodyMacroText);
  }

  private parseMacros(text: string): string {
    if (!text) return '';
    return text
      .replace(/\{\{machineName\}\}/g, 'DESKTOP-PRO-01')
      .replace(/\{\{metricName\}\}/g, 'CPU')
      .replace(/\{\{currentValue\}\}/g, '92%')
      .replace(/\{\{threshold\}\}/g, '80%')
      .replace(/\{\{timestamp\}\}/g, '2026-05-31 10:45:00 AM');
  }
}
