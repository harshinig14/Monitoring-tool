import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription, interval } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Device } from '../../core/models/device.model';
import { MetricSet } from '../../core/models/metrics.model';
import { ChartConfiguration, ChartDataset } from 'chart.js';
import { WebsocketService } from '../../core/services/websocket.service';

@Component({
  selector: 'app-realtime',
  templateUrl: './realtime.html',
  styleUrl: './realtime.scss',
  standalone: false,
})
export class Realtime implements OnInit, OnDestroy {
  public userId: number | null = null;
  public device: Device | null = null;
  
  // Current values
  public cpuUsage: number = 40;
  public memoryUsage: number = 49;
  public diskUsage: number = 40;
  public networkUsage: number = 1.8; // in MB/s
  
  // Active Tab
  public activeTab: string = '5m'; // '5m', '1h', '24h'

  // Chart configs
  public chartLabels: string[] = [];
  
  public cpuDatasets: ChartDataset[] = [{ data: [], label: 'CPU %', borderColor: '#00d2ff', backgroundColor: 'rgba(0, 210, 255, 0.05)', pointBackgroundColor: '#00d2ff' }];
  public memoryDatasets: ChartDataset[] = [{ data: [], label: 'Memory %', borderColor: '#00e676', backgroundColor: 'rgba(0, 230, 118, 0.05)', pointBackgroundColor: '#00e676' }];
  public diskDatasets: ChartDataset[] = [{ data: [], label: 'Disk %', borderColor: '#ff9100', backgroundColor: 'rgba(255, 145, 0, 0.05)', pointBackgroundColor: '#ff9100' }];
  public networkDatasets: ChartDataset[] = [{ data: [], label: 'Network MB/s', borderColor: '#af40ff', backgroundColor: 'rgba(175, 64, 255, 0.05)', pointBackgroundColor: '#af40ff' }];

  public chartOptions: ChartConfiguration['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { display: false } },
    scales: {
      x: {
        grid: { color: 'rgba(255, 255, 255, 0.03)' },
        ticks: { color: '#7d8b9b', font: { family: 'JetBrains Mono', size: 9 }, maxRotation: 0 }
      },
      y: {
        min: 0,
        max: 100,
        grid: { color: 'rgba(255, 255, 255, 0.03)' },
        ticks: { color: '#7d8b9b', font: { family: 'JetBrains Mono', size: 9 } }
      }
    },
    elements: {
      line: { tension: 0.3, borderWidth: 2, fill: true },
      point: { radius: 0, hoverRadius: 5 }
    }
  };

  public networkChartOptions: ChartConfiguration['options'] = {
    ...this.chartOptions,
    scales: {
      x: this.chartOptions!.scales!['x'],
      y: {
        min: 0,
        grid: { color: 'rgba(255, 255, 255, 0.03)' },
        ticks: { color: '#7d8b9b', font: { family: 'JetBrains Mono', size: 9 } }
      }
    }
  };

