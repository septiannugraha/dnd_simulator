'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Settings, Save } from 'lucide-react'
import { useToast } from '@/hooks/use-toast'
import api from '@/lib/api'

const settingsSchema = z.object({
  aiEnabled: z.boolean(),
  aiModel: z.string(),
  responseStyle: z.string(),
  contextPrompt: z.string().optional(),
  autoNarrate: z.boolean(),
  suggestionMode: z.boolean(),
})

type SettingsFormData = z.infer<typeof settingsSchema>

interface AIDMSettingsProps {
  sessionId: string
  isDM: boolean
}

export function AIDMSettings({ sessionId, isDM }: AIDMSettingsProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const { toast } = useToast()
  
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<SettingsFormData>({
    resolver: zodResolver(settingsSchema),
    defaultValues: {
      aiEnabled: true,
      aiModel: 'gpt-4',
      responseStyle: 'immersive',
      autoNarrate: true,
      suggestionMode: false,
    },
  })

  const aiEnabled = watch('aiEnabled')

  const responseStyles = [
    { value: 'immersive', label: 'Immersive Storytelling' },
    { value: 'concise', label: 'Concise & Direct' },
    { value: 'descriptive', label: 'Rich & Descriptive' },
    { value: 'dark', label: 'Dark & Gritty' },
    { value: 'humorous', label: 'Light & Humorous' },
  ]

  const aiModels = [
    { value: 'gpt-4', label: 'GPT-4 (Most Creative)' },
    { value: 'gpt-3.5-turbo', label: 'GPT-3.5 (Fast)' },
    { value: 'claude-3', label: 'Claude 3 (Balanced)' },
  ]

  const onSubmit = async (data: SettingsFormData) => {
    if (!isDM) {
      toast({
        title: 'Permission Denied',
        description: 'Only the DM can modify AI settings.',
        variant: 'destructive',
      })
      return
    }

    setIsSubmitting(true)
    try {
      await api.put(`/sessions/${sessionId}/ai/settings`, data)
      
      toast({
        title: 'Settings Saved',
        description: 'AI DM settings have been updated.',
      })
    } catch (error) {
      console.error('Failed to save settings:', error)
      toast({
        title: 'Error',
        description: 'Failed to save settings. Please try again.',
        variant: 'destructive',
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  if (!isDM) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Settings className="h-5 w-5" />
            AI DM Settings
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-gray-500">
            Only the Dungeon Master can configure AI settings.
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings className="h-5 w-5" />
          AI DM Settings
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="flex items-center justify-between">
            <Label htmlFor="aiEnabled">Enable AI Dungeon Master</Label>
            <Switch
              id="aiEnabled"
              checked={aiEnabled}
              onCheckedChange={(checked) => setValue('aiEnabled', checked)}
            />
          </div>

          {aiEnabled && (
            <>
              <div className="space-y-2">
                <Label htmlFor="aiModel">AI Model</Label>
                <Select 
                  onValueChange={(value) => setValue('aiModel', value)}
                  defaultValue="gpt-4"
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select AI model" />
                  </SelectTrigger>
                  <SelectContent>
                    {aiModels.map((model) => (
                      <SelectItem key={model.value} value={model.value}>
                        {model.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="responseStyle">Response Style</Label>
                <Select 
                  onValueChange={(value) => setValue('responseStyle', value)}
                  defaultValue="immersive"
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select response style" />
                  </SelectTrigger>
                  <SelectContent>
                    {responseStyles.map((style) => (
                      <SelectItem key={style.value} value={style.value}>
                        {style.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="contextPrompt">Custom Context Prompt</Label>
                <Textarea
                  id="contextPrompt"
                  placeholder="Add any special rules or context for the AI DM..."
                  className="min-h-[80px]"
                  {...register('contextPrompt')}
                />
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="autoNarrate">Auto-narrate player actions</Label>
                <Switch
                  id="autoNarrate"
                  {...register('autoNarrate')}
                />
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="suggestionMode">Show action suggestions</Label>
                <Switch
                  id="suggestionMode"
                  {...register('suggestionMode')}
                />
              </div>
            </>
          )}

          <Button 
            type="submit" 
            className="w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? (
              <>Saving...</>
            ) : (
              <>
                <Save className="h-4 w-4 mr-2" />
                Save Settings
              </>
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}