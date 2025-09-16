import React from 'react';
import { Stack } from 'expo-router';
import { useRouter } from 'expo-router';
import { Text, TouchableOpacity, View } from 'react-native';

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
            <View style={{ flexDirection: 'row' }}>
              <TouchableOpacity 
                onPress={() => router.push('/pos/reports')}
                style={{ padding: 10, marginRight: 10 }}
              >
                <Text style={{ color: '#fff', fontWeight: 'bold' }}>Reports</Text>
              </TouchableOpacity>
              <TouchableOpacity 
                onPress={() => router.push('/pos/checkout')}
                style={{ padding: 10 }}
              >
                <Text style={{ color: '#fff', fontWeight: 'bold' }}>Checkout</Text>
              </TouchableOpacity>
            </View>
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
        }} 
      />
      <Stack.Screen 
        name="reports" 
        options={{ 
          title: 'Reports & Analytics',
          headerShown: true,
          headerStyle: {
            backgroundColor: '#2e7d32',
          },
          headerTintColor: '#fff',
          headerTitleStyle: {
            fontWeight: 'bold',
          },
        }} 
      />
    </Stack>
  );
}