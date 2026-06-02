import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../core/services/api.service';

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
export class Alerts implements OnInit {
  public cpuTh: number = 80;
  public memTh: number = 80;
  public diskTh: number = 80;
  public netTh: number = 80;
  
  public saveSuccess: boolean = false;

  public alertsLog: AlertLog[] = [
    { severity: 'CRITICAL', metric: 'CPU', value: '92%', threshold: '80%', machine: 'DESKTOP-PRO-01', time: '2026-05-31 10:45' },
    { severity: 'WARNING', metric: 'Memory', value: '83%', threshold: '80%', machine: 'LINUX-PROD-SERVER-01', time: '2026-05-31 10:30' },
    { severity: 'WARNING', metric: 'Disk', value: '85%', threshold: '80%', machine: 'DESKTOP-PRO-01', time: '2026-05-31 09:15' }
  ];

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.apiService.getThresholds().subscribe({
      next: (res) => {
        if (res) {
          this.cpuTh = res.cpu_threshold;
          this.memTh = res.memory_threshold;
          this.diskTh = res.disk_threshold;
          this.netTh = res.network_threshold;
        }
      },
      error: (err) => {
        console.warn('Failed to load thresholds from DB, using fallback defaults', err);
      }
    });
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
