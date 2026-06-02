import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../core/services/api.service';

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
export class Configuration implements OnInit {
  public monitoringIntervals: DropdownOption[] = [
    { label: '10 sec', value: 10 },
    { label: '30 sec', value: 30 },
    { label: '1 min', value: 60 },
    { label: '5 min', value: 300 },
    { label: '10 min', value: 600 }
  ];
  
  public selectedInterval: number = 60;
  public saveSuccess: boolean = false;

  // Add Device Workflow states
  public showAddDevicePanel: boolean = false;
  public selectedOS: string = '';

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.apiService.getConfiguration().subscribe({
      next: (res) => {
        if (res && res.polling_frequency) {
          this.selectedInterval = res.polling_frequency;
        }
      },
      error: (err) => {
        console.warn('Failed to load live configuration, using default: 60s', err);
      }
    });
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
    this.showAddDevicePanel = true;
    this.selectedOS = '';
  }

  public selectOS(os: string): void {
    this.selectedOS = os;
    
    // Trigger agent download from backend
    const downloadUrl = this.apiService.getAgentDownloadUrl(os);
    window.location.href = downloadUrl;
  }
}
