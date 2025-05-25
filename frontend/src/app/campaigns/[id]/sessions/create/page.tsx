'use client';

import { useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { sessionApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { Loader2, ArrowLeft, Calendar } from 'lucide-react';

const createSessionSchema = z.object({
  name: z.string().min(3, 'Session name must be at least 3 characters'),
  description: z.string().optional(),
  scheduled_for: z.string().optional(),
});

type CreateSessionForm = z.infer<typeof createSessionSchema>;

export default function CreateSessionPage() {
  const router = useRouter();
  const params = useParams();
  const user = useAuthStore((state) => state.user);
  const [error, setError] = useState('');

  const campaignId = params.id as string;

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<CreateSessionForm>({
    resolver: zodResolver(createSessionSchema),
  });

  if (!user) {
    router.push('/login');
    return null;
  }

  const onSubmit = async (data: CreateSessionForm) => {
    try {
      setError('');
      const sessionData = {
        campaign_id: campaignId,
        name: data.name,
        description: data.description || '',
        scheduled_for: data.scheduled_for || undefined,
      };
      const response = await sessionApi.create(sessionData);
      router.push(`/sessions/${response.data.id}`);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create session');
    }
  };

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <Link
            href={`/campaigns/${campaignId}`}
            className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to campaign
          </Link>
        </div>

        <div className="bg-gray-800 rounded-lg p-8">
          <h1 className="text-3xl font-bold text-white mb-2">Create New Session</h1>
          <p className="text-gray-400 mb-8">Start a new game session for your campaign</p>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-2">
                Session Name
              </label>
              <input
                {...register('name')}
                type="text"
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="e.g., Chapter 1: The Beginning"
              />
              {errors.name && (
                <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-300 mb-2">
                Description (Optional)
              </label>
              <textarea
                {...register('description')}
                rows={3}
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Brief description of what will happen in this session..."
              />
            </div>

            <div>
              <label htmlFor="scheduled_for" className="block text-sm font-medium text-gray-300 mb-2">
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4" />
                  Schedule For (Optional)
                </div>
              </label>
              <input
                {...register('scheduled_for')}
                type="datetime-local"
                className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
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
                    Creating session...
                  </>
                ) : (
                  'Create Session'
                )}
              </button>
              <Link
                href={`/campaigns/${campaignId}`}
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