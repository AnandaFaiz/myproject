<script setup lang="ts">
definePageMeta({
  middleware: "auth",
});

import { ref, reactive, onMounted } from "vue";

// News data type
interface Berita {
  id: number;
  judul: string;
  slug: string;
  konten: string;
  penulis: string;
  kategori: string | null;
  gambar: string | null;
  createdAt: Date | string;
  updatedAt: Date | string;
}

const daftarBerita = ref<Berita[]>([]);
const loading = ref(false);
const pesanSukses = ref("");
const pesanError = ref("");
const { logout } = useAuth();

async function logoutAdmin() {
  await logout();
}

//State
const form = reactive({
  id: null as number | null,
  judul: "",
  konten: "",
  penulis: "",
  kategori: "",
  gambar: "",
});

const uploading = ref(false);
const previewGambar = ref("");
const modeEdit = ref(false);
const { user, fetchUser, fetchDenganAuth } = useAuth();

// Read : ambil semua berita
async function ambilBerita() {
  loading.value = true;
  try {
    const res = await $fetch<{ sukses: boolean; data: Berita[] }>(
      "/api/berita",
    );
    daftarBerita.value = res.data;
  } catch (err) {
    pesanError.value = "Gagal mengambil data berita";
  } finally {
    loading.value = false;
  }
}

//Create/Update
async function simpanBerita() {
  pesanSukses.value = "";
  pesanError.value = "";

  if (!form.judul || !form.konten || !form.penulis) {
    pesanError.value = "Judul, Konten, dan Penulis wajib diisi";
    return;
  }

  try {
    const body = {
      judul: form.judul,
      konten: form.konten,
      penulis: form.penulis,
      kategori: form.kategori,
      gambar: previewGambar.value || "",
    };

    if (modeEdit.value && form.id) {
      //Update
      await fetchDenganAuth(`/api/berita/${form.id}`, {
        method: "PUT",
        body,
      });
      pesanSukses.value = "Berita berhasil diperbarui";
    } else {
      //Create
      await fetchDenganAuth("/api/berita", {
        method: "POST",
        body,
      });
      pesanSukses.value = "Berita berhasil ditambahkan";
    }

    resetForm();

    await ambilBerita();
  } catch (err: any) {
    pesanError.value = err?.data?.pesan || "Gagal menyimpan berita";
  }
}

//Delete
async function hapusBerita(id: number) {
  if (!confirm("Yakin mau hapus berita ini?")) return;

  try {
    await fetchDenganAuth(`/api/berita/${id}`, { method: "DELETE" });
    pesanSukses.value = "Berita berhasil dihapus";
    await ambilBerita();
  } catch (err: any) {
    pesanError.value = err?.data?.pesan || "Gagal menghapus berita";
  }
}

//Upload Gambar
async function uploadGambar(e: Event) {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  uploading.value = true;
  try {
    const fd = new FormData();
    fd.append("file", file);

    //send to BE Go via Proxy
    const res = await fetchDenganAuth<{ sukses: boolean; url: string }>(
      "/api/upload",
      {
        method: "POST",
        body: fd,
      },
    );

    previewGambar.value = res.url;
  } catch (err: any) {
    pesanError.value = err?.data?.pesan || "Gagal upload gambar";
  } finally {
    uploading.value = false;
  }
}

//isi form untuk edit
function editBerita(item: Berita) {
  console.log("edit berita ID:", item.id);
  modeEdit.value = true;
  form.id = item.id;
  form.judul = item.judul;
  form.konten = item.konten;
  form.penulis = item.penulis;
  form.kategori = item.kategori ?? "";
  previewGambar.value = item.gambar ?? "";
  window.scrollTo({ top: 0, behavior: "smooth" });
}

//reset form
function resetForm() {
  modeEdit.value = false;
  form.id = null;
  form.judul = "";
  form.konten = "";
  form.penulis = "";
  form.kategori = "";
  form.gambar = "";
  previewGambar.value = "";
}

//ambil data saat halaman pertama dibuka
onMounted(() => {
  ambilBerita();
});
</script>

