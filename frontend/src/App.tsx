import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import DashboardLayout from './components/layouts/DashboardLayout';
import ScanPage from './pages/ScanPage';
import FindingsPage from './pages/FindingsPage';
import FlowGraphPage from './pages/FlowGraphPage';
import WorkspacePage from './pages/WorkspacePage';
import SettingsPage from './pages/SettingsPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<DashboardLayout />}>
          <Route index element={<ScanPage />} />
          <Route path="findings" element={<FindingsPage />} />
          <Route path="flow" element={<FlowGraphPage />} />
          <Route path="workspace" element={<WorkspacePage />} />
          <Route path="settings" element={<SettingsPage />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
