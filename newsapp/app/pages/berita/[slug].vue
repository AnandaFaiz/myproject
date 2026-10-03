<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";

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

interface Komentar {
  id: number;
  beritaId: number;
  userId: number;
  isi: string;
  createdAt: string;
  updatedAt: string;
  userNama: string;
}

//State
const route = useRoute();
const berita = ref<Berita | null>(null);
const semuaBerita = ref<Berita[]>([]);
const komentarList = ref<Komentar[]>([]);
const loadingKomentar = ref(false);
const isiKomentar = ref("");
const mengirim = ref(false);
const errorKomentar = ref("");
const hapusLoading = ref<number | null>(null);
const tersimpan = ref(false);
const loadingBookmark = ref(false);
const loading = ref(false);
const error = ref("");

function formatTanggal(tanggal: string) {
  return new Date(tanggal).toLocaleDateString("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

async function ambilKomentar(beritaId: number) {
  loadingKomentar.value = true;
  try {
    const res = await $fetch<{ sukses: boolean; data: Komentar[] }>(
      `/api/komentar/${beritaId}`,
    );
    komentarList.value = res.data;
  } catch (err: any) {
    console.error("Gagal mengambil komentar:", err);
  } finally {
    loadingKomentar.value = false;
  }
}

function formatWaktuRelatif(tanggal: string) {
  const sekarang = new Date().getTime();
  const waktu = new Date(tanggal).getTime();
  const selisihDetik = Math.floor((sekarang - waktu) / 1000);

  if (selisihDetik < 60) return `${selisihDetik} detik yang lalu`;
  if (selisihDetik < 3600)
    return `${Math.floor(selisihDetik / 60)} menit yang lalu`;
  if (selisihDetik < 86400)
    return `${Math.floor(selisihDetik / 3600)} jam yang lalu`;
  if (selisihDetik < 604800)
    return `${Math.floor(selisihDetik / 86400)} hari yang lalu`;
  return new Date(tanggal).toLocaleDateString("id-ID", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

async function kirimKomentar() {
  errorKomentar.value = "";
  if (!isiKomentar.value.trim()) {
    errorKomentar.value = "Komentar tidak boleh kosong";
    return;
  }

  if (!berita.value) return;

  mengirim.value = true;
  try {
    const token = useCookie<string | null>("auth_token");
    const res = await $fetch<{ sukses: boolean; data: Komentar }>(
      "/api/komentar",
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
        body: {
          beritaId: berita.value.id,
          isi: isiKomentar.value.trim(),
        },
      },
    );

    // Tambahkan komentar baru ke daftar komentar
    komentarList.value.unshift(res.data);

    //reset form
    isiKomentar.value = "";
  } catch (err: any) {
    errorKomentar.value = err?.data?.pesan || "Gagal mengirim komentar";
  } finally {
    mengirim.value = false;
  }
}

//ambil inisial dari nama
function inisial(nama: string) {
  return nama
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}

async function cekBookmark(beritaId: number) {
  if (!isLoggedIn.value) {
    tersimpan.value = false;
    return;
  }

  try {
    const token = useCookie<string | null>("auth_token");
    const res = await $fetch<{ sukses: boolean; tersimpan: boolean }>(
      `/api/bookmark/cek/${beritaId}`,
      {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      },
    );
    tersimpan.value = res.tersimpan;
  } catch (err) {
    console.error("Gagal cek bookmark:", err);
    tersimpan.value = false;
  }
}

async function toggleBookmark() {
  if (!berita.value) return;
  if (!isLoggedIn.value) {
    await navigateTo(`/login?redirect=${encodeURIComponent(route.fullPath)}`);
    return;
  }
  loadingBookmark.value = true;

  try {
    const token = useCookie<string | null>("auth_token");

    if (tersimpan.value) {
      await $fetch(`/api/bookmark/${berita.value.id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      });
      tersimpan.value = false;
    } else {
      await $fetch("/api/bookmark", {
        method: "POST",
        headers: {
          Authorization: `Bearer $(token.value)`,
        },
        body: {
          beritaId: berita.value.id,
        },
      });
      tersimpan.value = true;
    }
  } catch (err: any) {
    alert(err?.data?.pesan || "Gagal memperbarui bookmark");
  } finally {
    loadingBookmark.value = false;
  }
}

//hapus komentar
function bisaHapus(komentar: Komentar): boolean {
  if (!user.value) return false;
  return user.value.role === "admin" || user.value.id === komentar.userId;
}

async function hapusKomentar(id: number) {
  if (!confirm("Apakah Anda yakin ingin menghapus komentar ini?")) return;

  hapusLoading.value = id;
  try {
    const token = useCookie<string | null>("auth_token");

    await $fetch(`/api/komentar/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token.value}`,
      },
    });

    //hapus komentar dari daftar
    komentarList.value = komentarList.value.filter((k) => k.id !== id);
  } catch (err: any) {
    alert(err?.data?.pesan || "Gagal menghapus komentar");
  } finally {
    hapusLoading.value = null;
  }
}

//Filter : ambil berita lainnya selain berita yang dibaca user
const beritaLain = computed(() => {
  if (!berita.value) return [];
  return semuaBerita.value.filter((b) => b.id !== berita.value!.id).slice(0, 5);
});

const { user, isLoggedIn, isAdmin, fetchUser, logout } = useAuth();

onMounted(async () => {
  loading.value = true;
  try {
    if (!user.value) {
      await fetchUser();
    }

    //Fetch berita yang dibaca
    const res = await $fetch<{ sukses: boolean; data: Berita }>(
      `/api/berita/slug/${route.params.slug}`,
    );
    berita.value = res.data;

    const all = await $fetch<{ sukses: boolean; data: Berita[] }>(
      `/api/berita`,
    );
    semuaBerita.value = all.data;

    //ambil komentar setelah berita berhasil diambil
    if (berita.value) {
      await ambilKomentar(berita.value.id);
      await cekBookmark(berita.value.id);
    }
  } catch (err: any) {
    error.value = err?.data?.pesan || "Berita tidak ditemukan";
  } finally {
    loading.value = false;
  }
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
            <NuxtLink
              :to="`/login?redirect=${encodeURIComponent(route.fullPath)}`"
              class="nav-link primary"
              >Login</NuxtLink
            >
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

    <!--Main Content-->
    <main class="container">
      <p v-if="loading" class="info">Memuat...</p>
      <p v-else-if="error" class="info error">{{ error }}</p>

      <div v-else-if="berita" class="layout">
        <!--Berita Utama-->
        <article class="artikel">
          <img
            v-if="berita.gambar"
            :src="berita.gambar"
            :alt="berita.judul"
            class="gambar-utama"
          />
          <div v-if="berita.kategori" class="kategori-badge">
            {{ berita.kategori }}
          </div>

          <h1 class="judul">{{ berita.judul }}</h1>

          <div class="meta">
            <span
              >Oleh <strong>{{ berita.penulis }}</strong></span
            >
            <span class="pemisah">.</span>
            <span>{{ formatTanggal(berita.createdAt) }}</span>
          </div>

          <!--Bookmark-->
          <div class="aksi-berita">
            <button
              class="btn-bookmark"
              :class="{ tersimpan }"
              :disabled="loadingBookmark"
              @click="toggleBookmark"
            >
              <i
                class="loadingBookmark ? 'mdi mdi-loading mdi-spin' :tersimpan ? 'mdi mdi-bookmark' : 'mdi mdi-bookmark-outline'"
              ></i>
              <span>{{ tersimpan ? "Tersimpan" : "Bookmark" }}</span>
            </button>
          </div>

          <div class="konten">
            <p v-for="(paragraf, i) in berita.konten.split('\n')" :key="i">
              {{ paragraf }}
            </p>
          </div>

          <!--Comments Section-->
          <section class="komentar-section">
            <h2 class="komentar-judul">
              💬 Komentar
              <span class="komentar-count">{{ komentarList.length }}</span>
            </h2>

            <div class="komentar-form-wrapper">
              <!-- Kalau belum login -->
              <div v-if="!isLoggedIn" class="komentar-login-prompt">
                <p>Ingin berkomentar?</p>
                <NuxtLink
                  :to="`/login?redirect=${encodeURIComponent(route.fullPath)}`"
                  class="btn-login-komentar"
                >
                  <i class="mdi mdi-login"></i>
                  Login
                </NuxtLink>
              </div>

              <!-- Kalau sudah login -->
              <div v-else class="komentar-form">
                <div class="komentar-form-header">
                  <div class="komentar-avatar">
                    {{ inisial(user?.nama || "") }}
                  </div>
                  <span class="komentar-form-nama">{{ user?.nama }}</span>
                </div>

                <textarea
                  v-model="isiKomentar"
                  placeholder="Tulis komentar..."
                  rows="3"
                  maxlength="1000"
                  :disabled="mengirim"
                  @keydown.ctrl.enter="kirimKomentar"
                ></textarea>

                <div class="komentar-form-footer">
                  <span class="komentar-form-hint"
                    >Ctrl + Enter untuk kirim</span
                  >
                  <div class="komentar-form-actions">
                    <span
                      v-if="isiKomentar.length > 0"
                      class="komentar-char-count"
                    >
                      {{ isiKomentar.length }}/1000
                    </span>
                    <button
                      class="btn-kirim"
                      :disabled="mengirim || !isiKomentar.trim()"
                      @click="kirimKomentar"
                    >
                      <i v-if="mengirim" class="mdi mdi-loading mdi-spin"></i>
                      <i v-else class="mdi mdi-send"></i>
                      {{ mengirim ? "Mengirim..." : "Kirim" }}
                    </button>
                  </div>
                </div>

                <p v-if="errorKomentar" class="komentar-error">
                  {{ errorKomentar }}
                </p>
              </div>
            </div>

            <p v-if="loadingKomentar" class="komentar-info">
              Memuat komentar...
            </p>
            <p v-else-if="komentarList.length === 0" class="komentar-info">
              Belum ada komentar.
            </p>

            <div v-else class="komentar-list">
              <div
                v-for="kom in komentarList"
                :key="kom.id"
                class="komentar-item"
              >
                <div class="komentar-avatar">{{ inisial(kom.userNama) }}</div>
                <div class="komentar-body">
                  <div class="komentar-header">
                    <span class="komentar-nama">{{ kom.userNama }}</span>
                    <span class="komentar-waktu">{{
                      formatWaktuRelatif(kom.createdAt)
                    }}</span>

                    <!--Tombol Hapus Komentar-->
                    <button
                      v-if="bisaHapus(kom)"
                      class="btn-hapus-komentar"
                      :disabled="hapusLoading === kom.id"
                      @click="hapusKomentar(kom.id)"
                    >
                      <i
                        v-if="hapusLoading === kom.id"
                        class="mdi mdi-loading mdi-spin"
                      ></i>
                      <i v-else class="mdi mdi-delete"></i>
                    </button>
                  </div>
                  <p class="komentar-isi">{{ kom.isi }}</p>
                </div>
              </div>
            </div>
          </section>
        </article>

        <aside>
          <h2 class="sidebar-judul">Berita Lainnya</h2>
          <div v-if="beritaLain.length === 0" class="sidebar-kosong">
            Belum ada berita lain.
          </div>

          <NuxtLink
            v-for="item in beritaLain"
            :key="item.id"
            :to="`/berita/${item.slug}`"
            class="sidebar-item"
          >
            <div class="sidebar-thumb">
              <img v-if="item.gambar" :src="item.gambar" :alt="item.judul" />
              <div v-else class="sidebar-thumb-placeholder">
                <i class="mdi mdi-newspaper"></i>
              </div>
            </div>
            <div class="sidebar-isi">
              <span v-if="item.kategori" class="sidebar-kategori">
                {{ item.kategori }}
              </span>
              <h3 class="sidebar-item-judul">
                {{ item.judul }}
              </h3>
              <span class="sidebar-tanggal">
                {{ formatTanggal(item.createdAt) }}
              </span>
            </div>
          </NuxtLink>
        </aside>
      </div>
    </main>

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
  gap: 0.5rem;
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

.info {
  color: #6b7280;
  padding: 2rem 0;
  text-align: center;
}
.info.error {
  color: #dc2626;
}

.layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 900px) {
  .layout {
    grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
  }
}

.artikel {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}
.kembali {
  display: inline-block;
  color: #2563eb;
  text-decoration: none;
  font-size: 0.9rem;
  margin-bottom: 1.5rem;
}
.kembali:hover {
  text-decoration: underline;
}

.gambar-utama {
  width: 100%;
  max-height: 420px;
  object-fit: cover;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.kategori-badge {
  display: inline-block;
  background: #eff6ff;
  color: #2563eb;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 1rem;
}

.judul {
  font-size: 2rem;
  font-weight: 800;
  color: #111827;
  margin: 0 0 1rem;
  line-height: 1.25;
}

.meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.4rem;
  color: #6b7280;
  font-size: 0.9rem;
  padding-bottom: 1.25rem;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 1.5rem;
}
.meta .pemisah {
  color: #d1d5db;
}

/* Tombol bookmark */
.aksi-berita {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1.5rem;
}

.btn-bookmark {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-bookmark:hover:not(:disabled) {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.btn-bookmark:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-bookmark.tersimpan {
  background: #2563eb;
  border-color: #2563eb;
  color: white;
}

.btn-bookmark.tersimpan:hover:not(:disabled) {
  background: #1d4ed8;
  border-color: #1d4ed8;
}

.btn-bookmark i {
  font-size: 1.1rem;
}

.konten {
  font-size: 1rem;
  line-height: 1.8;
  color: #374151;
}
.konten p {
  margin: 0 0 1.25rem;
  text-align: justify;
  line-height: 1.8;
}

/* Section Komentar */
.komentar-section {
  margin-top: 2.5rem;
  padding-top: 2rem;
  border-top: 2px solid #e5e7eb;
}

.komentar-judul {
  font-size: 1.25rem;
  font-weight: 700;
  color: #111827;
  margin: 0 0 1.25rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.komentar-count {
  color: #6b7280;
  font-weight: 500;
  font-size: 1rem;
}

.komentar-info {
  color: #9ca3af;
  font-size: 0.9rem;
  padding: 1.5rem 0;
  text-align: center;
  font-style: italic;
}

.komentar-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.komentar-item {
  display: flex;
  gap: 0.75rem;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 8px;
}

.komentar-avatar {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #2563eb;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.85rem;
  letter-spacing: 0.5px;
}

.komentar-body {
  flex: 1;
  min-width: 0;
}

.komentar-header {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.35rem;
}

.komentar-nama {
  font-weight: 600;
  color: #111827;
  font-size: 0.95rem;
}

.komentar-waktu {
  color: #9ca3af;
  font-size: 0.8rem;
}

.komentar-isi {
  color: #374151;
  font-size: 0.95rem;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}
/* Form Komentar Wrapper */
.komentar-form-wrapper {
  margin-bottom: 2rem;
}

/* Prompt Login */
.komentar-login-prompt {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  margin-bottom: 1rem;
}
.komentar-login-prompt p {
  margin: 0;
  color: #1e40af;
  font-weight: 500;
}
.btn-login-komentar {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: #2563eb;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
}
.btn-login-komentar:hover {
  background: #1d4ed8;
}

.komentar-form {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}
.komentar-form-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}
.komentar-form-nama {
  font-weight: 600;
  color: #111827;
  font-size: 0.9rem;
}
.komentar-form textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
  resize: vertical;
  min-height: 80px;
  box-sizing: border-box;
}
.komentar-form textarea:focus {
  outline: 2px solid #2563eb;
  border-color: transparent;
}
.komentar-form textarea:disabled {
  background: #f3f4f6;
  cursor: not-allowed;
}
.komentar-form-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  margin-top: 0.75rem;
  flex-wrap: wrap;
}
.komentar-form-hint {
  font-size: 0.75rem;
  color: #9ca3af;
}
.komentar-form-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.komentar-char-count {
  font-size: 0.75rem;
  color: #9ca3af;
}
.btn-kirim {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: #2563eb;
  color: white;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
}
.btn-kirim:hover:not(:disabled) {
  background: #1d4ed8;
}
.btn-kirim:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}
.komentar-error {
  color: #dc2626;
  font-size: 0.85rem;
  margin: 0.75rem 0 0;
}
/* Tombol hapus komentar */
.btn-hapus-komentar {
  margin-left: auto; /* dorong ke kanan */
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 1.1rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.btn-hapus-komentar:hover:not(:disabled) {
  background: #fee2e2;
  color: #dc2626;
}

.btn-hapus-komentar:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

/* Loading spin */
.mdi-spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.sidebar {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  align-self: start;
  position: sticky;
  top: 80px;
  max-height: calc(100vh - 100px);
  overflow-y: auto;
}
.sidebar-judul {
  font-size: 1.1rem;
  font-weight: 700;
  color: #111827;
  margin: 0 0 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid #2563eb;
  display: inline-block;
}

.sidebar-kosong {
  color: #6b7280;
  font-size: 0.9rem;
  padding: 1rem 0;
}

.sidebar-item {
  display: flex;
  gap: 0.75rem;
  padding: 0.75rem 0;
  text-decoration: none;
  color: inherit;
  border-bottom: 1px solid #f3f4f6;
  transition: opacity 0.15s;
}
.sidebar-item:last-child {
  border-bottom: none;
}
.sidebar-item:hover {
  opacity: 0.75;
}

.sidebar-thumb {
  flex-shrink: 0;
  width: 80px;
  height: 60px;
  border-radius: 6px;
  overflow: hidden;
  background: #f3f4f6;
}
.sidebar-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.sidebar-thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  background: #e5e7eb;
}
.sidebar-thumb-placeholder i {
  font-size: 1.5rem;
}

.sidebar-isi {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.sidebar-kategori {
  font-size: 0.7rem;
  font-weight: 700;
  color: #2563eb;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 0.25rem;
}

.sidebar-item-judul {
  font-size: 0.9rem;
  font-weight: 600;
  color: #111827;
  margin: 0 0 0.35rem;
  line-height: 1.35;
  display: -webkit-box;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.sidebar-tanggal {
  font-size: 0.75rem;
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
