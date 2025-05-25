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

// Safe localStorage helper
const safeLocalStorage = {
  getItem: (key: string) => {
    if (typeof window === 'undefined') return null;
    try {
      return localStorage.getItem(key);
    } catch {
      return null;
    }
  },
  setItem: (key: string, value: string) => {
    if (typeof window === 'undefined') return;
    try {
      localStorage.setItem(key, value);
    } catch {
      // Ignore localStorage errors
    }
  },
  removeItem: (key: string) => {
    if (typeof window === 'undefined') return;
    try {
      localStorage.removeItem(key);
    } catch {
      // Ignore localStorage errors
    }
  },
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isLoading: true,

  login: async (username, password) => {
    const response = await authApi.login({ username, password });
    const { token, user } = response.data;
    
    safeLocalStorage.setItem('token', token);
    safeLocalStorage.setItem('user', JSON.stringify(user));
    
    set({ user, token, isLoading: false });
  },

  register: async (username, email, password) => {
    const response = await authApi.register({ username, email, password });
    const { token, user } = response.data;
    
    safeLocalStorage.setItem('token', token);
    safeLocalStorage.setItem('user', JSON.stringify(user));
    
    set({ user, token, isLoading: false });
  },

  logout: () => {
    safeLocalStorage.removeItem('token');
    safeLocalStorage.removeItem('user');
    set({ user: null, token: null });
  },

  checkAuth: async () => {
    const token = safeLocalStorage.getItem('token');
    const userStr = safeLocalStorage.getItem('user');
    
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
      safeLocalStorage.removeItem('token');
      safeLocalStorage.removeItem('user');
      set({ user: null, token: null, isLoading: false });
    }
  },
}));