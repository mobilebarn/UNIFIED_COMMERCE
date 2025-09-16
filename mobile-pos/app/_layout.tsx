import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { useFonts } from 'expo-font';
import React from 'react';
import { Stack } from 'expo-router';
import ApolloProviderWrapper from '../components/ApolloProviderWrapper';

export default function RootLayout() {
  return (
    <ApolloProviderWrapper>
      <Stack>
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
        <Stack.Screen name="auth/login" options={{ headerShown: false }} />
        <Stack.Screen name="+not-found" />
      </Stack>
    </ApolloProviderWrapper>
  );
}
