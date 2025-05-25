'use client'

import { useEffect, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Sparkles } from 'lucide-react'

interface AIResponse {
  id: string
  message: string
  context?: {
    scene?: string
    characters?: string[]
    recentActions?: string[]
  }
  createdAt: string
}

interface AIDMResponseProps {
  sessionId: string
  onNewResponse?: (response: AIResponse) => void
}

export function AIDMResponse({ sessionId, onNewResponse }: AIDMResponseProps) {
  const [responses, setResponses] = useState<AIResponse[]>([])
  const [isGenerating, setIsGenerating] = useState(false)

  useEffect(() => {
    // In a real implementation, this would listen for AI responses via WebSocket
    // For now, we'll simulate with mock data
    const mockResponse: AIResponse = {
      id: '1',
      message: 'You find yourselves standing at the entrance of an ancient tomb. The heavy stone doors are slightly ajar, revealing only darkness beyond. Strange symbols are carved into the weathered stone, and you can hear a faint whistling sound coming from within.',
      context: {
        scene: 'Ancient Tomb Entrance',
        characters: ['Thraldin Ironforge', 'Lyra Moonwhisper'],
        recentActions: ['Arrived at the tomb', 'Examined the entrance']
      },
      createdAt: new Date().toISOString()
    }
    
    setResponses([mockResponse])
  }, [sessionId])

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  }

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Sparkles className="h-5 w-5" />
          AI Dungeon Master
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4 max-h-[400px] overflow-y-auto">
          {responses.map((response) => (
            <div key={response.id} className="space-y-2">
              <div className="flex items-start gap-3">
                <div className="bg-purple-500 text-white rounded-full p-2">
                  <Sparkles className="h-4 w-4" />
                </div>
                <div className="flex-1">
                  <div className="bg-gray-100 dark:bg-gray-800 rounded-lg p-3">
                    <p className="text-sm whitespace-pre-wrap">{response.message}</p>
                  </div>
                  <p className="text-xs text-gray-500 mt-1">
                    {formatTime(response.createdAt)}
                  </p>
                </div>
              </div>
            </div>
          ))}
          
          {isGenerating && (
            <div className="flex items-center gap-2 text-gray-500">
              <div className="animate-pulse flex gap-1">
                <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse delay-75"></div>
                <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse delay-150"></div>
              </div>
              <span className="text-sm">AI is generating response...</span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}