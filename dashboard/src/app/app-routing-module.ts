import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainLayout } from './layout/main-layout/main-layout';
import { Devices } from './pages/devices/devices';
import { Realtime } from './pages/realtime/realtime';
import { Configuration } from './pages/configuration/configuration';
import { Alerts } from './pages/alerts/alerts';
import { MessageConfig } from './pages/message-config/message-config';
import { Reports } from './pages/reports/reports';

const routes: Routes = [
  {
    path: '',
    component: MainLayout,
    children: [
      { path: '', redirectTo: 'devices', pathMatch: 'full' },
      { path: 'devices', component: Devices },
      { path: 'realtime', component: Realtime },
      { path: 'realtime/:userId', component: Realtime },
      { path: 'configuration', component: Configuration },
      { path: 'alerts', component: Alerts },
      { path: 'message-config', component: MessageConfig },
      { path: 'reports', component: Reports }
    ]
  },
  { path: '**', redirectTo: 'devices' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
