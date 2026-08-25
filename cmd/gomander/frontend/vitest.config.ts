import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
	viteConfig,
	defineConfig({
		test: {
			// happy-dom over jsdom: jsdom 30 pulls undici 8, which needs
			// worker_threads.markAsUncloneable and so refuses to run on the Node 20 CI uses.
			environment: "happy-dom",
			include: ["src/**/*.test.ts", "src/**/*.test.tsx"],
			coverage: {
				provider: "v8",
				include: ["src/**"],
				exclude: ["src/testing/**", "src/contracts/adapters/inMemory.ts"],
			},
		},
	}),
);
