<template>
  <div class="cd-view">
    <div v-if="!isOwner" class="error-message">
      <p>{{ $t('deleteChar.notAuthorized') }}</p>
    </div>
    <div v-else>
      <p>Are you sure you want to delete {{ character.name }}?</p>
      <button @click="deleteCharacter" class="btn btn-danger">Yes, Delete</button>
      <button @click="$emit('cancel')" class="btn btn-secondary">Cancel</button>
    </div>
  </div>
</template>

<script>
import API from '../utils/api'

export default {
  name: "DeleteCharView",
  
  props: {
    character: {
      type: Object,
      required: true
    },
    isOwner: {
      type: Boolean,
      default: false
    }
  },
  
  emits: ['deleted', 'cancel'],
  
  methods: {
    async deleteCharacter() {
      try {
        const response = await API.delete(`/api/characters/${this.character.id}`)

        if (response.status === 200 || response.status === 204) {
          this.$emit('deleted')
          this.$router.push('/dashboard')
        } else {
          console.error('Failed to delete character')
        }
      } catch (error) {
        console.error('Error deleting character:', error)
      }
    }
  }
}
</script>

<style>
/* All common styles moved to main.css */
</style>

