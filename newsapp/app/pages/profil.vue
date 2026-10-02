<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";

definePageMeta({
  middleware: "auth",
});

const { user, fetchUser, fetchDenganAuth, logout } = useAuth();

// Form profil
const formProfil = reactive({
  nama: "",
  email: "",
});
const loadingProfil = ref(false);
const pesanProfil = ref("");
const errorProfil = ref("");

// Form password
const formPassword = reactive({
  passwordLama: "",
  passwordBaru: "",
  konfirmasi: "",
});
const loadingPassword = ref(false);
const pesanPassword = ref("");
const errorPassword = ref("");

onMounted(async () => {
  if (!user.value) {
    await fetchUser();
  }
  if (user.value) {
    formProfil.nama = user.value.nama;
    formProfil.email = user.value.email;
  }
});

async function simpanProfil() {
  pesanProfil.value = "";
  errorProfil.value = "";
  loadingProfil.value = true;

  try {
    await fetchDenganAuth("/api/auth/profil", {
      method: "PUT",
      body: { nama: formProfil.nama, email: formProfil.email },
    });
    await fetchUser();
    pesanProfil.value = "Profil berhasil diperbarui";
  } catch (err: any) {
    errorProfil.value = err?.data?.pesan || "Gagal memperbarui profil";
  } finally {
    loadingProfil.value = false;
  }
}

async function gantiPassword() {
  pesanPassword.value = "";
  errorPassword.value = "";

  if (formPassword.passwordBaru !== formPassword.konfirmasi) {
    errorPassword.value = "Konfirmasi password tidak cocok";
    return;
  }

  loadingPassword.value = true;

  try {
    await fetchDenganAuth("/api/auth/password", {
      method: "PUT",
      body: {
        passwordLama: formPassword.passwordLama,
        passwordBaru: formPassword.passwordBaru,
      },
    });
    pesanPassword.value = "Password berhasil diganti";
    formPassword.passwordLama = "";
    formPassword.passwordBaru = "";
    formPassword.konfirmasi = "";
  } catch (err: any) {
    errorPassword.value = err?.data?.pesan || "Gagal mengganti password";
  } finally {
    loadingPassword.value = false;
  }
}
</script>

<template>
  <div class="halaman">
    <header class="header">
      <div class="container header-inner">
        <NuxtLink to="/" class="brand">
          <i class="mdi mdi-newspaper mdi-24px"></i>
          <span class="logo">Portal Berita</span>
        </NuxtLink>

        <nav class="nav">
          <button
            class="nav-link logout"
            @click="logout"
            title="Logout"
            aria-label="Logout"
          >
            <i class="mdi mdi-logout-variant"></i>
            Logout
          </button>
        </nav>
      </div>
    </header>

    <main class="container main">
      <h1>👤 Profil Saya</h1>

      <div class="grid">
        <!-- Kartu Info -->
        <section class="kartu">
          <h2>Informasi Akun</h2>
          <div class="info-item">
            <span class="label">Nama</span>
            <span class="value">{{ user?.nama }}</span>
          </div>
          <div class="info-item">
            <span class="label">Email</span>
            <span class="value">{{ user?.email }}</span>
          </div>
          <div class="info-item">
            <span class="label">Role</span>
            <span class="badge" :class="user?.role">
              {{ user?.role === "admin" ? "👑 Admin" : "👤 User" }}
            </span>
          </div>
        </section>

        <!-- Form Edit Profil -->
        <section class="kartu">
          <h2>Edit Profil</h2>

          <div v-if="pesanProfil" class="notif sukses">{{ pesanProfil }}</div>
          <div v-if="errorProfil" class="notif error">{{ errorProfil }}</div>

          <form @submit.prevent="simpanProfil">
            <div class="field">
              <label>Nama</label>
              <input v-model="formProfil.nama" type="text" />
            </div>

            <div class="field">
              <label>Email</label>
              <input v-model="formProfil.email" type="email" />
            </div>

            <button type="submit" :disabled="loadingProfil" class="btn-primary">
              {{ loadingProfil ? "Menyimpan..." : "Simpan Perubahan" }}
            </button>
          </form>
        </section>

        <!-- Form Ganti Password -->
        <section class="kartu full">
          <h2>Ganti Password</h2>

          <div v-if="pesanPassword" class="notif sukses">
            {{ pesanPassword }}
          </div>
          <div v-if="errorPassword" class="notif error">
            {{ errorPassword }}
          </div>

          <form @submit.prevent="gantiPassword">
            <div class="field">
              <label>Password Lama</label>
              <input v-model="formPassword.passwordLama" type="password" />
            </div>

            <div class="field">
              <label>Password Baru</label>
              <input
                v-model="formPassword.passwordBaru"
                type="password"
                placeholder="Minimal 6 karakter"
              />
            </div>

            <div class="field">
              <label>Konfirmasi Password Baru</label>
              <input v-model="formPassword.konfirmasi" type="password" />
            </div>

            <button
              type="submit"
              :disabled="loadingPassword"
              class="btn-primary"
            >
              {{ loadingPassword ? "Memproses..." : "Ganti Password" }}
            </button>
          </form>
        </section>
      </div>
    </main>
  </div>
