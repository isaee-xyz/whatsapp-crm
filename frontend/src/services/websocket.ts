import { useContactsStore } from '@/stores/contacts'

// WebSocket message types
const WS_TYPE_NEW_MESSAGE = 'new_message'
const WS_TYPE_STATUS_UPDATE = 'status_update'
const WS_TYPE_SET_CONTACT = 'set_contact'
const WS_TYPE_PING = 'ping'
const WS_TYPE_PONG = 'pong'

interface WSMessage {
  type: string
  payload: any
}

class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private pingInterval: number | null = null
  private isConnected = false

  connect(token: string) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected')
      return
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const url = `${protocol}//${host}/ws?token=${token}`

    console.log('Connecting to WebSocket:', url)

    try {
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.isConnected = true
        this.reconnectAttempts = 0
        this.startPing()
      }

      this.ws.onmessage = (event) => {
        this.handleMessage(event.data)
      }

      this.ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason)
        this.isConnected = false
        this.stopPing()
        this.handleReconnect(token)
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
      this.handleReconnect(token)
    }
  }

  disconnect() {
    this.stopPing()
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.isConnected = false
    this.reconnectAttempts = this.maxReconnectAttempts // Prevent reconnect
  }

  private handleMessage(data: string) {
    try {
      const message: WSMessage = JSON.parse(data)
      console.log('WebSocket message received:', message.type)

      const store = useContactsStore()

      switch (message.type) {
        case WS_TYPE_NEW_MESSAGE:
          this.handleNewMessage(store, message.payload)
          break
        case WS_TYPE_STATUS_UPDATE:
          this.handleStatusUpdate(store, message.payload)
          break
        case WS_TYPE_PONG:
          // Pong received, connection is alive
          break
        default:
          console.log('Unknown message type:', message.type)
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error)
    }
  }

  private handleNewMessage(store: ReturnType<typeof useContactsStore>, payload: any) {
    // Check if this message is for the current contact
    const currentContact = store.currentContact
    if (currentContact && payload.contact_id === currentContact.id) {
      // Add message to the store
      store.addMessage({
        id: payload.id,
        contact_id: payload.contact_id,
        direction: payload.direction,
        message_type: payload.message_type,
        content: payload.content,
        status: payload.status,
        wamid: payload.wamid,
        error_message: payload.error_message,
        created_at: payload.created_at,
        updated_at: payload.updated_at
      })
    }

    // Update contacts list (for unread count, last message preview)
    store.fetchContacts()
  }

  private handleStatusUpdate(store: ReturnType<typeof useContactsStore>, payload: any) {
    store.updateMessageStatus(payload.message_id, payload.status)
  }

  private handleReconnect(token: string) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnect attempts reached')
      return
    }

    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)

    setTimeout(() => {
      this.connect(token)
    }, delay)
  }

  setCurrentContact(contactId: string | null) {
    this.send({
      type: WS_TYPE_SET_CONTACT,
      payload: { contact_id: contactId || '' }
    })
  }

  private send(message: WSMessage) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    }
  }

  private startPing() {
    this.stopPing()
    this.pingInterval = window.setInterval(() => {
      this.send({ type: WS_TYPE_PING, payload: {} })
    }, 30000) // Ping every 30 seconds
  }

  private stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }

  getIsConnected() {
    return this.isConnected
  }
}

// Export singleton instance
export const wsService = new WebSocketService()
