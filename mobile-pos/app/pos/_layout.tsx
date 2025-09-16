import React from 'react';
import { Stack, Text, TouchableOpacity } from 'expo-router';
import { useRouter } from 'expo-router';

export default function POSLayout() {
  const router = useRouter();
  
  return (
    <Stack>
      <Stack.Screen 
        name="index" 
        options={{ 
          title: 'Point of Sale',
          headerShown: true,
          headerStyle: {
            backgroundColor: '#2e7d32',
          },
          headerTintColor: '#fff',
          headerTitleStyle: {
            fontWeight: 'bold',
          },
          headerRight: () => (
            <TouchableOpacity 
              onPress={() => router.push('/pos/checkout')}
              style={{ padding: 10 }}
            >
              <Text style={{ color: '#fff', fontWeight: 'bold' }}>Checkout</Text>
            </TouchableOpacity>
          ),
        }} 
      />
      <Stack.Screen 
        name="checkout" 
        options={{ 
          title: 'Checkout',
          headerShown: true,
          headerStyle: {
            backgroundColor: '#2e7d32',
          },
          headerTintColor: '#fff',
          headerTitleStyle: {
            fontWeight: 'bold',
          },
          headerBackTitleVisible: false,
        }} 
      />
    </Stack>
  );
}