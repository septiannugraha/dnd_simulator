'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { campaignApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { Plus, Users, Calendar, ChevronRight, Loader2 } from 'lucide-react';
import Link from 'next/link';

interface Campaign {
  id: string;
  name: string;
  description: string;
  dm_id: string;
  dm_username: string;
  created_at: string;
  player_count: number;
  max_players: number;
  status: 'active' | 'paused' | 'completed';
}

export default function CampaignsPage() {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }
    loadCampaigns();
  }, [user, router]);

  const loadCampaigns = async () => {
    try {
      setLoading(true);
      const response = await campaignApi.list();
      setCampaigns(response.data.campaigns || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load campaigns');
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
            <h1 className="text-3xl font-bold text-white">Campaigns</h1>
            <p className="mt-2 text-gray-400">Join an adventure or create your own</p>
          </div>
          <Link
            href="/campaigns/create"
            className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors"
          >
            <Plus className="h-5 w-5" />
            Create Campaign
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
        ) : campaigns.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-400 mb-4">No campaigns found</p>
            <Link
              href="/campaigns/create"
              className="inline-flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors"
            >
              <Plus className="h-5 w-5" />
              Create your first campaign
            </Link>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {campaigns.map((campaign) => (
              <div
                key={campaign.id}
                className="bg-gray-800 rounded-lg p-6 hover:bg-gray-750 transition-colors"
              >
                <div className="mb-4">
                  <h3 className="text-xl font-semibold text-white mb-2">{campaign.name}</h3>
                  <p className="text-gray-400 text-sm line-clamp-2">{campaign.description}</p>
                </div>

                <div className="space-y-2 mb-4">
                  <div className="flex items-center gap-2 text-sm text-gray-400">
                    <Users className="h-4 w-4" />
                    <span>
                      {campaign.player_count} / {campaign.max_players} players
                    </span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-gray-400">
                    <Calendar className="h-4 w-4" />
                    <span>DM: {campaign.dm_username}</span>
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      campaign.status === 'active'
                        ? 'bg-green-500/20 text-green-400'
                        : campaign.status === 'paused'
                        ? 'bg-yellow-500/20 text-yellow-400'
                        : 'bg-gray-500/20 text-gray-400'
                    }`}
                  >
                    {campaign.status}
                  </span>
                  <Link
                    href={`/campaigns/${campaign.id}`}
                    className="flex items-center gap-1 text-purple-400 hover:text-purple-300 transition-colors"
                  >
                    View
                    <ChevronRight className="h-4 w-4" />
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}