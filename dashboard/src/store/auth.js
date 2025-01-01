import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import axios from 'axios'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import { jwtDecode } from 'jwt-decode'
import moment from 'moment'

export const useAuthStore = defineStore('auth_details', () => {
  const IsLoggedIn = ref(false)
  const AccessToken = ref('')
  const IsLoggingInProgress = ref(false)
  const currentUsername = ref('')
  let currentTime = ref(Date.now())

  function FetchBearerToken() {
    if (IsLoggedIn.value) {
      return 'Bearer ' + AccessToken.value
    }
    return ''
  }

  function SetCredential(token) {
    AccessToken.value = token
    localStorage.setItem('token', token)
    IsLoggedIn.value = true
    IsLoggingInProgress.value = true
    setTimeout(() => {
      IsLoggingInProgress.value = false
    }, 1000)
  }

  async function Login(username, password, totp) {
    let data = new FormData()
    data.append('username', username)
    data.append('password', password)
    data.append('totp', totp)

    // environment variable
    const HTTP_BASE_URL = getHttpBaseUrl()

    let config = {
      method: 'post',
      url: `${HTTP_BASE_URL}/auth/login`,
      data: data
    }

    try {
      const res = await axios.request(config)
      const resData = res.data
      SetCredential(resData.token)
      return {
        success: true,
        message: 'Logged in successfully !',
        totp_required: false
      }
    } catch (e) {
      if (e.response) {
        return {
          success: false,
          message: e.response.data.message || 'Unexpected error',
          totp_required: e.response.data.totp_required
        }
      } else {
        return {
          success: false,
          message: 'Failed to send request',
          totp_required: false
        }
      }
    }
  }

  function Logout() {
    // logout
    IsLoggedIn.value = false
    localStorage.clear()
    IsLoggingInProgress.value = true
    // redirect to /
    setTimeout(() => {
      this.$router.push({ name: 'Login' }).then(() => {
        this.$router.go(0)
      })
    }, 500)
  }

  async function CheckAuthStatus() {
    try {
      const token = localStorage.getItem('token')
      if (token) {
        const HTTP_BASE_URL = getHttpBaseUrl()
        let config = {
          method: 'get',
          url: `${HTTP_BASE_URL}/verify-auth`,
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
        const res = await axios.request(config)
        return res.status === 200
      }
    } catch (e) {
      if (e.isAxiosError) {
        if (e.message.includes('Network Error')) {
          return true
        }
      }
      return false
    }
  }

  async function logoutOnInvalidToken(callback) {
    if (!IsLoggedIn.value) {
      return
    }
    const isTokenValid = await CheckAuthStatus()
    if (!isTokenValid) {
      callback()
    }
  }

  function StartAuthChecker(callback) {
    setInterval(() => logoutOnInvalidToken(callback), 5000)
  }

  const sessionRelativeTimeoutStatus = computed(() => {
    if (IsLoggedIn.value) {
      try {
        const token = localStorage.getItem('token')
        if (token) {
          const decoded = jwtDecode(token)
          currentUsername.value = decoded.username ?? ''
          const exp = moment(new Date(decoded.exp * 1000))
          return moment.duration(exp.diff(currentTime.value)).humanize(true)
        }
      } catch (e) {
        return 'N/A'
      }
    }
    return ''
  })

  async function fetchSWVersion() {
    if (!IsLoggedIn.value) return '...'
    try {
      const HTTP_BASE_URL = getHttpBaseUrl()
      let config = {
        method: 'get',
        url: `${HTTP_BASE_URL}/version`,
        headers: {
          Authorization: FetchBearerToken()
        }
      }
      const res = await axios.request(config)
      return res.data
    } catch (e) {
      return 'N/A'
    }
  }

  setInterval(() => {
    currentTime.value = Date.now()
  }, 10000)
  return {
    IsLoggedIn,
    IsLoggingInProgress,
    FetchBearerToken,
    Login,
    Logout,
    SetCredential,
    StartAuthChecker,
    sessionRelativeTimeoutStatus,
    fetchSWVersion,
    currentUsername
  }
})
