<template>
  <div>
    <h2>Ausruestung</h2>
    <ul>
      <li v-for="item in ausruestung" :key="item.ausruestung_id">
        {{ item.name }} - {{ item.anzahl }} ({{ item.gewicht }}kg)
        <button @click="deleteItem(item.ausruestung_id)">Delete</button>
      </li>
    </ul>
    <AusruestungForm @added="fetchAusruestung" :characterId="characterId" />
  </div>
</template>

<script>
import API from '../utils/api'
import AusruestungForm from './AusruestungForm.vue'

export default {
  components: { AusruestungForm },
  props: ['characterId'],
  data() {
    return {
      ausruestung: [],
    }
  },
  async created() {
    this.fetchAusruestung()
  },
  methods: {
    async fetchAusruestung() {
      const response = await API.get(`/ausruestung/${this.characterId}`)
      this.ausruestung = response.data
    },
    async deleteItem(id) {
      await API.delete(`/ausruestung/${id}`)
      this.fetchAusruestung()
    },
  },
}
</script>

<style>
/* All common styles moved to main.css */
</style>
