'use client';

import { useEffect, useState, useRef } from 'react';
import { useParams } from 'next/navigation';
import { sessionApi } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { wsService, ChatMessage, DiceRollMessage } from '@/lib/websocket';
import ChatPanel from '@/components/chat-panel';
import DiceRoller from '@/components/dice-roller';
import CharacterStatus from '@/components/character-status';
import TurnTracker from '@/components/turn-tracker';
import SceneDescription from '@/components/scene-description';
import ActionPanel from '@/components/action-panel';
import { AIDMResponse } from '@/components/session/ai-dm-response';
import { AIActionPanel } from '@/components/session/ai-action-panel';
import { AIContextPanel } from '@/components/session/ai-context-panel';
import { AIDMSettings } from '@/components/session/ai-dm-settings';
import { Loader2, Users, Dice6, MessageSquare, Sword, Sparkles, Settings } from 'lucide-react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

interface SessionData {
  id: string;
  campaign_id: string;
  campaign_name: string;
  status: string;
  current_scene: string;
  scene_notes?: string;
  current_turn?: string;
  turn_order: Array<{
    character_id: string;
    character_name: string;
    initiative: number;
  }>;
  players: Array<{
    user_id: string;
    username: string;
    character_id?: string;
    character_name?: string;
    character_class?: string;
    character_level?: number;
    hit_points?: number;
    max_hit_points?: number;
  }>;
  dm_id: string;
  dm_username: string;
}

