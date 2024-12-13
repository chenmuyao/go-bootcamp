<template>
  <div class="login bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="login-page w-full max-w-2xl p-6 bg-white rounded-lg shadow-md">
      <h3 class="text-2xl font-bold text-center mb-6">Redirecting...</h3>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from "vue-router"

const router = useRouter()

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const token = urlParams.get('token');
  const refresh_token = urlParams.get('refresh_token');

  if (token) {
    localStorage.setItem('token', token);
  }

  if (refresh_token) {
    localStorage.setItem('refresh_token', refresh_token);
  }

  router.push('/user/profile');
})

</script>
