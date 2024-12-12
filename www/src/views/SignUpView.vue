<template>
  <div class="signup bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="signup-page w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h3 class="text-2xl font-bold text-center mb-6">Sign Up</h3>
      <form class="space-y-4" @submit.prevent="handleSignup">
        <div>
          <label for="email" class="block text-sm font-medium text-gumbo-700">Email:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500"
            type="email" id="email" v-model="form.email" required>
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gumbo-700">Password:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500"
            type="password" id="password" v-model="form.password" required>
        </div>
        <div>
          <label for="confirm-password" class="block text-sm font-medium text-gumbo-700">Confirm password:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500"
            type="password" id="confirm-password" v-model="form.confirmPassword" required>
        </div>
        <div v-if="errorMessage" class="p-2 mb-4 text-sm text-gumbo-800 bg-gumbo-50 rounded-lg">{{ errorMessage }}</div>
        <div>
          <button
            class="w-full px-4 py-2 font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-600 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50"
            type="submit" :disabled="!isFormValid">Sign Up</button>
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
const passwordRegexPattern = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$/

const router = useRouter()

const form = reactive({
  email: '',
  password: '',
  confirmPassword: '',
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

  if (form.password && !validatePassword(form.password)) {
    return 'Password should at least have letters, numbers and at least one special characters';
  }

  if (form.password && form.confirmPassword && form.password != form.confirmPassword) {
    return 'Passwords do not match';
  }

  return '';
}

const validateEmail = (val: string) => {
  return val.match(emailRegexPattern) != null;
}

const validatePassword = (val: string) => {
  return val.match(passwordRegexPattern) != null;
}

const handleSignup = () => {
  if (!isFormValid.value) return;

  axios.post('/user/signup', {
    'email': form.email,
    'password': form.password,
    'confirm_password': form.confirmPassword,
  })
    .then((response) => {
      console.log(response)
      if (response.status != 200) {
        alert(response.statusText);
        return
      }
      console.log(response.data);
      router.push("/login")
    })
    .catch((error) => {
      console.error(error)
    });
};
</script>

<style></style>
