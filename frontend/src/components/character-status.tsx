'use client';

import { Heart, Shield, User } from 'lucide-react';

interface CharacterStatusProps {
  player: {
    user_id: string;
    username: string;
    character_id?: string;
    character_name?: string;
    character_class?: string;
    character_level?: number;
    hit_points?: number;
    max_hit_points?: number;
  };
  isCurrentTurn: boolean;
  isSelected: boolean;
  onSelect: () => void;
}

export default function CharacterStatus({
  player,
  isCurrentTurn,
  isSelected,
  onSelect,
}: CharacterStatusProps) {
  const hpPercentage = player.hit_points && player.max_hit_points
    ? (player.hit_points / player.max_hit_points) * 100
    : 100;

  const hpColor = hpPercentage > 50 ? 'bg-green-500' : hpPercentage > 25 ? 'bg-yellow-500' : 'bg-red-500';

  return (
    <div
      onClick={onSelect}
      className={`p-3 rounded-lg cursor-pointer transition-all ${
        isSelected
          ? 'bg-purple-600/20 border border-purple-500'
          : 'bg-gray-700 hover:bg-gray-650 border border-transparent'
      } ${isCurrentTurn ? 'ring-2 ring-yellow-500' : ''}`}
    >
      <div className="flex items-start gap-3">
        <div className="p-2 bg-gray-600 rounded-full">
          <User className="h-5 w-5 text-gray-300" />
        </div>
        <div className="flex-1 min-w-0">
          <p className="font-medium text-white truncate">
            {player.character_name || player.username}
          </p>
          {player.character_class && (
            <p className="text-sm text-gray-400">
              Lvl {player.character_level} {player.character_class}
            </p>
          )}
          {player.hit_points !== undefined && (
            <div className="mt-2">
              <div className="flex items-center gap-1 text-xs text-gray-400 mb-1">
                <Heart className="h-3 w-3" />
                <span>
                  {player.hit_points} / {player.max_hit_points}
                </span>
              </div>
              <div className="w-full bg-gray-600 rounded-full h-2 overflow-hidden">
                <div
                  className={`h-full transition-all ${hpColor}`}
                  style={{ width: `${hpPercentage}%` }}
                />
              </div>
            </div>
          )}
        </div>
      </div>
      {isCurrentTurn && (
        <div className="mt-2 text-xs text-yellow-400 font-medium text-center">
          Current Turn
        </div>
      )}
    </div>
  );
}