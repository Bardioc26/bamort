<template>
  <div>
    <h2>Your Characters</h2>
    <ul>
      <li v-for="character in characters" :key="character.character_id" style="white-space: nowrap; /* Prevent line breaks inside list items */;">
        <!-- Link to Character Details -->
        <router-link :to="`/character/${character.id}`">View Details</router-link>
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
