export const parseError = (error: unknown, messagePrefix?: string) => {
	if (error instanceof Error) {
		return messagePrefix ? `${messagePrefix}: ${error.message}` : error.message;
	}

	if (typeof error === "string") {
		return messagePrefix ? `${messagePrefix}: ${error}` : error;
	}

	return "Unknown error";
};
