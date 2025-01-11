<template>
  <div class="spell-view">
    <h2>{{ character.name }}'s Skills</h2>
    <div class="spell-list">
        <table class="spell-table">
      <thead>
        <tr>
          <th>{{ $t('spell.name') }}</th>
          <th>{{ $t('spell.description') }}</th>
          <th>{{ $t('spell.bonus') }}</th>
          <th>{{ $t('spell.quelle') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-for="spell in character.zauber">
        <tr>
          <td>{{ spell.name || '-' }}</td>
          <td>{{ spell.beschreibung || '-' }}</td>
          <td>{{ spell.bonus || '0' }}</td>
          <td>{{ spell.quelle || '-' }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div> <!--- end spell-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>
.spell-view {
  color: #fff;
}
.spell-table {
  width: 100%;
  border-collapse: collapse;
  position: relative;
  margin-top: 1rem;
}

.spell-table thead {
  position: sticky;
  top: 0;
  background: #fff;
  z-index: 1;
}
.spell-table th,
.spell-table td {
  padding: 0.5rem;
  text-align: left;
  border: 1px solid #ddd;
}

.spell-table th {
  background-color: #f5f5f5;
  font-weight: bold;
}

.spell-table tr:nth-child(even) {
  background-color: #f9f9f9;
  color: #000000;
}

.spell-list {
  max-height: 400px;
  overflow-y: auto;
  overflow-x: hidden;
}

/* Custom scrollbar styling */
.spell-list::-webkit-scrollbar {
  width: 8px;
}

.spell-list::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.spell-list::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.spell-list::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>


<script>
export default {
  name: "SpellView",
  props: {
    character: {
      type: Object,
      required: true
    }
  }
};
</script>
