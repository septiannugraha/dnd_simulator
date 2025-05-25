'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { characterApi, campaignApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import {
  Loader2,
  ArrowLeft,
  Heart,
  Shield,
  Sword,
  Brain,
  Eye,
  Sparkles,
  Edit,
  Trash2,
} from 'lucide-react';

interface Character {
  id: string;
  user_id: string;
  name: string;
  race: string;
  class: string;
  level: number;
  background: string;
  alignment: string;
  hit_points: number;
  max_hit_points: number;
  armor_class: number;
  speed: number;
  initiative_modifier: number;
  proficiency_bonus: number;
  ability_scores: {
    strength: number;
    dexterity: number;
    constitution: number;
    intelligence: number;
    wisdom: number;
    charisma: number;
  };
  ability_modifiers: {
    strength: number;
    dexterity: number;
    constitution: number;
    intelligence: number;
    wisdom: number;
    charisma: number;
  };
  skills: Record<string, number>;
  saving_throws: Record<string, number>;
  campaign_id?: string;
  campaign_name?: string;
  created_at: string;
}

const ABILITY_ICONS = {
  strength: Sword,
  dexterity: Shield,
  constitution: Heart,
  intelligence: Brain,
  wisdom: Eye,
  charisma: Sparkles,
};

export default function CharacterDetailPage() {
  const router = useRouter();
  const params = useParams();
  const user = useAuthStore((state) => state.user);
  const [character, setCharacter] = useState<Character | null>(null);
  const [campaigns, setCampaigns] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAssignModal, setShowAssignModal] = useState(false);

  const characterId = params.id as string;
  const isOwner = user?.id === character?.user_id;

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }
    loadCharacter();
  }, [user, characterId, router]);

  const loadCharacter = async () => {
    try {
      setLoading(true);
      const response = await characterApi.get(characterId);
      setCharacter(response.data.character);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load character');
    } finally {
      setLoading(false);
    }
  };

  const loadCampaigns = async () => {
    try {
      const response = await campaignApi.list();
      setCampaigns(response.data.campaigns || []);
    } catch (err) {
      console.error('Failed to load campaigns:', err);
    }
  };

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this character?')) return;
    try {
      await characterApi.delete(characterId);
      router.push('/characters');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete character');
    }
  };

  const handleAssignToCampaign = async (campaignId: string) => {
    try {
      await characterApi.assignToCampaign(characterId, campaignId);
      setShowAssignModal(false);
      await loadCharacter();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to assign character');
    }
  };

  if (!user || loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-purple-500" />
      </div>
    );
  }

  if (error || !character) {
    return (
      <div className="min-h-screen bg-gray-900 p-8">
        <div className="max-w-3xl mx-auto">
          <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
            {error || 'Character not found'}
          </div>
          <Link
            href="/characters"
            className="inline-flex items-center gap-2 mt-4 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to characters
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8 flex justify-between items-center">
          <Link
            href="/characters"
            className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to characters
          </Link>
          {isOwner && (
            <div className="flex gap-2">
              <button
                onClick={() => router.push(`/characters/${characterId}/edit`)}
                className="p-2 text-gray-400 hover:text-white transition-colors"
              >
                <Edit className="h-5 w-5" />
              </button>
              <button
                onClick={handleDelete}
                className="p-2 text-red-400 hover:text-red-300 transition-colors"
              >
                <Trash2 className="h-5 w-5" />
              </button>
            </div>
          )}
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-8">
            <div className="bg-gray-800 rounded-lg p-6">
              <h1 className="text-3xl font-bold text-white mb-2">{character.name}</h1>
              <p className="text-gray-400 mb-4">
                Level {character.level} {character.race} {character.class}
              </p>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span className="text-gray-400">Background:</span>{' '}
                  <span className="text-white">{character.background}</span>
                </div>
                <div>
                  <span className="text-gray-400">Alignment:</span>{' '}
                  <span className="text-white">{character.alignment}</span>
                </div>
              </div>
            </div>

            <div className="bg-gray-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Ability Scores</h2>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {Object.entries(character.ability_scores).map(([ability, score]) => {
                  const Icon = ABILITY_ICONS[ability as keyof typeof ABILITY_ICONS];
                  const modifier = character.ability_modifiers[ability as keyof typeof character.ability_modifiers];
                  return (
                    <div key={ability} className="bg-gray-700 rounded-lg p-4 text-center">
                      <Icon className="h-6 w-6 mx-auto mb-2 text-purple-400" />
                      <p className="text-xs text-gray-400 uppercase mb-1">{ability.slice(0, 3)}</p>
                      <p className="text-2xl font-bold text-white">{score}</p>
                      <p className="text-sm text-gray-300">
                        {modifier >= 0 ? '+' : ''}{modifier}
                      </p>
                    </div>
                  );
                })}
              </div>
            </div>

            <div className="bg-gray-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Skills & Proficiencies</h2>
              <div className="grid grid-cols-2 gap-4">
                {Object.entries(character.skills).map(([skill, bonus]) => (
                  <div key={skill} className="flex justify-between">
                    <span className="text-gray-400 capitalize">{skill.replace(/_/g, ' ')}</span>
                    <span className="text-white">
                      {bonus >= 0 ? '+' : ''}{bonus}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="space-y-8">
            <div className="bg-gray-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Combat Stats</h2>
              <div className="space-y-4">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <Heart className="h-5 w-5 text-red-500" />
                    <span className="text-gray-400">Hit Points</span>
                  </div>
                  <p className="text-2xl font-bold text-white">
                    {character.hit_points} / {character.max_hit_points}
                  </p>
                </div>
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <Shield className="h-5 w-5 text-blue-500" />
                    <span className="text-gray-400">Armor Class</span>
                  </div>
                  <p className="text-2xl font-bold text-white">{character.armor_class}</p>
                </div>
                <div>
                  <p className="text-gray-400">Initiative</p>
                  <p className="text-lg text-white">
                    {character.initiative_modifier >= 0 ? '+' : ''}{character.initiative_modifier}
                  </p>
                </div>
                <div>
                  <p className="text-gray-400">Speed</p>
                  <p className="text-lg text-white">{character.speed} ft</p>
                </div>
                <div>
                  <p className="text-gray-400">Proficiency Bonus</p>
                  <p className="text-lg text-white">+{character.proficiency_bonus}</p>
                </div>
              </div>
            </div>

            <div className="bg-gray-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Campaign</h2>
              {character.campaign_name ? (
                <div>
                  <p className="text-gray-300">{character.campaign_name}</p>
                  <Link
                    href={`/campaigns/${character.campaign_id}`}
                    className="text-purple-400 hover:text-purple-300 text-sm mt-2 inline-block"
                  >
                    View Campaign →
                  </Link>
                </div>
              ) : (
                <div>
                  <p className="text-gray-400 mb-4">Not assigned to a campaign</p>
                  {isOwner && (
                    <button
                      onClick={() => {
                        loadCampaigns();
                        setShowAssignModal(true);
                      }}
                      className="w-full py-2 px-4 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors text-sm"
                    >
                      Assign to Campaign
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {showAssignModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <div className="bg-gray-800 rounded-lg p-6 max-w-md w-full">
            <h3 className="text-xl font-semibold text-white mb-4">Assign to Campaign</h3>
            <div className="space-y-2 max-h-60 overflow-y-auto">
              {campaigns.map((campaign) => (
                <button
                  key={campaign.id}
                  onClick={() => handleAssignToCampaign(campaign.id)}
                  className="w-full text-left p-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
                >
                  <p className="text-white font-medium">{campaign.name}</p>
                  <p className="text-sm text-gray-400">DM: {campaign.dm_username}</p>
                </button>
              ))}
            </div>
            <button
              onClick={() => setShowAssignModal(false)}
              className="mt-4 w-full py-2 px-4 border border-gray-600 rounded-lg text-gray-300 hover:bg-gray-700 transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}