  private routeSub!: Subscription;
  private wsSub!: Subscription;
  private mockIntervalSub: Subscription | null = null;
  private isMockMode: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private apiService: ApiService,
    private wsService: WebsocketService,
    private cdr: ChangeDetectorRef
  ) {}

  private logDebug(message: string, data?: any): void {
    console.log(message, data);
    this.apiService.postDebugLog(message, data || {}).subscribe({
      error: () => {}
    });
  }

  ngOnInit(): void {
    this.routeSub = this.route.params.subscribe(params => {
      if (params['userId']) {
        this.userId = Number(params['userId']);
      } else {
        const savedId = localStorage.getItem('selectedDeviceId');
        this.userId = savedId ? Number(savedId) : 1;
      }
      localStorage.setItem('selectedDeviceId', this.userId.toString());
      this.logDebug('Realtime component initialized for userId: ' + this.userId);
      this.loadDeviceDetails();
      this.loadMetricsHistory();
    });

    // Subscribe to WebSocket messages
    this.wsSub = this.wsService.getMessages().subscribe(msg => {
      if (msg) {
        if (msg.type === 'metric' && Number(msg.user_id) === this.userId) {
          // Live data arrived! Stop mock data tick if it is running
          if (this.isMockMode) {
            this.isMockMode = false;
            this.stopMockInterval();
          }

          this.cpuUsage = Math.round(msg.cpu);
          this.memoryUsage = Math.round(msg.memory);
          this.diskUsage = Math.round(msg.disk);
          this.networkUsage = Number((msg.network / (1024 * 1024)).toFixed(1));

          const metricsPayload: MetricSet = {
            cpu_usage: msg.cpu,
            memory_usage: msg.memory,
            disk_usage: msg.disk,
            network_usage: msg.network,
            trace_date: new Date().toISOString()
          };
          if (this.activeTab === '5m') {
            this.appendLatestToChart(metricsPayload);
          }
          this.cdr.detectChanges();
        } else if (msg.type === 'device_status' && Number(msg.user_id) === this.userId) {
          if (this.device) {
            this.device.status = (msg.status || '').toUpperCase();
            if (msg.status === 'OFFLINE') {
              this.stopMockInterval();
            }
            this.cdr.detectChanges();
          }
        }
      }
    });
  }

  ngOnDestroy(): void {
    if (this.routeSub) this.routeSub.unsubscribe();
    if (this.wsSub) this.wsSub.unsubscribe();
    this.stopMockInterval();
  }

  public changeTab(tab: string): void {
    this.activeTab = tab;
    this.loadMetricsHistory();
  }

  private loadDeviceDetails(): void {
    if (this.userId === null) return;
    this.logDebug('loadDeviceDetails starting for userId: ' + this.userId);
    
    this.apiService.getDeviceById(this.userId).subscribe({
      next: (dev) => {
        this.logDebug('loadDeviceDetails next: ', dev);
        if (dev) {
          this.device = { ...dev, status: (dev.status || '').toUpperCase() };
          this.logDebug('loadDeviceDetails device set to: ', this.device);
        } else {
          this.setupMockDevice();
        }
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.logDebug('loadDeviceDetails error: ', { message: err.message, status: err.status });
        this.setupMockDevice();
        this.cdr.detectChanges();
      }
    });
  }

  private setupMockDevice(): void {
    this.isMockMode = true;
    const mockDevices: { [key: number]: Device } = {
      1: { user_id: 1, machine_name: 'DESKTOP-PRO-01', os_type: 'windows', status: 'ONLINE', last_seen: 'Just now' },
      2: { user_id: 2, machine_name: 'LINUX-PROD-SERVER-01', os_type: 'linux', status: 'ONLINE', last_seen: 'Just now' },
      3: { user_id: 3, machine_name: 'MACBOOK-M3-DEVELOPMENT', os_type: 'darwin', status: 'ONLINE', last_seen: 'Just now' },
      4: { user_id: 4, machine_name: 'SQL-BACKUP-SERVER-02', os_type: 'linux', status: 'OFFLINE', last_seen: '3 hours ago' },
      5: { user_id: 5, machine_name: 'HR-LAPTOP-05', os_type: 'windows', status: 'ONLINE', last_seen: 'Just now' }
    };
    this.device = mockDevices[this.userId!] || mockDevices[1];
    this.logDebug('setupMockDevice done: ', this.device);
  }

  private loadMetricsHistory(): void {
    if (this.userId === null) return;
    this.logDebug('loadMetricsHistory starting for userId: ' + this.userId + ', tab: ' + this.activeTab);
    
    // Call proper endpoint depending on active tab
    const history$ = this.activeTab === '24h' 
      ? this.apiService.getDailyMetrics(this.userId) 
      : this.apiService.getHourlyMetrics(this.userId);

    history$.subscribe({
      next: (data) => {
        this.logDebug('loadMetricsHistory success: length = ' + (data ? data.length : 0));
        if (data && data.length > 0) {
          this.isMockMode = false;
          this.stopMockInterval();
          this.populateChartData(data);
        } else {
          this.generateMockHistory();
          this.startMockInterval();
        }
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.logDebug('loadMetricsHistory error: ', { message: err.message, status: err.status });
        this.generateMockHistory();
        this.startMockInterval();
        this.cdr.detectChanges();
      }
    });
  }

  private startMockInterval(): void {
    if (this.mockIntervalSub) return;
    this.mockIntervalSub = interval(5000).subscribe(() => {
      this.tickMockMetrics();
    });
  }

  private stopMockInterval(): void {
    if (this.mockIntervalSub) {
      this.mockIntervalSub.unsubscribe();
      this.mockIntervalSub = null;
    }
  }

  private downsample(list: MetricSet[], targetSize: number): MetricSet[] {
    if (list.length <= targetSize) {
      return list;
    }
    const result: MetricSet[] = [];
    const step = (list.length - 1) / (targetSize - 1);
    for (let i = 0; i < targetSize; i++) {
      const index = Math.round(i * step);
      result.push(list[index]);
    }
    return result;
  }

  private populateChartData(metricsList: MetricSet[]): void {
    if (!metricsList || metricsList.length === 0) return;

    let cleanList: MetricSet[] = [];

    if (this.activeTab === '5m') {
      // Filter to only include metrics from the last 5 minutes relative to the latest data point
      const latestDate = new Date(metricsList[metricsList.length - 1].trace_date);
      const fiveMinutesAgo = new Date(latestDate.getTime() - 5 * 60 * 1000);
      const filtered = metricsList.filter(m => new Date(m.trace_date) >= fiveMinutesAgo);
      // Downsample to 30 points if we have more
      cleanList = this.downsample(filtered, 30);
    } else if (this.activeTab === '1h') {
      // Downsample the entire 1-hour metrics list to 30 points
      cleanList = this.downsample(metricsList, 30);
    } else {
      // Downsample the entire 24-hour metrics list to 30 points
      cleanList = this.downsample(metricsList, 30);
    }

    const labels: string[] = [];
    const cpu: number[] = [];
    const mem: number[] = [];
    const disk: number[] = [];
    const net: number[] = [];

    cleanList.forEach(m => {
      labels.push(this.formatTimeLabel(m.trace_date));
      cpu.push(m.cpu_usage);
      mem.push(m.memory_usage);
      disk.push(m.disk_usage);
      net.push(Number((m.network_usage / (1024 * 1024)).toFixed(2)));
    });

    this.chartLabels = labels;
    this.cpuDatasets[0].data = cpu;
    this.memoryDatasets[0].data = mem;
    this.diskDatasets[0].data = disk;
    this.networkDatasets[0].data = net;
    
    // Set current values to latest raw point from metricsList
    const latest = metricsList[metricsList.length - 1];
    this.cpuUsage = Math.round(latest.cpu_usage);
    this.memoryUsage = Math.round(latest.memory_usage);
    this.diskUsage = Math.round(latest.disk_usage);
    this.networkUsage = Number((latest.network_usage / (1024 * 1024)).toFixed(1));
  }

  private appendLatestToChart(m: MetricSet): void {
    const label = this.formatTimeLabel(m.trace_date);
    
    // Maintain maximum 30 points
    if (this.chartLabels.length >= 30) {
      this.chartLabels.shift();
      this.cpuDatasets[0].data!.shift();
      this.memoryDatasets[0].data!.shift();
      this.diskDatasets[0].data!.shift();
      this.networkDatasets[0].data!.shift();
    }
    
    this.chartLabels.push(label);
    this.cpuDatasets[0].data!.push(m.cpu_usage);
    this.memoryDatasets[0].data!.push(m.memory_usage);
    this.diskDatasets[0].data!.push(m.disk_usage);
    this.networkDatasets[0].data!.push(Number((m.network_usage / (1024 * 1024)).toFixed(2)));
    
    // Refresh charts
    this.chartLabels = [...this.chartLabels];
  }

  private formatTimeLabel(dateStr: string): string {
    try {
      const date = new Date(dateStr);
      if (isNaN(date.getTime())) return dateStr;
      if (this.activeTab === '24h') {
        return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false });
      }
      return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });
    } catch {
      return dateStr;
    }
  }

  /* --- Premium Live Mock Generators for robust offline fallback --- */
  private generateMockHistory(): void {
    this.isMockMode = true;
    const pointsCount = this.activeTab === '5m' ? 30 : this.activeTab === '1h' ? 24 : 20;
    const labels: string[] = [];
    const cpu: number[] = [];
    const mem: number[] = [];
    const disk: number[] = [];
    const net: number[] = [];

    const now = new Date();
    // Setup initial base values depending on device selected
    let baseCpu = this.userId === 2 ? 60 : this.userId === 3 ? 30 : 40;
    let baseMem = this.userId === 2 ? 65 : this.userId === 3 ? 55 : 49;
    let baseDisk = this.userId === 2 ? 75 : 40;
    let baseNet = this.userId === 2 ? 4.5 : this.userId === 3 ? 2.1 : 1.8;

    for (let i = pointsCount - 1; i >= 0; i--) {
      const d = new Date(now.getTime() - i * (this.activeTab === '5m' ? 10000 : this.activeTab === '1h' ? 150000 : 43200000));
      labels.push(d.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }));
      
      // Add random walks
      baseCpu += (Math.random() - 0.5) * 8;
      baseCpu = Math.max(10, Math.min(95, baseCpu));
      
      baseMem += (Math.random() - 0.5) * 2;
      baseMem = Math.max(20, Math.min(90, baseMem));
      
      baseDisk += (Math.random() - 0.5) * 0.2;
      baseDisk = Math.max(10, Math.min(95, baseDisk));
      
      baseNet += (Math.random() - 0.5) * 0.8;
      baseNet = Math.max(0.1, baseNet);

      cpu.push(Math.round(baseCpu));
      mem.push(Math.round(baseMem));
      disk.push(Math.round(baseDisk));
      net.push(Number(baseNet.toFixed(1)));
    }

    this.chartLabels = labels;
    this.cpuDatasets[0].data = cpu;
    this.memoryDatasets[0].data = mem;
    this.diskDatasets[0].data = disk;
    this.networkDatasets[0].data = net;

    // Set current dashboard values to last element
    this.cpuUsage = cpu[cpu.length - 1];
    this.memoryUsage = mem[mem.length - 1];
    this.diskUsage = disk[disk.length - 1];
    this.networkUsage = net[net.length - 1];
  }

  private tickMockMetrics(): void {
    if (this.device?.status === 'OFFLINE') return;
    
    // Generate new random walk step
    this.cpuUsage = Math.round(Math.max(10, Math.min(95, this.cpuUsage + (Math.random() - 0.5) * 10)));
    this.memoryUsage = Math.round(Math.max(20, Math.min(90, this.memoryUsage + (Math.random() - 0.5) * 4)));
    this.diskUsage = Math.round(Math.max(10, Math.min(95, this.diskUsage + (Math.random() - 0.5) * 0.5)));
    this.networkUsage = Number(Math.max(0.1, this.networkUsage + (Math.random() - 0.5) * 0.6).toFixed(1));

    const now = new Date();
    const label = now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });

    if (this.chartLabels.length >= 30) {
      this.chartLabels.shift();
      this.cpuDatasets[0].data!.shift();
      this.memoryDatasets[0].data!.shift();
      this.diskDatasets[0].data!.shift();
      this.networkDatasets[0].data!.shift();
    }

    this.chartLabels.push(label);
    this.cpuDatasets[0].data!.push(this.cpuUsage);
    this.memoryDatasets[0].data!.push(this.memoryUsage);
    this.diskDatasets[0].data!.push(this.diskUsage);
    this.networkDatasets[0].data!.push(this.networkUsage);

    this.chartLabels = [...this.chartLabels];
  }
}
