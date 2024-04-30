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

import { localStorageModeKey } from "./constants.js";
import { AceEditor } from "./editor.js";
import { setEditorTheme } from "./theme.js";
import {
  getExprEditorValue,
  getInputEditorValue,
  getRunValues,
} from "./utils/editor.js";
import { getCurrentMode, getCurrentTheme } from "./utils/localStorage.js";

const shareButton = document.getElementById("share");
shareButton.addEventListener("click", share);

function share() {
  const values = getRunValues();

  const modeId = getCurrentMode();

  const str = JSON.stringify({
    ...values,
    mode: modeId,
  });
  var compressed_uint8array = pako.gzip(str);
  var b64encoded_string = btoa(
    String.fromCharCode.apply(null, compressed_uint8array)
  );

  const url = new URL(window.location.href);
  url.searchParams.set("content", b64encoded_string);
  window.history.pushState({}, "", url.toString());

  document.querySelector(".share-url__container").style.display = "flex";
  document.querySelector(".share-url__input").value = url.toString();
}

export const renderSharedContent = (
  mode,
  object,
  legacyObjectShared = false
) => {
  localStorage.setItem(localStorageModeKey, mode.id);

  try {
    if (legacyObjectShared) {
      const celEditor = new AceEditor("cel");
      celEditor.setValue(object.expression, -1);
      setEditorTheme(celEditor);

      const inputEditor = new AceEditor("dataInput");
      inputEditor.setValue(object.data, -1);
      setEditorTheme(inputEditor);
    } else {
      const exprEditor = new AceEditor(object.mode);
      exprEditor.setValue(object[object.mode], -1);
      setEditorTheme(exprEditor);
      document
        .querySelectorAll(".editor__input.data__input")
        ?.forEach((editor) => {
          const containerId = editor.id;
          const inputEditor = new AceEditor(containerId);
          inputEditor.setValue(object[containerId], -1);
          setEditorTheme(inputEditor);
        });
    }
  } catch (error) {
    console.error(error);
  }
};

let celCopyIcon = document.getElementById("cel-copy-icon");
let celCopyHover = document.getElementById("cel-copy-hover");
let celCopyClick = document.getElementById("cel-copy-click");
let celInput = document.getElementById("cel-cont");

celInput.addEventListener("mouseover", () => {
  celCopyIcon.style.display = "inline";
});
celInput.addEventListener("mouseleave", () => {
  celCopyIcon.style.display = "none";
});

celCopyIcon.addEventListener("click", () => {
  const exprEditorValue = getExprEditorValue();
  navigator.clipboard.writeText(exprEditorValue).catch(console.error);
  celCopyHover.style.display = "none";
  celCopyClick.style.display = "flex";
  setTimeout(() => {
    celCopyClick.style.display = "none";
  }, 1000);
});

celCopyIcon.addEventListener("mouseover", () => {
  celCopyHover.style.display = "flex";
});

celCopyIcon.addEventListener("mouseleave", () => {
  celCopyHover.style.display = "none";
});

let dataCopyIcon = document.getElementById("data-copy-icon");
let dataCopyHover = document.getElementById("data-copy-hover");
let dataCopyClick = document.getElementById("data-copy-click");
let dataInput = document.getElementById("data-cont");

dataInput.addEventListener("mouseover", () => {
  dataCopyIcon.style.display = "inline";
});

dataInput.addEventListener("mouseleave", () => {
  dataCopyIcon.style.display = "none";
});

dataCopyIcon.addEventListener("click", () => {
  const dataInputValue = getInputEditorValue();
  navigator.clipboard.writeText(dataInputValue);
  dataCopyHover.style.display = "none";
  dataCopyClick.style.display = "flex";
  setTimeout(() => {
    dataCopyClick.style.display = "none";
  }, 1000);
});

dataCopyIcon.addEventListener("mouseover", () => {
  dataCopyHover.style.display = "flex";
});

dataCopyIcon.addEventListener("mouseleave", () => {
  dataCopyHover.style.display = "none";
});

const copyButton = document.getElementById("copy");
copyButton.addEventListener("click", copy);

function copy() {
  const copyText = document.querySelector(".share-url__input");
  copyText.select();
  copyText.setSelectionRange(0, 99999);
  navigator.clipboard.writeText(copyText.value);
  window.getSelection().removeAllRanges();

  const tooltip = document.querySelector(".share-url__tooltip");
  tooltip.style.opacity = 1;
  setTimeout(() => {
    tooltip.style.opacity = 0;
  }, 3000);
}
