import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { MessageService } from 'primeng/api';
import { WebsocketService } from '../../core/services/websocket.service';

@Component({
  selector: 'app-main-layout',
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
  standalone: false,
})
export class MainLayout implements OnInit, OnDestroy {
  private wsSub!: Subscription;

  constructor(
    private wsService: WebsocketService,
    private messageService: MessageService
  ) {}

  ngOnInit(): void {
    this.wsSub = this.wsService.getMessages().subscribe(msg => {
      if (msg && msg.type === 'alert') {
        const isNet = msg.metric.toLowerCase().includes('network') || msg.metric.toLowerCase().includes('net');
        const unit = isNet ? ' MB/s' : '%';
        const formattedValue = Math.round(msg.value) + unit;
        
        const toastSeverity = msg.severity === 'CRITICAL' ? 'error' : 'warn';
        const toastSummary = msg.severity === 'CRITICAL' ? 'Critical System Alert' : 'System Warning';
        
        this.messageService.add({
          severity: toastSeverity,
          summary: toastSummary,
          detail: `${msg.machine}: ${msg.metric} usage is ${formattedValue}!`,
          life: 6000
        });
      }
    });
  }

  ngOnDestroy(): void {
    if (this.wsSub) {
      this.wsSub.unsubscribe();
    }
  }
}
