'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/auth-store';
import { Dice6, Loader2 } from 'lucide-react';

export default function HomePage() {
  const router = useRouter();
  const { user, isLoading, checkAuth } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    if (!isLoading) {
      if (user) {
        router.push('/campaigns');
      } else {
        router.push('/login');
      }
    }
  }, [user, isLoading, router]);

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center">
      <div className="text-center">
        <div className="flex justify-center mb-4">
          <Dice6 className="h-16 w-16 text-purple-500 animate-pulse" />
        </div>
        <h1 className="text-4xl font-bold text-white mb-2">D&D Simulator</h1>
        <div className="flex items-center justify-center text-gray-400">
          <Loader2 className="animate-spin mr-2 h-4 w-4" />
          Loading your adventure...
        </div>
      </div>
    </div>
  );
}