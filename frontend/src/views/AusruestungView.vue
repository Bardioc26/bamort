<template>
  <AusruestungList :characterId="$route.params.characterId" />
</template>

<script>
import AusruestungList from '../components/AusruestungList.vue'

export default {
  components: { AusruestungList },
}
</script>
