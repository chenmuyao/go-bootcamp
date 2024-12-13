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
import axios from "@/axios/axios";

const router = useRouter()

onMounted(() => {

  axios.post('/user/logout')
    .then((response) => {
      console.log(response)
      localStorage.removeItem('token');
      localStorage.removeItem('refresh_token');

      if (response.status != 200) {
        alert(response.statusText);
        return
      }
      console.log(response.data);
      router.push('/login');

    })
    .catch((error) => {
      console.error(error)
    });

})

</script>
