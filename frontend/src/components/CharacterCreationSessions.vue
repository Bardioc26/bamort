<template>
  <div v-if="sessions.length > 0" class="sessions-section">
    <div class="section-header">
      <h3>{{ $t('characters.list.continue_creation') }}</h3>
    </div>
    <div class="grid-container grid-2-columns">
      <div 
        v-for="session in sessions" 
        :key="session.session_id"
        class="card session-card"
        @click="continueSession(session.session_id)"
      >
        <div class="session-header">
          <h4 class="list-item-title">{{ session.name || $t('characters.list.unnamed_character') }}</h4>
          <span class="badge badge-primary progress-badge">{{ $t('characters.list.step') }} {{ session.current_step }}/{{ session.total_steps }}</span>
        </div>
        <div class="session-details">
          <p class="list-item-details"><strong>{{ $t('characters.list.race') }}:</strong> {{ session.rasse || $t('characters.list.not_selected') }}
          <span class="list-item-separator">|</span><strong>{{ $t('characters.list.class') }}:</strong> {{ session.typ || $t('characters.list.not_selected') }} </p> 
          <p class="list-item-details"><strong>{{ $t('characters.list.current_step') }}:</strong> {{ session.progress_text }} </p>
        </div>
        <div class="session-meta">
          <span class="session-date">{{ $t('characters.list.last_updated') }}: {{ formatDate(session.updated_at) }}</span>
          <span class="session-date">{{ $t('characters.list.expires') }}: {{ formatDate(session.expires_at) }}</span>
        </div>
        <div class="list-item-actions">
          <button @click.stop="deleteSession(session.session_id)" class="btn btn-danger btn-small">
            {{ $t('characters.list.delete_draft') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { formatDate } from '@/utils/dateUtils'

export default {
  name: 'CharacterCreationSessions',
  props: {
    sessions: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    continueSession(sessionId) {
      this.$emit('continue-session', sessionId)
    },
    
    deleteSession(sessionId) {
      this.$emit('delete-session', sessionId)
    },
    
    formatDate
  }
}
</script>

<style scoped>
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

.btn-small {
  padding: 5px 10px;
  font-size: 0.8rem;
}

/* Responsive Design */
@media (max-width: 768px) {
  .session-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
</style>
