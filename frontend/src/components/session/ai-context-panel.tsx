'use client'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Info, Users, MapPin, History } from 'lucide-react'

interface GameContext {
  currentScene: string
  activeCharacters: Array<{
    id: string
    name: string
    class: string
    hp: number
    maxHp: number
  }>
  recentEvents: Array<{
    id: string
    description: string
    timestamp: string
  }>
  sessionNotes?: string
}

interface AIContextPanelProps {
  context: GameContext
}

export function AIContextPanel({ context }: AIContextPanelProps) {
  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Info className="h-5 w-5" />
          Game Context
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Current Scene */}
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium">
            <MapPin className="h-4 w-4" />
            Current Scene
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            {context.currentScene}
          </p>
        </div>

        {/* Active Characters */}
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium">
            <Users className="h-4 w-4" />
            Active Characters
          </div>
          <div className="space-y-2">
            {context.activeCharacters.map((character) => (
              <div
                key={character.id}
                className="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-800 rounded-md"
              >
                <div>
                  <p className="font-medium text-sm">{character.name}</p>
                  <p className="text-xs text-gray-500">{character.class}</p>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant={character.hp < character.maxHp * 0.3 ? 'destructive' : 'default'}>
                    {character.hp}/{character.maxHp} HP
                  </Badge>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Events */}
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium">
            <History className="h-4 w-4" />
            Recent Events
          </div>
          <ScrollArea className="h-[150px]">
            <div className="space-y-2">
              {context.recentEvents.map((event) => (
                <div
                  key={event.id}
                  className="p-2 bg-gray-50 dark:bg-gray-800 rounded-md"
                >
                  <p className="text-sm">{event.description}</p>
                  <p className="text-xs text-gray-500 mt-1">
                    {new Date(event.timestamp).toLocaleTimeString()}
                  </p>
                </div>
              ))}
            </div>
          </ScrollArea>
        </div>

        {/* Session Notes */}
        {context.sessionNotes && (
          <div className="space-y-2">
            <p className="text-sm font-medium">Session Notes</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {context.sessionNotes}
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}