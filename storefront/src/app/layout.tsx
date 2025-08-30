import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Navigation } from "@/components/Navigation";
import { CartProvider } from "@/components/CartProvider";

const inter = Inter({
  subsets: ["latin"],
  display: 'swap',
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
        <CartProvider>
          <div className="min-h-screen bg-gray-50">
            <Navigation />
            <main className="pb-16">
              {children}
            </main>
            {/* Footer */}
            <footer className="bg-gray-900 text-gray-300">
              <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                  {/* Company Info */}
                  <div>
                    <h3 className="text-white text-lg font-bold mb-4">Unified Commerce</h3>
                    <p className="mb-4 text-sm">Your premium shopping destination for the latest products and innovations.</p>
                    <div className="flex space-x-4">
                      {/* Social Media Icons */}
                      <a href="#" className="text-gray-400 hover:text-white transition">
                        <svg className="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z" />
                        </svg>
                      </a>
                      <a href="#" className="text-gray-400 hover:text-white transition">
                        <svg className="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path d="M12 2C6.477 2 2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12c0-5.523-4.477-10-10-10z" />
                        </svg>
                      </a>
                      <a href="#" className="text-gray-400 hover:text-white transition">
                        <svg className="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.5 14c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zm-9 0c-.83 0-1.5-.67-1.5-1.5S6.67 13 7.5 13s1.5.67 1.5 1.5S8.33 16 7.5 16zm4.5-9C9.02 7 7 9.02 7 11.5h2c0-1.38 1.12-2.5 2.5-2.5s2.5 1.12 2.5 2.5h2c0-2.48-2.02-4.5-4.5-4.5z" />
                        </svg>
                      </a>
                    </div>
                  </div>
                  
                  {/* Quick Links */}
                  <div>
                    <h3 className="text-white text-lg font-bold mb-4">Quick Links</h3>
                    <ul className="space-y-2">
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Home</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Shop</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Categories</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">About Us</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Contact</a></li>
                    </ul>
                  </div>
                  
                  {/* Customer Service */}
                  <div>
                    <h3 className="text-white text-lg font-bold mb-4">Customer Service</h3>
                    <ul className="space-y-2">
                      <li><a href="#" className="text-gray-400 hover:text-white transition">My Account</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Order History</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Shipping Policy</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">Returns & Exchanges</a></li>
                      <li><a href="#" className="text-gray-400 hover:text-white transition">FAQ</a></li>
                    </ul>
                  </div>
                  
                  {/* Contact Info */}
                  <div>
                    <h3 className="text-white text-lg font-bold mb-4">Contact Us</h3>
                    <ul className="space-y-2">
                      <li className="flex items-start">
                        <svg className="h-5 w-5 mr-2 mt-0.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                        </svg>
                        <span>123 Commerce St, Business City, 10001</span>
                      </li>
                      <li className="flex items-start">
                        <svg className="h-5 w-5 mr-2 mt-0.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                        </svg>
                        <span>support@unifiedcommerce.com</span>
                      </li>
                      <li className="flex items-start">
                        <svg className="h-5 w-5 mr-2 mt-0.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
                        </svg>
                        <span>+1 (555) 123-4567</span>
                      </li>
                    </ul>
                  </div>
                </div>
                
                <div className="mt-12 border-t border-gray-800 pt-8 flex flex-col md:flex-row justify-between items-center">
                  <p className="text-sm">© 2025 Unified Commerce. All rights reserved.</p>
                  <div className="mt-4 md:mt-0 flex space-x-6">
                    <a href="#" className="text-gray-400 hover:text-white transition text-sm">Privacy Policy</a>
                    <a href="#" className="text-gray-400 hover:text-white transition text-sm">Terms of Service</a>
                    <a href="#" className="text-gray-400 hover:text-white transition text-sm">Cookie Policy</a>
                  </div>
                </div>
              </div>
            </footer>
          </div>
        </CartProvider>
      </body>
    </html>
  );
}
