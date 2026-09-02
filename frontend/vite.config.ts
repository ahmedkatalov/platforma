import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { "@": path.resolve(__dirname, "./src") },
  },
  build: {
    // Поддержка не самых свежих телефонов (iOS 14+, Android Chrome 87+):
    // современный синтаксис транспилируется, чтобы код не падал при разборе.
    target: ["es2020", "safari14", "chrome87", "firefox78"],
    rollupOptions: {
      output: {
        // Тяжёлые библиотеки — отдельными чанками, чтобы они кэшировались отдельно от кода приложения.
        manualChunks: {
          react: ["react", "react-dom", "react-router-dom"],
          redux: ["@reduxjs/toolkit", "react-redux"],
          charts: ["recharts"],
          markdown: ["react-markdown", "remark-gfm"],
        },
      },
    },
  },
  server: {
    host: "0.0.0.0",
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8090",
        changeOrigin: true,
      },
    },
  },
});
