import type { EnvironmentPath, UserConfig } from "@/contracts/types.ts";

export type UserConfigBuilder = {
	withLastOpenedProjectId: (lastOpenedProjectId: string) => UserConfigBuilder;
	withEnvironmentPaths: (
		...environmentPaths: EnvironmentPath[]
	) => UserConfigBuilder;
	withLocale: (locale: string) => UserConfigBuilder;
	build: () => UserConfig;
};

const builder = (data: UserConfig): UserConfigBuilder => ({
	withLastOpenedProjectId: (lastOpenedProjectId) =>
		builder({ ...data, lastOpenedProjectId }),
	withEnvironmentPaths: (...environmentPaths) =>
		builder({ ...data, environmentPaths }),
	withLocale: (locale) => builder({ ...data, locale }),
	build: () => ({ ...data, environmentPaths: [...data.environmentPaths] }),
});

export const newUserConfigBuilder = (): UserConfigBuilder =>
	builder({
		lastOpenedProjectId: "",
		environmentPaths: [],
		locale: "en",
	});
