// Utility functions for local storage
// Note: For mobile apps, you should use SecureStore from expo-secure-store for sensitive data

export const localStorage = {
  getItem: (key: string): string | null => {
    if (typeof window !== 'undefined') {
      return window.localStorage.getItem(key);
    }
    return null;
  },
  
  setItem: (key: string, value: string): void => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(key, value);
    }
  },
  
  removeItem: (key: string): void => {
    if (typeof window !== 'undefined') {
      window.localStorage.removeItem(key);
    }
  },
  
  clear: (): void => {
    if (typeof window !== 'undefined') {
      window.localStorage.clear();
    }
  }
};

// For mobile apps, you would typically use SecureStore:
/*
import * as SecureStore from 'expo-secure-store';

export const secureStorage = {
  getItem: async (key: string): Promise<string | null> => {
    return await SecureStore.getItemAsync(key);
  },
  
  setItem: async (key: string, value: string): Promise<void> => {
    await SecureStore.setItemAsync(key, value);
  },
  
  removeItem: async (key: string): Promise<void> => {
    await SecureStore.deleteItemAsync(key);
  }
};
*/