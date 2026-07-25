import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  Shield,
  Activity,
  Layers,
  Users,
  Settings,
  ChevronRight,
} from 'lucide-react';

interface SidebarProps {
  isOpen: boolean;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen }) => {
  const location = useLocation();

  const menuItems = [
    { label: 'Dashboard', path: '/', icon: Activity },
    { label: 'Findings', path: '/findings', icon: Shield },
    { label: 'Flow Graph', path: '/flow', icon: Layers },
    { label: 'Team', path: '/workspace', icon: Users },
    { label: 'Settings', path: '/settings', icon: Settings },
  ];

  return (
    <aside
      className={`${
        isOpen ? 'w-64' : 'w-20'
      } bg-gray-900 border-r border-gray-700 transition-all duration-300 flex flex-col`}
    >
      <div className="p-6 border-b border-gray-700 flex items-center justify-center">
        <Shield className="text-hlfa-primary" size={32} />
        {isOpen && <span className="ml-2 font-bold text-lg">HLFA</span>}
      </div>

      <nav className="flex-1 p-4 space-y-2">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname === item.path;

          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center px-4 py-3 rounded-lg transition-colors ${
                isActive
                  ? 'bg-hlfa-primary text-white'
                  : 'text-gray-300 hover:bg-gray-800'
              }`}
            >
              <Icon size={20} />
              {isOpen && <span className="ml-3">{item.label}</span>}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
};

export default Sidebar;
