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
    
    <div v-if="ownedCharacters.length === 0" class="empty-state">
      <h3>{{ $t('characters.list.no_characters') }}</h3>
      <p>{{ $t('characters.list.no_characters_description') }}</p>
    </div>
    
    <div v-else class="list-container horizontal-placement">
      <div class="charlist">
        {{ $t('characters.list.owned_characters_title') }}
        <div v-for="character in ownedCharacters" :key="character.character_id" class="list-item">
          <router-link :to="`/character/${character.id}`" class="list-item-content">
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
          </router-link>
        </div>
      </div>
      <div class="charlist">
        {{ $t('characters.list.shared_characters_title') }}
        <div v-for="character in sharedCharacters" :key="character.character_id" class="list-item">
          <router-link :to="`/character/${character.id}`" class="list-item-content">
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
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import API from '../utils/api'
import { formatDate } from '@/utils/dateUtils'
import CharacterCreationSessions from './CharacterCreationSessions.vue'

export default {
  components: {
    CharacterCreationSessions
  },
  data() {
    return {
      ownedCharacters: [],
      sharedCharacters: [],
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
        this.ownedCharacters = response.data.self_owned
        this.sharedCharacters = response.data.others
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
/* All common styles moved to main.css */

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

/* RouterLink als list-item-content styling */
.list-item-content {
  flex: 1;
  text-decoration: none;
  color: inherit;
  display: block;
}

.list-item-content:hover {
  text-decoration: none;
  color: inherit;
}

.list-item-content:hover .list-item-title {
  color: #007bff;
}

.horizontal-placement {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.charlist {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
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
