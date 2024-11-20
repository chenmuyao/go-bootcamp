<script setup lang="ts">
import LogIn from '../components/LogIn.vue'
</script>

<template>
  <div class="login">
    <LogIn />
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .login {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>

