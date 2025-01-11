<template>
  <div id="app">
    <Menu />
     <!-- Main Content Area -->
     <main class="main-content">
      <router-view />
    </main>
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

<style src="./assets/main.css"></style>
<style>
/* Global styles can go here */
</style>
