import React from 'react';
import { AlertTriangle, CheckCircle, Clock, AlertCircle } from 'lucide-react';

interface FindingCardProps {
  finding: any;
  onClick: (finding: any) => void;
}

const FindingCard: React.FC<FindingCardProps> = ({ finding, onClick }) => {
  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'CRITICAL':
        return 'bg-red-900 text-red-300';
      case 'HIGH':
        return 'bg-orange-900 text-orange-300';
      case 'MEDIUM':
        return 'bg-yellow-900 text-yellow-300';
      default:
        return 'bg-blue-900 text-blue-300';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'CONFIRMED':
        return <CheckCircle size={16} className="text-green-500" />;
      case 'OPEN':
        return <AlertTriangle size={16} className="text-yellow-500" />;
      case 'FALSE_POSITIVE':
        return <Clock size={16} className="text-gray-500" />;
      default:
        return <AlertCircle size={16} className="text-blue-500" />;
    }
  };

  return (
    <div
      onClick={() => onClick(finding)}
      className="bg-gray-800 border border-gray-700 rounded-lg p-4 cursor-pointer hover:border-hlfa-primary transition-colors"
    >
      <div className="flex items-start justify-between mb-2">
        <h3 className="font-semibold text-white flex-1">{finding.title}</h3>
        <span className={`px-2 py-1 rounded text-xs font-semibold ${getSeverityColor(finding.severity)}`}>
          {finding.severity}
        </span>
      </div>

      <p className="text-gray-400 text-sm mb-3">{finding.description}</p>

      <div className="flex items-center justify-between text-xs text-gray-500">
        <span className="px-2 py-1 bg-gray-900 rounded">{finding.category}</span>
        <div className="flex items-center space-x-1">{getStatusIcon(finding.status)}</div>
      </div>
    </div>
  );
};

export default FindingCard;
