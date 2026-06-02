export interface MetricSet {
  cpu_usage: number;
  memory_usage: number;
  disk_usage: number;
  network_usage: number;
  trace_date: string;
}

export interface MetricTrendResponse {
  cpu: number[];
  memory: number[];
  disk: number[];
  network: number[];
  timestamps: string[];
}
