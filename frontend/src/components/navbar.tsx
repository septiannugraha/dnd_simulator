'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/auth-store';
import { Dice6, User, LogOut, Menu, X } from 'lucide-react';
import { useState } from 'react';

export default function Navbar() {
  const router = useRouter();
  const { user, logout } = useAuthStore();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  if (!user) return null;

  return (
    <nav className="bg-gray-800 border-b border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <Link href="/campaigns" className="flex items-center gap-2">
              <Dice6 className="h-8 w-8 text-purple-500" />
              <span className="text-xl font-bold text-white">D&D Simulator</span>
            </Link>
            <div className="hidden md:block ml-10">
              <div className="flex items-baseline space-x-4">
                <Link
                  href="/campaigns"
                  className="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
                >
                  Campaigns
                </Link>
                <Link
                  href="/characters"
                  className="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
                >
                  Characters
                </Link>
              </div>
            </div>
          </div>

          <div className="hidden md:block">
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2 text-gray-300">
                <User className="h-5 w-5" />
                <span className="text-sm">{user.username}</span>
              </div>
              <button
                onClick={handleLogout}
                className="flex items-center gap-2 text-gray-300 hover:text-white transition-colors"
              >
                <LogOut className="h-5 w-5" />
                <span className="text-sm">Logout</span>
              </button>
            </div>
          </div>

          <div className="md:hidden">
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="text-gray-300 hover:text-white"
            >
              {mobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
            </button>
          </div>
        </div>
      </div>

      {mobileMenuOpen && (
        <div className="md:hidden bg-gray-800 border-t border-gray-700">
          <div className="px-2 pt-2 pb-3 space-y-1">
            <Link
              href="/campaigns"
              className="text-gray-300 hover:text-white block px-3 py-2 rounded-md text-base font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              Campaigns
            </Link>
            <Link
              href="/characters"
              className="text-gray-300 hover:text-white block px-3 py-2 rounded-md text-base font-medium"
              onClick={() => setMobileMenuOpen(false)}
            >
              Characters
            </Link>
          </div>
          <div className="pt-4 pb-3 border-t border-gray-700">
            <div className="flex items-center px-5 mb-3">
              <User className="h-5 w-5 text-gray-400" />
              <span className="ml-2 text-gray-300">{user.username}</span>
            </div>
            <button
              onClick={handleLogout}
              className="text-gray-300 hover:text-white block px-5 py-2 text-base font-medium"
            >
              Logout
            </button>
          </div>
        </div>
      )}
    </nav>
  );
}