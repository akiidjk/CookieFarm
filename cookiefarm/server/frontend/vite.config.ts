import { defineConfig } from "vite";
import react from '@vitejs/plugin-react-swc'
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  build: {
    target: "esnext",
    minify: 'esbuild',
    terserOptions: {
      compress: {
        drop_console: true,
      },
    },
    sourcemap: false,
    chunkSizeWarningLimit: 1200,
  },
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  server: {
    host: "0.0.0.0",
    port: 5173,
    proxy: {
      "/api/v1": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
      },
    },
  },
});
