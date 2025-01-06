<template>
  <div>
    <h1>Character Details</h1>
    <p>ID: {{ character.id }}</p>
    <p>Name: {{ character.name }}</p>
    <p>Race: {{ character.rasse }}</p>
    <p>Typ: {{ character.typ }}</p>
    <p>Grad: {{ character.grad }}</p>
    <p>Besitzer: {{ character.owner }}</p>
    <p>öffentlich: {{ character.public }} </p>
    <router-link to="/dashboard">Back to Character List</router-link>
  </div>
</template>

<script>
import API from '../utils/api'

export default {
  name: "CharacterDetails",
  props: ["id"], // Receive the route parameter as a prop
  data() {
    return {
      character: {},
    };
  },
  async created() {
    const token = localStorage.getItem('token')
    const response = await API.get(`/api/characters/${this.id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    this.character = response.data
  },
};
</script>
