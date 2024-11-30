<template>
  <div class="signup">
    <div class="signup-page">
      <h3>Sign Up</h3>
      <form class="w3-container" @submit.prevent="handleSignup">
        <div>
          <label for="email">Email:</label>
          <input class="w3-input" type="email" id="email" v-model="form.email" required>
        </div>
        <div>
          <label for="password">Password:</label>
          <input class="w3-input" type="password" id="password" v-model="form.password" required>
        </div>
        <div>
          <label for="confirm-password">Confirm password:</label>
          <input class="w3-input" type="password" id="confirm-password" v-model="form.confirmPassword" required>
        </div>
        <div v-if="errorMessage" class="w3-panel w3-red">{{ errorMessage }}</div>
        <div>
          <button class="w3-margin-top w3-btn w3-ripple w3-green" type="submit" :disabled="!isFormValid">Sign
            Up</button>
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

<style>
@media (min-width: 1024px) {
  .signup {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
