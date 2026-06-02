import { NgModule, provideBrowserGlobalErrorListeners } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing-module';
import { App } from './app';
import { Sidebar } from './layout/sidebar/sidebar';
import { Header } from './layout/header/header';
import { MainLayout } from './layout/main-layout/main-layout';
import { Devices } from './pages/devices/devices';
import { Realtime } from './pages/realtime/realtime';
import { Configuration } from './pages/configuration/configuration';
import { Alerts } from './pages/alerts/alerts';
import { MessageConfig } from './pages/message-config/message-config';
import { Reports } from './pages/reports/reports';

// PrimeNG Config and Theme Preset
import { providePrimeNG } from 'primeng/config';
import Aura from '@primeng/themes/aura';

// PrimeNG Modules
import { TableModule } from 'primeng/table';
import { SelectModule } from 'primeng/select';
import { InputTextModule } from 'primeng/inputtext';
import { InputNumberModule } from 'primeng/inputnumber';
import { ButtonModule } from 'primeng/button';
import { ProgressBarModule } from 'primeng/progressbar';
import { DatePickerModule } from 'primeng/datepicker';

// Chart.js Module
import { BaseChartDirective, provideCharts, withDefaultRegisterables } from 'ng2-charts';

@NgModule({
  declarations: [
    App,
    Sidebar,
    Header,
    MainLayout,
    Devices,
    Realtime,
    Configuration,
    Alerts,
    MessageConfig,
    Reports,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    AppRoutingModule,
    
    // PrimeNG Modules
    TableModule,
    SelectModule,
    InputTextModule,
    InputNumberModule,
    ButtonModule,
    ProgressBarModule,
    DatePickerModule,
    
    // Charts Directive
    BaseChartDirective
  ],
  providers: [
    provideBrowserGlobalErrorListeners(),
    providePrimeNG({
      theme: {
        preset: Aura
      }
    }),
    provideCharts(withDefaultRegisterables())
  ],
  bootstrap: [App],
})
export class AppModule {}
