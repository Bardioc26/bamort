<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
    <!-- Add search input -->
    <div class="search-box">
      <input
        type="text"
        v-model="searchTerm"
        :placeholder="`${$t('search')} ${$t('Spell')}...`"
      />
    </div>
  </div>
  
  <!-- Import CSV Section -->
  <div class="import-section">
    <h3>{{ $t('spell.import') || 'Import Spells' }}</h3>
    <div class="file-upload">
      <input
        ref="fileInput"
        type="file"
        accept=".csv"
        @change="handleFileSelect"
        style="display: none;"
      />
      <button 
        @click="$refs.fileInput.click()"
        :disabled="importing"
        class="upload-btn"
      >
        {{ importing ? ($t('spell.importing') || 'Importing...') : ($t('spell.selectCsv') || 'Select CSV File') }}
      </button>
      <span v-if="selectedFile" class="file-name">{{ selectedFile.name }}</span>
    </div>
    <button 
      v-if="selectedFile && !importing"
      @click="importSpells"
      class="import-btn"
    >
      {{ $t('spell.import') || 'Import Spells' }}
    </button>
    <div v-if="importResult" class="import-result" :class="importResult.success ? 'success' : 'error'">
      {{ importResult.message }}
      <span v-if="importResult.total_spells"> ({{ importResult.total_spells }} spells total)</span>
    </div>
  </div>

  <div class="cd-view">
    <div class="cd-list">
      <div class="tables-container">
          <table class="cd-table">
            <thead>
              <tr>
                <th class="cd-table-header">{{ $t('spell.id') }}</th>
                <th class="cd-table-header">{{ $t('spell.category') }}<button @click="sortBy('category')">-{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('spell.name') }} <button @click="sortBy('name')">-{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('spell.level') }}</th>
                <th class="cd-table-header">{{ $t('spell.apverbrauch') }}</th>
                <th class="cd-table-header">{{ $t('spell.zauberdauer') }}</th>
                <th class="cd-table-header">{{ $t('spell.reichweite') }}</th>
                <th class="cd-table-header">{{ $t('spell.wirkungsziel') }}</th>
                <th class="cd-table-header">{{ $t('spell.wirkungsbereich') }}</th>
                <th class="cd-table-header">{{ $t('spell.wirkungsdauer') }}</th>
                <th class="cd-table-header">{{ $t('spell.ursprung') }}</th>
                <th class="cd-table-header">{{ $t('spell.description') }}</th>
                <th class="cd-table-header">{{ $t('spell.quelle') }}</th>
                <th class="cd-table-header">{{ $t('spell.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(dtaItem, index) in filteredAndSortedSpells" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <td>{{ dtaItem.category|| '-' }}</td>
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.level || '0' }}</td>
                  <td>{{ dtaItem.ap || '0' }}</td>
                  <td>{{ dtaItem.zauberdauer || '-' }}</td>
                  <td>{{ dtaItem.reichweite || '0' }}</td>
                  <td>{{ dtaItem.wirkungsziel || '-' }}</td>
                  <td>{{ dtaItem.wirkungsbereich || '-' }}</td>
                  <td>{{ dtaItem.wirkungsdauer || '-' }}</td>
                  <td>{{ dtaItem.ursprung || '-' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ dtaItem.quelle || '-' }}</td>
                  <td>{{ dtaItem.system || 'midgard' }}</td>
                  <td>
                    <button @click="startEdit(index)">Edit</button>
                  </td>
                </tr>
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;"/></td>
                  <td><select v-model="editedItem.category" style="width:80px;">
                    <option v-for="category in mdata['spellcategories']"
                            :key="category"
                            :value="category">
                      {{ category }}
                    </option>
                  </select></td>
                  <td><input v-model="editedItem.name"/></td>
                  <td><input v-model.number="editedItem.level" type="number" style="width:40px;"/></td>
                  <td><input v-model="editedItem.ap" style="width:40px;"/></td>
                  <td><input v-model="editedItem.zauberdauer" /></td>
                  <td><input v-model="editedItem.reichweite" style="width:40px;"/></td>
                  <td><input v-model="editedItem.wirkungsziel" /></td>
                  <td><input v-model="editedItem.wirkungsbereich" /></td>
                  <td><input v-model="editedItem.wirkungsdauer" /></td>
                  <td><input v-model="editedItem.ursprung" /></td>
                  <td><input v-model="editedItem.beschreibung" /></td>
                  <td><input v-model="editedItem.quelle"  style="width:80px;"/></td>
                  <td><input v-model="editedItem.system"  style="width:80px;"/></td>
                  <td>
                    <button @click="saveEdit(index)">Save</button>
                    <button @click="cancelEdit">Cancel</button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
      </div>
    </div> <!--- end cd-list-->
  </div> <!--- end character -datasheet-->
</template>

<!-- <style scoped> -->
<style>
.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  height: fit-content;
  padding: 0.5rem;
  border-bottom: 1px solid #ddd;
}

