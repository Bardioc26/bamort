<template>
  <div class="equipment-view">
    <h2>{{ character.name }}'s Ausrüstung</h2>
    <div class="equipment-list">
      <table class="equipment-table">
      <thead>
        <tr>
          <th>{{ $t('equipment.name') }}</th>
          <th>{{ $t('equipment.description') }}</th>
          <th>{{ $t('equipment.weight') }}</th>
          <th>{{ $t('equipment.value') }}</th>
          <th>{{ $t('equipment.amount') }}</th>
          <th>{{ $t('equipment.contained_in') }}</th>
          <th>{{ $t('equipment.bonus') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-for="equipment in character.ausruestung">
        <tr>
          <td>{{ equipment.name || '-' }}</td>
          <td>{{ equipment.beschreibung || '-' }}</td>
          <td>{{ equipment.gewicht || '-' }}</td>
          <td>{{ equipment.wert || '-' }}</td>
          <td>{{ equipment.anzahl || '-' }}</td>
          <td>{{ equipment.beinhaltet_in || '-' }}</td>
          <td>{{ equipment.bonus || '-' }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div> <!--- end equipment-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>.equipment-view {
  color: #fff;
  }
  .equipment-table {
    width: 100%;
    border-collapse: collapse;
    position: relative;
    margin-top: 1rem;
  }

  .equipment-table thead {
    position: sticky;
    top: 0;
    background: #fff;
    z-index: 1;
  }
  .equipment-table th,
  .equipment-table td {
    padding: 0.5rem;
    text-align: left;
    border: 1px solid #ddd;
  }

  .equipment-table th {
    background-color: #f5f5f5;
    font-weight: bold;
  }

  .equipment-table tr:nth-child(even) {
    background-color: #f9f9f9;
    color: #000000;
  }

  .equipment-list {
    max-height: 400px;
    overflow-y: auto;
    overflow-x: hidden;
  }

  /* Custom scrollbar styling */
  .equipment-list::-webkit-scrollbar {
    width: 8px;
  }

  .equipment-list::-webkit-scrollbar-track {
    background: #f1f1f1;
  }

  .equipment-list::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 4px;
  }

  .equipment-list::-webkit-scrollbar-thumb:hover {
    background: #555;
  }
</style>


<script>
export default {
name: "EquipmentView",
  props: {
    character: {
      type: Object,
      required: true
    }
  }
};
</script>
