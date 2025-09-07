import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { User } from '@/graphql/queries';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
  setUser: (user: User) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: (user, token) => {
        // Store the token in localStorage
        if (typeof window !== 'undefined') {
          localStorage.setItem('auth-token', token);
        }
        set({ user, isAuthenticated: true });
      },
      logout: () => {
        // Remove the token from localStorage
        if (typeof window !== 'undefined') {
          localStorage.removeItem('auth-token');
        }
        set({ user: null, isAuthenticated: false });
      },
      setUser: (user) => set({ user }),
    }),
    {
      name: 'auth-storage', // name of the item in the storage (must be unique)
      partialize: (state) => ({ 
        user: state.user, 
        isAuthenticated: state.isAuthenticated 
      }),
    }
  )
);