.search-box {
  margin-left: auto;
}

.search-box input {
  padding: 0.5rem;
  width: 250px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.import-section {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
  border: 2px solid #1da766;
  border-radius: 8px;
  background-color: #f8fcfa;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.import-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.file-upload {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.upload-btn, .import-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #1da766;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s;
}

.upload-btn:hover, .import-btn:hover {
  background-color: #166d4a;
}

.upload-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.file-name {
  font-style: italic;
  color: #666;
}

.import-result {
  padding: 0.5rem;
  border-radius: 4px;
  margin-top: 1rem;
}

.import-result.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.import-result.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
.tables-container {
  display: flex;
  gap: 1rem;
  width: 100%;
}

.table-wrapper-left {
  flex: 6;
  min-width: 0; /* Prevent table from overflowing */
}
.table-wrapper-right {
  flex: 4;
  min-width: 0; /* Prevent table from overflowing */
}

.cd-table {
  width: 100%;
}
.cd-table-header {
  background-color: #1da766;
  font-weight: bold;
}
</style>


<script>
import API from '../../utils/api'
export default {
  name: "SpellView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        spells: [],
        spellcategories: []
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
      selectedFile: null,
      importing: false,
      importResult: null
    }
  },
  computed: {
    filteredAndSortedSpells() {
      if (!this.mdata?.spells) return [];

      return [...this.mdata.spells]
        .filter(spell => {
          const searchLower = this.searchTerm.toLowerCase();
          return !this.searchTerm ||
            spell.name?.toLowerCase().includes(searchLower) ||
            spell.category?.toLowerCase().includes(searchLower);
        })
        .sort((a, b) => {
          const aValue = (a[this.sortField] || '').toLowerCase();
          const bValue = (b[this.sortField] || '').toLowerCase();
          return this.sortAsc
            ? aValue.localeCompare(bValue)
            : bValue.localeCompare(aValue);
        });
    },
    sortedSpells() {
      return [...this.mdata.spells].sort((a, b) => {
        const aValue = (a[this.sortField] || '').toLowerCase();
        const bValue = (b[this.sortField] || '').toLowerCase();
        return this.sortAsc
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      });
    }
  },
  methods: {
    startEdit(index) {
      this.editingIndex = index;
      this.editedItem = { ...this.filteredAndSortedSpells[index] };
    },
    saveEdit(index) {
      //this.$emit('update-spell', { index, spell: this.editedItem });
      this.handleSpellUpdate( { index, spell: this.editedItem });
      this.editingIndex = -1;
      this.editedItem = null;
    },
    cancelEdit() {
      this.editingIndex = -1;
      this.editedItem = null;
    },
    sortBy(field) {
      if (this.sortField === field) {
        this.sortAsc = !this.sortAsc;
      } else {
        this.sortField = field;
        this.sortAsc = true;
      }
    },
    async handleSpellUpdate({ index, spell }) {
      try {
          const response = await API.put(
            `/api/maintenance/spells/${spell.id}`, spell,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}` ,
                'Content-Type': 'application/json'
              }
            }
          )
          if (!response.statusText== "OK") throw new Error('Update failed');
          const updatedSkill = response.data;
          // Update the spell in mdata
          this.mdata.spells = this.mdata.spells.map(s =>
            s.id === updatedSkill.id ? updatedSkill : s
          );
        } catch (error) {
          console.error('Failed to update spell:', error);
        }
      },
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (file && file.type === 'text/csv') {
        this.selectedFile = file;
        this.importResult = null;
      } else {
        this.selectedFile = null;
        this.importResult = {
          success: false,
          message: 'Please select a valid CSV file.'
        };
      }
    },
    async importSpells() {
      if (!this.selectedFile) {
        this.importResult = {
          success: false,
          message: 'Please select a CSV file first.'
        };
        return;
      }

      this.importing = true;
      this.importResult = null;

      try {
        const formData = new FormData();
        formData.append('file', this.selectedFile);

        const response = await API.post('/api/importer/spells/csv', formData, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'multipart/form-data'
          }
        });

        this.importResult = {
          success: true,
          message: response.data.message,
          total_spells: response.data.total_spells
        };

        // Refresh the spells data after successful import
        await this.refreshSpellsData();

      } catch (error) {
        console.error('Failed to import spells:', error);
        this.importResult = {
          success: false,
          message: error.response?.data?.message || 'Import failed. Please try again.'
        };
      } finally {
        this.importing = false;
      }
    },
    async refreshSpellsData() {
      try {
        const token = localStorage.getItem('token');
        const response = await API.get('/api/maintenance', {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        // Update the spells data
        if (response.data.spells) {
          this.mdata.spells = response.data.spells;
        }
      } catch (error) {
        console.error('Failed to refresh spells data:', error);
      }
    }
  }
};
</script>
