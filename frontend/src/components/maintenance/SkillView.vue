<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
    <div class="search-box">
      <input
        type="text"
        v-model="searchTerm"
        :placeholder="`${$t('search')} ${$t('Skill')}...`"
      />
      <button @click="startCreate" class="btn-primary">{{ $t('newSkill') }}</button>
    </div>
  </div>

  <div class="cd-view">
    <div class="cd-list">
      <!-- Filter Row -->
      <div class="filter-row">
        <div class="filter-item">
          <label>{{ $t('skill.category') }}:</label>
          <select v-model="filterCategory">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="cat in availableCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('skill.difficulty') }}:</label>
          <select v-model="filterDifficulty">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="diff in availableDifficulties" :key="diff" :value="diff">{{ diff }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('skill.improvable') }}:</label>
          <select v-model="filterImprovable">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option value="true">{{ $t('yes') || 'Yes' }}</option>
            <option value="false">{{ $t('no') || 'No' }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('skill.innateskill') }}:</label>
          <select v-model="filterInnateskill">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option value="true">{{ $t('yes') || 'Yes' }}</option>
            <option value="false">{{ $t('no') || 'No' }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('skill.bonusskill') }}:</label>
          <select v-model="filterBonuseigenschaft">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="bonus in availableBonuseigenschaften" :key="bonus" :value="bonus">{{ bonus }}</option>
          </select>
        </div>
        <button @click="clearFilters" class="btn-clear-filters">{{ $t('clearFilters') || 'Clear Filters' }}</button>
      </div>

      <div class="tables-container">
        <table class="cd-table">
          <thead>
            <tr>
              <th class="cd-table-header">{{ $t('skill.id') }}</th>
              <th class="cd-table-header">
                {{ $t('skill.category') }}
                <button @click="sortBy('category')">{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '-' }}</button>
              </th>
              <th class="cd-table-header">
                {{ $t('skill.name') }}
                <button @click="sortBy('name')">{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '-' }}</button>
              </th>
              <th class="cd-table-header">{{ $t('skill.difficulty') }}</th>
              <th class="cd-table-header">{{ $t('skill.initialwert') }}</th>
              <th class="cd-table-header">{{ $t('skill.basiswert') }}</th>
              <th class="cd-table-header">{{ $t('skill.improvable') }}</th>
              <th class="cd-table-header">{{ $t('skill.innateskill') }}</th>
              <th class="cd-table-header">{{ $t('skill.description') }}</th>
              <th class="cd-table-header">{{ $t('skill.bonusskill') }}</th>
              <th class="cd-table-header">{{ $t('skill.quelle') }}&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>
              <th class="cd-table-header">{{ $t('skill.system') }}</th>
              <th class="cd-table-header"> </th>
            </tr>
          </thead>
          <tbody>
            <!-- Create New Skill Row -->
            <tr v-if="creatingNew">
              <td>New</td>
              <td colspan="11">
                <!-- Create form -->
                <div class="edit-form">
                  <div class="edit-row">
                    <div class="edit-field">
                      <label>{{ $t('skill.name') }}:</label>
                      <input v-model="editedItem.name" />
                    </div>
                    <div class="edit-field">
                      <label>{{ $t('skill.initialwert') }}:</label>
                      <input v-model.number="editedItem.initialwert" type="number" style="width:60px;" />
                    </div>
                    <div class="edit-field">
                      <label>{{ $t('skill.basiswert') }}:</label>
                      <input v-model.number="editedItem.basiswert" type="number" style="width:60px;" />
                    </div>
                    <div class="edit-field">
                      <label>{{ $t('skill.bonusskill') }}:</label>
                      <select v-model="editedItem.bonuseigenschaft" style="width:80px;">
                        <option value="">-</option>
                        <option value="St">St</option>
                        <option value="Gs">Gs</option>
                        <option value="Gw">Gw</option>
                        <option value="Ko">Ko</option>
                        <option value="In">In</option>
                        <option value="Zt">Zt</option>
                        <option value="Au">Au</option>
                        <option value="pA">pA</option>
                        <option value="Wk">Wk</option>
                        <option value="B">B</option>
                      </select>
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field">
                      <label>{{ $t('skill.improvable') }}:</label>
                      <input type="checkbox" v-model="editedItem.improvable" />
                    </div>
                    <div class="edit-field">
                      <label>{{ $t('skill.innateskill') }}:</label>
                      <input type="checkbox" v-model="editedItem.innateskill" />
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field full-width">
                      <label>{{ $t('skill.description') }}:</label>
                      <input v-model="editedItem.beschreibung" />
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field">
                      <label>{{ $t('skill.quelle') }}:</label>
                      <select v-model="editedItem.sourceCode" style="width:100px;">
                        <option v-for="source in availableSources" :key="source.code" :value="source.code">
                          {{ source.code }}
                        </option>
                      </select>
                    </div>
                    <div class="edit-field">
                      <label>{{ $t('skill.page') || 'Page' }}:</label>
                      <input v-model.number="editedItem.page_number" type="number" style="width:60px;" />
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field">
                      <label>{{ $t('skill.system') }}</label>
                      <select v-model.number="selectedSystemId" style="width:140px;">
                        <option value="">-</option>
                        <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                          {{ system.label }}
                        </option>
                      </select>
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field">
                      <label>{{ $t('skill.system') }}</label>
                      <select v-model.number="createSelectedSystemId" style="width:140px;">
                        <option value="">-</option>
                        <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                          {{ system.label }}
                        </option>
                      </select>
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field full-width">
                      <label>{{ $t('skill.categories') || 'Categories' }}:</label>
                      <div class="category-checkboxes">
                        <div v-for="category in mdata.categories" :key="category.id" class="category-checkbox">
                          <input 
                            type="checkbox" 
                            :value="category.id"
                            v-model="editedItem.selectedCategories"
                            @change="onCategoryToggle(category.id)"
                          />
                          <label>{{ category.name }}</label>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="edit-row">
                    <div class="edit-field full-width">
                      <label>{{ $t('skill.difficulty') || 'Difficulty' }}:</label>
                      <div class="difficulty-selects">
                        <div v-for="catId in editedItem.selectedCategories" :key="catId" class="difficulty-select">
                          <span>{{ getCategoryName(catId) }}:</span>
                          <select v-model="editedItem.categoryDifficulties[catId]" style="width:120px;">
                            <option v-for="diff in mdata.difficulties" :key="diff.id" :value="diff.id">
                              {{ diff.name }}
                            </option>
                          </select>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="edit-actions">
                    <button @click="saveCreate" class="btn-save">{{ $t('createSkill') }}</button>
                    <button @click="cancelCreate" class="btn-cancel">Cancel</button>
                  </div>
                </div>
              </td>
            </tr>

            <template v-for="(dtaItem, index) in filteredAndSortedSkills" :key="dtaItem.id">
              <!-- Display Mode -->
              <tr v-if="editingIndex !== index">
                <td>{{ dtaItem.id || '' }}</td>
                <td>{{ formatCategories(dtaItem.categories) }}</td>
                <td>{{ dtaItem.name || '-' }}</td>
                <td>{{ formatDifficulties(dtaItem.difficulties) }}</td>
                <td>{{ dtaItem.initialwert || '0' }}</td>
                <td>{{ dtaItem.basiswert || '0' }}</td>
                <td><input type="checkbox" :checked="dtaItem.improvable" disabled /></td>
                <td><input type="checkbox" :checked="dtaItem.innateskill" disabled /></td>
                <td>{{ dtaItem.beschreibung || '-' }}</td>
                <td>{{ dtaItem.bonuseigenschaft || '-' }}</td>
                <td>{{ formatQuelle(dtaItem) }}</td>
                <td>{{ getSystemCodeById(dtaItem.game_system_id, dtaItem.game_system || 'midgard') }}</td>
                <td>
                  <button @click="startEdit(index)">Edit</button>
                </td>
              </tr>
              
              <!-- Edit Mode -->
              <tr v-else>
                <td><input v-model="editedItem.id" style="width:20px;" disabled /></td>
                <td colspan="11">
                  <!-- Expanded edit form -->
                  <div class="edit-form">
                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('skill.name') }}:</label>
                        <input v-model="editedItem.name" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('skill.initialwert') }}:</label>
                        <input v-model.number="editedItem.initialwert" type="number" style="width:60px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('skill.basiswert') }}:</label>
                        <input v-model.number="editedItem.basiswert" type="number" style="width:60px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('skill.bonusskill') }}:</label>
                        <select v-model="editedItem.bonuseigenschaft" style="width:80px;">
                          <option value="">-</option>
                          <option value="St">St</option>
                          <option value="Gs">Gs</option>
                          <option value="Gw">Gw</option>
                          <option value="Ko">Ko</option>
                          <option value="In">In</option>
                          <option value="Zt">Zt</option>
                          <option value="Au">Au</option>
                          <option value="pA">pA</option>
                          <option value="Wk">Wk</option>
                          <option value="B">B</option>
                        </select>
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('skill.improvable') }}:</label>
                        <input type="checkbox" v-model="editedItem.improvable" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('skill.innateskill') }}:</label>
                        <input type="checkbox" v-model="editedItem.innateskill" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field full-width">
                        <label>{{ $t('skill.description') }}:</label>
                        <input v-model="editedItem.beschreibung" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('skill.quelle') }}:</label>
                        <select v-model="editedItem.sourceCode" style="width:100px;">
                          <option v-for="source in availableSources" :key="source.code" :value="source.code">
                            {{ source.code }}
                          </option>
                        </select>
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('skill.page') || 'Page' }}:</label>
                        <input v-model.number="editedItem.page_number" type="number" style="width:60px;" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field full-width">
                        <label>{{ $t('skill.categories') || 'Categories' }}:</label>
                        <div class="category-checkboxes">
                          <div v-for="category in mdata.categories" :key="category.id" class="category-checkbox">
                            <input 
                              type="checkbox" 
                              :value="category.id"
                              v-model="editedItem.selectedCategories"
                              @change="onCategoryToggle(category.id)"
                            />
                            <label>{{ category.name }}</label>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field full-width">
                        <label>{{ $t('skill.difficulty') || 'Difficulty' }}:</label>
                        <div class="difficulty-selects">
                          <div v-for="catId in editedItem.selectedCategories" :key="catId" class="difficulty-select">
                            <span>{{ getCategoryName(catId) }}:</span>
                            <select v-model="editedItem.categoryDifficulties[catId]" style="width:120px;">
                              <option v-for="diff in mdata.difficulties" :key="diff.id" :value="diff.id">
                                {{ diff.name }}
                              </option>
                            </select>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div class="edit-actions">
                      <button @click="saveEdit(index)" class="btn-save">Save</button>
                      <button @click="cancelEdit" class="btn-cancel">Cancel</button>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style>
