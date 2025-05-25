'use client';

import { useState } from 'react';
import { Sword, Shield, Sparkles, MessageCircle } from 'lucide-react';

interface ActionPanelProps {
  characterId: string;
  sessionId: string;
  onAction: (action: string, actionType: string, target?: string) => void;
}

const ACTION_TEMPLATES = {
  combat: [
    { label: 'Attack', icon: Sword, prompt: 'I attack [target] with my [weapon]' },
    { label: 'Defend', icon: Shield, prompt: 'I take a defensive stance' },
    { label: 'Cast Spell', icon: Sparkles, prompt: 'I cast [spell] at [target]' },
  ],
  roleplay: [
    { label: 'Speak', icon: MessageCircle, prompt: 'I say to [character]: "[dialogue]"' },
    { label: 'Examine', icon: MessageCircle, prompt: 'I examine [object/area]' },
    { label: 'Interact', icon: MessageCircle, prompt: 'I interact with [object/character]' },
  ],
  exploration: [
    { label: 'Move', icon: MessageCircle, prompt: 'I move to [location]' },
    { label: 'Search', icon: MessageCircle, prompt: 'I search for [item/clue]' },
    { label: 'Use Item', icon: MessageCircle, prompt: 'I use [item]' },
  ],
};

export default function ActionPanel({ characterId, sessionId, onAction }: ActionPanelProps) {
  const [actionType, setActionType] = useState<'combat' | 'roleplay' | 'exploration'>('combat');
  const [customAction, setCustomAction] = useState('');
  const [target, setTarget] = useState('');

  const handleTemplateAction = (template: any) => {
    setCustomAction(template.prompt);
  };

  const handleSubmitAction = (e: React.FormEvent) => {
    e.preventDefault();
    if (customAction.trim()) {
      onAction(customAction.trim(), actionType, target || undefined);
      setCustomAction('');
      setTarget('');
    }
  };

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold text-white">Actions</h3>

      {/* Action Type Tabs */}
      <div className="flex gap-2">
        {(['combat', 'roleplay', 'exploration'] as const).map((type) => (
          <button
            key={type}
            onClick={() => setActionType(type)}
            className={`flex-1 px-3 py-2 rounded text-sm capitalize transition-colors ${
              actionType === type
                ? 'bg-purple-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            }`}
          >
            {type}
          </button>
        ))}
      </div>

      {/* Action Templates */}
      <div className="space-y-2">
        <p className="text-sm text-gray-400">Quick Actions</p>
        {ACTION_TEMPLATES[actionType].map((template, index) => {
          const Icon = template.icon;
          return (
            <button
              key={index}
              onClick={() => handleTemplateAction(template)}
              className="w-full flex items-center gap-3 px-3 py-2 bg-gray-700 hover:bg-gray-600 rounded text-left transition-colors"
            >
              <Icon className="h-4 w-4 text-purple-400" />
              <span className="text-sm text-white">{template.label}</span>
            </button>
          );
        })}
      </div>

      {/* Custom Action Form */}
      <form onSubmit={handleSubmitAction} className="space-y-3">
        <div>
          <label className="block text-sm text-gray-400 mb-1">
            Target (optional)
          </label>
          <input
            type="text"
            value={target}
            onChange={(e) => setTarget(e.target.value)}
            placeholder="e.g., goblin, door, NPC name"
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
          />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">
            Action Description
          </label>
          <textarea
            value={customAction}
            onChange={(e) => setCustomAction(e.target.value)}
            placeholder="Describe your action..."
            rows={3}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
          />
        </div>
        <button
          type="submit"
          disabled={!customAction.trim()}
          className="w-full py-2 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Submit Action
        </button>
      </form>
    </div>
  );
}