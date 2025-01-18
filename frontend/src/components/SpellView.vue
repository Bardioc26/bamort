<template>
  <div class="cd-view">
    <h2>{{ character.name }}'s Zauber</h2>
    <div class="cd-list">
        <table class="cd-table">
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
    </div> <!--- end cd-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>
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
