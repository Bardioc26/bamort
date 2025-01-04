<template>
  <div>
    <h2>Your Characters</h2>
    <ul>
      <li v-for="character in characters" :key="character.character_id">
        {{ character.name }} ({{ character.rasse }}, {{ character.typ }}, {{ character.grad }}, {{ character.owner }}, {{ character.public }} )
        <button @click="goToAusruestung(character.character_id)">Manage Equipment</button>
      </li>
    </ul>
  </div>
</template>

<script>
import API from '../utils/api'

export default {
  data() {
    return {
      characters: [],
    }
  },
  async created() {
    const token = localStorage.getItem('token')
    const response = await API.get('/api/characters', {
      headers: { Authorization: `Bearer ${token}` },
    })
    this.characters = response.data
  },
  methods: {
    goToAusruestung(characterId) {
      this.$router.push(`/api/ausruestung/${characterId}`)
    },
  },
}
</script>
