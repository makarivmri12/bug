import React from 'react';
import { Settings as SettingsIcon, Save } from 'lucide-react';

const SettingsPage: React.FC = () => {
  const [settings, setSettings] = React.useState({
    autoScan: true,
    notifications: true,
    darkMode: true,
    apiEndpoint: 'http://localhost:8080',
  });

  const handleSave = () => {
    console.log('Settings saved:', settings);
  };

  return (
    <div className="space-y-6 max-w-2xl">
      <h1 className="text-3xl font-bold flex items-center space-x-2">
        <SettingsIcon size={32} />
        <span>Settings</span>
      </h1>

      <div className="bg-gray-800 rounded-lg border border-gray-700 p-6 space-y-6">
        {/* Auto Scan */}
        <div className="flex items-center justify-between">
          <div>
            <p className="font-semibold">Auto Scan</p>
            <p className="text-sm text-gray-400">Automatically run scans on new targets</p>
          </div>
          <input
            type="checkbox"
            checked={settings.autoScan}
            onChange={(e) => setSettings({ ...settings, autoScan: e.target.checked })}
            className="w-5 h-5 rounded"
          />
        </div>

        {/* Notifications */}
        <div className="border-t border-gray-700 pt-6 flex items-center justify-between">
          <div>
            <p className="font-semibold">Notifications</p>
            <p className="text-sm text-gray-400">Receive alerts on new findings</p>
          </div>
          <input
            type="checkbox"
            checked={settings.notifications}
            onChange={(e) => setSettings({ ...settings, notifications: e.target.checked })}
            className="w-5 h-5 rounded"
          />
        </div>

        {/* API Endpoint */}
        <div className="border-t border-gray-700 pt-6">
          <p className="font-semibold mb-2">API Endpoint</p>
          <input
            type="text"
            value={settings.apiEndpoint}
            onChange={(e) => setSettings({ ...settings, apiEndpoint: e.target.value })}
            className="w-full bg-gray-900 border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-hlfa-primary"
          />
        </div>

        {/* Save Button */}
        <div className="border-t border-gray-700 pt-6 flex justify-end">
          <button
            onClick={handleSave}
            className="bg-hlfa-primary hover:bg-blue-600 text-white font-semibold px-6 py-2 rounded-lg flex items-center space-x-2 transition-colors"
          >
            <Save size={20} />
            <span>Save Changes</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;
