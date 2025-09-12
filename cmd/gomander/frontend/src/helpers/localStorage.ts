export const removeKeyFromLocalStorage = (key: string) => {
  try {
    localStorage.removeItem(key);
  } catch (e) {
    // Ignore write errors
  }
};
