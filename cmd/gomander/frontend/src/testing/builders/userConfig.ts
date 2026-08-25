import type { EnvironmentPath, UserConfig } from "@/contracts/types.ts";

export class UserConfigBuilder {
	private data: UserConfig = {
		lastOpenedProjectId: "",
		environmentPaths: [],
		locale: "en",
	};

	withLastOpenedProjectId(lastOpenedProjectId: string): this {
		this.data.lastOpenedProjectId = lastOpenedProjectId;
		return this;
	}

	withEnvironmentPaths(...environmentPaths: EnvironmentPath[]): this {
		this.data.environmentPaths = environmentPaths;
		return this;
	}

	withLocale(locale: string): this {
		this.data.locale = locale;
		return this;
	}

	build(): UserConfig {
		return {
			...this.data,
			environmentPaths: [...this.data.environmentPaths],
		};
	}
}
