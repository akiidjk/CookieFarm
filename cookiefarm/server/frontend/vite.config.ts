import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [["babel-plugin-react-compiler", {}]],
      },
    }),
    tailwindcss(),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes("node_modules")) {
            return undefined;
          }

          if (id.includes("echarts")) {
            return "vendor-echarts";
          }

          if (id.includes("@cloudflare/kumo")) {
            return "vendor-kumo";
          }

          if (id.includes("react-router")) {
            return "vendor-router";
          }

          if (id.includes("react")) {
            return "vendor-react";
          }

          if (id.includes("msw") || id.includes("ansi-to-html")) {
            return "vendor-ops";
          }

          return "vendor";
        },
      },
    },
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
