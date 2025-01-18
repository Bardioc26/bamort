<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
      <!-- Add search input -->
      <div class="search-box">
        <input
          type="text"
          v-model="searchTerm"
          :placeholder="`${$t('search')} ${$t('Skill')}...`"
        />
      </div>
    </div>

  <div class="cd-view">
    <div class="cd-list">
      <div class="tables-container">
          <table class="cd-table">
            <thead>
              <tr>
                <th class="cd-table-header">{{ $t('skill.id') }}</th>
                <th class="cd-table-header">{{ $t('skill.category') }}<button @click="sortBy('category')">-{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('skill.name') }} <button @click="sortBy('name')">-{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('skill.initialwert') }}</th>
                <th class="cd-table-header">{{ $t('skill.improvable') }}</th>
                <th class="cd-table-header">{{ $t('skill.innateskill') }}</th>
                <th class="cd-table-header">{{ $t('skill.description') }}</th>
                <th class="cd-table-header">{{ $t('skill.bonusskill') }}</th>
                <th class="cd-table-header">{{ $t('skill.quelle') }}</th>
                <th class="cd-table-header">{{ $t('skill.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(dtaItem, index) in filteredAndSortedSkills" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <td>{{ dtaItem.category|| '-' }}</td>
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.initialwert || '0' }}</td>
                  <td>{{ dtaItem.improvable || '0' }}</td>
                  <td>{{ dtaItem.innateskill || '0' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ dtaItem.bonuseigenschaft || '-' }}</td>
                  <td>{{ dtaItem.quelle || '-' }}</td>
                  <td>{{ dtaItem.system || 'midgard' }}</td>
                  <td>
                    <button @click="startEdit(index)">Edit</button>
                  </td>
                </tr>
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;"/></td>
                  <td><select v-model="editedItem.category" style="width:80px;">
                    <option v-for="category in mdata['skillcategories']"
                            :key="category"
                            :value="category">
                      {{ category }}
                    </option>
                  </select></td>
                  <td><input v-model="editedItem.name"/></td>
                  <td><input v-model.number="editedItem.initialwert" type="number" style="width:40px;"/></td>
                  <td><input type="checkbox" :checked="true" v-model="editedItem.improvable" style="width:50px;"/></td>
                  <td><input type="checkbox" :checked="true" v-model="editedItem.innateskill" style="width:50px;"/></td>
                  <td><input v-model="editedItem.beschreibung" /></td>
                  <td><input v-model="editedItem.bonuseigenschaft" style="width:80px;" ></td>
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
  name: "SkillView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        skills: [],
        skillcategories: []
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
    filteredAndSortedSkills() {
      if (!this.mdata?.skills) return [];

      return [...this.mdata.skills]
        .filter(skill => {
          const searchLower = this.searchTerm.toLowerCase();
          return !this.searchTerm ||
            skill.name?.toLowerCase().includes(searchLower) ||
            skill.category?.toLowerCase().includes(searchLower);
        })
        .sort((a, b) => {
          const aValue = (a[this.sortField] || '').toLowerCase();
          const bValue = (b[this.sortField] || '').toLowerCase();
          return this.sortAsc
            ? aValue.localeCompare(bValue)
            : bValue.localeCompare(aValue);
        });
    },
    sortedSkills() {
      return [...this.mdata.skills].sort((a, b) => {
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
      this.editedItem = { ...this.filteredAndSortedSkills[index] };
    },
    saveEdit(index) {
      //this.$emit('update-skill', { index, skill: this.editedItem });
      this.handleSkillUpdate({ index, skill: this.editedItem });
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
    async handleSkillUpdate({ index, skill }) {
      try {
        const response = await API.put(
          `/api/maintenance/skills/${skill.id}`, skill,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}` ,
              'Content-Type': 'application/json'
            }
          }
        )
        if (!response.statusText== "OK") throw new Error('Update failed');
        const updatedSkill = response.data;
        // Update the skill in mdata
        this.mdata.skills = this.mdata.skills.map(s =>
          s.id === updatedSkill.id ? updatedSkill : s
        );
      } catch (error) {
        console.error('Failed to update skill:', error);
      }
    }
  }
};
</script>
