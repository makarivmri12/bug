import React from 'react';
import ReactFlow, { Background, Controls } from 'reactflow';
import 'reactflow/dist/style.css';

const FlowGraphPage: React.FC = () => {
  const [nodes, setNodes] = React.useState([
    {
      id: '1',
      data: { label: 'Login' },
      position: { x: 0, y: 0 },
      style: {
        background: '#3b82f6',
        color: '#fff',
        border: '1px solid #1e293b',
        borderRadius: '8px',
        padding: '10px',
        minWidth: '100px',
      },
    },
    {
      id: '2',
      data: { label: 'Dashboard' },
      position: { x: 250, y: 0 },
      style: {
        background: '#3b82f6',
        color: '#fff',
        border: '1px solid #1e293b',
        borderRadius: '8px',
        padding: '10px',
        minWidth: '100px',
      },
    },
  ]);

  const [edges] = React.useState([
    { id: 'e1-2', source: '1', target: '2', animated: true },
  ]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold mb-2">Application Flow</h1>
        <p className="text-gray-400">Visualize the application flow and endpoint relationships</p>
      </div>

      <div className="bg-gray-800 rounded-lg border border-gray-700 overflow-hidden" style={{ height: '600px' }}>
        <ReactFlow nodes={nodes} edges={edges} fitView>
          <Background color="#1e293b" gap={16} />
          <Controls />
        </ReactFlow>
      </div>
    </div>
  );
};

export default FlowGraphPage;
