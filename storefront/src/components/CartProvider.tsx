'use client';

import { ReactNode } from 'react';

interface CartProviderProps {
  children: ReactNode;
}

export function CartProvider({ children }: CartProviderProps) {
  return <>{children}</>;
}
