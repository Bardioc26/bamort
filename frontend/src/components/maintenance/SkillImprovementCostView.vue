<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }} - {{ $t('skillimprovement.title') }}</h2>
    <button class="btn-primary" @click="startCreate">{{ $t('newEntry') }}</button>
  </div>

  <div v-if="error" class="error-box">{{ error }}</div>

  <div class="cd-view">
    <div class="cd-list">
      <table class="cd-table">
        <thead>
          <tr>
            <th class="cd-table-header">{{ $t('skillimprovement.id') }}</th>
            <th class="cd-table-header">{{ $t('skillimprovement.level') }}</th>
            <th class="cd-table-header">{{ $t('skillimprovement.te') }}</th>
            <th class="cd-table-header">{{ $t('skillimprovement.category') }}</th>
            <th class="cd-table-header">{{ $t('skillimprovement.difficulty') }}</th>
            <th class="cd-table-header"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="6">{{ $t('common.loading') }}</td>
          </tr>
          <tr v-if="creatingNew">
            <td>New</td>
            <td colspan="5">
              <div class="edit-form">
                <div class="edit-row">
                  <label>{{ $t('skillimprovement.level') }}</label>
                  <input v-model.number="newItem.current_level" type="number" />
                  <label class="inline-label">{{ $t('skillimprovement.te') }}</label>
                  <input v-model.number="newItem.te_required" type="number" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('skillimprovement.category') }}</label>
                  <select v-model.number="newItem.category_id">
                    <option v-for="cat in categoryOptions" :key="cat.id" :value="cat.id">
                      {{ cat.label }}
                    </option>
                  </select>
                  <label class="inline-label">{{ $t('skillimprovement.difficulty') }}</label>
                  <select v-model.number="newItem.difficulty_id">
                    <option v-for="diff in difficultyOptions" :key="diff.id" :value="diff.id">
                      {{ diff.label }}
                    </option>
                  </select>
                </div>
                <div class="edit-actions">
                  <button class="btn-primary btn-save" :disabled="isSaving" @click="saveCreate">
                    <span v-if="!isSaving">{{ $t('common.save') }}</span>
                    <span v-else>{{ $t('common.saving') }}</span>
                  </button>
                  <button class="btn-cancel" :disabled="isSaving" @click="cancelCreate">
                    {{ $t('common.cancel') }}
                  </button>
                </div>
              </div>
            </td>
          </tr>
          <template v-for="cost in costs" :key="cost.id">
            <tr v-if="editingId !== cost.id">
              <td>{{ cost.id }}</td>
              <td>{{ cost.current_level }}</td>
              <td>{{ cost.te_required }}</td>
              <td>{{ displayCategory(cost) }}</td>
              <td>{{ displayDifficulty(cost) }}</td>
              <td><button @click="startEdit(cost)">{{ $t('common.edit') }}</button></td>
            </tr>
            <tr v-else>
              <td>{{ cost.id }}</td>
              <td colspan="5">
                <div class="edit-form">
                  <div class="edit-row">
                    <label>{{ $t('skillimprovement.level') }}</label>
                    <input v-model.number="editedItem.current_level" type="number" />
                    <label class="inline-label">{{ $t('skillimprovement.te') }}</label>
                    <input v-model.number="editedItem.te_required" type="number" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('skillimprovement.category') }}</label>
                    <select v-model.number="editedItem.category_id">
                      <option v-for="cat in categoryOptions" :key="cat.id" :value="cat.id">
                        {{ cat.label }}
                      </option>
                    </select>
                    <label class="inline-label">{{ $t('skillimprovement.difficulty') }}</label>
                    <select v-model.number="editedItem.difficulty_id">
                      <option v-for="diff in difficultyOptions" :key="diff.id" :value="diff.id">
                        {{ diff.label }}
                      </option>
                    </select>
                  </div>
                  <div class="edit-actions">
                    <button class="btn-primary btn-save" :disabled="isSaving" @click="saveEdit">
                      <span v-if="!isSaving">{{ $t('common.save') }}</span>
                      <span v-else>{{ $t('common.saving') }}</span>
                    </button>
                    <button class="btn-cancel" :disabled="isSaving" @click="cancelEdit">
                      {{ $t('common.cancel') }}
                    </button>
                  </div>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.error-box {
  margin: 10px 0;
  padding: 10px 12px;
  background: #ffe3e3;
  color: #8a1c1c;
  border: 1px solid #f5c2c2;
  border-radius: 6px;
}
.edit-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.edit-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}
.edit-actions {
  display: flex;
  gap: 10px;
}
.inline-label {
  margin-left: 10px;
}
</style>

