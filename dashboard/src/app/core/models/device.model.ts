export interface Device {
  user_id: number;
  machine_name: string;
  os_type: string;
  status: string; // 'ONLINE' or 'OFFLINE'
  last_seen: string;
}