<template>
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

  <div class="container">
    <div class="header-admin">
      <h1>Admin - Kelola Berita</h1>
    </div>

    <!--Notif-->
    <div v-if="pesanSukses" class="notif sukses">{{ pesanSukses }}</div>
    <div v-if="pesanError" class="notif error">{{ pesanError }}</div>

    <!--Form-->
    <section class="form-section">
      <h2>{{ modeEdit ? "Edit Berita" : "Tambah Berita" }}</h2>
      <form @submit.prevent="simpanBerita">
        <div class="field">
          <label for="">Judul</label>
          <input v-model="form.judul" type="text" placeholder="Judul berita" />
        </div>

        <div class="field">
          <label for="">Penulis</label>
          <input
            v-model="form.penulis"
            type="text"
            placeholder="Nama penulis"
          />
        </div>

        <div class="field">
          <label for="">Kategori</label>
          <input
            v-model="form.kategori"
            type="text"
            placeholder="Kategori (opsional)"
          />
        </div>

        <div class="field">
          <label for="">Gambar</label>
          <input
            type="file"
            accept="image/*"
            @change="uploadGambar"
            :disabled="uploading"
          />
          <p v-if="uploading" class="info-upload">Mengupload...</p>
          <img
            v-if="previewGambar"
            :src="previewGambar"
            alt="Preview"
            class="preview"
          />
        </div>

        <div class="field">
          <label for="">Konten</label>
          <textarea
            v-model="form.konten"
            rows="6"
            placeholder="Isi berita"
          ></textarea>
        </div>

        <div class="actions">
          <button type="submit" class="btn-primary">
            {{ modeEdit ? "Update" : "Simpan" }}
          </button>
          <button
            v-if="modeEdit"
            type="button"
            class="btn-secondary"
            @click="resetForm"
          >
            Batal
          </button>
        </div>
      </form>
    </section>

    <!--LIST-->
    <section class="list-section">
      <h2>Daftar Berita</h2>

      <p v-if="loading">Memuat...</p>
      <p v-else-if="daftarBerita.length === 0">Belum ada berita</p>

      <table v-else>
        <thead>
          <tr>
            <th>Gambar</th>
            <th>Judul</th>
            <th>Penulis</th>
            <th>Kategori</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in daftarBerita" :key="item.id">
            <td>
              <img v-if="item.gambar" :src="item.gambar" alt="" class="thumb" />
              <div v-else class="thumb-placeholder"></div>
            </td>
            <td>{{ item.judul }}</td>
            <td>{{ item.penulis }}</td>
            <td>{{ item.kategori }}</td>
            <td>
              <button class="btn-edit" @click="editBerita(item)">Edit</button>
              <button class="btn-hapus" @click="hapusBerita(item.id)">
                Hapus
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>

<style scoped>
.container {
  max-width: 1300px;
  margin: 0 auto;
  padding: 0 1.25rem;
  font-family: system-ui, sans-serif;
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
h1 {
  margin-bottom: 1rem;
}
h2 {
  margin: 1.5rem 0 1rem;
  font-size: 1.2rem;
}
.notif {
  padding: 0.75rem 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
}
.notif.sukses {
  background: #d1fae5;
  color: #065f46;
}
.notif.error {
  background: #fee2e2;
  color: #991b1b;
}

.form-section,
.list-section {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.25rem;
  margin-bottom: 1.5rem;
}

.field {
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
}
.field label {
  font-weight: 600;
  margin-bottom: 0.25rem;
}
.field input,
.field textarea {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
}

.preview {
  margin-top: 0.5rem;
  max-width: 200px;
  border-radius: 8px;
}

.thumb {
  width: 60px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
}

.thumb-placeholder {
  width: 60px;
  height: 40px;
  background: #e5e7eb;
  border-radius: 4px;
}

.info-upload {
  font-size: 0.85rem;
  color: #6b7280;
  margin: 0.25rem 0 0;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 0.95rem;
}
.btn-primary {
  background: #2563eb;
  color: white;
}
.btn-primary:hover {
  background: #1d4ed8;
}
.btn-secondary {
  background: #e5e7eb;
  color: #374151;
}
.btn-edit {
  background: #f59e0b;
  color: white;
  margin-right: 0.5rem;
}
.btn-hapus {
  background: #dc2626;
  color: white;
}

table {
  width: 100%;
  border-collapse: collapse;
}
th,
td {
  text-align: left;
  padding: 0.5rem;
  border-bottom: 1px solid #e5e7eb;
}
th {
  background: #f3f4f6;
}
</style>
