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

  postDebugLog(message: string, data: any): Observable<any> {
    return this.http.post(`${this.baseUrl}/debug-log`, { message, data });
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
    return this.http.get<{ polling_frequency: number }>(`${this.baseUrl}/configuration?t=` + new Date().getTime());
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
    return this.http.get<any>(`${this.baseUrl}/email-config`);
  }

  updateEmailConfig(config: any): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/email-config`, config);
  }

  sendTestEmail(): Observable<any> {
    return this.http.post<any>(`${this.baseUrl}/email/test`, {});
  }

  getEmailLogs(): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/email/logs`);
  }

  // 8. Agent Download URLs
  getAgentDownloadUrl(os: string): string {
    return `${this.baseUrl}/download-agent/${os}`;
  }

  requestAddDevice(osType: string): Observable<{ download_url: string }> {
    return this.http.post<{ download_url: string }>(`${this.baseUrl}/devices/request`, { os_type: osType });
  }

  deactivateDevice(id: number): Observable<any> {
    return this.http.post<any>(`${this.baseUrl}/devices/${id}/deactivate`, {});
  }

  activateDevice(id: number): Observable<any> {
    return this.http.post<any>(`${this.baseUrl}/devices/${id}/activate`, {});
  }

  removeDevice(id: number): Observable<any> {
    return this.http.delete<any>(`${this.baseUrl}/devices/${id}`);
  }

  getDeviceHistory(id: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/devices/${id}/history`);
  }

  // 8.5. Get Recent Alerts
  getAlerts(): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/alerts/recent`);
  }

  // 9. Reports Data
  getReportsDaily(deviceId?: number): Observable<any> {
    const url = deviceId ? `${this.baseUrl}/reports/daily?deviceId=${deviceId}` : `${this.baseUrl}/reports/daily`;
    return this.http.get<any>(url);
  }

  getReportsWeekly(deviceId?: number): Observable<any> {
    const url = deviceId ? `${this.baseUrl}/reports/weekly?deviceId=${deviceId}` : `${this.baseUrl}/reports/weekly`;
    return this.http.get<any>(url);
  }

  getReportsMonthly(deviceId?: number): Observable<any> {
    const url = deviceId ? `${this.baseUrl}/reports/monthly?deviceId=${deviceId}` : `${this.baseUrl}/reports/monthly`;
    return this.http.get<any>(url);
  }

  getReportsHistory(limit: number = 30, deviceId?: number): Observable<any[]> {
    const url = deviceId ? `${this.baseUrl}/reports/history?limit=${limit}&deviceId=${deviceId}` : `${this.baseUrl}/reports/history?limit=${limit}`;
    return this.http.get<any[]>(url);
  }

  getCSVExportUrl(from: string, to: string, type: string, deviceId?: number): string {
    let url = `${this.baseUrl}/reports/export/csv?from=${from}&to=${to}&type=${type}`;
    if (deviceId) {
      url += `&deviceId=${deviceId}`;
    }
    return url;
  }
}
