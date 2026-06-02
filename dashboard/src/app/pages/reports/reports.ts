import { Component, OnInit } from '@angular/core';
import { ChartConfiguration, ChartDataset } from 'chart.js';
import { ApiService } from '../../core/services/api.service';
interface ReportLog {
  period: string;
  avgCpu: string;
  avgMem: string;
  avgDisk: string;
  avgNet: string;
}

@Component({
  selector: 'app-reports',
  templateUrl: './reports.html',
  styleUrl: './reports.scss',
  standalone: false,
})
export class Reports implements OnInit {
  public activeTab: string = 'daily'; // 'daily', 'weekly', 'monthly'
  
  // Averaged stats shown in cards
  public avgCpu: number = 46;
  public avgMem: number = 62;
  public avgDisk: number = 75;
  public avgNet: string = '5.6 MB/s';

  // Grouped bar chart configurations
  public chartLabels: string[] = [];
  
  public chartDatasets: ChartDataset[] = [
    { data: [], label: 'CPU %', backgroundColor: '#00d2ff', hoverBackgroundColor: '#00b0d6', borderRadius: 4 },
    { data: [], label: 'Memory %', backgroundColor: '#00e676', hoverBackgroundColor: '#00c862', borderRadius: 4 },
    { data: [], label: 'Disk %', backgroundColor: '#ff9100', hoverBackgroundColor: '#e07f00', borderRadius: 4 },
    { data: [], label: 'Network MB/s', backgroundColor: '#af40ff', hoverBackgroundColor: '#962be6', borderRadius: 4 }
  ];

  public chartOptions: ChartConfiguration['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: 'top',
        labels: {
          color: '#7d8b9b',
          font: { family: 'Inter', size: 10, weight: 'bold' },
          boxWidth: 15,
          padding: 15
        }
      }
    },
    scales: {
      x: {
        grid: { color: 'rgba(255, 255, 255, 0.02)' },
        ticks: { color: '#7d8b9b', font: { family: 'Inter', size: 10 } }
      },
      y: {
        min: 0,
        max: 100,
        grid: { color: 'rgba(255, 255, 255, 0.02)' },
        ticks: { color: '#7d8b9b', font: { family: 'Inter', size: 10 } }
      }
    }
  };

  public currentLogs: ReportLog[] = [];

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.loadReportData();
  }

  public changeTab(tab: string): void {
    this.activeTab = tab;
    this.loadReportData();
  }

  private loadReportData(): void {
    let report$;
    if (this.activeTab === 'daily') {
      report$ = this.apiService.getReportsDaily();
    } else if (this.activeTab === 'weekly') {
      report$ = this.apiService.getReportsWeekly();
    } else {
      report$ = this.apiService.getReportsMonthly();
    }

    report$.subscribe({
      next: (res) => {
        if (res) {
          this.avgCpu = res.avgCpu;
          this.avgMem = res.avgMem;
          this.avgDisk = res.avgDisk;
          this.avgNet = res.avgNet;
          this.chartLabels = res.chartLabels;
          
          this.chartDatasets[0].data = res.cpuData;
          this.chartDatasets[1].data = res.memoryData;
          this.chartDatasets[2].data = res.diskData;
          this.chartDatasets[3].data = res.networkData;
          
          this.currentLogs = res.logs;
          
          // Trigger chart update
          this.chartLabels = [...this.chartLabels];
        }
      },
      error: (err) => {
        console.error('Failed to load reports from backend:', err);
      }
    });
  }

  public downloadCSV(): void {
    let csvContent = 'data:text/csv;charset=utf-8,';
    csvContent += 'Period,Avg CPU,Avg Memory,Avg Disk,Avg Network\r\n';
    
    this.currentLogs.forEach(row => {
      csvContent += `"${row.period}","${row.avgCpu}","${row.avgMem}","${row.avgDisk}","${row.avgNet}"\r\n`;
    });

    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', `SysPulse_${this.activeTab}_report.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
}
