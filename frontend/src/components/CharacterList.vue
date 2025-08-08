<template>
  <div class="fullwidth-container">
    <div class="page-header">
      <h2>Your Characters</h2>
    </div>
    
    <!-- Create New Character Button -->
    <div class="create-character-section">
      <button @click="createNewCharacter" class="btn btn-success btn-large">Create New Character</button>
    </div>

    <!-- Active Character Creation Sessions -->
    <div v-if="creationSessions.length > 0" class="sessions-section">
      <div class="section-header">
        <h3>Continue Character Creation</h3>
      </div>
      <div class="grid-container grid-2-columns">
        <div 
          v-for="session in creationSessions" 
          :key="session.session_id"
          class="card session-card"
          @click="continueSession(session.session_id)"
        >
          <div class="session-header">
            <h4 class="list-item-title">{{ session.name || 'Unnamed Character' }}</h4>
            <span class="badge badge-primary progress-badge">Step {{ session.current_step }}/{{ session.total_steps }}</span>
          </div>
          <div class="session-details">
            <p class="list-item-details"><strong>Race:</strong> {{ session.rasse || 'Not selected' }}
            <span class="list-item-separator">|</span><strong>Class:</strong> {{ session.typ || 'Not selected' }} 
            <span class="list-item-separator">|</span><strong>Current step:</strong> {{ session.progress_text }}</p>
          </div>
          <div class="session-meta">
            <span class="session-date">Last updated: {{ formatDate(session.updated_at) }}</span>
            <span class="session-date">Expires: {{ formatDate(session.expires_at) }}</span>
          </div>
          <div class="list-item-actions">
            <button @click.stop="deleteSession(session.session_id)" class="btn btn-danger btn-small">
              Delete Draft
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="characters.length === 0" class="empty-state">
      <h3>No Characters Yet</h3>
      <p>Create your first character to get started!</p>
    </div>
    
    <div v-else class="list-container">
      <div v-for="character in characters" :key="character.character_id" class="list-item">
        <div class="list-item-content">
          <h4 class="list-item-title">{{ character.name }}</h4>
          <div class="list-item-details">
            {{ character.rasse }} <span class="list-item-separator">|</span>
            {{ character.typ }} <span class="list-item-separator">|</span>
            {{ character.grad }} <span class="list-item-separator">|</span>
            {{ character.owner }} <span class="list-item-separator">|</span>
            <span class="badge" :class="character.public ? 'badge-success' : 'badge-secondary'">
              {{ character.public ? 'Public' : 'Private' }}
            </span>
          </div>
        </div>
        <div class="list-item-actions">
          <router-link :to="`/character/${character.id}`" class="btn btn-primary">View Details</router-link>
          <button @click="goToAusruestung(character.character_id)" class="btn btn-secondary">Manage Equipment</button>
        </div>
      </div>
    </div>
  </div>
</template><script>
import API from '../utils/api'

export default {
  data() {
    return {
      characters: [],
      creationSessions: [],
    }
  },
  async created() {
    await this.loadCharacters()
    await this.loadCreationSessions()
  },
  methods: {
    async loadCharacters() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.get('/api/characters', {
          headers: { Authorization: `Bearer ${token}` },
        })
        this.characters = response.data
      } catch (error) {
        console.error('Error loading characters:', error)
      }
    },
    
    async loadCreationSessions() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.get('/api/characters/create-sessions', {
          headers: { Authorization: `Bearer ${token}` },
        })
        this.creationSessions = response.data.sessions || []
      } catch (error) {
        console.error('Error loading creation sessions:', error)
        // Don't show error to user since this is a new feature
        this.creationSessions = []
      }
    },
    
    continueSession(sessionId) {
      this.$router.push(`/character/create/${sessionId}`)
    },
    
    async deleteSession(sessionId) {
      if (confirm('Are you sure you want to delete this character draft?')) {
        try {
          const token = localStorage.getItem('token')
          await API.delete(`/api/characters/create-session/${sessionId}`, {
            headers: { Authorization: `Bearer ${token}` },
          })
          
          // Reload sessions after deletion
          await this.loadCreationSessions()
        } catch (error) {
          console.error('Error deleting session:', error)
          alert('Error deleting character draft')
        }
      }
    },
    
    goToAusruestung(characterId) {
      this.$router.push(`/api/ausruestung/${characterId}`)
    },
    async createNewCharacter() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.post('/api/characters/create-session', {}, {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        const sessionId = response.data.session_id
        this.$router.push(`/character/create/${sessionId}`)
      } catch (error) {
        console.error('Error creating character session:', error)
        alert('Fehler beim Erstellen der Charakter-Session')
      }
    },
    
    formatDate(dateString) {
      if (!dateString) return 'Unknown'
      return new Date(dateString).toLocaleDateString()
    },
  },
}
</script>

<style scoped>
/* Spezifische Styles nur für CharacterList */
.create-character-section {
  margin-bottom: 30px;
  padding: 20px;
  border: 2px dashed #28a745;
  border-radius: 8px;
  text-align: center;
  background-color: #f8fff9;
}

.btn-large {
  padding: 12px 24px;
  font-size: 16px;
}

.btn-small {
  padding: 5px 10px;
  font-size: 0.8rem;
}

.sessions-section {
  margin-bottom: 30px;
}

.session-card {
  cursor: pointer;
}

.session-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.session-details {
  margin-bottom: 10px;
}

.session-meta {
  display: flex;
  flex-direction: column;
  gap: 5px;
  margin-bottom: 10px;
  padding-top: 10px;
  border-top: 1px solid #eee;
}

.session-date {
  font-size: 0.8rem;
  color: #888;
}

/* Responsive Design */
@media (max-width: 768px) {
  .session-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
  
  .list-item {
    flex-direction: column;
    gap: 10px;
  }
  
  .list-item-actions {
    align-self: stretch;
    justify-content: flex-start;
  }
}
</style>
