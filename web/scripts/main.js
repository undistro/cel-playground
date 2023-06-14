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
import {
  EDITOR_ELEMENTS,
  EXAMPLES,
  SAMPLE_DATA,
  WASM_URL,
} from "./constants.js";

const selectInstance = NiceSelect.bind(document.getElementById("examples"));

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

(async function loadAndRunGoWasm() {
  const go = new Go();

  const buffer = pako.ungzip(await (await fetch(WASM_URL)).arrayBuffer());

  // A fetched response might be decompressed twice on Firefox.
  // See https://bugzilla.mozilla.org/show_bug.cgi?id=610679
  if (buffer[0] === 0x1f && buffer[1] === 0x8b) {
    buffer = pako.ungzip(buffer);
  }

  WebAssembly.instantiate(buffer, go.importObject)
    .then((result) => {
      go.run(result.instance);
      document.getElementById("run").disabled = false;
    })
    .catch((err) => {
      console.error(err);
    });
})();

const runButton = document.getElementById("run");
runButton.addEventListener("click", run);
document.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.code === "Enter") {
    run();
  }
});

fetch("../assets/data.json")
  .then((response) => response.json())
  .then(({ examples, versions }) => {

    // Dynamically set the CEL Go version
    document.getElementById("version").innerText = versions["cel-go"];

    // Load the examples into the select element
    const examplesList = document.getElementById("examples");
    examples.forEach((example) => {
      const option = document.createElement("option");
      option.value = example.name;
      option.innerText = example.name;
      examplesList.appendChild(option);
    });

    selectInstance.update();

    examplesList.addEventListener("change", (event) => {
      const example = examples.find(
        (example) => example.name === event.target.value
      );
      celEditor.setValue(example.cel, -1);
      dataEditor.setValue(example.data, -1);
    });
  })
  .catch((err) => {
    console.error(err);
  });