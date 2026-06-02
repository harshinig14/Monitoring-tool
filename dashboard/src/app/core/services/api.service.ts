import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Device } from '../models/device.model';
import { MetricSet } from '../models/metrics.model';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private baseUrl = 'http://localhost:8081/api/v1';

  constructor(private http: HttpClient) {}

  // 1. Devices API
  getDevices(): Observable<Device[]> {
    return this.http.get<Device[]>(`${this.baseUrl}/devices`);
  }

  getDeviceById(id: number): Observable<Device> {
    return this.http.get<Device>(`${this.baseUrl}/devices/${id}`);
  }

  // 2. Real-Time Metrics API (userId is user_id)
  getRealtimeMetrics(userId: number): Observable<MetricSet> {
    return this.http.get<MetricSet>(`${this.baseUrl}/metrics/realtime/${userId}`);
  }

  // 3. Hourly Trend Metrics API
  getHourlyMetrics(userId: number): Observable<MetricSet[]> {
    return this.http.get<MetricSet[]>(`${this.baseUrl}/metrics/hourly/${userId}`);
  }

  // 4. Daily Trend Metrics API
  getDailyMetrics(userId: number): Observable<MetricSet[]> {
    return this.http.get<MetricSet[]>(`${this.baseUrl}/metrics/daily/${userId}`);
  }

  // 5. Configuration Settings
  getConfiguration(): Observable<{ polling_frequency: number }> {
    return this.http.get<{ polling_frequency: number }>(`${this.baseUrl}/configuration`);
  }

  updateConfiguration(pollingFrequency: number): Observable<{ success: boolean }> {
    return this.http.put<{ success: boolean }>(`${this.baseUrl}/configuration`, { polling_frequency: pollingFrequency });
  }

  // 6. Threshold Configuration
  getThresholds(): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/alerts/thresholds`);
  }

  updateThresholds(thresholds: any): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/alerts/thresholds`, thresholds);
  }

  // 7. Email / Message Configuration
  getEmailConfig(): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/alerts/email-config`);
  }

  updateEmailConfig(config: any): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/alerts/email-config`, config);
  }

  // 8. Agent Download URLs
  getAgentDownloadUrl(os: string): string {
    return `${this.baseUrl}/download-agent/${os}`;
  }

  // 9. Reports Data
  getReportsDaily(): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/reports/daily`);
  }

  getReportsWeekly(): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/reports/weekly`);
  }

  getReportsMonthly(): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/reports/monthly`);
  }
}
