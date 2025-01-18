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
  margin-bottom: 0.3rem;
  height: fit-content;
  padding: 0.5rem;
}
.search-box {
  margin-bottom: 1rem;
}
.search-box input {
  padding: 0.2rem;
  width: 200px;
  border: 1px solid #ddd;
  border-radius: 4px;
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
      editedItem: null
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
      }
  }
};
</script>
