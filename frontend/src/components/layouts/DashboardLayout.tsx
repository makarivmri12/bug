import React from 'react';
import { Outlet, Link } from 'react-router-dom';
import { Menu, Settings, LogOut } from 'lucide-react';
import Sidebar from '../Sidebar';
import TopBar from '../TopBar';

const DashboardLayout: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = React.useState(true);

  return (
    <div className="flex h-screen bg-hlfa-dark">
      <Sidebar isOpen={sidebarOpen} />
      <div className="flex-1 flex flex-col">
        <TopBar onMenuClick={() => setSidebarOpen(!sidebarOpen)} />
        <main className="flex-1 overflow-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default DashboardLayout;
