/// <reference types="vitest" />
import path from "path"
import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"

import tailwind from "tailwindcss"
import autoprefixer from "autoprefixer"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()],
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
    extensions: ['.ts', '.vue', '.json']
  },
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      provider: 'istanbul',
      exclude: [
        "tailwind.config.js",
        "src/lib/routes",
        "src/main.ts",
        "src/App.vue",
        "src/components/NavigationBar.vue",
        "src/components/ui/calendar/Calendar.vue",
        "src/components/ui/calendar/index.ts",
        "src/components/DateRangePicker.vue",
      ],
    },
  },
})
