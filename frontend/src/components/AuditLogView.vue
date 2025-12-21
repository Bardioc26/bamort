<template>
  <div class="audit-log-view">
    <h4>{{ $t('audit.title', 'Änderungsprotokoll') }}</h4>
    
    <!-- Filter Controls -->
    <div class="filter-controls">
      <div class="filter-group">
        <label>{{ $t('audit.filter_by_field', 'Filter nach Feld') }}:</label>
        <select v-model="selectedField" @change="loadAuditLog" class="filter-select">
          <option value="">{{ $t('audit.all_fields', 'Alle Felder') }}</option>
          <option value="experience_points">{{ $t('audit.experience_points', 'Erfahrungspunkte') }}</option>
          <option value="gold">{{ $t('audit.gold', 'Gold') }}</option>
          <option value="silver">{{ $t('audit.silver', 'Silber') }}</option>
          <option value="copper">{{ $t('audit.copper', 'Kupfer') }}</option>
        </select>
      </div>
      
      <div class="filter-group">
        <label>{{ $t('audit.filter_by_date', 'Zeitraum') }}:</label>
        <select v-model="selectedDateRange" @change="loadAuditLog" class="filter-select">
          <option value="">{{ $t('audit.all_time', 'Alle Zeit') }}</option>
          <option value="today">{{ $t('audit.today', 'Heute') }}</option>
          <option value="week">{{ $t('audit.this_week', 'Diese Woche') }}</option>
          <option value="month">{{ $t('audit.this_month', 'Dieser Monat') }}</option>
          <option value="custom">{{ $t('audit.custom_range', 'Benutzerdefiniert') }}</option>
        </select>
      </div>
      
      <div v-if="selectedDateRange === 'custom'" class="date-range-group">
        <input 
          v-model="customDateFrom" 
          type="date" 
          @change="loadAuditLog"
          class="date-input"
        />
        <span>bis</span>
        <input 
          v-model="customDateTo" 
          type="date" 
          @change="loadAuditLog"
          class="date-input"
        />
      </div>
      
      <div class="filter-group">
        <label>
          <input 
            v-model="groupByDate" 
            type="checkbox"
            class="checkbox-input"
          />
          {{ $t('audit.group_by_date', 'Nach Datum gruppieren') }}
        </label>
      </div>
      
      <button @click="loadAuditLog" class="btn-refresh" :disabled="isLoading">
        <span v-if="isLoading">⏳</span>
        <span v-else>🔄</span>
        {{ $t('audit.refresh', 'Aktualisieren') }}
      </button>
    </div>

    <!-- Statistics -->
    <div v-if="stats" class="stats-section">
      <h5>{{ $t('audit.statistics', 'Statistiken') }}</h5>
      <div class="stats-grid">
        <div class="stat-item">
          <span class="stat-label">{{ $t('audit.total_changes', 'Gesamte Änderungen') }}:</span>
          <span class="stat-value">{{ stats.total_changes }}</span>
        </div>
        <div class="stat-item ep-stat">
          <span class="stat-label">{{ $t('audit.ep_spent', 'EP ausgegeben') }}:</span>
          <span class="stat-value negative">{{ stats.total_ep_spent }}</span>
        </div>
        <div class="stat-item ep-stat">
          <span class="stat-label">{{ $t('audit.ep_gained', 'EP erhalten') }}:</span>
          <span class="stat-value positive">{{ stats.total_ep_gained }}</span>
        </div>
        <div class="stat-item gold-stat">
          <span class="stat-label">{{ $t('audit.gold_spent', 'Gold ausgegeben') }}:</span>
          <span class="stat-value negative">{{ stats.total_gold_spent }}</span>
        </div>
        <div class="stat-item gold-stat">
          <span class="stat-label">{{ $t('audit.gold_gained', 'Gold erhalten') }}:</span>
          <span class="stat-value positive">{{ stats.total_gold_gained }}</span>
        </div>
      </div>
    </div>

    <!-- Audit Log Entries -->
    <div class="audit-entries">
      <div v-if="isLoading" class="loading">
        {{ $t('audit.loading', 'Lädt...') }}
      </div>
      
      <div v-else-if="auditEntries.length === 0" class="no-entries">
        {{ $t('audit.no_entries', 'Keine Änderungen gefunden') }}
      </div>
      
      <div v-else>
        <div v-if="groupByDate">
          <div v-for="(entries, date) in groupedEntries" :key="date" class="date-group">
            <h6 class="date-group-header">{{ formatDateHeader(date) }}</h6>
            <div 
              v-for="entry in entries" 
              :key="entry.id" 
              class="audit-entry"
              :class="[
                entry.difference > 0 ? 'positive-change' : 'negative-change',
                `field-${entry.field_name}`
              ]"
            >
              <div class="entry-header">
                <div class="entry-field">
                  <span class="field-icon">{{ getFieldIcon(entry.field_name) }}</span>
                  <span class="field-name">{{ getFieldDisplayName(entry.field_name) }}</span>
                </div>
                <div class="entry-timestamp">
                  <div class="timestamp-time">{{ formatTime(entry.timestamp) }}</div>
                  <div class="timestamp-relative">{{ formatRelativeTime(entry.timestamp) }}</div>
                </div>
              </div>
              
              <div class="entry-content">
                <div class="value-change">
                  <span class="old-value">{{ entry.old_value }}</span>
                  <span class="arrow">→</span>
                  <span class="new-value">{{ entry.new_value }}</span>
                  <span class="difference" :class="entry.difference > 0 ? 'positive' : 'negative'">
                    ({{ entry.difference > 0 ? '+' : '' }}{{ entry.difference }})
                  </span>
                </div>
                
                <div class="entry-reason">
                  <span class="reason-label">{{ $t('audit.reason', 'Grund') }}:</span>
                  <span class="reason-value">{{ getReasonDisplayName(entry.reason) }}</span>
                </div>
                
                <div v-if="entry.notes" class="entry-notes">
                  <span class="notes-label">{{ $t('audit.notes', 'Notizen') }}:</span>
                  <span class="notes-value">{{ entry.notes }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div v-else>
          <div 
            v-for="entry in auditEntries" 
            :key="entry.id" 
            class="audit-entry"
            :class="[
              entry.difference > 0 ? 'positive-change' : 'negative-change',
              `field-${entry.field_name}`
            ]"
          >
            <div class="entry-header">
              <div class="entry-field">
                <span class="field-icon">{{ getFieldIcon(entry.field_name) }}</span>
                <span class="field-name">{{ getFieldDisplayName(entry.field_name) }}</span>
              </div>
              <div class="entry-timestamp">
                <div class="timestamp-date">{{ formatDate(entry.timestamp) }}</div>
                <div class="timestamp-time">{{ formatTime(entry.timestamp) }}</div>
                <div class="timestamp-relative">{{ formatRelativeTime(entry.timestamp) }}</div>
              </div>
            </div>
            
            <div class="entry-content">
              <div class="value-change">
                <span class="old-value">{{ entry.old_value }}</span>
                <span class="arrow">→</span>
                <span class="new-value">{{ entry.new_value }}</span>
                <span class="difference" :class="entry.difference > 0 ? 'positive' : 'negative'">
                  ({{ entry.difference > 0 ? '+' : '' }}{{ entry.difference }})
                </span>
              </div>
              
              <div class="entry-reason">
                <span class="reason-label">{{ $t('audit.reason', 'Grund') }}:</span>
                <span class="reason-value">{{ getReasonDisplayName(entry.reason) }}</span>
              </div>
              
              <div v-if="entry.notes" class="entry-notes">
                <span class="notes-label">{{ $t('audit.notes', 'Notizen') }}:</span>
                <span class="notes-value">{{ entry.notes }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import API from '@/utils/api'

export default {
  name: "AuditLogView",
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      auditEntries: [],
      stats: null,
      selectedField: '',
      selectedDateRange: '',
      customDateFrom: '',
      customDateTo: '',
      groupByDate: true,
      isLoading: false
    };
  },
  
  computed: {
    groupedEntries() {
      if (!this.groupByDate) return {};
      
      const grouped = {};
      this.auditEntries.forEach(entry => {
        const date = new Date(entry.timestamp).toLocaleDateString('de-DE');
        if (!grouped[date]) {
          grouped[date] = [];
        }
        grouped[date].push(entry);
      });
      
      return grouped;
    }
  },
  created() {
    this.$api = API;
    this.loadAuditLog();
    this.loadStats();
  },
  methods: {
    async loadAuditLog() {
      if (!this.character?.id) return;
      
      this.isLoading = true;
      try {
        let url = `/api/characters/${this.character.id}/audit-log`;
        const params = new URLSearchParams();
        
        if (this.selectedField) {
          params.append('field', this.selectedField);
        }
        
        // Datumsfilter
        if (this.selectedDateRange) {
          const now = new Date();
          let fromDate, toDate;
          
          switch (this.selectedDateRange) {
            case 'today':
              fromDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
              toDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59);
              break;
            case 'week':
              const dayOfWeek = now.getDay();
              const mondayOffset = dayOfWeek === 0 ? 6 : dayOfWeek - 1; // Montag als Wochenbeginn
              fromDate = new Date(now.getFullYear(), now.getMonth(), now.getDate() - mondayOffset);
              toDate = new Date();
              break;
            case 'month':
              fromDate = new Date(now.getFullYear(), now.getMonth(), 1);
              toDate = new Date();
              break;
            case 'custom':
              if (this.customDateFrom) {
                fromDate = new Date(this.customDateFrom);
              }
              if (this.customDateTo) {
                toDate = new Date(this.customDateTo);
                toDate.setHours(23, 59, 59, 999); // Ende des Tages
              }
              break;
          }
          
          if (fromDate) {
            params.append('from', fromDate.toISOString());
          }
          if (toDate) {
            params.append('to', toDate.toISOString());
          }
        }
        
        if (params.toString()) {
          url += '?' + params.toString();
        }
          
        const response = await this.$api.get(url);
        this.auditEntries = response.data.entries || [];
      } catch (error) {
        console.error('Fehler beim Laden des Audit-Logs:', error);
        this.auditEntries = [];
      } finally {
        this.isLoading = false;
      }
    },
    
    async loadStats() {
      if (!this.character?.id) return;
      
      try {
        const response = await this.$api.get(`/api/characters/${this.character.id}/audit-log/stats`);
        this.stats = response.data.stats;
      } catch (error) {
        console.error('Fehler beim Laden der Statistiken:', error);
        this.stats = null;
      }
    },
    
    getFieldIcon(fieldName) {
      const icons = {
        'experience_points': '⭐',
        'gold': '💰',
        'silver': '🥈',
        'copper': '🥉'
      };
      return icons[fieldName] || '📝';
    },
    
    getFieldDisplayName(fieldName) {
      const names = {
        'experience_points': 'Erfahrungspunkte',
        'gold': 'Goldstücke',
        'silver': 'Silberstücke',
        'copper': 'Kupferstücke'
      };
      return names[fieldName] || fieldName;
    },
    
    getReasonDisplayName(reason) {
      const reasons = {
        'manual': 'Manuell',
        'skill_learning': 'Fertigkeit lernen',
        'skill_improvement': 'Fertigkeit verbessern',
        'spell_learning': 'Zauber lernen',
        'spell_improvement': 'Zauber verbessern',
        'equipment': 'Ausrüstung',
        'reward': 'Belohnung',
        'correction': 'Korrektur',
        'import': 'Import'
      };
      return reasons[reason] || reason;
    },
    
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleString('de-DE', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    },
    
    formatDate(timestamp) {
      return new Date(timestamp).toLocaleDateString('de-DE', {
        weekday: 'short',
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      });
    },
    
    formatTime(timestamp) {
      return new Date(timestamp).toLocaleTimeString('de-DE', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    },
    
    formatRelativeTime(timestamp) {
      const now = new Date();
      const past = new Date(timestamp);
      const diffInSeconds = Math.floor((now - past) / 1000);
      
      if (diffInSeconds < 60) {
        return 'gerade eben';
      } else if (diffInSeconds < 3600) {
        const minutes = Math.floor(diffInSeconds / 60);
        return `vor ${minutes} Min.`;
      } else if (diffInSeconds < 86400) {
        const hours = Math.floor(diffInSeconds / 3600);
        return `vor ${hours} Std.`;
      } else if (diffInSeconds < 604800) {
        const days = Math.floor(diffInSeconds / 86400);
        return `vor ${days} Tag${days > 1 ? 'en' : ''}`;
      } else {
        const weeks = Math.floor(diffInSeconds / 604800);
        return `vor ${weeks} Woche${weeks > 1 ? 'n' : ''}`;
      }
    },
    
    formatDateHeader(dateString) {
      const date = new Date(dateString.split('.').reverse().join('-')); // Umwandlung von dd.mm.yyyy zu yyyy-mm-dd
      const today = new Date();
      const yesterday = new Date(today);
      yesterday.setDate(yesterday.getDate() - 1);
      
      if (date.toLocaleDateString('de-DE') === today.toLocaleDateString('de-DE')) {
        return 'Heute (' + dateString + ')';
      } else if (date.toLocaleDateString('de-DE') === yesterday.toLocaleDateString('de-DE')) {
        return 'Gestern (' + dateString + ')';
      } else {
        return date.toLocaleDateString('de-DE', {
          weekday: 'long',
          day: '2-digit',
          month: 'long',
          year: 'numeric'
        });
      }
    }
  },
  
  watch: {
    'character.id'() {
      this.loadAuditLog();
      this.loadStats();
    }
  }
};
</script>

