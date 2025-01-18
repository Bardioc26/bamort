<template>
  <div class="cd-view">
    <div class="cd-list">
      <div class="tables-container">
          <table class="cd-table">
            <thead>
              <tr>
                <th class="cd-table-header">{{ $t('skill.id') }}</th>
                <th class="cd-table-header">{{ $t('skill.category') }}</th>
                <th class="cd-table-header">{{ $t('skill.name') }}</th>
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
              <template v-for="(dtaItem, index) in mdata['skills']" :key="index">
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
                  <td><input v-model.boolean="editedItem.improvable" style="width:50px;"/></td>
                  <td><input v-model.boolean="editedItem.innateskill" style="width:50px;"/></td>
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

<style>
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
export default {
  name: "SkillView",
  props: {
    mdata: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      editingIndex: -1,
      editedItem: null
    }
  },
  methods: {
    startEdit(index) {
      this.editingIndex = index;
      this.editedItem = { ...this.mdata['skills'][index] };
    },
    saveEdit(index) {
      this.$emit('update-skill', { index, skill: this.editedItem });
      this.editingIndex = -1;
      this.editedItem = null;
    },
    cancelEdit() {
      this.editingIndex = -1;
      this.editedItem = null;
    }
  }
};
</script>
