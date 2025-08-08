<template>
  <div class="fullwidth-container cd-view">
    <h2>{{ character.name }}'s Ausrüstung</h2>
    <div class="cd-list">
      <table class="cd-table">
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
    </div> <!--- end cd-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>
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
