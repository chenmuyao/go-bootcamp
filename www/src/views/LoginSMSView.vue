<template>
  <div class="login bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="login-page w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h3 class="text-2xl font-bold text-center mb-6">Login</h3>
      <form class="space-y-4" @submit.prevent="handleLogin">
        <div>
          <label for="phone" class="block text-sm font-medium text-gumbo-700">Phone Number:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" type="tel" id="phone" v-model="form.phone" required>
        </div>
        <div v-if="smsSent">
          <label for="code">Code:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" type="text" id="code" v-model="form.code" required>
        </div>
        <div v-if="errorMessage" class="p-2 mb-4 text-sm text-gumbo-800 bg-gumbo-50 rounded-lg">{{ errorMessage }}</div>
        <div v-if="!smsSent">
          <button class="w-full px-4 py-2 font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50" type="button" :disabled="!validatePhone(form.phone)"
            @click="sendSMS">Send
            SMS</button>
        </div>
        <div v-else class="flex space-x-4">
          <button class="block mt-4 text-center text-gumbo-500 hover:underline"
            type="button"
            :disabled="!validatePhone(form.phone)" @click="sendSMS">Resend SMS</button>
          <button class="flex-grow px-4 py-2 font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50"
            type="submit"
            :disabled="!validateCode(form.code) || !validatePhone(form.phone)">Login</button>
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
</style>
