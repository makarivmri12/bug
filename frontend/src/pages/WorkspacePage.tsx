import React from 'react';
import { Users, Plus, Trash2 } from 'lucide-react';

const WorkspacePage: React.FC = () => {
  const [members, setMembers] = React.useState([
    { id: 1, name: 'You', email: 'user@example.com', role: 'Admin' },
  ]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Team & Workspace</h1>
        <button className="bg-hlfa-primary hover:bg-blue-600 text-white font-semibold px-6 py-2 rounded-lg flex items-center space-x-2 transition-colors">
          <Plus size={20} />
          <span>Add Member</span>
        </button>
      </div>

      <div className="bg-gray-800 rounded-lg border border-gray-700 overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-700">
          <h2 className="text-lg font-bold flex items-center space-x-2">
            <Users size={20} />
            <span>Team Members</span>
          </h2>
        </div>
        <div className="divide-y divide-gray-700">
          {members.map((member) => (
            <div key={member.id} className="px-6 py-4 flex items-center justify-between hover:bg-gray-700 transition-colors">
              <div>
                <p className="font-semibold">{member.name}</p>
                <p className="text-sm text-gray-400">{member.email}</p>
              </div>
              <div className="flex items-center space-x-4">
                <span className="px-3 py-1 bg-gray-900 rounded text-sm font-medium">{member.role}</span>
                {member.id !== 1 && (
                  <button className="p-2 hover:bg-gray-600 rounded transition-colors text-red-500">
                    <Trash2 size={16} />
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WorkspacePage;
