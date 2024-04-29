/**
 * Copyright 2024 Undistro Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

export function getRunValues() {
  const currentMode = getCurrentMode();
  const exprEditor = new AceEditor(currentMode);
  setEditorTheme(exprEditor);
  let values = {
    [currentMode]: exprEditor.getValue(),
  };

  document.querySelectorAll(".editor__input.data__input").forEach((editor) => {
    const containerId = editor.id;
    const dataEditor = new AceEditor(containerId);
    setEditorTheme(dataEditor);
    values = {
      ...values,
      [containerId]: dataEditor.getValue(),
    };
  });

  return values;
}
