'use client';

import { useState } from 'react';
import { Dice1, Dice2, Dice3, Dice4, Dice5, Dice6 } from 'lucide-react';

interface DiceRollerProps {
  onRoll: (dice: string) => void;
}

const COMMON_DICE = [
  { label: 'd4', value: '1d4', icon: Dice1 },
  { label: 'd6', value: '1d6', icon: Dice2 },
  { label: 'd8', value: '1d8', icon: Dice3 },
  { label: 'd10', value: '1d10', icon: Dice4 },
  { label: 'd12', value: '1d12', icon: Dice5 },
  { label: 'd20', value: '1d20', icon: Dice6 },
];

const QUICK_ROLLS = [
  { label: 'Attack', value: '1d20' },
  { label: 'Damage', value: '1d8' },
  { label: 'Initiative', value: '1d20' },
  { label: 'Skill Check', value: '1d20' },
  { label: 'Saving Throw', value: '1d20' },
];

export default function DiceRoller({ onRoll }: DiceRollerProps) {
  const [customDice, setCustomDice] = useState('');

  const handleCustomRoll = (e: React.FormEvent) => {
    e.preventDefault();
    if (customDice.trim()) {
      onRoll(customDice.trim());
      setCustomDice('');
    }
  };

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold text-white">Dice Roller</h3>
      
      {/* Common Dice */}
      <div>
        <p className="text-sm text-gray-400 mb-2">Common Dice</p>
        <div className="grid grid-cols-3 gap-2">
          {COMMON_DICE.map((dice) => {
            const Icon = dice.icon;
            return (
              <button
                key={dice.value}
                onClick={() => onRoll(dice.value)}
                className="flex flex-col items-center justify-center p-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              >
                <Icon className="h-6 w-6 text-purple-400 mb-1" />
                <span className="text-sm text-white">{dice.label}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Quick Rolls */}
      <div>
        <p className="text-sm text-gray-400 mb-2">Quick Rolls</p>
        <div className="grid grid-cols-2 gap-2">
          {QUICK_ROLLS.map((roll) => (
            <button
              key={roll.label}
              onClick={() => onRoll(roll.value)}
              className="px-3 py-2 bg-gray-700 hover:bg-gray-600 rounded text-sm text-white transition-colors"
            >
              {roll.label}
            </button>
          ))}
        </div>
      </div>

      {/* Custom Dice */}
      <form onSubmit={handleCustomRoll}>
        <p className="text-sm text-gray-400 mb-2">Custom Roll</p>
        <div className="flex gap-2">
          <input
            type="text"
            value={customDice}
            onChange={(e) => setCustomDice(e.target.value)}
            placeholder="e.g., 2d6+3"
            className="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
          />
          <button
            type="submit"
            disabled={!customDice.trim()}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Roll
          </button>
        </div>
      </form>
    </div>
  );
}