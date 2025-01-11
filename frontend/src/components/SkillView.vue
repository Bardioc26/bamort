<template>
  <div class="skill-view">
    <h2>{{ character.name }}'s Skills</h2>
    <div class="skill-list">
        <table class="skill-table">
      <thead>
        <tr>
          <th>{{ $t('skill.name') }}</th>
          <th>{{ $t('skill.description') }}</th>
          <th>{{ $t('skill.value') }}</th>
          <th>{{ $t('skill.bonus') }}</th>
          <th>{{ $t('skill.pp') }}</th>
          <th>{{ $t('skill.note') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-for="skill in character.fertigkeiten">
        <tr>
          <td>{{ skill.name || '-' }}</td>
          <td>{{ skill.beschreibung || '-' }}</td>
          <td>{{ skill.fertigkeitswert || '-' }}</td>
          <td>{{ skill.bonus || '0' }}</td>
          <td>{{ skill.pp || '0' }}</td>
          <td>{{ skill.bemerkung || '-' }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div> <!--- end skill-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>
.skill-view {
  color: #fff;
}
.skill-table {
  width: 100%;
  border-collapse: collapse;
  position: relative;
  margin-top: 1rem;
}

.skill-table thead {
  position: sticky;
  top: 0;
  background: #fff;
  z-index: 1;
}
.skill-table th,
.skill-table td {
  padding: 0.5rem;
  text-align: left;
  border: 1px solid #ddd;
}

.skill-table th {
  background-color: #f5f5f5;
  font-weight: bold;
}

.skill-table tr:nth-child(even) {
  background-color: #f9f9f9;
  color: #000000;
}

.skill-list {
  max-height: 400px;
  overflow-y: auto;
  overflow-x: hidden;
}

/* Custom scrollbar styling */
.skill-list::-webkit-scrollbar {
  width: 8px;
}

.skill-list::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.skill-list::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.skill-list::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>


<script>
export default {
  name: "SkillView",
  props: {
    character: {
      type: Object,
      required: true
    }
  }
};
</script>
