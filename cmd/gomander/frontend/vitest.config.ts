import path from "node:path";
import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
	viteConfig,
	defineConfig({
		resolve: {
			alias: {
				"@wailsio/runtime": path.resolve(
					__dirname,
					"src/testing/wailsRuntime.ts",
				),
			},
		},
		test: {
			environment: "happy-dom",
			setupFiles: ["src/testing/setup.ts"],
			include: ["src/**/*.test.ts", "src/**/*.test.tsx"],
			coverage: {
				provider: "v8",
				include: ["src/**"],
				exclude: [
					"src/testing/**",
					"src/contracts/adapters/inMemory.ts",
					"src/commandOutput/adapters/recording.ts",
				],
			},
		},
	}),
);