<style scoped>
.audit-log-view {
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-top: 20px;
}

.audit-log-view h4 {
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #007bff;
}

.filter-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.date-range-group {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: 10px;
}

.filter-select, .date-input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
}

.btn-refresh {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
}

.btn-refresh:hover:not(:disabled) {
  background: #0056b3;
}

.stats-section {
  margin-bottom: 25px;
  padding: 15px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.stats-section h5 {
  margin-bottom: 15px;
  color: #555;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
}

.stat-label {
  font-weight: 500;
  color: #666;
}

.stat-value {
  font-weight: bold;
  font-size: 1.1em;
}

.stat-value.positive {
  color: #28a745;
}

.stat-value.negative {
  color: #dc3545;
}

.audit-entries {
  background: white;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.loading, .no-entries {
  padding: 40px;
  text-align: center;
  color: #666;
  font-style: italic;
}

.audit-entry {
  padding: 15px;
  border-bottom: 1px solid #e9ecef;
  transition: background-color 0.2s;
}

.audit-entry:last-child {
  border-bottom: none;
}

.audit-entry:hover {
  background-color: #f8f9fa;
}

.audit-entry.positive-change {
  border-left: 4px solid #28a745;
}

.audit-entry.negative-change {
  border-left: 4px solid #dc3545;
}

.entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.entry-field {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
}

.field-icon {
  font-size: 1.2em;
}

.entry-timestamp {
  color: #666;
  font-size: 0.9em;
  text-align: right;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.timestamp-date {
  font-weight: 500;
  color: #555;
}

.timestamp-time {
  font-size: 0.85em;
  color: #888;
  font-family: monospace;
}

.timestamp-relative {
  font-size: 0.8em;
  color: #999;
  font-style: italic;
}

.entry-content {
  margin-left: 20px;
}

.value-change {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  font-family: monospace;
  font-size: 1.1em;
}

.old-value {
  color: #666;
}

.arrow {
  color: #007bff;
  font-weight: bold;
}

.new-value {
  font-weight: bold;
}

.difference {
  font-weight: bold;
  font-size: 0.9em;
}

.difference.positive {
  color: #28a745;
}

.difference.negative {
  color: #dc3545;
}

.entry-reason, .entry-notes {
  display: flex;
  gap: 8px;
  margin-bottom: 5px;
  font-size: 0.9em;
}

.reason-label, .notes-label {
  font-weight: 500;
  color: #666;
  min-width: 50px;
}

.reason-value {
  background: #e9ecef;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.85em;
  font-weight: 500;
}

.notes-value {
  color: #555;
  font-style: italic;
}

.date-group {
  margin-bottom: 25px;
}

.date-group-header {
  background: #007bff;
  color: white;
  padding: 8px 15px;
  margin: 0 0 10px 0;
  border-radius: 4px;
  font-weight: 500;
  font-size: 0.9em;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.checkbox-input {
  margin-right: 8px;
}
</style>
