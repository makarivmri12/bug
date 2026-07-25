import { useEffect, useRef } from 'react';
import io, { Socket } from 'socket.io-client';
import { useScanStore } from '../stores/scanStore';

export const useWebSocket = () => {
  const socketRef = useRef<Socket | null>(null);
  const { addFinding, updateScanProgress } = useScanStore();

  useEffect(() => {
    socketRef.current = io(window.location.origin, {
      path: '/socket.io',
      reconnection: true,
      reconnectionDelay: 1000,
      reconnectionDelayMax: 5000,
      reconnectionAttempts: 5,
    });

    socketRef.current.on('finding:discovered', (data) => {
      addFinding(data);
    });

    socketRef.current.on('scan:progress', (data) => {
      updateScanProgress(data.scanId, data.progress);
    });

    socketRef.current.on('scan:completed', (data) => {
      console.log('Scan completed:', data);
    });

    return () => {
      if (socketRef.current) {
        socketRef.current.disconnect();
      }
    };
  }, [addFinding, updateScanProgress]);

  return socketRef.current;
};
