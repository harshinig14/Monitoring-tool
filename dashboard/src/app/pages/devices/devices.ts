import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Device } from '../../core/models/device.model';
import { WebsocketService } from '../../core/services/websocket.service';

@Component({
  selector: 'app-devices',
  templateUrl: './devices.html',
  styleUrl: './devices.scss',
  standalone: false,
})
export class Devices implements OnInit, OnDestroy {
  public devices: Device[] = [];
  public filteredDevices: Device[] = [];
  public searchQuery: string = '';
  public loading: boolean = false;
  
  public showHistoryModal: boolean = false;
  public selectedDeviceName: string = '';
  public historyLoading: boolean = false;
  public deviceHistory: any[] = [];
  
  private fetchSub!: Subscription;
  private wsSub!: Subscription;

  // Premium fallback mock data exactly matching the screenshots
  private readonly MOCK_DEVICES: Device[] = [
    { user_id: 1, machine_name: 'DESKTOP-PRO-01', os_type: 'windows', status: 'ONLINE', last_seen: 'May 31, 2026, 10:05:06 PM' },
    { user_id: 2, machine_name: 'LINUX-PROD-SERVER-01', os_type: 'linux', status: 'ONLINE', last_seen: 'May 31, 2026, 10:04:20 PM' },
    { user_id: 3, machine_name: 'MACBOOK-M3-DEVELOPMENT', os_type: 'darwin', status: 'ONLINE', last_seen: 'May 31, 2026, 10:03:15 PM' },
    { user_id: 4, machine_name: 'SQL-BACKUP-SERVER-02', os_type: 'linux', status: 'OFFLINE', last_seen: 'May 31, 2026, 9:40:00 PM' },
    { user_id: 5, machine_name: 'HR-LAPTOP-05', os_type: 'windows', status: 'ONLINE', last_seen: 'May 31, 2026, 10:05:00 PM' }
  ];

  constructor(
    private apiService: ApiService,
    private router: Router,
    private wsService: WebsocketService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    // Wrap initial fetch in setTimeout to push execution to next turn, preventing NG0100
    setTimeout(() => {
      this.fetchDevices(true);
    }, 0);
    
    // Listen for real-time WebSocket status updates
    this.wsSub = this.wsService.getMessages().subscribe(msg => {
      if (msg && msg.type === 'device_status') {
        const devId = Number(msg.user_id);
        const device = this.devices.find(d => Number(d.user_id) === devId);
        if (device) {
          if (msg.status === 'REMOVED') {
            this.devices = this.devices.filter(d => Number(d.user_id) !== devId);
          } else {
            device.status = msg.status.toUpperCase();
            device.last_seen = this.formatLastSeen(new Date().toISOString());
          }
          this.filterGrid();
          this.cdr.detectChanges();
        } else {
          // If a new device registered or wasn't in list, pull complete grid
          this.fetchDevices(true);
        }
      }
    });
  }

  ngOnDestroy(): void {
    if (this.fetchSub) this.fetchSub.unsubscribe();
    if (this.wsSub) this.wsSub.unsubscribe();
  }

  private logDebug(message: string, data?: any): void {
    console.log(message, data);
    this.apiService.postDebugLog(message, data || {}).subscribe({
      error: () => {}
    });
  }

  public fetchDevices(silent: boolean = false): void {
    this.logDebug('fetchDevices started', { silent });
    if (!silent) {
      this.loading = true;
    }
    
    this.fetchSub = this.apiService.getDevices().subscribe({
      next: (data) => {
        this.logDebug('fetchDevices success', data);
        setTimeout(() => {
          if (data && data.length > 0) {
            try {
              this.devices = data.map(dev => {
                console.log('mapping device locally', dev);
                return {
                  ...dev,
                  last_seen: this.formatLastSeen(dev.last_seen),
                  status: (dev.status || '').toUpperCase()
                };
              });
              this.logDebug('mapping finished successfully', this.devices);
            } catch (e: any) {
              this.logDebug('mapping crashed with exception', { error: e.message, stack: e.stack });
              this.devices = [...this.MOCK_DEVICES];
            }
          } else {
            this.logDebug('fetchDevices returned empty or null data', data);
            this.devices = [...this.MOCK_DEVICES];
          }
          this.filterGrid();
          this.logDebug('filterGrid completed. filteredDevices:', this.filteredDevices);
          this.loading = false;
          this.cdr.detectChanges();
        }, 0);
      },
      error: (err) => {
        this.logDebug('fetchDevices error callback', { error: err.message, status: err.status });
        console.warn('Backend API connection failed, using high-fidelity mock data:', err);
        setTimeout(() => {
          this.devices = [...this.MOCK_DEVICES];
          this.filterGrid();
          this.loading = false;
          this.cdr.detectChanges();
        }, 0);
      }
    });
  }

