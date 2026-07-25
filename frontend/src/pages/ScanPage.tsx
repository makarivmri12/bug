import React from 'react';
import { Activity, AlertTriangle, CheckCircle, Clock } from 'lucide-react';
import { useScanStore } from '../stores/scanStore';

const ScanPage: React.FC = () => {
  const { scans, isScanning, findings } = useScanStore();
  const [targetUrl, setTargetUrl] = React.useState('');

  const handleStartScan = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Starting scan for:', targetUrl);
    // API call to start scan
  };

  const stats = [
    {
      label: 'Total Findings',
      value: findings.length,
      icon: AlertTriangle,
      color: 'text-yellow-500',
    },
    {
      label: 'Critical Issues',
      value: findings.filter((f) => f.severity === 'CRITICAL').length,
      icon: AlertTriangle,
      color: 'text-red-500',
    },
    {
      label: 'Confirmed',
      value: findings.filter((f) => f.status === 'CONFIRMED').length,
      icon: CheckCircle,
      color: 'text-green-500',
    },
    {
      label: 'Active Scans',
      value: scans.filter((s) => s.status === 'running').length,
      icon: Activity,
      color: 'text-blue-500',
    },
  ];

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>

      {/* Quick Scan Form */}
      <div className="bg-gray-800 rounded-lg border border-gray-700 p-6">
        <h2 className="text-xl font-bold mb-4">Quick Scan</h2>
        <form onSubmit={handleStartScan} className="flex gap-2">
          <input
            type="url"
            value={targetUrl}
            onChange={(e) => setTargetUrl(e.target.value)}
            placeholder="Enter target URL (e.g., https://example.com)"
            className="flex-1 bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-hlfa-primary"
            required
          />
          <button
            type="submit"
            disabled={isScanning}
            className="bg-hlfa-primary hover:bg-blue-600 disabled:bg-gray-700 text-white font-semibold px-6 py-2 rounded-lg transition-colors"
          >
            {isScanning ? 'Scanning...' : 'Start Scan'}
          </button>
        </form>
      </div>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <div
              key={stat.label}
              className="bg-gray-800 rounded-lg border border-gray-700 p-6 flex items-center space-x-4"
            >
              <Icon size={32} className={stat.color} />
              <div>
                <p className="text-gray-400 text-sm">{stat.label}</p>
                <p className="text-2xl font-bold">{stat.value}</p>
              </div>
            </div>
          );
        })}
      </div>

      {/* Recent Scans */}
      <div className="bg-gray-800 rounded-lg border border-gray-700 p-6">
        <h2 className="text-xl font-bold mb-4">Recent Scans</h2>
        <div className="space-y-2">
          {scans.length === 0 ? (
            <p className="text-gray-400">No scans yet</p>
          ) : (
            scans.slice(0, 5).map((scan) => (
              <div key={scan.id} className="flex items-center justify-between p-3 bg-gray-900 rounded">
                <div>
                  <p className="font-semibold">{scan.target}</p>
                  <p className="text-sm text-gray-400">{scan.status}</p>
                </div>
                <div className="w-24 bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-hlfa-primary h-2 rounded-full transition-all"
                    style={{ width: `${scan.progress}%` }}
                  ></div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
};

export default ScanPage;
