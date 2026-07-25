import React from 'react';
import { Bell, Menu, LogOut, User } from 'lucide-react';

interface TopBarProps {
  onMenuClick: () => void;
}

const TopBar: React.FC<TopBarProps> = ({ onMenuClick }) => {
  return (
    <header className="bg-gray-900 border-b border-gray-700 px-6 py-4 flex items-center justify-between">
      <button
        onClick={onMenuClick}
        className="p-2 hover:bg-gray-800 rounded-lg transition-colors"
      >
        <Menu size={24} className="text-gray-300" />
      </button>

      <div className="flex items-center space-x-4">
        <button className="p-2 hover:bg-gray-800 rounded-lg transition-colors relative">
          <Bell size={20} className="text-gray-300" />
          <span className="absolute top-1 right-1 w-2 h-2 bg-hlfa-danger rounded-full"></span>
        </button>

        <div className="flex items-center space-x-2 pl-4 border-l border-gray-700">
          <div className="w-8 h-8 bg-hlfa-primary rounded-full flex items-center justify-center">
            <User size={16} />
          </div>
          <span className="text-sm text-gray-300">User</span>
        </div>

        <button className="p-2 hover:bg-gray-800 rounded-lg transition-colors">
          <LogOut size={20} className="text-gray-300" />
        </button>
      </div>
    </header>
  );
};

export default TopBar;
