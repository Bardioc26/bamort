<template>
  <div class="cd-view">
    <div class="cd-list">
      <div class="tables-container">
          <h2 style="line-height: 1.5; margin-top: 5px;"><!-- {{ character.name }}'s -->Fertigkeiten</h2>
          <table class="cd-table">
            <thead>
              <tr>
                <th class="cd-table-header" width="60%">{{ $t('skill.name') }}</th>
                <th class="cd-table-header" width="35">{{ $t('skill.value') }}</th>
                <th class="cd-table-header" width="35">{{ $t('skill.bonus') }}</th>
                <th class="cd-table-header" width="35">{{ $t('skill.pp') }}</th>
                <!-- <th class="cd-table-header">{{ $t('skill.description') }}</th>-->
                <th class="cd-table-header" width="30%">{{ $t('skill.note') }}</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="skills,categorie in character.categorizedskills">
              <tr>
                <td colspan="6">{{ categorie || '-' }}</td>
              </tr>
              <template v-for="skill in skills">
                <tr>
                  <td>{{ skill.name || '-' }}</td>
                  <td>{{ skill.fertigkeitswert || '-' }}</td>
                  <td>{{ skill.bonus || '0' }}</td>
                  <td>{{ skill.pp || '0' }}</td>
                  <!-- <td>{{ skill.beschreibung || '-' }}</td>-->
                  <td>{{ skill.bemerkung || '-' }}</td>
                </tr>
              </template>
              </template>
              <template v-for="skill in character.waffenfertigkeiten">
                <tr>
                  <td>{{ skill.name || '-' }}</td>
                  <td>{{ skill.fertigkeitswert || '-' }}</td>
                  <td>{{ skill.bonus || '0' }}</td>
                  <td>{{ skill.pp || '0' }}</td>
                  <!-- <td>{{ skill.beschreibung || '-' }}</td> -->
                  <td>{{ skill.bemerkung || '-' }}</td>
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
  name: "WeaponView",
  props: {
    character: {
      type: Object,
      required: true
    }
  }
};
</script>
