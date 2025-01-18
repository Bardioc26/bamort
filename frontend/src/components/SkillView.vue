<template>
  <div class="cd-view">
    <h2>{{ character.name }}'s Fertigkeiten</h2>
    <div class="cd-list">
        <table class="cd-table">
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
    </div> <!--- end cd-list-->
  </div> <!--- end character -datasheet-->

</template>

<style>
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
