import { create } from 'zustand';
import { authApi } from './api';

interface User {
  id: string;
  username: string;
  email: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isLoading: true,

  login: async (username, password) => {
    const response = await authApi.login({ username, password });
    const { token, user } = response.data;
    
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    
    set({ user, token, isLoading: false });
  },

  register: async (username, email, password) => {
    const response = await authApi.register({ username, email, password });
    const { token, user } = response.data;
    
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    
    set({ user, token, isLoading: false });
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    set({ user: null, token: null });
  },

  checkAuth: async () => {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    
    if (!token || !userStr) {
      set({ isLoading: false });
      return;
    }

    try {
      const user = JSON.parse(userStr);
      set({ user, token, isLoading: false });
      
      // Verify token is still valid
      await authApi.me();
    } catch (error) {
      // Token is invalid
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      set({ user: null, token: null, isLoading: false });
    }
  },
}));