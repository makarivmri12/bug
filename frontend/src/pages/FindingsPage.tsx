import React from 'react';
import FindingCard from '../components/FindingCard';
import FindingDetail from '../components/FindingDetail';
import { useScanStore } from '../stores/scanStore';

const FindingsPage: React.FC = () => {
  const { findings } = useScanStore();
  const [selectedFinding, setSelectedFinding] = React.useState(null);
  const [filterSeverity, setFilterSeverity] = React.useState('ALL');

  const filteredFindings = React.useMemo(() => {
    if (filterSeverity === 'ALL') return findings;
    return findings.filter((f) => f.severity === filterSeverity);
  }, [findings, filterSeverity]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Findings</h1>
        <div className="flex space-x-2">
          {['ALL', 'CRITICAL', 'HIGH', 'MEDIUM', 'LOW'].map((severity) => (
            <button
              key={severity}
              onClick={() => setFilterSeverity(severity)}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                filterSeverity === severity
                  ? 'bg-hlfa-primary text-white'
                  : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
              }`}
            >
              {severity}
            </button>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {filteredFindings.map((finding) => (
          <FindingCard
            key={finding.id}
            finding={finding}
            onClick={setSelectedFinding}
          />
        ))}
      </div>

      {filteredFindings.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-400 text-lg">No findings found</p>
        </div>
      )}

      {selectedFinding && (
        <FindingDetail
          finding={selectedFinding}
          onClose={() => setSelectedFinding(null)}
          onStatusChange={(status) => {
            console.log('Update status to:', status);
          }}
        />
      )}
    </div>
  );
};

export default FindingsPage;
