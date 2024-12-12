<template>
  <div class="login bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="login-page w-full max-w-2xl p-6 bg-white rounded-lg shadow-md">
      <h3 v-if="isLoading" class="text-2xl font-bold text-center mb-6">Loading...</h3>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from "@/axios/axios";
import { onMounted, ref } from 'vue'

const isLoading = ref(true);

const setLoading = (set: bool) => {
  isLoading.value = set;
}


onMounted(() => {
  axios.get('/oauth2/gitea/authurl')
    .then((res) => res.data)
    .then((data) => {
      setLoading(false)
      if (data && data.data) {
        window.location.href = data.data
      }
    })
})

</script>
