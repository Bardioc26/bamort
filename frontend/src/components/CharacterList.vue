<template>
  <div>
    <h2>Your Characters</h2>
    <ul>
      <li v-for="character in characters" :key="character.character_id">
        {{ character.name }} ({{ character.rasse }})
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
    const response = await API.get('/characters', {
      headers: { Authorization: `Bearer ${token}` },
    })
    this.characters = response.data
  },
  methods: {
    goToAusruestung(characterId) {
      this.$router.push(`/ausruestung/${characterId}`)
    },
  },
}
</script>