</template>

<style scoped>
.halaman {
  min-height: 100vh;
  background: #f9fafb;
  font-family: system-ui, sans-serif;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.25rem;
  width: 100%;
}

/* Header */
.header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
}
.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 1rem;
  padding-bottom: 1rem;
}
.brand {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  text-decoration: none;
  color: inherit;
}
.logo {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
  color: #111827;
}
.nav {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.nav-link {
  color: #374151;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 0.4rem 0.9rem;
  border-radius: 6px;
  border: none;
  background: transparent;
  cursor: pointer;
}
.nav-link:hover {
  background: #f3f4f6;
}
.nav-link.primary {
  background: #2563eb;
  color: white;
}
.nav-link.primary:hover {
  background: #1d4ed8;
}
.nav-link.logout {
  background: #dc2626;
  color: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.2rem 0.5rem;
  gap: 0.5rem;
  border-radius: 6px;
  border: none;
  cursor: pointer;
}
.nav-link.logout i {
  font-size: 1.25rem;
}
.nav-user {
  color: #6b7280;
  font-size: 0.9rem;
  text-decoration: none;
  padding: 0.4rem 0.75rem;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  line-height: 1;
}
.nav-user:hover {
  background: #f3f4f6;
  color: #111827;
}

/* Main */
.main {
  padding-top: 2rem;
  padding-bottom: 3rem;
}
h1 {
  margin: 0 0 0.25rem;
  font-size: 1.75rem;
  color: #111827;
}
.sub {
  color: #6b7280;
  margin: 0 0 2rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
}
.kartu {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}
.kartu.full {
  grid-column: 1 / -1;
}
.kartu h2 {
  margin: 0 0 1.25rem;
  font-size: 1.1rem;
  color: #111827;
}

/* Info */
.info-item {
  display: flex;
  justify-content: space-between;
  padding: 0.6rem 0;
  border-bottom: 1px solid #f3f4f6;
}
.info-item:last-child {
  border-bottom: none;
}
.label {
  color: #6b7280;
  font-size: 0.9rem;
}
.value {
  color: #111827;
  font-weight: 500;
}

.badge {
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
}
.badge.admin {
  background: #fef3c7;
  color: #92400e;
}
.badge.user {
  background: #dbeafe;
  color: #1e40af;
}

/* Form */
.field {
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
}
.field label {
  font-weight: 600;
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
  color: #374151;
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

.btn-primary {
  background: #2563eb;
  color: white;
  padding: 0.6rem 1.25rem;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  cursor: pointer;
}
.btn-primary:hover {
  background: #1d4ed8;
}
.btn-primary:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}

/* Notif */
.notif {
  padding: 0.75rem 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}
.notif.sukses {
  background: #d1fae5;
  color: #065f46;
}
.notif.error {
  background: #fee2e2;
  color: #991b1b;
}
</style>
