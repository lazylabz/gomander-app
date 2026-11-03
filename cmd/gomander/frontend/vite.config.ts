import path from "node:path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react-swc";
import wails from "@wailsio/runtime/plugins/vite";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
	// Wails proxies the dev window to 127.0.0.1, which "localhost" may not bind.
	server: {
		host: "127.0.0.1",
		port: Number(process.env.WAILS_VITE_PORT) || 9245,
		strictPort: true,
	},
	// Under Vitest @wailsio/runtime resolves to a stub the typed-events plugin can't extend.
	plugins: [
		react(),
		tailwindcss(),
		...(process.env.VITEST ? [] : [wails("./bindings")]),
	],
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src"),
		},
	},
});
