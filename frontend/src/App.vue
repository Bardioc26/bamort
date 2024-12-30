<template>
  <div id="app">
    <Menu />
    <router-view />
  </div>
</template>

<script>
import Menu from "./components/Menu.vue";

export default {
  components: {
    Menu,
  },
};
</script>

<style>
/* Global styles can go here */
</style>
