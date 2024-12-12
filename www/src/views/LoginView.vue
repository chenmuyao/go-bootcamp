<template>
  <div class="login bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="login-page w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h3 class="text-2xl font-bold text-center mb-6">Login</h3>
      <form class="space-y-4" @submit.prevent="handleLogin">
        <div>
          <label for="email" class="block text-sm font-medium text-gumbo-700">Email:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" type="email"
            id="email" v-model="form.email" required>
        </div>
        <div>
          <label for="password">Password:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500"
            type="password" id="password" v-model="form.password" required>
        </div>
        <div v-if="errorMessage" class="p-2 mb-4 text-sm text-gumbo-800 bg-gumbo-50 rounded-lg">{{ errorMessage }}</div>
        <div>
          <button
            class="w-full px-4 py-2 font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50"
            type="submit" :disabled="!isFormValid">Login</button>
          <a class="block mt-4 text-center text-gumbo-500 hover:underline" href="/login_sms">Signup or Login with SMS</a>
          <a class="block mt-4 text-center text-gumbo-500 hover:underline" href="/login_gitea">Signup or Login with Gitea</a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue';
import axios from "@/axios/axios";
import { useRouter } from "vue-router"

const emailRegexPattern = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/

const router = useRouter()

const form = reactive({
  email: '',
  password: '',
});

const errorMessage = computed(() => {
  return validateForm();
});

const isFormValid = computed(() => {
  return '' === validateForm();
});

const validateForm = () => {
  if (form.email && !validateEmail(form.email)) {
    return 'Please enter a valid email';
  }

  if (!form.password) {
    return 'Empty password is not allowed';
  }

  return '';
}

const validateEmail = (val: string) => {
  return val.match(emailRegexPattern) != null;
}

const handleLogin = () => {
  if (!isFormValid.value) return;

  axios.post('/user/login', {
    'email': form.email,
    'password': form.password,
  })
    .then((response) => {
      console.log(response)
      if (response.status != 200) {
        alert(response.statusText);
        return
      }
      console.log(response.data);
      router.push({ path: "/user/profile" });

    })
    .catch((error) => {
      console.error(error)
    });

};
</script>


<style></style>
