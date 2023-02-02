import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { resolve as pathResolve } from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  optimizeDeps: { exclude: ["svelte-navigator"] },
  resolve: {
    alias: {
      "~": pathResolve(__dirname, "src"),
    },
  },
});
