<template>
  <form @submit.prevent="addItem">
    <input v-model="name" placeholder="Name" required />
    <input v-model.number="anzahl" type="number" placeholder="Anzahl" required />
    <input v-model.number="gewicht" type="number" placeholder="Gewicht" required />
    <button type="submit">Add Item</button>
  </form>
</template>

<script>
import API from '../utils/api'

export default {
  props: ['characterId'],
  data() {
    return {
      name: '',
      anzahl: 1,
      gewicht: 0,
    }
  },
  methods: {
    async addItem() {
      await API.post('/ausruestung', {
        character_id: this.characterId,
        name: this.name,
        anzahl: this.anzahl,
        gewicht: this.gewicht,
      })
      this.$emit('added')
      this.name = ''
      this.anzahl = 1
      this.gewicht = 0
    },
  },
}
</script>
