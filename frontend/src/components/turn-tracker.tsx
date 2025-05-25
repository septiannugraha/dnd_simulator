'use client';

import { ChevronRight, SkipForward } from 'lucide-react';
import { sessionApi } from '@/lib/api';

interface TurnTrackerProps {
  turnOrder: Array<{
    character_id: string;
    character_name: string;
    initiative: number;
  }>;
  currentTurn?: string;
  isDM: boolean;
  sessionId: string;
}

export default function TurnTracker({
  turnOrder,
  currentTurn,
  isDM,
  sessionId,
}: TurnTrackerProps) {
  const handleAdvanceTurn = async () => {
    try {
      await sessionApi.advanceTurn(sessionId);
    } catch (error) {
      console.error('Failed to advance turn:', error);
    }
  };

  const sortedTurnOrder = [...turnOrder].sort((a, b) => b.initiative - a.initiative);

  return (
    <div className="space-y-2">
      {sortedTurnOrder.map((character, index) => (
        <div
          key={character.character_id}
          className={`flex items-center justify-between p-2 rounded ${
            character.character_id === currentTurn
              ? 'bg-yellow-500/20 border border-yellow-500'
              : 'bg-gray-700'
          }`}
        >
          <div className="flex items-center gap-2">
            {character.character_id === currentTurn && (
              <ChevronRight className="h-4 w-4 text-yellow-500" />
            )}
            <div>
              <p className="text-sm text-white">{character.character_name}</p>
              <p className="text-xs text-gray-400">Initiative: {character.initiative}</p>
            </div>
          </div>
        </div>
      ))}

      {isDM && turnOrder.length > 0 && (
        <button
          onClick={handleAdvanceTurn}
          className="w-full flex items-center justify-center gap-2 px-3 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors text-sm"
        >
          <SkipForward className="h-4 w-4" />
          Next Turn
        </button>
      )}
    </div>
  );
}