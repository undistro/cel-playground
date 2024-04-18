/**
 * Copyright 2023 Undistro Authors
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

import { setCost } from "./utils/render-functions.js";
import { AceEditor } from "./editor.js";

// Add the following polyfill for Microsoft Edge 17/18 support:
// <script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script>
// (see https://caniuse.com/#feat=textencoder)
if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

const output = document.getElementById("output");

function getRunValues() {
  const currentMode = localStorage.getItem(localStorageModeKey) ?? "cel";

  const exprEditor = new AceEditor(currentMode);
  let values = {
    expression: exprEditor.getValue(),
  };

  document.querySelectorAll(".editor__input.data__input").forEach((editor) => {
    const containerId = editor.id;
    const dataEditor = new AceEditor(containerId);
    values = {
      ...values,
      [containerId]: dataEditor.getValue(),
    };
  });

  return values;
}

function run() {
  const cost = document.getElementById("cost");
  const values = getRunValues();
  output.value = "Evaluating...";
  setCost("");
  console.log({ values: getRunValues() });
  const result = eval(values);

  const { output: resultOutput, isError } = result;

  if (isError) {
    output.value = resultOutput;
    output.style.color = "red";
  } else {
    output.value = JSON.stringify(result);
    output.style.color = "white";
    setCost(cost);
  }
}

function share() {
  const values = getRunValues();

  const str = JSON.stringify({
    ...values,
    mode: localStorage.getItem(localStorageModeKey) ?? "cel",
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

var urlParams = new URLSearchParams(window.location.search);
if (urlParams.has("content")) {
  const content = urlParams.get("content");
  try {
    const decodedUint8Array = new Uint8Array(
      atob(content)
        .split("")
        .map(function (char) {
          return char.charCodeAt(0);
        })
    );

    const decompressedData = pako.ungzip(decodedUint8Array, { to: "string" });
    if (!decompressedData) {
      throw new Error("Invalid content parameter");
    }
    const obj = JSON.parse(decompressedData);
    localStorage.setItem(localStorageModeKey, obj.mode);
    new AceEditor(obj.mode).setValue(obj.expression, -1);
    document
      .querySelectorAll(".editor__input.data__input")
      .forEach((editor) => {
        const containerId = editor.id;
        new AceEditor(containerId).setValue(obj[containerId], -1);
      });
  } catch (error) {
    console.error(error);
  }
}

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
  let value = celEditor.editor.getValue();
  navigator.clipboard.writeText(value).catch(console.error);
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
  let value = dataEditor.editor.getValue();
  navigator.clipboard.writeText(value);
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

(async function loadAndRunGoWasm() {
  const go = new Go();

  const buffer = pako.ungzip(
    await (await fetch("assets/main.wasm.gz")).arrayBuffer()
  );

  // A fetched response might be decompressed twice on Firefox.
  // See https://bugzilla.mozilla.org/show_bug.cgi?id=610679
  if (buffer[0] === 0x1f && buffer[1] === 0x8b) {
    buffer = pako.ungzip(buffer);
  }

  WebAssembly.instantiate(buffer, go.importObject)
    .then((result) => {
      go.run(result.instance);
      document.getElementById("run").disabled = false;
      document.getElementById("output").placeholder =
        "Press 'Run' to evaluate your CEL expression.";
    })
    .catch((err) => {
      console.error(err);
    });
})();

const runButton = document.getElementById("run");
const shareButton = document.getElementById("share");
const copyButton = document.getElementById("copy");

runButton.addEventListener("click", run);
shareButton.addEventListener("click", share);
copyButton.addEventListener("click", copy);
document.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.code === "Enter") {
    run();
  }
});

fetch("../assets/examples/cel.json")
  .then((response) => response.json())
  .then(({ versions }) => {
    document.getElementById("version").innerText = versions["cel-go"];
  })
  .catch((err) => {
    console.error(err);
  });
