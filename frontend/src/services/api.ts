import axios from 'axios';

const API_BASE_URL = '/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const scanAPI = {
  createScan: (data: any) => apiClient.post('/scans', data),
  getScan: (scanId: string) => apiClient.get(`/scans/${scanId}`),
  listScans: (workspaceId: string) => apiClient.get(`/scans?workspace=${workspaceId}`),
  getFindings: (scanId: string) => apiClient.get(`/scans/${scanId}/findings`),
};

export const findingAPI = {
  getFinding: (findingId: string) => apiClient.get(`/findings/${findingId}`),
  updateFinding: (findingId: string, data: any) => apiClient.put(`/findings/${findingId}`, data),
  confirmFinding: (findingId: string) => apiClient.put(`/findings/${findingId}`, { status: 'CONFIRMED' }),
};

export const heuristicAPI = {
  setBaseline: (data: any) => axios.post('http://localhost:5000/api/v1/differential/baseline', data),
  analyze: (data: any) => axios.post('http://localhost:5000/api/v1/differential/analyze', data),
  profileEndpoint: (data: any) => axios.post('http://localhost:5000/api/v1/profile/endpoint', data),
};

export default apiClient;
