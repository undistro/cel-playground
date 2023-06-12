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

import { AceEditor } from "./editor.js";
import { EDITOR_ELEMENTS, EXAMPLES, SAMPLE_DATA, WASM_URL } from "./constants.js";

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

const [celEditor, dataEditor] = EDITOR_ELEMENTS.map((id) => new AceEditor(id));

dataEditor.setValue(SAMPLE_DATA, -1);

function run() {
  const data = dataEditor.getValue();
  const expression = celEditor.getValue();
  const output = document.getElementById("output");

  output.value = "Evaluating...";
  const result = eval(expression, data);

  const { output: resultOutput, isError } = result;
  output.value = `${resultOutput}`;
  output.style.color = isError ? "red" : "white";
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject)
  .then((result) => {
    go.run(result.instance);
    document.getElementById("run").disabled = false;
  })
  .catch((err) => {
    console.error(err);
  });

const runButton = document.getElementById("run");
runButton.addEventListener("click", run);
document.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.code === "Enter") {
    run();
  }
});

const list = document.createElement("div");

EXAMPLES.forEach((example) => {
  const button = document.createElement("button");
  button.classList.add("example-item");
  button.innerText = example.name;
  list.appendChild(button);
});

tippy("#examples", {
  content: list.innerHTML,
  triggerTarget: document.getElementById("examples"),
  placement: "bottom-end",
  trigger: "click",
  animation: "shift-away",
  interactive: true,
  interactiveBorder: 10,
  interactiveDebounce: 0,
  aria: null,
});
