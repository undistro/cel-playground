import { AceEditor } from "../editor.js";
import { setEditorTheme } from "../theme.js";
import { getCurrentMode } from "./localStorage.js";

export function getExprEditorValue() {
  const exprEditor = new AceEditor(getCurrentMode());
  const exprEditorValue = exprEditor.getValue();
  setEditorTheme(exprEditor);
  return exprEditorValue;
}

export function getInputEditorValue() {
  const editorsInputEl = document.querySelectorAll(
    ".editor__input.data__input"
  );
  const tabsEL = document.getElementById("tabs");
  const currentTabActiveIndex = Number(tabsEL.getAttribute("data-tab-active"));
  const editor = editorsInputEl[currentTabActiveIndex];
  const inputEditor = new AceEditor(editor.id);
  setEditorTheme(inputEditor);
  return inputEditor.getValue();
}
