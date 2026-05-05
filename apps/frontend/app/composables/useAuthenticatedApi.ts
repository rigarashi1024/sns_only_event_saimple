export function useAuthenticatedApi() {
  const config = useRuntimeConfig()

  function redirectToLogin() {
    window.alert('セッションの有効期限が切れました。ログイン画面へ移動します。')
    window.location.assign('/login')
  }

  async function fetchWithSession<T>(path: string, method: 'GET' | 'POST' = 'GET') {
    const response = await $fetch.raw<T>(path, {
      baseURL: config.public.apiBaseUrl,
      method,
      credentials: 'include',
      headers: {
        'X-Requested-With': 'XMLHttpRequest',
      },
      ignoreResponseError: true,
    })

    if (response.status === 401) {
      redirectToLogin()
      return null
    }

    if (response.status >= 400) {
      throw new Error(`${path} request failed with status ${response.status}`)
    }

    return response._data ?? null
  }

  return {
    fetchWithSession,
  }
}