export default function SessionPage() {
  const params = useParams();
  const { user } = useAuthStore();
  const [session, setSession] = useState<SessionData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [diceRolls, setDiceRolls] = useState<DiceRollMessage[]>([]);
  const [selectedCharacterId, setSelectedCharacterId] = useState<string>('');
  const [aiResponses, setAiResponses] = useState<any[]>([]);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const sessionId = params.id as string;
  const isDM = user?.id === session?.dm_id;
  const playerData = session?.players?.find(p => p.user_id === user?.id);

  useEffect(() => {
    if (!user) return;
    loadSession();
  }, [user, sessionId]);

  useEffect(() => {
    if (!user || !session) return;
    
    // Connect WebSocket
    const token = localStorage.getItem('token');
    if (token) {
      const myPlayer = session.players?.find(p => p.user_id === user?.id);
      wsService.connect(sessionId, token, myPlayer?.character_id);
      
      // Set up event listeners
      wsService.on('chat_message', handleChatMessage);
      wsService.on('dice_roll', handleDiceRoll);
      wsService.on('session_update', handleSessionUpdate);
      wsService.on('player_joined', handlePlayerJoined);
      wsService.on('player_left', handlePlayerLeft);
      wsService.on('turn_update', handleTurnUpdate);
      wsService.on('ai_response', handleAIResponse);
    }
    
    return () => {
      wsService.disconnect();
    };
  }, [user, session, sessionId]);

  useEffect(() => {
    // Auto-scroll to bottom of messages
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const loadSession = async () => {
    try {
      setLoading(true);
      const response = await sessionApi.get(sessionId);
      setSession(response.data);
      
      // Set default selected character
      if (response.data.players) {
        const myPlayer = response.data.players?.find(
          (p: any) => p.user_id === user?.id
        );
        if (myPlayer?.character_id) {
          setSelectedCharacterId(myPlayer.character_id);
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load session');
    } finally {
      setLoading(false);
    }
  };

  const handleChatMessage = (data: ChatMessage) => {
    setMessages(prev => [...prev, data]);
  };

  const handleDiceRoll = (data: DiceRollMessage) => {
    setDiceRolls(prev => [...prev, data]);
    // Also add to chat
    const rollMessage: ChatMessage = {
      user_id: '',
      username: data.character_name,
      message: `rolled ${data.dice}: ${data.result}${data.details ? ` (${data.details.join(', ')})` : ''}`,
      timestamp: new Date().toISOString(),
    };
    setMessages(prev => [...prev, rollMessage]);
  };

  const handleSessionUpdate = (data: any) => {
    setSession(prev => prev ? { ...prev, ...data } : null);
  };

  const handlePlayerJoined = (data: any) => {
    // Update player list
    loadSession();
  };

  const handlePlayerLeft = (data: any) => {
    // Update player list
    loadSession();
  };

  const handleTurnUpdate = (data: any) => {
    setSession(prev => prev ? {
      ...prev,
      current_turn: data.current_turn,
      turn_order: data.turn_order || prev.turn_order,
    } : null);
  };

  const handleSendMessage = (message: string) => {
    wsService.sendChatMessage(message, selectedCharacterId);
  };

  const handleRollDice = (dice: string) => {
    if (selectedCharacterId) {
      wsService.rollDice(dice, selectedCharacterId);
    }
  };

  const handleAction = (action: string, actionType: string, target?: string) => {
    if (selectedCharacterId) {
      wsService.sendAction(selectedCharacterId, action, actionType, target);
    }
  };

  const handleAIResponse = (data: any) => {
    setAiResponses(prev => [...prev, {
      id: Date.now().toString(),
      message: data.message,
      context: data.context,
      createdAt: new Date().toISOString(),
    }]);
    
    // Also add to chat as a special AI message
    const aiMessage: ChatMessage = {
      user_id: 'ai-dm',
      username: 'AI Dungeon Master',
      message: data.message,
      timestamp: new Date().toISOString(),
    };
    setMessages(prev => [...prev, aiMessage]);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-purple-500" />
      </div>
    );
  }

  if (error || !session) {
    return (
      <div className="min-h-screen bg-gray-900 p-8">
        <div className="max-w-3xl mx-auto">
          <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-3 rounded-lg">
            {error || 'Session not found'}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="h-screen flex flex-col">
        {/* Header */}
        <div className="bg-gray-800 border-b border-gray-700 px-4 py-3">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-xl font-bold text-white">{session.campaign_name}</h1>
              <p className="text-sm text-gray-400">
                DM: {session.dm_username} • Status: {session.status}
              </p>
            </div>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2 text-gray-400">
                <Users className="h-5 w-5" />
                <span>{session.players?.length || 0} players</span>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 flex overflow-hidden">
          {/* Left Sidebar - Players */}
          <div className="w-64 bg-gray-800 border-r border-gray-700 p-4 overflow-y-auto">
            <h2 className="text-lg font-semibold text-white mb-4">Players</h2>
            <div className="space-y-3">
              {session.players?.map((player) => (
                <CharacterStatus
                  key={player.user_id}
                  player={player}
                  isCurrentTurn={session.current_turn === player.character_id}
                  onSelect={() => setSelectedCharacterId(player.character_id || '')}
                  isSelected={selectedCharacterId === player.character_id}
                />
              ))}
            </div>
            
            {session.turn_order.length > 0 && (
              <div className="mt-6">
                <h3 className="text-lg font-semibold text-white mb-4">Turn Order</h3>
                <TurnTracker
                  turnOrder={session.turn_order}
                  currentTurn={session.current_turn}
                  isDM={isDM}
                  sessionId={sessionId}
                />
              </div>
            )}
          </div>

          {/* Center - Main Game Area */}
          <div className="flex-1 flex flex-col">
            {/* Scene Description */}
            <div className="bg-gray-850 border-b border-gray-700 p-4">
              <SceneDescription
                scene={session.current_scene}
                notes={session.scene_notes}
                isDM={isDM}
                sessionId={sessionId}
              />
            </div>

            {/* Chat and Actions */}
            <div className="flex-1 flex">
              <div className="flex-1 flex flex-col">
                <ChatPanel
                  messages={messages}
                  onSendMessage={handleSendMessage}
                  currentUserId={user?.id || ''}
                />
                <div ref={messagesEndRef} />
              </div>
              
              {/* Right Panel - Actions */}
              <div className="w-96 bg-gray-800 border-l border-gray-700">
                <Tabs defaultValue="actions" className="h-full">
                  <TabsList className={`grid w-full ${isDM ? 'grid-cols-4' : 'grid-cols-3'}`}>
                    <TabsTrigger value="actions">Actions</TabsTrigger>
                    <TabsTrigger value="ai-dm">AI DM</TabsTrigger>
                    <TabsTrigger value="context">Context</TabsTrigger>
                    {isDM && <TabsTrigger value="settings">Settings</TabsTrigger>}
                  </TabsList>
                  
                  <TabsContent value="actions" className="p-4 space-y-4">
                    <DiceRoller onRoll={handleRollDice} />
                    
                    {selectedCharacterId && (
                      <ActionPanel
                        characterId={selectedCharacterId}
                        sessionId={sessionId}
                        onAction={handleAction}
                      />
                    )}
                  </TabsContent>
                  
                  <TabsContent value="ai-dm" className="p-4 space-y-4">
                    <AIDMResponse 
                      sessionId={sessionId} 
                      onNewResponse={(response) => {
                        setAiResponses(prev => [...prev, response]);
                      }}
                    />
                    
                    {selectedCharacterId && (
                      <AIActionPanel
                        sessionId={sessionId}
                        characterId={selectedCharacterId}
                        onActionSubmit={(action) => {
                          console.log('AI Action submitted:', action);
                        }}
                      />
                    )}
                  </TabsContent>
                  
                  <TabsContent value="context" className="p-4">
                    <AIContextPanel
                      context={{
                        currentScene: session.current_scene,
                        activeCharacters: session.players
                          .filter(p => p.character_id)
                          .map(p => ({
                            id: p.character_id!,
                            name: p.character_name || 'Unknown',
                            class: p.character_class || 'Unknown',
                            hp: p.hit_points || 0,
                            maxHp: p.max_hit_points || 0,
                          })),
                        recentEvents: messages.slice(-5).map((msg, idx) => ({
                          id: String(idx),
                          description: `${msg.username}: ${msg.message}`,
                          timestamp: msg.timestamp,
                        })),
                        sessionNotes: session.scene_notes,
                      }}
                    />
                  </TabsContent>
                  
                  {isDM && (
                    <TabsContent value="settings" className="p-4">
                      <AIDMSettings sessionId={sessionId} isDM={isDM} />
                    </TabsContent>
                  )}
                </Tabs>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}