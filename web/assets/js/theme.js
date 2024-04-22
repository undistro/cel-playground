import {
  ACE_EDITOR,
  localStorageModeKey,
  localStorageThemeKey,
} from "./constants.js";
import { AceEditor } from "./editor.js";
import { getCurrentTheme } from "./utils/localStorage.js";
const toggleBtn = document.getElementsByClassName("toggle-theme")[0];

const { theme: aceEditorTheme } = ACE_EDITOR;

window.addEventListener("load", () => {
  let theme = localStorage.getItem(localStorageThemeKey);
  if (theme === "dark") {
    toggleTheme("dark");
  }
});

toggleBtn.addEventListener("click", function () {
  let currTheme = localStorage.getItem(localStorageThemeKey);
  if (currTheme === "dark") toggleTheme("light");
  else toggleTheme("dark");
});

export function applyThemeToEditors() {
  const theme = localStorage.getItem(localStorageThemeKey);
  const exprEditor = new AceEditor(
    localStorage.getItem(localStorageModeKey) ?? "cel"
  );
  const editorsInputEl = document.querySelectorAll(
    ".editor__input.data__input"
  );
  setEditorTheme(exprEditor);
  editorsInputEl.forEach((editor) => {
    const inputEditor = new AceEditor(editor.id);
    setEditorTheme(inputEditor);
  });
}

function toggleTheme(theme) {
  let toggleIcon = document.getElementsByClassName("toggle-theme__icon")[0];
  let celLogo = document.getElementsByClassName("cel-logo")[0];
  let copyIcon = document.querySelectorAll(".editor-copy-icon");

  if (theme === "dark") {
    document.body.classList.add("dark");
    toggleIcon.src = "./assets/img/moon.svg";
    celLogo.src = "./assets/img/logo-dark.svg";
    copyIcon[0].src = "./assets/img/copy-dark.svg";
    copyIcon[1].src = "./assets/img/copy-dark.svg";
  } else {
    document.body.classList.remove("dark");
    toggleIcon.src = "./assets/img/sun.svg";
    celLogo.src = "./assets/img/logo.svg";
    copyIcon[0].src = "./assets/img/copy.svg";
    copyIcon[1].src = "./assets/img/copy.svg";
  }
  localStorage.setItem(localStorageThemeKey, theme);
  applyThemeToEditors();
}

function getEditorByTheme(currentTheme) {
  return aceEditorTheme[currentTheme];
}

export function setEditorTheme({ editor }) {
  const theme = getCurrentTheme();
  const editorTheme = getEditorByTheme(theme);
  editor.setTheme(editorTheme);
}
