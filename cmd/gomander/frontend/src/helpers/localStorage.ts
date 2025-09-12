export const removeKeyFromLocalStorage = (key: string) => {
  try {
    localStorage.removeItem(key);
  } catch (e) {
    console.log(e);
    // Ignore write errors
  }
};