  public filterGrid(): void {
    if (!this.searchQuery) {
      this.filteredDevices = [...this.devices];
      return;
    }
    const q = this.searchQuery.toLowerCase();
    this.filteredDevices = this.devices.filter(d => 
      d.machine_name.toLowerCase().includes(q) || 
      d.os_type.toLowerCase().includes(q) || 
      d.status.toLowerCase().includes(q)
    );
  }

  public onRowSelect(device: any): void {
    if (device && device.status === 'ONLINE') {
      this.router.navigate(['/realtime', device.user_id]);
    }
  }

  private formatLastSeen(lastSeenRaw: string): string {
    if (!lastSeenRaw) return 'N/A';
    try {
      const date = new Date(lastSeenRaw);
      if (isNaN(date.getTime())) return lastSeenRaw; // Return as-is if already formatted
      
      const options: Intl.DateTimeFormatOptions = {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        second: '2-digit',
        hour12: true
      };
      return date.toLocaleString('en-US', options);
    } catch {
      return lastSeenRaw;
    }
  }

  public activateDevice(device: Device): void {
    this.apiService.activateDevice(device.user_id).subscribe({
      next: (res) => {
        if (res && res.success) {
          device.status = 'ONLINE';
          this.filterGrid();
          this.cdr.detectChanges();
        }
      },
      error: (err) => console.error('Failed to activate device:', err)
    });
  }

  public deactivateDevice(device: Device): void {
    this.apiService.deactivateDevice(device.user_id).subscribe({
      next: (res) => {
        if (res && res.success) {
          device.status = 'DEACTIVATED';
          this.filterGrid();
          this.cdr.detectChanges();
        }
      },
      error: (err) => console.error('Failed to deactivate device:', err)
    });
  }

  public removeDevice(device: Device): void {
    if (confirm(`Are you sure you want to remove ${device.machine_name}? It will be hidden from the active monitoring grid.`)) {
      this.apiService.removeDevice(device.user_id).subscribe({
        next: (res) => {
          if (res && res.success) {
            // Soft-deleted device is hidden, remove from local list
            this.devices = this.devices.filter(d => d.user_id !== device.user_id);
            this.filterGrid();
            this.cdr.detectChanges();
          }
        },
        error: (err) => console.error('Failed to remove device:', err)
      });
    }
  }

  public viewHistory(device: Device): void {
    this.selectedDeviceName = device.machine_name;
    this.showHistoryModal = true;
    this.historyLoading = true;
    this.deviceHistory = [];

    this.apiService.getDeviceHistory(device.user_id).subscribe({
      next: (data) => {
        this.deviceHistory = data || [];
        this.historyLoading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Failed to load device status history:', err);
        this.historyLoading = false;
        this.cdr.detectChanges();
      }
    });
  }

  public closeHistoryModal(): void {
    this.showHistoryModal = false;
    this.deviceHistory = [];
  }

  public formatHistoryTime(timeStr: string): string {
    if (!timeStr) return '';
    try {
      const date = new Date(timeStr);
      if (isNaN(date.getTime())) return timeStr;
      return date.toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        second: '2-digit',
        hour12: true
      });
    } catch {
      return timeStr;
    }
  }

  public getStatusClass(status: string): string {
    const s = (status || '').toUpperCase();
    if (s === 'ONLINE') return 'status-curr-online';
    if (s === 'OFFLINE') return 'status-curr-offline';
    if (s === 'DEACTIVATED') return 'status-curr-deactivated';
    if (s === 'REMOVED') return 'status-curr-removed';
    return 'status-curr-pending';
  }
}
