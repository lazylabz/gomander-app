import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
	viteConfig,
	defineConfig({
		test: {
			environment: "jsdom",
			include: ["src/**/*.test.ts", "src/**/*.test.tsx"],
			coverage: {
				provider: "v8",
				include: ["src/**"],
				exclude: ["src/testing/**", "src/contracts/adapters/inMemory.ts"],
			},
		},
	}),
);
