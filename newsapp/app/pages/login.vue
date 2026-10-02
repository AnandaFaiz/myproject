<script setup lang="ts">
import { ref } from "vue";

const route = useRoute();
const email = ref("");
const password = ref("");
const pesanError = ref("");
const loading = ref(false);

const { fetchUser, setToken } = useAuth();

const redirectTo = computed(() => {
  const r = route.query.redirect;
  return typeof r === "string" ? r : null;
});

async function login() {
  pesanError.value = "";
  loading.value = true;

  try {
    const res = await $fetch<{
      sukses: boolean;
      token: string;
      user: { role: string };
    }>("/api/auth/login", {
      method: "POST",
      body: { email: email.value, password: password.value },
    });

    setToken(res.token);
    await fetchUser();

    if (redirectTo.value) {
      await navigateTo(redirectTo.value);
    } else if (res.user.role === "admin") {
      await navigateTo("/admin/berita");
    } else {
      await navigateTo("/");
    }
  } catch (err: any) {
    pesanError.value =
      err?.data?.pesan || err?.data?.statusMessage || "Login gagal";
  } finally {
    loading.value = false;
  }
}
</script>
<template>
  <div class="halaman">
    <div class="kartu">
      <h1>Login</h1>

      <form @submit.prevent="login">
        <div class="field">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="email@example.com" />
        </div>

        <div class="field">
          <label>Password</label>
          <input v-model="password" type="password" placeholder="******" />
        </div>

        <p v-if="pesanError" class="error">{{ pesanError }}</p>

        <button type="submit" :disabled="loading">
          <i class="mdi mdi-login"></i>
          {{ loading ? "Memproses..." : "Masuk" }}
        </button>
      </form>

      <p class="daftar-link">
        Belum punya akun?
        <NuxtLink to="/register">Daftar di sini</NuxtLink>
      </p>
      <NuxtLink to="/" class="kembali">Kembali ke beranda</NuxtLink>
    </div>
  </div>
</template>

<style scoped>
.halaman {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  font-family: system-ui, sans-serif;
  padding: 1rem;
}
.kartu {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  width: 100%;
  max-width: 400px;
}
h1 {
  margin: 0 0 0.25rem;
  font-size: 1.5rem;
  color: #111827;
}
.sub {
  color: #6b7280;
  margin: 0 0 1.5rem;
  font-size: 0.9rem;
}

.field {
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
}
.field label {
  font-weight: 600;
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
}
.field input {
  padding: 0.6rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
}
.field input:focus {
  outline: 2px solid #2563eb;
  border-color: transparent;
}

.error {
  color: #dc2626;
  font-size: 0.85rem;
  margin: 0.5rem 0;
}
button {
  width: 100%;
  padding: 0.7rem;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  margin-top: 0.5rem;
}
button:hover {
  background: #1d4ed8;
}
button:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}

.daftar-link {
  text-align: center;
  font-size: 0.9rem;
  color: #6b7280;
  margin-top: 1rem;
}
.daftar-link a {
  color: #2563eb;
  text-decoration: none;
  font-weight: 500;
}
.daftar-link a:hover {
  text-decoration: underline;
}
.kembali {
  display: block;
  text-align: center;
  margin-top: 1rem;
  color: #6b7280;
  text-decoration: none;
  font-size: 0.85rem;
}
.kembali:hover {
  color: #2563eb;
}
</style>
