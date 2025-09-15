import { useState } from "react";

export const useLocalStorageState = <T>(key: string, defaultValue: T) => {
  const [state, setState] = useState<T>(() => {
    try {
      const storedValue = window.localStorage.getItem(key);
      return storedValue ? (JSON.parse(storedValue) as T) : defaultValue;
    } catch {
      return defaultValue;
    }
  });

  const setLocalStorageState = (value: T | ((prevState: T) => T)) => {
    setState((prevState) => {
      const newValue = value instanceof Function ? value(prevState) : value;
      try {
        window.localStorage.setItem(key, JSON.stringify(newValue));
      } catch {
        // Ignore write errors
      }
      return newValue;
    });
  };

  return [state, setLocalStorageState] as const;
};
