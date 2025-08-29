export const EXPECTED_VALIDATION_ERROR = "EXPECTED_VALIDATION_ERROR";

export const parseError = (error: unknown, messagePrefix?: string) => {
  if (error instanceof Error) {
    if (error.cause === EXPECTED_VALIDATION_ERROR) {
      return error.message;
    }
    return messagePrefix ? `${messagePrefix}: ${error.message}` : error.message;
  }

  if (typeof error === "string") {
    return messagePrefix ? `${messagePrefix}: ${error}` : error;
  }

  return "Unknown error";
};
