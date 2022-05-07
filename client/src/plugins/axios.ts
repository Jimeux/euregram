import axios from "axios"

const baseURL = (import.meta.env.VITE_BASE_URL)?.toString()
axios.defaults.baseURL = baseURL

interface Auth {
  access_token: string
}

// auth instance to avoid default headers or interceptor loops
const authInstance = axios.create({baseURL})
// cache auth
let auth: Auth | null = null

axios.interceptors.request.use(async (config) => {
  if (auth == null) {
    const str = localStorage.getItem("auth")
    if (str) auth = JSON.parse(str)
  }

  if (auth != null) {
    const bearer = `Bearer ${auth.access_token}`
    if (config.headers) {
      config.headers['Authorization'] = bearer
    } else {
      config.headers = {"Authorization": bearer}
    }
  }

  // allow requests without auth to be sent (triggers Google redirect on response)
  return config
}, (error) => {
  return Promise.reject(error)
})

axios.interceptors.response.use((response) => response, async (error) => {
  if (window.location.search.includes("www.googleapis.com")) {
    const searchParams = new URLSearchParams(window.location.search)
    const code = searchParams.get("code")
    const state = searchParams.get("state")
    const response = await authInstance.post("/auth/confirm", {code, state})
    localStorage.setItem("auth", JSON.stringify(response.data))

    window.location.href = "/"
  } else if (error.response.status === 401 || error.response.status === 403) {
    await redirectToGoogleAuth()
  } else {
    return Promise.reject(error)
  }
})

async function redirectToGoogleAuth() {
  localStorage.removeItem("auth")
  const response = await authInstance.get("/auth/init")
  window.location.href = response.data.redirect_url
}
