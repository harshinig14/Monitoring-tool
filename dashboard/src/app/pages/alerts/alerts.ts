import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { Subscription } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { WebsocketService } from '../../core/services/websocket.service';

interface AlertLog {
  severity: 'CRITICAL' | 'WARNING';
  metric: string;
  value: string;
  threshold: string;
  machine: string;
  time: string;
}

@Component({
  selector: 'app-alerts',
  templateUrl: './alerts.html',
  styleUrl: './alerts.scss',
  standalone: false,
})
export class Alerts implements OnInit, OnDestroy {
  public cpuTh: number = 80;
  public memTh: number = 80;
  public diskTh: number = 80;
  public netTh: number = 80;
  
  public saveSuccess: boolean = false;
  private wsSub!: Subscription;

  private readonly MOCK_ALERTS: AlertLog[] = [
    { severity: 'CRITICAL', metric: 'CPU', value: '92%', threshold: '80%', machine: 'DESKTOP-PRO-01', time: '2026-05-31 10:45' },
    { severity: 'WARNING', metric: 'Memory', value: '83%', threshold: '80%', machine: 'LINUX-PROD-SERVER-01', time: '2026-05-31 10:30' },
    { severity: 'WARNING', metric: 'Disk', value: '85%', threshold: '80%', machine: 'DESKTOP-PRO-01', time: '2026-05-31 09:15' }
  ];

  public alertsLog: AlertLog[] = [];

  constructor(
    private apiService: ApiService,
    private wsService: WebsocketService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    // Wrap database queries in setTimeout to push execution to next turn, preventing NG0100
    setTimeout(() => {
      // 1. Fetch thresholds
      this.apiService.getThresholds().subscribe({
        next: (res) => {
          setTimeout(() => {
            if (res) {
              this.cpuTh = res.cpu_threshold;
              this.memTh = res.memory_threshold;
              this.diskTh = res.disk_threshold;
              this.netTh = res.network_threshold;
              this.cdr.detectChanges();
            }
          }, 0);
        },
        error: (err) => {
          console.warn('Failed to load thresholds from DB, using fallback defaults', err);
        }
      });

      // 2. Fetch live alert logs
      this.fetchAlerts();
    }, 0);

    // 3. Listen for WebSocket alert broadcasts
    this.wsSub = this.wsService.getMessages().subscribe(msg => {
      if (msg && msg.type === 'alert') {
        const isNet = msg.metric.toLowerCase().includes('network') || msg.metric.toLowerCase().includes('net');
        let thresholdVal = '80%';
        if (msg.metric.toLowerCase().includes('cpu')) thresholdVal = `${this.cpuTh}%`;
        else if (msg.metric.toLowerCase().includes('mem') || msg.metric.toLowerCase().includes('ram')) thresholdVal = `${this.memTh}%`;
        else if (msg.metric.toLowerCase().includes('disk')) thresholdVal = `${this.diskTh}%`;
        else if (isNet) thresholdVal = `${this.netTh} MB/s`;

        const newAlert: AlertLog = {
          severity: msg.severity as 'CRITICAL' | 'WARNING',
          metric: msg.metric,
          value: Math.round(msg.value) + (isNet ? ' MB/s' : '%'),
          threshold: thresholdVal,
          machine: msg.machine,
          time: this.formatRelativeTime(new Date().toISOString())
        };

        setTimeout(() => {
          this.alertsLog = [newAlert, ...this.alertsLog];
          this.cdr.detectChanges();
        }, 0);
      }
    });
  }

  ngOnDestroy(): void {
    if (this.wsSub) {
      this.wsSub.unsubscribe();
    }
  }

  public fetchAlerts(): void {
    this.apiService.getAlerts().subscribe({
      next: (data) => {
        setTimeout(() => {
          if (data && data.length > 0) {
            this.alertsLog = data.map(a => ({
              severity: a.severity,
              metric: a.metric_name,
              value: Math.round(a.current_value) + (a.metric_name === 'Network' ? ' MB/s' : '%'),
              threshold: Math.round(a.threshold_value) + (a.metric_name === 'Network' ? ' MB/s' : '%'),
              machine: a.machine_name,
              time: this.formatRelativeTime(a.created_at)
            }));
          } else {
            this.alertsLog = [...this.MOCK_ALERTS];
          }
          this.cdr.detectChanges();
        }, 0);
      },
      error: () => {
        setTimeout(() => {
          this.alertsLog = [...this.MOCK_ALERTS];
          this.cdr.detectChanges();
        }, 0);
      }
    });
  }

  private formatRelativeTime(dateStr: string): string {
    if (!dateStr) return '';
    try {
      const date = new Date(dateStr);
      if (isNaN(date.getTime())) return dateStr;
      const yyyy = date.getFullYear();
      const mm = String(date.getMonth() + 1).padStart(2, '0');
      const dd = String(date.getDate()).padStart(2, '0');
      const hh = String(date.getHours()).padStart(2, '0');
      const min = String(date.getMinutes()).padStart(2, '0');
      return `${yyyy}-${mm}-${dd} ${hh}:${min}`;
    } catch {
      return dateStr;
    }
  }

  public adjust(metric: string, amount: number): void {
    if (metric === 'cpu') this.cpuTh = Math.max(0, Math.min(100, this.cpuTh + amount));
    if (metric === 'mem') this.memTh = Math.max(0, Math.min(100, this.memTh + amount));
    if (metric === 'disk') this.diskTh = Math.max(0, Math.min(100, this.diskTh + amount));
    if (metric === 'net') this.netTh = Math.max(0, Math.min(1000, this.netTh + amount));
  }

  public saveThresholds(): void {
    const payload = {
      cpu_threshold: this.cpuTh,
      memory_threshold: this.memTh,
      disk_threshold: this.diskTh,
      network_threshold: this.netTh
    };

    this.apiService.updateThresholds(payload).subscribe({
      next: (res) => {
        if (res && res.success) {
          this.saveSuccess = true;
          setTimeout(() => {
            this.saveSuccess = false;
          }, 3000);
        }
      },
      error: (err) => {
        console.error('Failed to save thresholds to DB:', err);
      }
    });
  }
}
