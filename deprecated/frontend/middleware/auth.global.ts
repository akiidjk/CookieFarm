import { HOST } from '@/lib/config'

export default defineNuxtRouteMiddleware(async (to, _from) => {
  // const token = useCookie('token')

  // if (!token.value && to.path !== '/login') {
  //   return navigateTo('/login')
  // }

  // try {
  //   const res = await fetch(HOST + '/api/v1/auth/verify', {
  //     method: 'POST',
  //     body: JSON.stringify({ token: token.value }),
  //     headers: {
  //       'Content-Type': 'application/json',
  //     },
  //   })

  //   if (!res.ok) {
  //     token.value = null
  //     if (to.path !== '/login') {
  //       return navigateTo('/login')
  //     }
  //   }
  // } catch (error) {
  //   token.value = null
  //   if (to.path !== '/login') {
  //     return navigateTo('/login')
  //   }
  //   console.error("Error on: ", error)
  // }

  // if (token.value && to.path === '/login') {
  //   return navigateTo('/')
  // }
  //
  return Promise.resolve()
})
