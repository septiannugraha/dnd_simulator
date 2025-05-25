'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { characterApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { Plus, Sword, Shield, Heart, Loader2 } from 'lucide-react';

interface Character {
  id: string;
  user_id: string;
  name: string;
  race: string;
  class: string;
  level: number;
  hit_points: number;
  max_hit_points: number;
  armor_class: number;
  campaign_id?: string;
  campaign_name?: string;
  created_at: string;
}

export default function CharactersPage() {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }
    loadCharacters();
  }, [user, router]);

  const loadCharacters = async () => {
    try {
      setLoading(true);
      const response = await characterApi.list();
      setCharacters(response.data.characters || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load characters');
    } finally {
      setLoading(false);
    }
  };

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-white">My Characters</h1>
            <p className="mt-2 text-gray-400">Manage your heroes and adventurers</p>
          </div>
          <Link
            href="/characters/create"
            className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors"
          >
            <Plus className="h-5 w-5" />
            Create Character
          </Link>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <Loader2 className="h-8 w-8 animate-spin text-purple-500" />
          </div>
        ) : error ? (
          <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
            {error}
          </div>
        ) : characters.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-400 mb-4">No characters yet</p>
            <Link
              href="/characters/create"
              className="inline-flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors"
            >
              <Plus className="h-5 w-5" />
              Create your first character
            </Link>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {characters.map((character) => (
              <Link
                key={character.id}
                href={`/characters/${character.id}`}
                className="bg-gray-800 rounded-lg p-6 hover:bg-gray-750 transition-colors"
              >
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-xl font-semibold text-white">{character.name}</h3>
                    <p className="text-gray-400">
                      Level {character.level} {character.race} {character.class}
                    </p>
                  </div>
                </div>

                <div className="grid grid-cols-3 gap-4 mb-4">
                  <div className="text-center">
                    <div className="flex items-center justify-center mb-1">
                      <Heart className="h-5 w-5 text-red-500" />
                    </div>
                    <p className="text-sm text-gray-400">HP</p>
                    <p className="text-white font-semibold">
                      {character.hit_points}/{character.max_hit_points}
                    </p>
                  </div>
                  <div className="text-center">
                    <div className="flex items-center justify-center mb-1">
                      <Shield className="h-5 w-5 text-blue-500" />
                    </div>
                    <p className="text-sm text-gray-400">AC</p>
                    <p className="text-white font-semibold">{character.armor_class}</p>
                  </div>
                  <div className="text-center">
                    <div className="flex items-center justify-center mb-1">
                      <Sword className="h-5 w-5 text-yellow-500" />
                    </div>
                    <p className="text-sm text-gray-400">Level</p>
                    <p className="text-white font-semibold">{character.level}</p>
                  </div>
                </div>

                {character.campaign_name && (
                  <div className="pt-4 border-t border-gray-700">
                    <p className="text-sm text-gray-400">
                      Campaign: <span className="text-gray-300">{character.campaign_name}</span>
                    </p>
                  </div>
                )}
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}