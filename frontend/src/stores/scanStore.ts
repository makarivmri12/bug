import { create } from 'zustand';

interface ScanStore {
  scans: any[];
  currentScan: any | null;
  isScanning: boolean;
  findings: any[];
  addScan: (scan: any) => void;
  setCurrentScan: (scan: any) => void;
  setScanning: (isScanning: boolean) => void;
  addFinding: (finding: any) => void;
  updateScanProgress: (scanId: string, progress: number) => void;
}

export const useScanStore = create<ScanStore>((set) => ({
  scans: [],
  currentScan: null,
  isScanning: false,
  findings: [],
  addScan: (scan) =>
    set((state) => ({
      scans: [...state.scans, scan],
    })),
  setCurrentScan: (scan) =>
    set({
      currentScan: scan,
    }),
  setScanning: (isScanning) =>
    set({
      isScanning,
    }),
  addFinding: (finding) =>
    set((state) => ({
      findings: [...state.findings, finding],
    })),
  updateScanProgress: (scanId, progress) =>
    set((state) => ({
      scans: state.scans.map((s) =>
        s.id === scanId ? { ...s, progress } : s,
      ),
    })),
}));
