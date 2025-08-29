import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { CartProvider } from "@/components/CartProvider";
import { Navigation } from "@/components/Navigation";

const inter = Inter({
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Unified Commerce - Your Premium Shopping Destination",
  description: "Discover amazing products with our modern e-commerce platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} antialiased`}>
        <div className="min-h-screen bg-gray-50">
          <Navigation />
          <main className="pb-16">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
