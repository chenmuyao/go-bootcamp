<template>
  <div class="login">
    <div class="login-page">
      <h3>Login</h3>
      <form class="w3-container" @submit.prevent="handleLogin">
        <div>
          <label for="phone">Phone Number:</label>
          <input class="w3-input" type="tel" id="phone" v-model="form.phone" required>
        </div>
        <div v-if="smsSent">
          <label for="code">Code:</label>
          <input class="w3-input" type="text" id="code" v-model="form.code" required>
        </div>
        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
        <div v-if="!smsSent">
          <button class="w3-margin-top w3-button w3-green" type="button" :disabled="!validatePhone(form.phone)"
            @click="sendSMS">Send
            SMS</button>
        </div>
        <div v-else>
          <button class="w3-margin w3-button w3-green" type="button"
            :disabled="!validatePhone(form.phone)" @click="sendSMS">Resend SMS</button>
          <button class="w3-margin w3-button w3-green" type="submit"
            :disabled="!validateCode(form.code) || !validatePhone(form.phone)">Login</button>
          <a class="w3-margin" href="/login">Login with password</a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, ref } from 'vue';
import axios from "@/axios/axios";
import { useRouter } from "vue-router"

const phoneRegexPattern = /^\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}$/
const codeRegexPattern = /^\d{6}$/

const router = useRouter()

const smsSent = ref(false);

const form = reactive({
  phone: '',
  code: '',
});

const errorMessage = computed(() => {
  return validateForm();
});

// TODO: Add a timer to disable resend button for 1 min

const validateForm = () => {
  if (form.phone && !validatePhone(form.phone)) {
    return 'Please enter a valid phone number';
  }

  if (form.code && !validateCode(form.code)) {
    return 'Please enter a valid code';
  }

  return '';
}

const validateCode = (val: string) => {
  return val.match(codeRegexPattern) != null;
}

const validatePhone = (val: string) => {
  return val.match(phoneRegexPattern) != null;
}

const sendSMS = () => {
  if (!validatePhone(form.phone)) return;

  axios.post('/user/login_sms/code/send', {
    'phone': form.phone,
  })
    .then((response) => {
      console.log(response)
      if (response.status != 200) {
        alert(response.statusText);
        return
      }
      console.log(response.data);
      smsSent.value = true;
    })
    .catch((error) => {
      console.error(error)
    });
}

const handleLogin = () => {
  if (!validateCode(form.code) || !validatePhone(form.phone)) return;

  axios.post('/user/login_sms', {
    'phone': form.phone,
    'code': form.code.toString(),
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


<style>
@media (min-width: 1024px) {
  .login {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
