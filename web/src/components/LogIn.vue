<template>
  <div class="login-page">
    <h3>Login</h3>
    <form class="w3-container" @submit.prevent="handleLogin">
      <div>
        <label for="email">Email:</label>
        <input class="w3-input" type="email" id="email" v-model="form.email" required>
      </div>
      <div>
        <label for="password">Password:</label>
        <input class="w3-input" type="password" id="password" v-model="form.password" required>
      </div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div>
        <button class="w3-margin-top w3-button w3-green" type="submit" :disabled="!isFormValid">Login</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue';
import axios from 'axios';

const emailRegexPattern = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/

const emit = defineEmits([ 'login' ]);

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

  // Emit an event
  emit('login', {
    email: form.email,
    password: form.password,
  });

  axios.post('http://localhost:7779/user/login', {
    email: form.email,
    password: form.password,
  })
    .then((response) => {
      // get the jwt token from header (x-jwt-token)
      // store it somewhere for further usage
      // and redirect to user/profile page
    })
    .catch((error) => {
      console.error(error)
    });
};
</script>

