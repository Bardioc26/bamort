<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
      <!-- Add search input -->
      <div class="search-box">
        <input
          type="text"
          v-model="searchTerm"
          :placeholder="`${$t('search')} ${$t('Weapons')}...`"
        />
      </div>
    </div>

  <div class="cd-view">
    <div class="cd-list">
      <div class="tables-container">
          <table class="cd-table">
            <thead>
              <tr>
                <th class="cd-table-header">{{ $t('weapon.id') }}</th>
                <!-- <th class="cd-table-header">{{ $t('weapon.category') }}<button @click="sortBy('category')">-{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '' }}</button></th> -->
                <th class="cd-table-header">{{ $t('weapon.name') }} <button @click="sortBy('name')">-{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('weapon.gewicht') }}</th>
                <th class="cd-table-header">{{ $t('weapon.wert') }}</th>
                <th class="cd-table-header">{{ $t('weapon.damage') }}</th>
                <th class="cd-table-header">{{ $t('weapon.description') }}</th>
                <th class="cd-table-header">{{ $t('weapon.quelle') }}</th>
                <th class="cd-table-header">{{ $t('weapon.personal_item') }}</th>
                <th class="cd-table-header">{{ $t('weapon.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(dtaItem, index) in filteredAndSortedWeaponss" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <!-- <td>{{ dtaItem.category|| '-' }}</td> -->
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.gewicht || '-' }}</td>
                  <td>{{ dtaItem.wert || '-' }}</td>
                  <td>{{ dtaItem.damage || '-' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ dtaItem.quelle || '-' }}</td>
                  <td>{{ dtaItem.personal_item || '0' }}</td>
                  <td>{{ dtaItem.system || 'midgard' }}</td>
                  <td>
                    <button @click="startEdit(index)">Edit</button>
                  </td>
                </tr>
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;"/></td>
                  <!-- <td><select v-model="editedItem.category" style="width:80px;">
                    <option v-for="category in mdata['weaponscategories']"
                            :key="category"
                            :value="category">
                      {{ category }}
                    </option>
                  </select></td> -->
                  <td><input v-model="editedItem.name"/></td>
                  <td><input v-model.number="editedItem.gewicht" type="number" style="width:40px;"/></td>
                  <td><input v-model="editedItem.wert" /></td>
                  <td><input v-model="editedItem.damage" /></td>
                  <td><input v-model="editedItem.beschreibung" /></td>
                  <td><input v-model="editedItem.quelle"  style="width:80px;"/></td>
                  <td><input type="checkbox" :checked="true" v-model="editedItem.personal_item" style="width:50px;"/></td>
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
  name: "WeaponView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        weaponss: [],
        weaponscategories: []
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
    filteredAndSortedWeaponss() {
      if (!this.mdata?.weapons) return [];

      return [...this.mdata.weapons]
        .filter(weapons => {
          const searchLower = this.searchTerm.toLowerCase();
          return !this.searchTerm ||
            weapons.name?.toLowerCase().includes(searchLower)
            //|| weapons.category?.toLowerCase().includes(searchLower);
        })
        .sort((a, b) => {
          const aValue = (a[this.sortField] || '').toLowerCase();
          const bValue = (b[this.sortField] || '').toLowerCase();
          return this.sortAsc
            ? aValue.localeCompare(bValue)
            : bValue.localeCompare(aValue);
        });
    },
    sortedWeaponss() {
      return [...this.mdata.weapons].sort((a, b) => {
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
      this.editedItem = { ...this.filteredAndSortedWeaponss[index] };
    },
    saveEdit(index) {
      //this.$emit('update-weapons', { index, weapons: this.editedItem });
      this.handleWeaponsUpdate( { index, weapons: this.editedItem });
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
    async handleWeaponsUpdate({ index, weapons }) {
      try {
          const response = await API.put(
            `/api/maintenance/weapons/${weapons.id}`, weapons,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}` ,
                'Content-Type': 'application/json'
              }
            }
          )
          if (response.status !== 200) throw new Error('Update failed');
          const updatedSkill = response.data;
          // Update the weapons in mdata
          this.mdata.weapons = this.mdata.weapons.map(s =>
            s.id === updatedSkill.id ? updatedSkill : s
          );
        } catch (error) {
          console.error('Failed to update weapons:', error);
        }
      }
  }
};
</script>
