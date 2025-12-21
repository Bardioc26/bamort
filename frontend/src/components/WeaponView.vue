<template>
  <div class="cd-view">
    <div class="cd-list">
      <table class="cd-table">
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
          <td>{{ weapon.description || '-' }}</td>
          <td>{{ weapon.weight || '-' }}</td>
          <td>{{ weapon.value || '-' }}</td>
          <td>{{ weapon.amount || '-' }}</td>
          <td>{{ weapon.contained_in || '-' }}</td>
          <td>{{ weapon.Anb || '-' }}/{{ weapon.Abwb || '-' }}</td>
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
name: "WeaponView",
  props: {
    character: {
      type: Object,
      required: true
    }
  }
};
</script>
