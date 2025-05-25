interface WebSocketMessage {
  type: string;
  payload: any;
  timestamp?: string;
}

interface DiceRollMessage {
  character_id: string;
  character_name: string;
  dice: string;
  result: number;
  details?: number[];
}

interface ChatMessage {
  user_id: string;
  username: string;
  character_id?: string;
  character_name?: string;
  message: string;
  timestamp: string;
}

interface SessionUpdate {
  status: string;
  current_scene?: string;
  current_turn?: string;
  turn_order?: Array<{
    character_id: string;
    character_name: string;
    initiative: number;
  }>;
}

class WebSocketService {
  private ws: WebSocket | null = null;
  private sessionId: string | null = null;
  private token: string | null = null;
  private characterId: string | null = null;
  private listeners: Map<string, Set<Function>> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1 second

  connect(sessionId: string, token: string, characterId?: string) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return;
    }

    this.sessionId = sessionId;
    this.token = token;
    this.characterId = characterId || null;
    this.reconnectAttempts = 0;

    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080';
    // The backend expects session_id and token as query parameters
    let url = `${wsUrl}/api/sessions/${sessionId}/ws?session_id=${sessionId}&token=${token}`;
    if (characterId) {
      url += `&character_id=${characterId}`;
    }

    try {
      // Note: Browser WebSocket API doesn't support custom headers.
      // The auth token is sent as a query parameter
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;
        this.emit('connected', {});
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket disconnected', event);
        this.emit('disconnected', {});
        this.attemptReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.emit('error', error);
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          this.handleMessage(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      this.emit('error', error);
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

    console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`);

    setTimeout(() => {
      if (this.sessionId && this.token) {
        this.connect(this.sessionId, this.token, this.characterId || undefined);
      }
    }, delay);
  }

  private handleMessage(message: any) {
    // Handle different message types based on backend structure
    switch (message.type) {
      case 'player_joined':
        this.emit('player_joined', message.data);
        break;
      case 'player_left':
        this.emit('player_left', message.data);
        break;
      case 'chat_ic':
      case 'chat_ooc':
        this.emit('chat_message', {
          type: message.type,
          user_id: message.user_id,
          username: message.username,
          character_id: message.data?.character_id,
          message: message.data?.content,
          timestamp: message.timestamp,
        });
        break;
      case 'dice_result':
        this.emit('dice_roll', {
          user_id: message.user_id,
          username: message.username,
          character_id: message.data?.character_id,
          dice: message.data?.dice,
          result: message.data?.result,
          total: message.data?.total,
          purpose: message.data?.purpose,
          special_message: message.data?.special_message,
        });
        break;
      case 'session_update':
        this.emit('session_update', message.data);
        break;
      case 'character_update':
        this.emit('character_update', message.data);
        break;
      case 'ai_response':
        this.emit('ai_response', message.data);
        break;
      case 'turn_update':
        this.emit('turn_update', message.data);
        break;
      default:
        console.warn('Unknown message type:', message.type);
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
      this.sessionId = null;
      this.token = null;
      this.characterId = null;
    }
  }

  send(type: string, data: any) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');
      return;
    }

    const message = {
      type,
      data,
      timestamp: new Date().toISOString(),
    };

    this.ws.send(JSON.stringify(message));
  }

  sendChatMessage(message: string, characterId?: string, isIC: boolean = false) {
    const type = isIC ? 'chat_ic' : 'chat_ooc';
    this.send(type, {
      content: message,
      character_id: characterId,
      is_ic: isIC,
    });
  }

  rollDice(dice: string, characterId: string, purpose?: string) {
    this.send('dice_roll', {
      dice,
      character_id: characterId,
      purpose,
    });
  }

  updateCharacter(characterId: string, field: string, value: any, oldValue?: any) {
    this.send('character_update', {
      character_id: characterId,
      field,
      value,
      old_value: oldValue,
    });
  }

  sendAction(characterId: string, action: string, actionType: string, target?: string) {
    this.send('player_action', {
      character_id: characterId,
      action,
      action_type: actionType,
      target,
    });
  }

  // Event listener management
  on(event: string, callback: Function) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event)!.add(callback);
  }

  off(event: string, callback: Function) {
    if (this.listeners.has(event)) {
      this.listeners.get(event)!.delete(callback);
    }
  }

  private emit(event: string, data: any) {
    if (this.listeners.has(event)) {
      this.listeners.get(event)!.forEach(callback => {
        callback(data);
      });
    }
  }
}

// Export singleton instance
export const wsService = new WebSocketService();

// Export types
export type {
  WebSocketMessage,
  DiceRollMessage,
  ChatMessage,
  SessionUpdate,
};