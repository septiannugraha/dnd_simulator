'use client';

import { useState } from 'react';
import { Edit, Save, X } from 'lucide-react';
import { sessionApi } from '@/lib/api';

interface SceneDescriptionProps {
  scene: string;
  notes?: string;
  isDM: boolean;
  sessionId: string;
}

export default function SceneDescription({
  scene,
  notes,
  isDM,
  sessionId,
}: SceneDescriptionProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editedScene, setEditedScene] = useState(scene);
  const [editedNotes, setEditedNotes] = useState(notes || '');

  const handleSave = async () => {
    try {
      await sessionApi.updateScene(sessionId, editedScene, editedNotes);
      setIsEditing(false);
    } catch (error) {
      console.error('Failed to update scene:', error);
    }
  };

  const handleCancel = () => {
    setEditedScene(scene);
    setEditedNotes(notes || '');
    setIsEditing(false);
  };

  if (isEditing && isDM) {
    return (
      <div className="space-y-3">
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">
            Scene Description
          </label>
          <textarea
            value={editedScene}
            onChange={(e) => setEditedScene(e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            rows={3}
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">
            DM Notes (private)
          </label>
          <textarea
            value={editedNotes}
            onChange={(e) => setEditedNotes(e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            rows={2}
          />
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleSave}
            className="flex items-center gap-2 px-3 py-1 bg-green-600 hover:bg-green-700 text-white rounded transition-colors text-sm"
          >
            <Save className="h-4 w-4" />
            Save
          </button>
          <button
            onClick={handleCancel}
            className="flex items-center gap-2 px-3 py-1 bg-gray-600 hover:bg-gray-700 text-white rounded transition-colors text-sm"
          >
            <X className="h-4 w-4" />
            Cancel
          </button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-start mb-2">
        <h2 className="text-lg font-semibold text-white">Current Scene</h2>
        {isDM && (
          <button
            onClick={() => setIsEditing(true)}
            className="p-1 text-gray-400 hover:text-white transition-colors"
          >
            <Edit className="h-4 w-4" />
          </button>
        )}
      </div>
      <p className="text-gray-300">{scene}</p>
      {isDM && notes && (
        <div className="mt-3 p-3 bg-gray-700/50 rounded">
          <p className="text-xs text-gray-400 mb-1">DM Notes (private)</p>
          <p className="text-sm text-gray-300">{notes}</p>
        </div>
      )}
    </div>
  );
}