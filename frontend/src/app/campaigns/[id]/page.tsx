'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { campaignApi, sessionApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import {
  Loader2,
  ArrowLeft,
  Users,
  Calendar,
  Settings,
  Play,
  Copy,
  Check,
  Plus,
} from 'lucide-react';

interface Campaign {
  id: string;
  name: string;
  description: string;
  setting: string;
  dm_id: string;
  dm_username: string;
  max_players: number;
  homebrew_rules?: string;
  invite_code: string;
  created_at: string;
  players: Array<{
    user_id: string;
    username: string;
    character_id?: string;
    character_name?: string;
    joined_at: string;
  }>;
  status: 'active' | 'paused' | 'completed';
}

interface Session {
  id: string;
  name: string;
  status: string;
  scheduled_for?: string;
  created_at: string;
}

export default function CampaignDetailPage() {
  const router = useRouter();
  const params = useParams();
  const user = useAuthStore((state) => state.user);
  const [campaign, setCampaign] = useState<Campaign | null>(null);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  const campaignId = params.id as string;
  const isDM = user?.id === campaign?.dm_id;
  const isPlayer = campaign?.players.some((p) => p.user_id === user?.id);

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }
    loadCampaign();
  }, [user, campaignId, router]);

  const loadCampaign = async () => {
    try {
      setLoading(true);
      const [campaignRes, sessionsRes] = await Promise.all([
        campaignApi.get(campaignId),
        campaignApi.getSessions(campaignId),
      ]);
      setCampaign(campaignRes.data.campaign);
      setSessions(sessionsRes.data.sessions || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load campaign');
    } finally {
      setLoading(false);
    }
  };

  const handleJoinCampaign = async () => {
    try {
      await campaignApi.join(campaignId);
      await loadCampaign();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to join campaign');
    }
  };

  const handleLeaveCampaign = async () => {
    if (!confirm('Are you sure you want to leave this campaign?')) return;
    try {
      await campaignApi.leave(campaignId);
      router.push('/campaigns');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to leave campaign');
    }
  };

  const copyInviteCode = () => {
    if (campaign?.invite_code) {
      navigator.clipboard.writeText(campaign.invite_code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  if (!user || loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-purple-500" />
      </div>
    );
  }

  if (error || !campaign) {
    return (
      <div className="min-h-screen bg-gray-900 p-8">
        <div className="max-w-3xl mx-auto">
          <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
            {error || 'Campaign not found'}
          </div>
          <Link
            href="/campaigns"
            className="inline-flex items-center gap-2 mt-4 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to campaigns
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <Link
            href="/campaigns"
            className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to campaigns
          </Link>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-8">
            <div className="bg-gray-800 rounded-lg p-6">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h1 className="text-3xl font-bold text-white mb-2">{campaign.name}</h1>
                  <p className="text-gray-400">{campaign.description}</p>
                </div>
                {isDM && (
                  <Link
                    href={`/campaigns/${campaignId}/settings`}
                    className="p-2 text-gray-400 hover:text-white transition-colors"
                  >
                    <Settings className="h-5 w-5" />
                  </Link>
                )}
              </div>

              <div className="space-y-3 mb-6">
                <div className="flex items-center gap-2 text-sm text-gray-400">
                  <Calendar className="h-4 w-4" />
                  <span>Setting: {campaign.setting}</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-400">
                  <Users className="h-4 w-4" />
                  <span>
                    {campaign.players.length} / {campaign.max_players} players
                  </span>
                </div>
              </div>

              {campaign.homebrew_rules && (
                <div className="border-t border-gray-700 pt-4">
                  <h3 className="text-sm font-medium text-gray-300 mb-2">Homebrew Rules</h3>
                  <p className="text-gray-400 text-sm">{campaign.homebrew_rules}</p>
                </div>
              )}
            </div>

            <div className="bg-gray-800 rounded-lg p-6">
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-xl font-semibold text-white">Sessions</h2>
                {isDM && (
                  <Link
                    href={`/campaigns/${campaignId}/sessions/create`}
                    className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors text-sm"
                  >
                    <Plus className="h-4 w-4" />
                    New Session
                  </Link>
                )}
              </div>

              {sessions.length === 0 ? (
                <p className="text-gray-400 text-center py-8">No sessions yet</p>
              ) : (
                <div className="space-y-4">
                  {sessions.map((session) => (
                    <div
                      key={session.id}
                      className="border border-gray-700 rounded-lg p-4 hover:bg-gray-750 transition-colors"
                    >
                      <div className="flex justify-between items-start">
                        <div>
                          <h3 className="font-medium text-white mb-1">{session.name}</h3>
                          <p className="text-sm text-gray-400">
                            {new Date(session.created_at).toLocaleDateString()}
                          </p>
                        </div>
                        <Link
                          href={`/sessions/${session.id}`}
                          className="flex items-center gap-2 px-3 py-1 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors text-sm"
                        >
                          <Play className="h-4 w-4" />
                          {session.status === 'active' ? 'Continue' : 'View'}
                        </Link>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <div className="space-y-8">
            <div className="bg-gray-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Dungeon Master</h2>
              <p className="text-gray-300">{campaign.dm_username}</p>
            </div>

            <div className="bg-gray-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Players</h2>
              {campaign.players.length === 0 ? (
                <p className="text-gray-400">No players yet</p>
              ) : (
                <div className="space-y-3">
                  {campaign.players.map((player) => (
                    <div key={player.user_id} className="flex justify-between items-center">
                      <div>
                        <p className="text-gray-300">{player.username}</p>
                        {player.character_name && (
                          <p className="text-sm text-gray-500">{player.character_name}</p>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {isDM && (
              <div className="bg-gray-800 rounded-lg p-6">
                <h2 className="text-xl font-semibold text-white mb-4">Invite Code</h2>
                <div className="flex items-center gap-2">
                  <code className="flex-1 px-3 py-2 bg-gray-700 rounded text-gray-300 text-sm">
                    {campaign.invite_code}
                  </code>
                  <button
                    onClick={copyInviteCode}
                    className="p-2 text-gray-400 hover:text-white transition-colors"
                  >
                    {copied ? <Check className="h-5 w-5" /> : <Copy className="h-5 w-5" />}
                  </button>
                </div>
              </div>
            )}

            {!isDM && !isPlayer && (
              <button
                onClick={handleJoinCampaign}
                className="w-full py-3 px-4 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors"
              >
                Join Campaign
              </button>
            )}

            {!isDM && isPlayer && (
              <button
                onClick={handleLeaveCampaign}
                className="w-full py-3 px-4 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
              >
                Leave Campaign
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}