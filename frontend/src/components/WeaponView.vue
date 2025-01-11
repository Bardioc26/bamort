<template>
  <div class="weapon-view">
    <h2>{{ character.name }}'s Ausrüstung</h2>
    <div class="weapon-list">
      <table class="weapon-table">
      <thead>
        <tr>
          <th>{{ $t('weapon.name') }}</th>
          <th>{{ $t('weapon.description') }}</th>
          <th>{{ $t('weapon.weight') }}</th>
          <th>{{ $t('weapon.value') }}</th>
          <th>{{ $t('weapon.amount') }}</th>
          <th>{{ $t('weapon.contained_in') }}</th>
          <th>{{ $t('weapon.bonus') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-for="weapon in character.waffen">
        <tr>
          <td>{{ weapon.name || '-' }}</td>
          <td>{{ weapon.beschreibung || '-' }}</td>
          <td>{{ weapon.gewicht || '-' }}</td>
          <td>{{ weapon.wert || '-' }}</td>
          <td>{{ weapon.anzahl || '-' }}</td>
          <td>{{ weapon.beinhaltet_in || '-' }}</td>
          <td>{{ weapon.Anb || '-' }}/{{ weapon.Abwb || '-' }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div> <!--- end weapon-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>.weapon-view {
  color: #fff;
  }
  .weapon-table {
    width: 100%;
    border-collapse: collapse;
    position: relative;
    margin-top: 1rem;
  }

  .weapon-table thead {
    position: sticky;
    top: 0;
    background: #fff;
    z-index: 1;
  }
  .weapon-table th,
  .weapon-table td {
    padding: 0.5rem;
    text-align: left;
    border: 1px solid #ddd;
  }

  .weapon-table th {
    background-color: #f5f5f5;
    font-weight: bold;
  }

  .weapon-table tr:nth-child(even) {
    background-color: #f9f9f9;
    color: #000000;
  }


  .weapon-list {
    max-height: 400px;
    overflow-y: auto;
    overflow-x: hidden;
  }

  /* Custom scrollbar styling */
  .weapon-list::-webkit-scrollbar {
    width: 8px;
  }

  .weapon-list::-webkit-scrollbar-track {
    background: #f1f1f1;
  }

  .weapon-list::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 4px;
  }

  .weapon-list::-webkit-scrollbar-thumb:hover {
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
