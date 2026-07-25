import React, { useState } from 'react';
import { X, Copy, Check } from 'lucide-react';

interface FindingDetailProps {
  finding: any;
  onClose: () => void;
  onStatusChange: (status: string) => void;
}

const FindingDetail: React.FC<FindingDetailProps> = ({ finding, onClose, onStatusChange }) => {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (!finding) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg max-w-2xl w-full mx-4 max-h-96 overflow-y-auto">
        <div className="flex items-center justify-between p-6 border-b border-gray-700 sticky top-0 bg-gray-800">
          <h2 className="text-xl font-bold">{finding.title}</h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-700 rounded transition-colors"
          >
            <X size={20} />
          </button>
        </div>

        <div className="p-6 space-y-4">
          <div>
            <label className="text-sm text-gray-400 block mb-1">Description</label>
            <p className="text-white">{finding.description}</p>
          </div>

          <div>
            <label className="text-sm text-gray-400 block mb-1">Evidence</label>
            <pre className="bg-gray-900 p-3 rounded text-xs overflow-x-auto text-gray-300">
              {JSON.stringify(finding.evidence, null, 2)}
            </pre>
          </div>

          <div>
            <label className="text-sm text-gray-400 block mb-1">PoC</label>
            <div className="flex items-center space-x-2">
              <code className="bg-gray-900 p-2 rounded flex-1 text-xs text-gray-300 overflow-x-auto">
                {finding.poc_curl}
              </code>
              <button
                onClick={() => copyToClipboard(finding.poc_curl)}
                className="p-2 hover:bg-gray-700 rounded transition-colors"
              >
                {copied ? <Check size={16} /> : <Copy size={16} />}
              </button>
            </div>
          </div>

          <div>
            <label className="text-sm text-gray-400 block mb-2">Status</label>
            <div className="flex space-x-2">
              {['OPEN', 'CONFIRMED', 'FALSE_POSITIVE', 'FIXED'].map((status) => (
                <button
                  key={status}
                  onClick={() => onStatusChange(status)}
                  className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                    finding.status === status
                      ? 'bg-hlfa-primary text-white'
                      : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  }`}
                >
                  {status}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default FindingDetail;
