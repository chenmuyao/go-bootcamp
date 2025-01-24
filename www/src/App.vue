<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header class="bg-gumbo-600 p-4">
    <div class="container mx-auto">
      <nav class="flex justify-between items-center">
        <div class="flex space-x-4">
          <RouterLink to="/" class="text-white hover:text-gumbo-300">Home</RouterLink>
          <RouterLink to="/about" class="text-white hover:text-gumbo-300">About</RouterLink>
          <RouterLink to="/signup" class="text-white hover:text-gumbo-300">Sign Up</RouterLink>
          <RouterLink to="/login" class="text-white hover:text-gumbo-300">Login</RouterLink>
          <RouterLink to="/login_sms" class="text-white hover:text-gumbo-300">SignUp/Login SMS</RouterLink>
          <RouterLink to="/logout" class="text-white hover:text-gumbo-300">Logout</RouterLink>
          <RouterLink to="/user/profile" class="text-white hover:text-gumbo-300">User Profile</RouterLink>
          <RouterLink to="/articles/list" class="text-white hover:text-gumbo-300">Article list</RouterLink>
        </div>
        <!-- Add additional navigation items here if needed -->
      </nav>
    </div>
  </header>

  <RouterView />
</template>
