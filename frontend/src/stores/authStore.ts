import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

interface AuthStore {
  user: any | null;
  workspace: any | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  setWorkspace: (workspace: any) => void;
}

export const useAuthStore = create<AuthStore>(
  devtools(
    persist(
      (set) => ({
        user: null,
        workspace: null,
        isAuthenticated: false,
        login: async (email: string, password: string) => {
          // API call would go here
          set({ isAuthenticated: true });
        },
        logout: () => {
          set({ user: null, isAuthenticated: false });
        },
        setWorkspace: (workspace: any) => {
          set({ workspace });
        },
      }),
      { name: 'auth-store' },
    ),
  ),
);
