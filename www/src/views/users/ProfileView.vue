<template>
  <div class="profile bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="profile-page w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h1 class="text-3xl font-bold text-center mb-6">User Profile</h1>
      <form class="space-y-4" @submit.prevent="saveProfile">
        <div>
          <label for="name" class="block text-sm font-medium text-gumbo-700">Name:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" id="name" v-model="userProfile.name" type="text" :disabled="!isEditing" />
        </div>
        <div>
          <label for="email" class="block text-sm font-medium text-gumbo-700">Email:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" id="email" v-model="userProfile.email" type="email" disabled />
        </div>
        <div>
          <label for="phone" class="block text-sm font-medium text-gumbo-700">Phone:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" id="phone" v-model="userProfile.phone" type="tel" disabled />
        </div>
        <div>
          <label for="birthday" class="block text-sm font-medium text-gumbo-700">Birthday:</label>
          <input class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" id="birthday" v-model="userProfile.birthday" type="date" :disabled="!isEditing" />
        </div>
        <div>
          <label for="profile" class="block text-sm font-medium text-gumbo-700">About Me:</label>
          <textarea class="w-full p-2 border border-gumbo-300 rounded focus:outline-none focus:border-gumbo-500" id="profile" v-model="userProfile.profile" :disabled="!isEditing"></textarea>
        </div>
        <div class="flex justify-between space-x-2">
          <button v-if="isEditing" type="submit" class="px-4 py-2 font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50">Save</button>
          <button v-if="!isEditing" type="button" @click="toggleEdit" class="px-4 py-2 font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50">Edit</button>
          <button v-else type="button" @click="toggleEdit" class="px-4 py-2 font-semibold text-white bg-red-500 rounded hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-600 focus:ring-opacity-50">Cancel</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from "@/axios/axios";

interface UserProfile {
  name: string;
  email: string;
  phone: string;
  birthday: string;
  profile: string;
}

const userProfile = ref<UserProfile>({
  name: '',
  email: '',
  phone: '',
  birthday: '',
  profile: ''
});

const isEditing = ref(false);

const fetchProfile = async () => {
  axios.get('/user/profile')
    .then((response) => {
      console.log(response)
      userProfile.value = response.data.data;
    })
    .catch((error) => {
      console.error(error)
    });
};

const saveProfile = async () => {
  console.log(userProfile.value);
  axios.post('/user/edit', {
    name: userProfile.value.name,
    birthday: userProfile.value.birthday,
    profile: userProfile.value.profile,
  })
    .then((response) => {
      console.log(response)
      alert('Profile saved successfully!');
      toggleEdit(); // Go back to view mode after saving
    })
    .catch((error) => {
      console.error(error)
    });
};

const toggleEdit = () => {
  isEditing.value = !isEditing.value;
};

onMounted(() => {
  fetchProfile();
});

</script>

<style>
</style>
