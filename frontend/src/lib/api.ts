import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Create axios instance with default config
export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear token and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authApi = {
  register: (data: { username: string; email: string; password: string }) =>
    api.post('/auth/register', data),
  
  login: (data: { username: string; password: string }) =>
    api.post('/auth/login', data),
  
  me: () => api.get('/auth/me'),
};

// Campaign API
export const campaignApi = {
  create: (data: any) => api.post('/campaigns', data),
  list: () => api.get('/campaigns'),
  get: (id: string) => api.get(`/campaigns/${id}`),
  update: (id: string, data: any) => api.put(`/campaigns/${id}`, data),
  delete: (id: string) => api.delete(`/campaigns/${id}`),
  join: (id: string) => api.post(`/campaigns/${id}/join`),
  leave: (id: string) => api.post(`/campaigns/${id}/leave`),
  getSessions: (id: string) => api.get(`/campaigns/${id}/sessions`),
};

// Character API
export const characterApi = {
  create: (data: any) => api.post('/characters', data),
  list: () => api.get('/characters'),
  get: (id: string) => api.get(`/characters/${id}`),
  update: (id: string, data: any) => api.put(`/characters/${id}`, data),
  delete: (id: string) => api.delete(`/characters/${id}`),
  assignToCampaign: (id: string, campaignId: string) =>
    api.post(`/characters/${id}/assign`, { campaign_id: campaignId }),
};

// Session API
export const sessionApi = {
  create: (data: any) => api.post('/sessions', data),
  get: (id: string) => api.get(`/sessions/${id}`),
  join: (id: string, characterId: string) =>
    api.post(`/sessions/${id}/join`, { character_id: characterId }),
  leave: (id: string) => api.post(`/sessions/${id}/leave`),
  start: (id: string) => api.post(`/sessions/${id}/start`),
  pause: (id: string) => api.post(`/sessions/${id}/pause`),
  resume: (id: string) => api.post(`/sessions/${id}/resume`),
  end: (id: string, status: 'completed' | 'cancelled') =>
    api.post(`/sessions/${id}/end`, { status }),
  updateScene: (id: string, scene: string, notes?: string) =>
    api.put(`/sessions/${id}/scene`, { scene, notes }),
  getStatus: (id: string) => api.get(`/sessions/${id}/status`),
  setInitiative: (id: string, characterId: string, initiative: number) =>
    api.post(`/sessions/${id}/initiative`, { character_id: characterId, initiative }),
  advanceTurn: (id: string, force?: boolean) =>
    api.post(`/sessions/${id}/turn/advance`, { force }),
};

// D&D Data API
export const dndApi = {
  getRaces: () => api.get('/dnd/races'),
  getClasses: () => api.get('/dnd/classes'),
  getBackgrounds: () => api.get('/dnd/backgrounds'),
};

// AI DM API
export const aiApi = {
  sendAction: (sessionId: string, data: {
    character_id: string;
    action: string;
    action_type: 'combat' | 'roleplay' | 'exploration';
    target?: string;
  }) => api.post(`/sessions/${sessionId}/action`, data),
  
  getNarrative: (sessionId: string) => api.get(`/sessions/${sessionId}/narrative`),
  getEvents: (sessionId: string, type: string) => api.get(`/sessions/${sessionId}/events/${type}`),
};

export default api;