<script>
import API from '../../utils/api'

export default {
  name: 'SkillImprovementCostView',
  data() {
    return {
      costs: [],
      editingId: null,
      editedItem: null,
      creatingNew: false,
      newItem: null,
      isLoading: false,
      isSaving: false,
      error: '',
    }
  },
  async created() {
    await this.loadCosts()
  },
  computed: {
    categoryOptions() {
      const seen = new Map()
      this.costs.forEach(c => {
        const id = c.category_id ?? c.skillCategoryId
        const name = c.category_name || c.skillCategoryName
        if (id != null && !seen.has(id)) {
          seen.set(id, name ? `${name} (${id})` : `${id}`)
        }
      })
      return Array.from(seen.entries()).map(([id, label]) => ({ id, label }))
    },
    difficultyOptions() {
      const seen = new Map()
      this.costs.forEach(c => {
        const id = c.difficulty_id ?? c.skillDifficultyId
        const name = c.difficulty_name || c.skillDifficultyName
        if (id != null && !seen.has(id)) {
          seen.set(id, name ? `${name} (${id})` : `${id}`)
        }
      })
      return Array.from(seen.entries()).map(([id, label]) => ({ id, label }))
    },
  },
  methods: {
    displayCategory(cost) {
      return cost.category_name || cost.skillCategoryName || cost.skillCategoryId || cost.category_id
    },
    displayDifficulty(cost) {
      return cost.difficulty_name || cost.skillDifficultyName || cost.skillDifficultyId || cost.difficulty_id
    },
    async loadCosts() {
      this.isLoading = true
      this.error = ''
      try {
        const resp = await API.get('/api/maintenance/skill-improvement-cost2')
        this.costs = resp.data?.costs || []
      } catch (err) {
        console.error('Failed to load costs:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isLoading = false
      }
    },
    startEdit(cost) {
      this.editingId = cost.id
      this.editedItem = {
        ...cost,
        category_id: cost.category_id ?? cost.skillCategoryId,
        difficulty_id: cost.difficulty_id ?? cost.skillDifficultyId,
      }
    },
    cancelEdit() {
      this.editingId = null
      this.editedItem = null
    },
    startCreate() {
      this.cancelEdit()
      const defaultCategory = this.categoryOptions[0]?.id ?? null
      const defaultDifficulty = this.difficultyOptions[0]?.id ?? null
      this.newItem = {
        current_level: 0,
        te_required: 0,
        category_id: defaultCategory,
        difficulty_id: defaultDifficulty,
      }
      this.creatingNew = true
    },
    cancelCreate() {
      this.creatingNew = false
      this.newItem = null
    },
    async saveEdit() {
      if (!this.editedItem) return
      const payload = {
        current_level: this.editedItem.current_level,
        te_required: this.editedItem.te_required,
        category_id: this.editedItem.category_id,
        difficulty_id: this.editedItem.difficulty_id,
      }
      this.isSaving = true
      try {
        const resp = await API.put(`/api/maintenance/skill-improvement-cost2/${this.editingId}`, payload)
        const idx = this.costs.findIndex(c => c.id === this.editingId)
        if (idx !== -1) this.costs.splice(idx, 1, resp.data)
        this.cancelEdit()
      } catch (err) {
        console.error('Failed to save cost:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
    async saveCreate() {
      if (!this.newItem) return
      const payload = {
        current_level: this.newItem.current_level,
        te_required: this.newItem.te_required,
        category_id: this.newItem.category_id,
        difficulty_id: this.newItem.difficulty_id,
      }
      this.isSaving = true
      try {
        const resp = await API.post('/api/maintenance/skill-improvement-cost2', payload)
        this.costs.push(resp.data)
        this.cancelCreate()
      } catch (err) {
        console.error('Failed to create cost:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
  },
}
</script>
