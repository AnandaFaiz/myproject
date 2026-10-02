interface User {
    id: number
    nama: string
    email: string
    role: string
}

export const useAuth = () =>{
    const user = useState<User | null>('auth-user',() => null)
    const loading = useState<boolean>('auth-loading',() => false)

    //token disimpan di cookie
    const token = useCookie<string | null>('auth_token', {
        maxAge: 60 * 60 * 24 * 7,
        sameSite: 'lax',
        path: '/',
    })

    //helper: fetch dengan auth header
    function fetchDenganAuth<T = any>(url: string, options: any ={}): Promise<T>{
        const headers: Record <string,string> = {
            ...(options.headers || {}),
        }

        if (token.value){
            headers.Authorization = `Bearer ${token.value}`
        }

        //kalau body bukan FormData, set Content-Type JSON
        if (options.body && !(options.body instanceof FormData)){
            headers['Content-Type'] = 'application/json'
        }

        return $fetch<T>(url, { ...options,headers})
    }

    async function fetchUser() {
        if (!token.value){
            user.value = null
            return
        }

        loading.value = true
        try{
            const res = await $fetch<{sukses: boolean; user: User | null}>('/api/auth/me', {
                headers:{
                    Authorization: `Bearer ${token.value}`,
                },
            })
            user.value = res.user
        }catch{
            user.value = null
            token.value = null
        }finally{
            loading.value = false
        }
    }

    async function logout() {
        // await $fetch('/api/auth/logout', {method:'POST'})

        token.value = null
        user.value = null
        await navigateTo('/')
    }

    //simpan token setelah login
    function setToken(newToken: string){
        token.value = newToken
    }

    const isLoggedIn = computed(() => !!user.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    return{
        user,
        token,
        loading,
        isLoggedIn,
        isAdmin,
        fetchUser,
        logout,
        setToken,
        fetchDenganAuth,
    }
}