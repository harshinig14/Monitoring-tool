import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { Subscription, interval } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Device } from '../../core/models/device.model';

interface DropdownOption {
  label: string;
  value: number;
}

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.html',
  styleUrl: './configuration.scss',
  standalone: false,
})
export class Configuration implements OnInit, OnDestroy {
  public monitoringIntervals: DropdownOption[] = [
    { label: '10 seconds', value: 10 },
    { label: '30 seconds', value: 30 },
    { label: '1 minute', value: 60 },
    { label: '5 minutes', value: 300 },
    { label: '10 minutes', value: 600 }
  ];
  
  public selectedInterval: number = 60;
  public saveSuccess: boolean = false;

  // Add Device Workflow states
  public showAddDevicePanel: boolean = false;
  public selectedOS: string = '';

  constructor(private apiService: ApiService, private cdr: ChangeDetectorRef) {}

  ngOnInit(): void {
    // Load config
    this.apiService.postDebugLog('Config ngOnInit started', {}).subscribe();
    this.apiService.getConfiguration().subscribe({
      next: (res) => {
        this.apiService.postDebugLog('Config ngOnInit loaded res', res).subscribe();
        if (res && res.polling_frequency) {
          this.selectedInterval = res.polling_frequency;
          this.apiService.postDebugLog('Config ngOnInit set selectedInterval to', { selectedInterval: this.selectedInterval }).subscribe();
          this.cdr.detectChanges();
        } else {
          this.apiService.postDebugLog('Config ngOnInit: res or polling_frequency missing/falsy', res || {}).subscribe();
        }
      },
      error: (err) => {
        console.warn('Failed to load live configuration, using default: 60s', err);
        this.apiService.postDebugLog('Config ngOnInit load error', err).subscribe();
      }
    });
  }

  ngOnDestroy(): void {
  }

  public saveConfig(): void {
    this.apiService.updateConfiguration(this.selectedInterval).subscribe({
      next: (res) => {
        if (res && res.success) {
          this.saveSuccess = true;
          setTimeout(() => {
            this.saveSuccess = false;
          }, 3000);
        }
      },
      error: (err) => {
        console.error('Failed to save configuration to DB:', err);
      }
    });
  }

  public openAddDevice(): void {
    this.showAddDevicePanel = !this.showAddDevicePanel; // Toggle panel state
    this.selectedOS = '';
    
    if (this.showAddDevicePanel) {
      setTimeout(() => {
        const element = document.getElementById('add-device-wizard');
        if (element) {
          element.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      }, 150);
    }
  }

  public selectOS(os: string): void {
    this.selectedOS = os;
    
    // Call the device request API to log the request and create a placeholder device
    this.apiService.requestAddDevice(os).subscribe({
      next: (res) => {
        if (res && res.download_url) {
          // Trigger file download using the returned url path (running on port 8081)
          window.location.href = 'http://localhost:8081' + res.download_url;
        }
      },
      error: (err) => {
        console.error('Failed to request device placeholder:', err);
        // Fallback to direct download url if API fails
        window.location.href = this.apiService.getAgentDownloadUrl(os);
      }
    });
  }
}
