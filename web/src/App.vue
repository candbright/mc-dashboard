<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<style>
*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
}

#app {
  width: 100vw;
  min-height: 100vh;
}
</style>
