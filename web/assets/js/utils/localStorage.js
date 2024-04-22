import { localStorageModeKey, localStorageThemeKey } from "../constants.js";

export function getCurrentTheme() {
  return localStorage.getItem(localStorageThemeKey) ?? "light";
}

export function getCurrentMode() {
  return localStorage.getItem(localStorageModeKey) ?? "cel";
}
