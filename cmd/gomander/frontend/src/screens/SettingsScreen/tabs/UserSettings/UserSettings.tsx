import { Route, WandSparkles } from "lucide-react";
import { useTranslation } from "react-i18next";

import { type Theme, useTheme } from "@/contexts/theme.tsx";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/design-system/components/ui/card.tsx";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { Input } from "@/design-system/components/ui/input.tsx";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/design-system/components/ui/select";
import { useSettingsContext } from "@/screens/SettingsScreen/context/settingsContext.tsx";
import { EnvironmentPathsField } from "@/screens/SettingsScreen/tabs/UserSettings/components/EnvironmentPathsField.tsx";
import { EnvironmentPathsInfoDialog } from "@/screens/SettingsScreen/tabs/UserSettings/components/EnvironmentPathsInfoDialog.tsx";

export const UserSettings = () => {
	const { t } = useTranslation();
	const { userSettingsForm, supportedLanguages } = useSettingsContext();
	const { rawTheme, setRawTheme } = useTheme();

	return (
		<Form {...userSettingsForm}>
			<form className="w-full h-full flex flex-col justify-between">
				<div className="flex flex-col gap-2">
					<Card>
						<CardHeader>
							<CardTitle className="flex items-center">
								<Route size={20} />
								<span className="ml-2 mr-2">
									{t("userSettingsForm.envPathsTitle")}
								</span>
								<EnvironmentPathsInfoDialog />
							</CardTitle>
							<CardDescription>
								{t("userSettingsForm.envPathsDescription")}
							</CardDescription>
						</CardHeader>
						<CardContent>
							<EnvironmentPathsField />
						</CardContent>
					</Card>
					<Card>
						<CardHeader>
							<CardTitle className="flex items-center space-x-2">
								<WandSparkles size={20} />
								<span>{t("userSettingsForm.preferencesTitle")}</span>
							</CardTitle>
							<CardDescription>
								{t("userSettingsForm.preferencesDescription")}
							</CardDescription>
						</CardHeader>
						<CardContent className="space-y-3">
							<FormField
								control={userSettingsForm.control}
								name="locale"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("userSettingsForm.languageLabel")}</FormLabel>
										<FormControl>
											<Select
												onValueChange={field.onChange}
												defaultValue={field.value}
											>
												<FormControl>
													<SelectTrigger>
														<SelectValue
															placeholder={t(
																"userSettingsForm.languagePlaceholder",
															)}
														/>
													</SelectTrigger>
												</FormControl>
												<SelectContent>
													{supportedLanguages.map((language) => (
														<SelectItem
															key={language.value}
															value={language.value}
														>
															{language.label}
														</SelectItem>
													))}
												</SelectContent>
											</Select>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
							<FormItem>
								<FormLabel>{t("userSettingsForm.themeLabel")}</FormLabel>
								<Select
									onValueChange={(value) => {
										setRawTheme(value as Theme);
									}}
									value={rawTheme}
								>
									<SelectTrigger>
										<SelectValue
											placeholder={t("userSettingsForm.themePlaceholder")}
										/>
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="system">
											{t("userSettingsForm.themeSystem")}
										</SelectItem>
										<SelectItem value="light">
											{t("userSettingsForm.themeLight")}
										</SelectItem>
										<SelectItem value="dark">
											{t("userSettingsForm.themeDark")}
										</SelectItem>
									</SelectContent>
								</Select>
								<FormDescription className="text-xs">
									{t("userSettingsForm.themeDescription")}
								</FormDescription>
							</FormItem>
							<FormField
								control={userSettingsForm.control}
								name="logLineLimit"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("userSettingsForm.logLimitLabel")}</FormLabel>
										<FormControl>
											<Input
												type="number"
												min={1}
												max={5000}
												{...field}
												onChange={(e) => field.onChange(Number(e.target.value))}
											/>
										</FormControl>
										<FormDescription className="text-xs">
											{t("userSettingsForm.logLimitDescription")}
										</FormDescription>
										<FormMessage />
									</FormItem>
								)}
							/>
						</CardContent>
					</Card>
				</div>
			</form>
		</Form>
	);
};
