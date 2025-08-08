<template>
  <div class="fullwidth-container">
    <div class="page-header">
      <h2>{{ $t('characters.list.title') }}</h2>
    </div>
    
    <!-- Create New Character Button -->
    <div class="create-character-section">
      <button @click="createNewCharacter" class="btn btn-success btn-large">{{ $t('characters.list.create_new') }}</button>
    </div>

    <!-- Active Character Creation Sessions -->
    <CharacterCreationSessions 
      :sessions="creationSessions"
      @continue-session="continueSession"
      @delete-session="handleDeleteSession"
    />
    
    <div v-if="characters.length === 0" class="empty-state">
      <h3>{{ $t('characters.list.no_characters') }}</h3>
      <p>{{ $t('characters.list.no_characters_description') }}</p>
    </div>
    
    <div v-else class="list-container">
      <div v-for="character in characters" :key="character.character_id" class="list-item">
        <div class="list-item-content">
          <h4 class="list-item-title">{{ character.name }}</h4>
          <div class="list-item-details">
            {{ character.rasse }} <span class="list-item-separator">|</span>
            {{ character.typ }} <span class="list-item-separator">|</span>
            {{ $t('characters.list.grade') }}: {{ character.grad }} <span class="list-item-separator">|</span>
            {{ $t('characters.list.owner') }}: {{ character.owner }} <span class="list-item-separator">|</span>
            <span class="badge" :class="character.public ? 'badge-success' : 'badge-secondary'">
              {{ character.public ? $t('characters.list.public') : $t('characters.list.private') }}
            </span>
          </div>
        </div>
        <div class="list-item-actions">
          <router-link :to="`/character/${character.id}`" class="btn btn-primary">{{ $t('characters.list.view_details') }}</router-link>
          <button @click="goToAusruestung(character.character_id)" class="btn btn-secondary">{{ $t('characters.list.manage_equipment') }}</button>
        </div>
      </div>
    </div>
  </div>
</template><script>
import API from '../utils/api'
import { formatDate } from '@/utils/dateUtils'
import CharacterCreationSessions from './CharacterCreationSessions.vue'

export default {
  components: {
    CharacterCreationSessions
  },
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
    
    handleDeleteSession(sessionId) {
      this.deleteSession(sessionId)
    },
    
    async deleteSession(sessionId) {
      if (confirm(this.$t('characters.list.delete_draft_confirm'))) {
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
    
    formatDate
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

/* Responsive Design */
@media (max-width: 768px) {
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
