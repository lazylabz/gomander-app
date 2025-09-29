export const removeKeyFromLocalStorage = (key: string) => {
  try {
    localStorage.removeItem(key);
  } catch {
    // Ignore write errors
  }
};
