'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Send, Wand2 } from 'lucide-react'
import { useToast } from '@/hooks/use-toast'
import api from '@/lib/api'

const actionSchema = z.object({
  actionType: z.string().min(1, 'Please select an action type'),
  description: z.string().min(10, 'Please describe your action in detail'),
  target: z.string().optional(),
})

type ActionFormData = z.infer<typeof actionSchema>

interface AIActionPanelProps {
  sessionId: string
  characterId: string
  onActionSubmit?: (action: ActionFormData) => void
}

export function AIActionPanel({ sessionId, characterId, onActionSubmit }: AIActionPanelProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const { toast } = useToast()
  
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<ActionFormData>({
    resolver: zodResolver(actionSchema),
  })

  const actionTypes = [
    { value: 'explore', label: 'Explore' },
    { value: 'interact', label: 'Interact' },
    { value: 'combat', label: 'Combat Action' },
    { value: 'skill', label: 'Use Skill' },
    { value: 'speak', label: 'Speak/Persuade' },
    { value: 'custom', label: 'Custom Action' },
  ]

  const onSubmit = async (data: ActionFormData) => {
    setIsSubmitting(true)
    try {
      // Send action to AI endpoint
      const response = await api.post(`/sessions/${sessionId}/ai/action`, {
        characterId,
        ...data,
      })

      toast({
        title: 'Action Submitted',
        description: 'The AI DM is processing your action...',
      })

      if (onActionSubmit) {
        onActionSubmit(data)
      }

      reset()
    } catch (error) {
      console.error('Failed to submit action:', error)
      toast({
        title: 'Error',
        description: 'Failed to submit action. Please try again.',
        variant: 'destructive',
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  const quickActions = [
    { label: 'Look Around', action: 'I carefully examine my surroundings, looking for anything unusual or noteworthy.' },
    { label: 'Search', action: 'I search the area thoroughly for hidden items, doors, or clues.' },
    { label: 'Listen', action: 'I stop and listen carefully for any sounds or movements nearby.' },
    { label: 'Check for Traps', action: 'I cautiously check the area for any traps or dangers.' },
  ]

  const handleQuickAction = (action: string) => {
    setValue('description', action)
    setValue('actionType', 'explore')
  }

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Wand2 className="h-5 w-5" />
          Submit Action to AI DM
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="actionType">Action Type</Label>
            <Select onValueChange={(value) => setValue('actionType', value)}>
              <SelectTrigger>
                <SelectValue placeholder="Select action type" />
              </SelectTrigger>
              <SelectContent>
                {actionTypes.map((type) => (
                  <SelectItem key={type.value} value={type.value}>
                    {type.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.actionType && (
              <p className="text-sm text-red-500">{errors.actionType.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Describe Your Action</Label>
            <Textarea
              id="description"
              placeholder="Describe what your character does..."
              className="min-h-[100px]"
              {...register('description')}
            />
            {errors.description && (
              <p className="text-sm text-red-500">{errors.description.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="target">Target (Optional)</Label>
            <input
              type="text"
              id="target"
              placeholder="e.g., 'the mysterious door', 'the goblin leader'"
              className="w-full px-3 py-2 border rounded-md"
              {...register('target')}
            />
          </div>

          <div className="space-y-2">
            <Label>Quick Actions</Label>
            <div className="flex flex-wrap gap-2">
              {quickActions.map((quick, index) => (
                <Button
                  key={index}
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => handleQuickAction(quick.action)}
                >
                  {quick.label}
                </Button>
              ))}
            </div>
          </div>

          <Button 
            type="submit" 
            className="w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? (
              <>Submitting...</>
            ) : (
              <>
                <Send className="h-4 w-4 mr-2" />
                Submit Action
              </>
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}