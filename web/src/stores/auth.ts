import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

interface UserProfile {
  name: string;
  birthday: string;
  profile: string;
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<UserProfile | null>(null)
  const isAuthenticated = ref(false)
  const error = ref<string | null>(null)

  const router = useRouter()

  const login = (email: string, password: string) => {
    try {
      // Reset previous errors
      error.value = null

      // Make login API request
      axios.post('http:localhost:7779/user/login', {
        email,
        password
      })
        .then((response) => {
          console.log(response)
          // Store JWT token
          token.value = response.headers['x-jwt-token']

          // Set authentication state
          isAuthenticated.value = true

          // Fetch user profile
          fetchUserProfile()

          // Redirect to profile page
          router.push('/user/profile')
        })
        .catch((error) => {
          console.error(error)
          throw error
        });
    } catch (err: any) {
      // Handle login errors
      isAuthenticated.value = false
      error.value = err.response?.data?.message || 'Login failed'
      token.value = null
    }
  }

  const fetchUserProfile = () => {
    try {
      // Ensure we have a token before fetching profile
      if (!token.value) {
        throw new Error('No authentication token')
      }

      // Fetch user profile with JWT token
      const response = axios.get('http:localhost:7779/user/profile', {
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      })

      // Store user data
      user.value = response.data
    } catch (err: any) {
      // Handle profile fetch error
      error.value = err.response?.data?.message || 'Failed to fetch profile'
      logout()
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    isAuthenticated.value = false
    router.push('/login')
  }

  return {
    token,
    user,
    isAuthenticated,
    error,
    login,
    logout,
    fetchUserProfile
  }
})
