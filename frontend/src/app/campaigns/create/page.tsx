'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { campaignApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { Loader2, ArrowLeft } from 'lucide-react';

const createCampaignSchema = z.object({
  name: z.string().min(3, 'Campaign name must be at least 3 characters'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  setting: z.string().min(3, 'Setting must be at least 3 characters'),
  max_players: z.number().min(2).max(8),
  homebrew_rules: z.string().optional(),
});

type CreateCampaignForm = z.infer<typeof createCampaignSchema>;

export default function CreateCampaignPage() {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const [error, setError] = useState('');

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<CreateCampaignForm>({
    resolver: zodResolver(createCampaignSchema),
    defaultValues: {
      max_players: 5,
    },
  });

  if (!user) {
    router.push('/login');
    return null;
  }

  const onSubmit = async (data: CreateCampaignForm) => {
    try {
      setError('');
      const response = await campaignApi.create(data);
      router.push(`/campaigns/${response.data.campaign.id}`);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create campaign');
    }
  };

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <Link
            href="/campaigns"
            className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to campaigns
          </Link>
        </div>

        <div className="bg-gray-800 rounded-lg p-8">
          <h1 className="text-3xl font-bold text-white mb-2">Create New Campaign</h1>
          <p className="text-gray-400 mb-8">Set up your world and invite players to join</p>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-2">
                Campaign Name
              </label>
              <input
                {...register('name')}
                type="text"
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="e.g., The Lost Mines of Phandelver"
              />
              {errors.name && (
                <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-300 mb-2">
                Description
              </label>
              <textarea
                {...register('description')}
                rows={4}
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Describe your campaign's story and what players can expect..."
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-500">{errors.description.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="setting" className="block text-sm font-medium text-gray-300 mb-2">
                Setting
              </label>
              <input
                {...register('setting')}
                type="text"
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="e.g., Forgotten Realms, Eberron, Homebrew World"
              />
              {errors.setting && (
                <p className="mt-1 text-sm text-red-500">{errors.setting.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="max_players" className="block text-sm font-medium text-gray-300 mb-2">
                Maximum Players
              </label>
              <input
                {...register('max_players', { valueAsNumber: true })}
                type="number"
                min="2"
                max="8"
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              />
              {errors.max_players && (
                <p className="mt-1 text-sm text-red-500">{errors.max_players.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="homebrew_rules" className="block text-sm font-medium text-gray-300 mb-2">
                Homebrew Rules (Optional)
              </label>
              <textarea
                {...register('homebrew_rules')}
                rows={3}
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Any custom rules or modifications to standard D&D 5e..."
              />
            </div>

            <div className="flex gap-4">
              <button
                type="submit"
                disabled={isSubmitting}
                className="flex-1 flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="animate-spin -ml-1 mr-3 h-5 w-5" />
                    Creating campaign...
                  </>
                ) : (
                  'Create Campaign'
                )}
              </button>
              <Link
                href="/campaigns"
                className="flex-1 flex justify-center py-3 px-4 border border-gray-600 rounded-lg text-gray-300 hover:bg-gray-700 transition-colors"
              >
                Cancel
              </Link>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}