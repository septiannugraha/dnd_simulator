'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { characterApi, dndApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { Loader2, ArrowLeft, Dice6 } from 'lucide-react';

const createCharacterSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  race: z.string().min(1, 'Please select a race'),
  class: z.string().min(1, 'Please select a class'),
  background: z.string().min(1, 'Please select a background'),
  alignment: z.string().min(1, 'Please select an alignment'),
  ability_scores: z.object({
    strength: z.number().min(3).max(20),
    dexterity: z.number().min(3).max(20),
    constitution: z.number().min(3).max(20),
    intelligence: z.number().min(3).max(20),
    wisdom: z.number().min(3).max(20),
    charisma: z.number().min(3).max(20),
  }),
});

type CreateCharacterForm = z.infer<typeof createCharacterSchema>;

interface DndData {
  races: string[];
  classes: string[];
  backgrounds: string[];
}

const ALIGNMENTS = [
  'Lawful Good',
  'Neutral Good',
  'Chaotic Good',
  'Lawful Neutral',
  'True Neutral',
  'Chaotic Neutral',
  'Lawful Evil',
  'Neutral Evil',
  'Chaotic Evil',
];

export default function CreateCharacterPage() {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const [error, setError] = useState('');
  const [dndData, setDndData] = useState<DndData>({ races: [], classes: [], backgrounds: [] });
  const [rolledScores, setRolledScores] = useState<number[]>([]);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch,
  } = useForm<CreateCharacterForm>({
    resolver: zodResolver(createCharacterSchema),
    defaultValues: {
      ability_scores: {
        strength: 10,
        dexterity: 10,
        constitution: 10,
        intelligence: 10,
        wisdom: 10,
        charisma: 10,
      },
    },
  });

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }
    loadDndData();
  }, [user, router]);

  const loadDndData = async () => {
    try {
      const [racesRes, classesRes, backgroundsRes] = await Promise.all([
        dndApi.getRaces(),
        dndApi.getClasses(),
        dndApi.getBackgrounds(),
      ]);
      setDndData({
        races: racesRes.data.races || [],
        classes: classesRes.data.classes || [],
        backgrounds: backgroundsRes.data.backgrounds || [],
      });
    } catch (err: any) {
      setError('Failed to load character creation data');
    }
  };

  const rollAbilityScores = () => {
    const scores = [];
    for (let i = 0; i < 6; i++) {
      const rolls = Array(4)
        .fill(0)
        .map(() => Math.floor(Math.random() * 6) + 1)
        .sort((a, b) => b - a);
      scores.push(rolls.slice(0, 3).reduce((a, b) => a + b, 0));
    }
    setRolledScores(scores.sort((a, b) => b - a));
  };

  const assignRolledScore = (ability: keyof CreateCharacterForm['ability_scores'], index: number) => {
    if (rolledScores[index]) {
      setValue(`ability_scores.${ability}`, rolledScores[index]);
      const newScores = [...rolledScores];
      newScores.splice(index, 1);
      setRolledScores(newScores);
    }
  };

  if (!user) return null;

  const onSubmit = async (data: CreateCharacterForm) => {
    try {
      setError('');
      const response = await characterApi.create(data);
      router.push(`/characters/${response.data.character.id}`);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create character');
    }
  };

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <Link
            href="/characters"
            className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to characters
          </Link>
        </div>

        <div className="bg-gray-800 rounded-lg p-8">
          <h1 className="text-3xl font-bold text-white mb-2">Create New Character</h1>
          <p className="text-gray-400 mb-8">Build your hero following D&D 5e rules</p>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
            {error && (
              <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-2">
                  Character Name
                </label>
                <input
                  {...register('name')}
                  type="text"
                  className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  placeholder="e.g., Aragorn"
                />
                {errors.name && (
                  <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="race" className="block text-sm font-medium text-gray-300 mb-2">
                  Race
                </label>
                <select
                  {...register('race')}
                  className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option value="">Select a race</option>
                  {dndData.races.map((race) => (
                    <option key={race} value={race}>
                      {race}
                    </option>
                  ))}
                </select>
                {errors.race && (
                  <p className="mt-1 text-sm text-red-500">{errors.race.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="class" className="block text-sm font-medium text-gray-300 mb-2">
                  Class
                </label>
                <select
                  {...register('class')}
                  className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option value="">Select a class</option>
                  {dndData.classes.map((cls) => (
                    <option key={cls} value={cls}>
                      {cls}
                    </option>
                  ))}
                </select>
                {errors.class && (
                  <p className="mt-1 text-sm text-red-500">{errors.class.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="background" className="block text-sm font-medium text-gray-300 mb-2">
                  Background
                </label>
                <select
                  {...register('background')}
                  className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option value="">Select a background</option>
                  {dndData.backgrounds.map((bg) => (
                    <option key={bg} value={bg}>
                      {bg}
                    </option>
                  ))}
                </select>
                {errors.background && (
                  <p className="mt-1 text-sm text-red-500">{errors.background.message}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label htmlFor="alignment" className="block text-sm font-medium text-gray-300 mb-2">
                  Alignment
                </label>
                <select
                  {...register('alignment')}
                  className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option value="">Select an alignment</option>
                  {ALIGNMENTS.map((alignment) => (
                    <option key={alignment} value={alignment}>
                      {alignment}
                    </option>
                  ))}
                </select>
                {errors.alignment && (
                  <p className="mt-1 text-sm text-red-500">{errors.alignment.message}</p>
                )}
              </div>
            </div>

            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium text-white">Ability Scores</h3>
                <button
                  type="button"
                  onClick={rollAbilityScores}
                  className="flex items-center gap-2 px-3 py-1 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors text-sm"
                >
                  <Dice6 className="h-4 w-4" />
                  Roll Scores
                </button>
              </div>

              {rolledScores.length > 0 && (
                <div className="mb-4 p-4 bg-gray-700 rounded-lg">
                  <p className="text-sm text-gray-300 mb-2">Rolled scores (click to assign):</p>
                  <div className="flex gap-2 flex-wrap">
                    {rolledScores.map((score, index) => (
                      <span
                        key={index}
                        className="px-3 py-1 bg-purple-600 text-white rounded cursor-pointer hover:bg-purple-700 transition-colors"
                      >
                        {score}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {Object.entries({
                  strength: 'STR',
                  dexterity: 'DEX',
                  constitution: 'CON',
                  intelligence: 'INT',
                  wisdom: 'WIS',
                  charisma: 'CHA',
                }).map(([ability, label]) => (
                  <div key={ability}>
                    <label className="block text-sm font-medium text-gray-300 mb-1">
                      {label}
                    </label>
                    <input
                      {...register(`ability_scores.${ability as keyof CreateCharacterForm['ability_scores']}`, {
                        valueAsNumber: true,
                      })}
                      type="number"
                      min="3"
                      max="20"
                      className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-center focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    />
                  </div>
                ))}
              </div>
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
                    Creating character...
                  </>
                ) : (
                  'Create Character'
                )}
              </button>
              <Link
                href="/characters"
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