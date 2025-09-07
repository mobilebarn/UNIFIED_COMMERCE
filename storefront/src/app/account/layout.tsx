'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function AccountLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  
  const navigation = [
    { name: 'Account Dashboard', href: '/account' },
    { name: 'Order History', href: '/account/orders' },
    { name: 'Saved Items', href: '/account/saved' },
    { name: 'Address Book', href: '/account/addresses' },
    { name: 'Payment Methods', href: '/account/payments' },
  ];

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-lg shadow-md overflow-hidden">
          <div className="grid grid-cols-1 lg:grid-cols-4">
            {/* Sidebar Navigation */}
            <div className="border-r border-gray-200 bg-gray-50">
              <div className="p-6">
                <h2 className="text-lg font-medium text-gray-900">Account</h2>
                <p className="text-gray-600 text-sm">Manage your account settings</p>
              </div>
              <nav className="border-t border-gray-200">
                <ul className="space-y-1 py-4">
                  {navigation.map((item) => (
                    <li key={item.name}>
                      <Link
                        href={item.href}
                        className={`block px-6 py-3 text-sm font-medium ${
                          pathname === item.href
                            ? 'bg-blue-50 text-blue-600 border-l-4 border-blue-600'
                            : 'text-gray-700 hover:bg-gray-100'
                        }`}
                      >
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </nav>
            </div>
            
            {/* Main Content */}
            <div className="lg:col-span-3">
              <div className="p-6">
                {children}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}