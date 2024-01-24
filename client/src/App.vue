<script setup lang="ts">
  import NavigationBar from '@/components/NavigationBar.vue';
</script>
<template>
  <NavigationBar />
  <main class="container m-auto max-w-screen-lg py-2">
    <RouterView />
  </main>
</template>

<style scoped>
</style>
