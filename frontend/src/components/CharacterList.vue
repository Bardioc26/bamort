<template>
  <div>
    <h2>Your Characters</h2>
    
    <!-- Create New Character Button -->
    <div class="create-character-section">
      <button @click="createNewCharacter" class="create-btn">Create New Character</button>
    </div>

    <!-- Active Character Creation Sessions -->
    <div v-if="creationSessions.length > 0" class="creation-sessions-section">
      <h3>Continue Character Creation</h3>
      <div class="sessions-grid">
        <div 
          v-for="session in creationSessions" 
          :key="session.session_id"
          class="session-card"
          @click="continueSession(session.session_id)"
        >
          <div class="session-header">
            <h4>{{ session.name || 'Unnamed Character' }}</h4>
            <span class="session-progress">Step {{ session.current_step }}/{{ session.total_steps }}</span>
          </div>
          <div class="session-details">
            <p><strong>Race:</strong> {{ session.rasse || 'Not selected' }}</p>
            <p><strong>Class:</strong> {{ session.typ || 'Not selected' }}</p>
            <p><strong>Current step:</strong> {{ session.progress_text }}</p>
          </div>
          <div class="session-meta">
            <span class="last-updated">Last updated: {{ formatDate(session.updated_at) }}</span>
            <span class="expires">Expires: {{ formatDate(session.expires_at) }}</span>
          </div>
          <div class="session-actions">
            <button @click.stop="deleteSession(session.session_id)" class="delete-session-btn">
              Delete Draft
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <ul>
      <li v-for="character in characters" :key="character.character_id" style="white-space: nowrap; /* Prevent line breaks inside list items */;">
        <!-- Link to Character Details -->
        <router-link :to="`/character/${character.id}`">View Details</router-link>
        {{ character.name }} ({{ character.rasse }}, {{ character.typ }}, {{ character.grad }}, {{ character.owner }}, {{ character.public }} )
        <button @click="goToAusruestung(character.character_id)">Manage Equipment</button>
      </li>
    </ul>
  </div>
</template>

<script>
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
.create-character-section {
  margin-bottom: 20px;
  padding: 15px;
  border: 2px dashed #ccc;
  border-radius: 8px;
  text-align: center;
  background-color: #f9f9f9;
}

.create-btn {
  background-color: #4CAF50;
  color: white;
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.create-btn:hover {
  background-color: #45a049;
}

.creation-sessions-section {
  margin-bottom: 30px;
}

.creation-sessions-section h3 {
  color: #333;
  margin-bottom: 15px;
  font-size: 1.2rem;
}

.sessions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
  margin-bottom: 20px;
}

.session-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.session-card:hover {
  border-color: #007bff;
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
  transform: translateY(-2px);
}

.session-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.session-header h4 {
  margin: 0;
  color: #333;
  font-size: 1.1rem;
}

.session-progress {
  background: #007bff;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: bold;
}

.session-details {
  margin-bottom: 10px;
}

.session-details p {
  margin: 5px 0;
  font-size: 0.9rem;
  color: #666;
}

.session-meta {
  display: flex;
  flex-direction: column;
  gap: 5px;
  margin-bottom: 10px;
  padding-top: 10px;
  border-top: 1px solid #eee;
}

.last-updated,
.expires {
  font-size: 0.8rem;
  color: #888;
}

.session-actions {
  display: flex;
  justify-content: flex-end;
}

.delete-session-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.delete-session-btn:hover {
  background: #c82333;
}

/* Responsive Design */
@media (max-width: 768px) {
  .sessions-grid {
    grid-template-columns: 1fr;
  }
  
  .session-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
</style>
