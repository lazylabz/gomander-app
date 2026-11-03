import { LOGS_FONT_SIZE } from "@/constants/logsFontSize.ts";
import type { EnvironmentPath, UserConfig } from "@/contracts/types.ts";

export class UserConfigBuilder {
	private data: UserConfig = {
		lastOpenedProjectId: "",
		environmentPaths: [],
		locale: "en",
		logsFontSize: LOGS_FONT_SIZE.default,
		zoomLevel: 1,
		notificationsEnabled: true,
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

	withLogsFontSize(logsFontSize: number): this {
		this.data.logsFontSize = logsFontSize;
		return this;
	}

	withZoomLevel(zoomLevel: number): this {
		this.data.zoomLevel = zoomLevel;
		return this;
	}

	withNotificationsEnabled(notificationsEnabled: boolean): this {
		this.data.notificationsEnabled = notificationsEnabled;
		return this;
	}

	build(): UserConfig {
		return {
			...this.data,
			environmentPaths: [...this.data.environmentPaths],
		};
	}
}
