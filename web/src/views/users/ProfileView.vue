<template>
  <div class="profile">
    <h1 class="w3-margin">User Profile</h1>
    <div class="w3-margin">
      <form class="w3-form" @submit.prevent="saveProfile">
        <div>
          <label for="name">Name:</label>
          <input class="w3-input" id="name" v-model="userProfile.name" type="text" :disabled="!isEditing" />
        </div>
        <div>
          <label for="email">Email:</label>
          <input class="w3-input" id="email" v-model="userProfile.email" type="email" disabled />
        </div>
        <div>
          <label for="phone">Phone:</label>
          <input class="w3-input" id="phone" v-model="userProfile.phone" type="tel" disabled />
        </div>
        <div>
          <label for="birthday">Birthday:</label>
          <input class="w3-input" id="birthday" v-model="userProfile.birthday" type="date" :disabled="!isEditing" />
        </div>
        <div>
          <label for="profile">About Me:</label>
          <textarea class="w3-input" id="profile" v-model="userProfile.profile" :disabled="!isEditing"></textarea>
        </div>
        <button class="w3-button w3-blue w3-margin" type="submit" v-if="isEditing">Save</button>
        <button class="w3-button w3-green w3-margin" type="button" v-if="!isEditing" @click="toggleEdit">Edit</button>
        <button class="w3-button w3-green w3-margin" type="button" v-else @click="toggleEdit">Cancel</button>
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
      userProfile.value = response.data;
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
@media (min-width: 1024px) {
  .profile {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
