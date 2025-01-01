import { defineStore } from 'pinia'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import axios from 'axios'
import { ref } from 'vue'

export const useSystemConfigStore = defineStore('system_config', () => {
  const HTTP_BASE_URL = getHttpBaseUrl()
  const isSubmitting = ref(false)
  const isUpdating = ref(false)
  const isLoadingConfig = ref(false)

  async function submit(data) {
    isSubmitting.value = true
    let config = {
      method: 'post',
      url: `${HTTP_BASE_URL}/setup`,
      data: data
    }
    try {
      let res = await axios.request(config)
      isSubmitting.value = false
      return {
        success: true,
        message: res.data.message
      }
    } catch (e) {
      isSubmitting.value = false
      if (e.response) {
        return {
          success: false,
          message: e.response.data.message || 'Unexpected error'
        }
      } else {
        return {
          success: false,
          message: 'Failed to send request'
        }
      }
    }
  }

  function getAuthToken() {
    return 'Bearer ' + (localStorage.getItem('token') || '')
  }

  async function fetch() {
    isLoadingConfig.value = true
    let config = {
      method: 'get',
      url: `${HTTP_BASE_URL}/config/system`,
      headers: {
        Authorization: getAuthToken()
      }
    }
    try {
      let res = await axios.request(config)
      isLoadingConfig.value = false
      return {
        success: true,
        data: res.data
      }
    } catch (e) {
      isLoadingConfig.value = false
      if (e.response) {
        return {
          success: false,
          message: e.response.data.message || 'Unexpected error'
        }
      } else {
        return {
          success: false,
          message: 'Failed to send request'
        }
      }
    }
  }

  async function update(data) {
    isUpdating.value = true
    let config = {
      method: 'put',
      url: `${HTTP_BASE_URL}/config/system`,
      data: data,
      headers: {
        Authorization: getAuthToken()
      }
    }
    try {
      let res = await axios.request(config)
      isUpdating.value = false
      return {
        success: true,
        message: res.data.message
      }
    } catch (e) {
      isUpdating.value = false
      if (e.response) {
        return {
          success: false,
          message: e.response.data.message || 'Unexpected error'
        }
      } else {
        return {
          success: false,
          message: 'Failed to send request'
        }
      }
    }
  }

  return {
    submit,
    isSubmitting,
    fetch,
    isLoadingConfig,
    update,
    isUpdating
  }
})