/* All styles moved to main.css as per project conventions */
</style>

<script>
import API from '../../utils/api'

export default {
  name: "SkillView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        skills: [],
        categories: [],
        difficulties: [],
        sources: []
      })
    }
  },
  data() {
    return {
      searchTerm: '',
      sortField: 'name',
      sortAsc: true,
      editingIndex: -1,
      editedItem: null,
      creatingNew: false,
      filterCategory: '',
      filterDifficulty: '',
      filterImprovable: '',
      filterInnateskill: '',
      filterBonuseigenschaft: '',
      enhancedSkills: [],
      availableSources: [],
      gameSystems: [],
      selectedSystemId: null,
      createSelectedSystemId: null
    }
  },
  async created() {
    await Promise.all([
      this.loadGameSystems(),
      this.loadEnhancedSkills()
    ])
  },
  computed: {
    availableCategories() {
      const categories = new Set()
      this.enhancedSkills.forEach(skill => {
        if (skill.categories) {
          skill.categories.forEach(cat => categories.add(cat.category_name))
        }
      })
      return Array.from(categories).sort()
    },
    availableDifficulties() {
      const difficulties = new Set()
      this.enhancedSkills.forEach(skill => {
        if (skill.difficulties) {
          skill.difficulties.forEach(diff => difficulties.add(diff))
        }
      })
      return Array.from(difficulties).sort()
    },
    availableBonuseigenschaften() {
      const bonuses = new Set()
      this.enhancedSkills.forEach(skill => {
        if (skill.bonuseigenschaft) {
          bonuses.add(skill.bonuseigenschaft)
        }
      })
      return Array.from(bonuses).sort()
    },
    filteredAndSortedSkills() {
      let filtered = [...this.enhancedSkills]

      // Apply search filter
      if (this.searchTerm) {
        const searchLower = this.searchTerm.toLowerCase()
        filtered = filtered.filter(skill => 
          skill.name?.toLowerCase().includes(searchLower) ||
          this.formatCategories(skill.categories).toLowerCase().includes(searchLower)
        )
      }

      // Apply category filter
      if (this.filterCategory) {
        filtered = filtered.filter(skill =>
          skill.categories && skill.categories.some(cat => cat.category_name === this.filterCategory)
        )
      }

      // Apply difficulty filter
      if (this.filterDifficulty) {
        filtered = filtered.filter(skill =>
          skill.difficulties && skill.difficulties.includes(this.filterDifficulty)
        )
      }

      // Apply improvable filter
      if (this.filterImprovable !== '') {
        const improvableValue = this.filterImprovable === 'true'
        filtered = filtered.filter(skill => skill.improvable === improvableValue)
      }

      // Apply innateskill filter
      if (this.filterInnateskill !== '') {
        const innateskillValue = this.filterInnateskill === 'true'
        filtered = filtered.filter(skill => skill.innateskill === innateskillValue)
      }

      // Apply bonuseigenschaft filter
      if (this.filterBonuseigenschaft) {
        filtered = filtered.filter(skill => skill.bonuseigenschaft === this.filterBonuseigenschaft)
      }

      // Apply sorting
      filtered.sort((a, b) => {
        let aValue, bValue

        if (this.sortField === 'category') {
          aValue = this.formatCategories(a.categories).toLowerCase()
          bValue = this.formatCategories(b.categories).toLowerCase()
        } else {
          aValue = (a[this.sortField] || '').toString().toLowerCase()
          bValue = (b[this.sortField] || '').toString().toLowerCase()
        }

        return this.sortAsc ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue)
      })

      return filtered
    },
    systemOptions() {
      return this.gameSystems.map(gs => ({
        id: gs.id,
        label: gs.name ? `${gs.code} (${gs.name})` : gs.code
      }))
    }
  },
  methods: {
    async loadGameSystems() {
      try {
        const resp = await API.get('/api/maintenance/game-systems')
        this.gameSystems = (resp.data?.game_systems || []).map(this.normalizeSystem)
      } catch (error) {
        console.error('Failed to load game systems:', error)
      }
    },
    normalizeSystem(gs) {
      return {
        id: gs.id ?? gs.ID,
        code: gs.code ?? gs.Code,
        name: gs.name ?? gs.Name,
        is_active: gs.is_active ?? gs.IsActive,
      }
    },
    async loadEnhancedSkills() {
      try {
        const response = await API.get('/api/maintenance/skills-enhanced')
        this.enhancedSkills = response.data.skills || []
        this.availableSources = response.data.sources || []
        // Also update mdata for compatibility
        if (response.data.categories) {
          this.mdata.categories = response.data.categories
        }
        if (response.data.difficulties) {
          this.mdata.difficulties = response.data.difficulties
        }
      } catch (error) {
        console.error('Failed to load enhanced skills:', error)
      }
    },
    formatCategories(categories) {
      if (!categories || categories.length === 0) return '-'
      return categories.map(cat => cat.category_name).join(', ')
    },
    formatDifficulties(difficulties) {
      if (!difficulties || difficulties.length === 0) return '-'
      return difficulties.join(', ')
    },
    formatQuelle(skill) {
      if (skill.source_id && this.availableSources.length > 0) {
        const source = this.availableSources.find(s => s.id === skill.source_id)
        if (source) {
          if (skill.page_number) {
            return `${source.code}:${skill.page_number}`
          } else {
            // No page number - show code and quelle if available
            const quelle = skill.quelle ? ` (${skill.quelle})` : ''
            return `${source.code}${quelle}`
          }
        }
      }
      return skill.quelle || '-'
    },
    getCategoryName(categoryId) {
      const category = this.mdata.categories.find(c => c.id === categoryId)
      return category ? category.name : `ID:${categoryId}`
    },
    startEdit(index) {
      const skill = this.filteredAndSortedSkills[index]
      
      // Initialize edit object
      this.editedItem = {
        ...skill,
        selectedCategories: skill.categories ? skill.categories.map(cat => cat.category_id) : [],
        categoryDifficulties: {},
        sourceCode: this.getSourceCode(skill.source_id),
      }

      // Map category difficulties
      if (skill.categories) {
        skill.categories.forEach(cat => {
          this.editedItem.categoryDifficulties[cat.category_id] = cat.difficulty_id
        })
      }

      this.selectedSystemId = skill.game_system_id ?? this.findSystemIdByCode(skill.game_system)

      this.editingIndex = index
    },
    getSourceCode(sourceId) {
      if (!sourceId || !this.availableSources.length) return ''
      const source = this.availableSources.find(s => s.id === sourceId)
      return source ? source.code : ''
    },
    findSystemIdByCode(code) {
      if (!code) return null
      const match = this.gameSystems.find(gs => gs.code === code)
      return match ? match.id : null
    },
    getSystemCodeById(systemId, fallback = '') {
      if (!systemId) return fallback
      const sys = this.gameSystems.find(gs => gs.id === systemId)
      return sys ? sys.code : fallback
    },
    onCategoryToggle(categoryId) {
      // If category was removed, also remove its difficulty setting
      if (!this.editedItem.selectedCategories.includes(categoryId)) {
        delete this.editedItem.categoryDifficulties[categoryId]
      } else {
        // Set a default difficulty if not already set
        if (!this.editedItem.categoryDifficulties[categoryId] && this.mdata.difficulties.length > 0) {
          // Find "normal" difficulty or use first one
          const normalDiff = this.mdata.difficulties.find(d => d.name.toLowerCase() === 'normal')
          this.editedItem.categoryDifficulties[categoryId] = normalDiff ? normalDiff.id : this.mdata.difficulties[0].id
        }
      }
    },
    async saveEdit(index) {
      try {
        // Find source ID from code
        const source = this.availableSources.find(s => s.code === this.editedItem.sourceCode)
        const selectedSystem = this.gameSystems.find(gs => gs.id === this.selectedSystemId)
        
        // Build category_difficulties array
        const categoryDifficulties = this.editedItem.selectedCategories.map(catId => ({
          category_id: catId,
          difficulty_id: this.editedItem.categoryDifficulties[catId]
        }))

        const updateData = {
          id: this.editedItem.id,
          name: this.editedItem.name,
          beschreibung: this.editedItem.beschreibung,
          game_system: selectedSystem ? selectedSystem.code : (this.editedItem.game_system || 'midgard'),
          game_system_id: selectedSystem ? selectedSystem.id : (this.editedItem.game_system_id ?? null),
          initialwert: this.editedItem.initialwert,
          basiswert: this.editedItem.basiswert || 0,
          bonuseigenschaft: this.editedItem.bonuseigenschaft,
          improvable: this.editedItem.improvable,
          innateskill: this.editedItem.innateskill,
          source_id: source ? source.id : null,
          page_number: this.editedItem.page_number || 0,
          category_difficulties: categoryDifficulties
        }

        const response = await API.put(
          `/api/maintenance/skills-enhanced/${this.editedItem.id}`,
          updateData
        )

        // Update the skill in the list using splice for proper reactivity
        const skillIndex = this.enhancedSkills.findIndex(s => s.id === this.editedItem.id)
        if (skillIndex !== -1) {
          this.enhancedSkills.splice(skillIndex, 1, response.data)
        }

        this.editingIndex = -1
        this.editedItem = null
        this.selectedSystemId = null
      } catch (error) {
        console.error('Failed to update skill:', error)
        alert('Failed to update skill: ' + (error.response?.data?.error || error.message))
      }
    },
    cancelEdit() {
      this.editingIndex = -1
      this.editedItem = null
    },
    sortBy(field) {
      if (this.sortField === field) {
        this.sortAsc = !this.sortAsc
      } else {
        this.sortField = field
        this.sortAsc = true
      }
    },
    clearFilters() {
      this.searchTerm = ''
      this.filterCategory = ''
      this.filterDifficulty = ''
      this.filterImprovable = ''
      this.filterInnateskill = ''
      this.filterBonuseigenschaft = ''
    },
    startCreate() {
      // Initialize new skill object with defaults
      const defaultSystem = this.gameSystems.find(gs => gs.is_active) || this.gameSystems[0] || null
      this.createSelectedSystemId = defaultSystem ? defaultSystem.id : null
      this.editedItem = {
        name: '',
        beschreibung: '',
        game_system: defaultSystem ? defaultSystem.code : 'midgard',
        game_system_id: defaultSystem ? defaultSystem.id : null,
        initialwert: 5,
        basiswert: 0,
        bonuseigenschaft: '',
        improvable: true,
        innateskill: false,
        sourceCode: this.availableSources.length > 0 ? this.availableSources[0].code : '',
        page_number: 0,
        selectedCategories: [],
        categoryDifficulties: {}
      }
      this.creatingNew = true
    },
    async saveCreate() {
      try {
        // Validate required fields
        if (!this.editedItem.name) {
          alert('Skill name is required')
          return
        }

        // Find source ID from code
        const source = this.availableSources.find(s => s.code === this.editedItem.sourceCode)
        const selectedSystem = this.gameSystems.find(gs => gs.id === this.createSelectedSystemId)
        
        // Build category_difficulties array
        const categoryDifficulties = this.editedItem.selectedCategories.map(catId => ({
          category_id: catId,
          difficulty_id: this.editedItem.categoryDifficulties[catId]
        }))

        const createData = {
          name: this.editedItem.name,
          beschreibung: this.editedItem.beschreibung,
          game_system: selectedSystem ? selectedSystem.code : (this.editedItem.game_system || 'midgard'),
          game_system_id: selectedSystem ? selectedSystem.id : (this.editedItem.game_system_id ?? null),
          initialwert: this.editedItem.initialwert,
          basiswert: this.editedItem.basiswert || 0,
          bonuseigenschaft: this.editedItem.bonuseigenschaft,
          improvable: this.editedItem.improvable,
          innateskill: this.editedItem.innateskill,
          source_id: source ? source.id : null,
          page_number: this.editedItem.page_number || 0,
          category_difficulties: categoryDifficulties
        }

        const response = await API.post(
          '/api/maintenance/skills-enhanced',
          createData
        )

        // Add the new skill to the list
        this.enhancedSkills.push(response.data)

        // Hide the create dialog
        this.creatingNew = false
        this.editedItem = null
        this.createSelectedSystemId = null
      } catch (error) {
        console.error('Failed to create skill:', error)
        alert('Failed to create skill: ' + (error.response?.data?.error || error.message))
        // Don't close dialog on error so user can fix the issue
      }
    },
    cancelCreate() {
      this.creatingNew = false
      this.editedItem = null
      this.createSelectedSystemId = null
    }
  }
}
</script>
