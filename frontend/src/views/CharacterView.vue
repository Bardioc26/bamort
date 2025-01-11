<template>
  <CharacterForm :characterId="$route.params.characterId" />
</template>

<script>
import CharacterForm from '../components/CharacterDetails.vue'

export default {
  components: { CharacterForm },
}
</script>
