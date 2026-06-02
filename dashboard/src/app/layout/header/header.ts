import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { Subscription, interval } from 'rxjs';
import { filter, map } from 'rxjs/operators';

@Component({
  selector: 'app-header',
  templateUrl: './header.html',
  standalone: false,
  styleUrl: './header.scss',
})
export class Header implements OnInit, OnDestroy {
  public pageTitle: string = 'Connected Devices';
  public currentTime: string = '';
  
  private routerSub!: Subscription;
  private clockSub!: Subscription;

  constructor(private router: Router) {
    this.updateClock();
  }

  ngOnInit(): void {
    // 1. Subscribe to Route Changes to set Dynamic Page Title
    this.routerSub = this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map(() => this.router.url)
    ).subscribe(url => {
      this.setPageTitle(url);
    });
    
    // Set initial title
    this.setPageTitle(this.router.url);

    // 2. Start Live Digital Clock Ticking
    this.clockSub = interval(1000).subscribe(() => {
      this.updateClock();
    });
  }

  ngOnDestroy(): void {
    if (this.routerSub) this.routerSub.unsubscribe();
    if (this.clockSub) this.clockSub.unsubscribe();
  }

  private setPageTitle(url: string): void {
    if (url.startsWith('/devices')) {
      this.pageTitle = 'Connected Devices';
    } else if (url.startsWith('/realtime')) {
      this.pageTitle = 'Real-Time Stats';
    } else if (url.startsWith('/configuration')) {
      this.pageTitle = 'Configuration';
    } else if (url.startsWith('/alerts')) {
      this.pageTitle = 'Alerts';
    } else if (url.startsWith('/message-config')) {
      this.pageTitle = 'Message Config';
    } else if (url.startsWith('/reports')) {
      this.pageTitle = 'Reports';
    } else {
      this.pageTitle = 'SysPulse Dashboard';
    }
  }

  private updateClock(): void {
    const now = new Date();
    // Format: Jun 2, 2026 10:00:36 AM
    const options: Intl.DateTimeFormatOptions = {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
      second: '2-digit',
      hour12: true
    };
    this.currentTime = now.toLocaleDateString('en-US', options);
  }
}
