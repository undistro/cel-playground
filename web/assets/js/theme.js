import { AceEditor } from "./editor.js";
const toggleBtn = document.getElementsByClassName("toggle-theme")[0];

window.addEventListener("load", () => {
  let theme = localStorage.getItem("theme");
  if (theme === "dark") {
    toggleTheme("dark");
  }
});

toggleBtn.addEventListener("click", function () {
  let currTheme = localStorage.getItem("theme");
  if (currTheme === "dark") toggleTheme("light");
  else toggleTheme("dark");
});

export function applyThemeToEditors() {
  const theme = localStorage.getItem("theme");
  const exprEditor = new AceEditor(
    localStorage.getItem(localStorageModeKey) ?? "cel"
  );
  const editorsInputEl = document.querySelectorAll(
    ".editor__input.data__input"
  );

  if (theme === "dark") {
    exprEditor.editor.setTheme("ace/theme/tomorrow_night");
    editorsInputEl.forEach((editor) => {
      new AceEditor(editor.id).editor.setTheme("ace/theme/tomorrow_night");
    });
  } else {
    exprEditor.editor.setTheme("ace/theme/clouds");
    editorsInputEl.forEach((editor) => {
      new AceEditor(editor.id);
    });
  }
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
  localStorage.setItem("theme", theme);
  applyThemeToEditors();
}
