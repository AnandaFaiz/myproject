export default defineNuxtRouteMiddleware(async (to) => {
  const perluLogin =
    to.path.startsWith('/admin') || to.path.startsWith('/profil')

  if (!perluLogin) return

  const { fetchUser, isLoggedIn, isAdmin } = useAuth()

  await fetchUser()

  if (!isLoggedIn.value) {
    return navigateTo({
      path:'/login',
      query: { redirect: to.fullPath },
    })
  }

  // Khusus /admin, role admin
  if (to.path.startsWith('/admin') && !isAdmin.value) {
    return navigateTo('/')
  }
})