<template>
    <div class="cd-view">
      <DeleteCharView
      v-if="showDeleteDialog"
      :character="character"
      @deleted="handleDeleted"
      @cancel="handleCancel"
    />
      <p>Are you sure you want to delete {{ character.name }}?</p>
      <button @click="deleteCharacter">Yes</button>
      <button @click="$emit('cancel')">No</button>
    </div>
  </template>

  <script>
  export default {
    name: "DeleteCharView",

    props: {
      character: {
        type: Object,
        required: true
      }
    },
    methods: {
      async deleteCharacter() {
        try {
          const response = await fetch(`/api/characters/${this.character.id}`, {
            method: 'DELETE'
          });
          if (response.ok) {
            this.$emit('deleted');
          } else {
            console.error('Failed to delete character');
          }
        } catch (error) {
          console.error('Error:', error);
        }
      },
      handleCancel() {
      this.showDeleteDialog = false;
      // Optional: Go back in router history
      this.$router.go(-1);
    },
    handleDeleted() {
      this.$router.push('/characters');
    }
    }
  };
  </script>

  <style>
    /*
  .cd-view {
    text-align: center;
  }
  button {
    margin: 5px;
  }*/
  </style>

