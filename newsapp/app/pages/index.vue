<script setup lang="ts">
import { ref, onMounted } from "vue";

interface Berita {
  id: number;
  judul: string;
  slug: string;
  konten: string;
  penulis: string;
  kategori: string | null;
  gambar: string | null;
  createdAt: string;
}

const daftarBerita = ref<Berita[]>([]);
const loading = ref(false);
const error = ref("");

async function ambilBerita() {
  loading.value = true;
  try {
    const res = await $fetch<{ sukses: boolean; data: Berita[] }>(
      "/api/berita",
    );
    daftarBerita.value = res.data;
  } catch (err: any) {
    error.value = err?.data?.statusMessage || "Gagal memuat berita";
  } finally {
    loading.value = false;
  }
}

//function ambil cuplikan konten
function cuplikan(teks: string, panjang = 120) {
  if (teks.length <= panjang) return teks;
  return teks.slice(0, panjang).trim() + "...";
}

//function format tanggal
function formatTanggal(tanggal: string) {
  return new Date(tanggal).toLocaleDateString("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

//warna placeholder per kategori
function warnaKategori(kategori: string | null) {
  const peta: Record<string, string> = {
    Teknologi: "#3b82f6",
    Tutorial: "#10b981",
    Review: "#f59e0b",
    Umum: "8b5cf6",
  };
  return peta[kategori ?? ""] || "#6b7280";
}

const { user, isLoggedIn, isAdmin, fetchUser, logout } = useAuth();

onMounted(async () => {
  await fetchUser();
  ambilBerita();
});
</script>

<template>
  <div class="halaman">
    <!--Header-->
    <header class="header">
      <div class="container header-inner">
        <NuxtLink to="/" class="brand">
          <i class="mdi mdi-newspaper mdi-24px"></i>
          <span class="logo">Portal Berita</span>
        </NuxtLink>

        <nav class="nav">
          <!--Login-->
          <template v-if="!isLoggedIn">
            <NuxtLink to="/login" class="nav-link primary"> Login </NuxtLink>
          </template>

          <template v-else>
            <NuxtLink v-if="isAdmin" to="/admin/berita" class="nav-link">
              Dashboard Admin
            </NuxtLink>
            <NuxtLink to="/profil" class="nav-link nav-user">
              <i class="mdi mdi-account mdi-24px"></i>
              {{ user?.nama }}
            </NuxtLink>
            <button
              class="nav-link logout"
              @click="logout"
              title="Logout"
              aria-label="Logout"
            >
              <i class="mdi mdi-logout-variant"></i>
              Logout
            </button>
          </template>
        </nav>
      </div>
    </header>

    <!--Main-->
    <main class="container">
      <h2 class="judul-section">Berita Terbaru</h2>
      <p v-if="loading" class="info">Memuat berita...</p>
      <p v-else-if="error" class="info error">{{ error }}</p>
      <p v-else-if="daftarBerita.length === 0" class="info">
        Belum ada berita.
        <NuxtLink to="/admin/berita"> Tambah disini </NuxtLink>
      </p>
      <div v-else class="grid-berita">
        <NuxtLink
          v-for="item in daftarBerita"
          :key="item.id"
          :to="`/berita/${item.slug}`"
          class="kartu"
        >
          <!--Placeholder gambar-->
          <div
            class="kartu-gambar"
            :style="
              item.gambar ? {} : { background: warnaKategori(item.kategori) }
            "
          >
            <img
              v-if="item.gambar"
              :src="item.gambar"
              :alt="item.judul"
              class="kartu-img"
            />
            <span class="kartu-kategori">{{ item.kategori || "Umum" }}</span>
          </div>
          <div class="kartu-isi">
            <h3 class="kartu-judul">{{ item.judul }}</h3>
            <p class="kartu-cuplikan">{{ cuplikan(item.konten) }}</p>
            <div class="kartu-meta">
              <span>{{ item.penulis }}</span>
              <span>.</span>
              <span>{{ formatTanggal(item.createdAt) }}</span>
            </div>
          </div>
        </NuxtLink>
      </div>
    </main>

    <!--Footer-->
    <footer class="footer">
      <div class="container">
        <p>© 2026 Portal Berita - WkwkLand</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.halaman {
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
  background: #f9fafb;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.25rem;
  width: 100%;
}
.header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 10;
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
  color: #111827;
  margin: 0;
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
main {
  flex: 1;
  padding-top: 2rem;
  padding-bottom: 3rem;
}

.judul-section {
  font-size: 1.5rem;
  margin: 0 0 1.5rem;
  color: #111827;
}

.info {
  color: #6b7280;
  padding: 2rem 0;
  text-align: center;
}
.info.error {
  color: #dc2626;
}
.grid-berita {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}
.kartu {
  background: white;
  border-radius: 10px;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition:
    transform 0.2s,
    box-shadow 0.2s;
  display: flex;
  flex-direction: column;
}
.kartu:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
}

.kartu-gambar {
  position: relative;
  height: 200px;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding: 0.75rem;
}
.kartu-img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.kartu-kategori {
  background: rgba(255, 255, 255, 0.95);
  color: #111827;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.kartu-isi {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  flex: 1;
}
.kartu-judul {
  font-size: 1.05rem;
  margin: 0 0 0.5rem;
  color: #111827;
  line-height: 1.4;
}
.kartu-cuplikan {
  color: #6b7280;
  font-size: 0.9rem;
  line-height: 1.5;
  margin: 0 0 1rem;
  flex: 1;
}
.kartu-meta {
  display: flex;
  gap: 0.4rem;
  font-size: 0.8rem;
  color: #9ca3af;
}
.footer {
  background: white;
  border-top: 1px solid #e5e7eb;
  padding: 1.25rem 0;
  color: #6b7280;
  font-size: 0.85rem;
  text-align: center;
}
</style>
