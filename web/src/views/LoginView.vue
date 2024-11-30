<template>
  <div class="login">
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
          <a class="w3-margin" href="/login_sms">Signup or Login with SMS</a>
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
      if(response.status != 200) {
        alert(response.statusText);
        return
      }
      console.log(response.data);
      router.push({path: "/about"});

    })
    .catch((error) => {
      console.error(error)
    });

};
</script>


<style>
@media (min-width: 1024px) {
  .login {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
