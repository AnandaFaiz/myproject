import { fileURLToPath } from 'node:url'

const rootDir = fileURLToPath(new URL('.', import.meta.url))
const appDir = fileURLToPath(new URL('./app', import.meta.url))

export default defineNuxtConfig({
  alias: {
    '~~': rootDir,
    '@@': rootDir,
    '~': appDir,
    '@': appDir,
  },

  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  
  css: [
    '@mdi/font/css/materialdesignicons.min.css',
  ],

  modules: ['@nuxtjs/tailwindcss', '@nuxt/icon'],

  nitro:{
    devProxy:{
      '/api':{
        target: 'http://localhost:8080/api',
        changeOrigin: true,
      },
      '/uploads':{
        target: 'http://localhost:8080/uploads',
        changeOrigin: true,
      },
    },
  },
})