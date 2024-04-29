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
import {
  handleRenderAccordions,
  hideAccordions,
} from "./components/accordions/result.js";
import { getRunValues } from "./utils/editor.js";
import { getCurrentMode } from "./utils/localStorage.js";

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

function run() {
  const values = getRunValues();
  output.value = "Evaluating...";
  setCost("");

  try {
    const mode = getCurrentMode();
    const result = eval(mode, values);
    const { output: resultOutput, isError } = result;
    if (isError) {
      output.value = resultOutput;
      output.style.color = "red";
      hideAccordions();
    } else {
      const obj = JSON.parse(resultOutput);
      const resultCost = obj?.cost;
      delete obj.cost;

      if ("result" in obj) {
        output.value = JSON.stringify(obj.result);
        output.style.color = "white";
      } else {
        handleRenderAccordions(obj);
      }
      setCost(resultCost);
    }
  } catch (error) {
    output.value = "";
    setCost("");
    console.log(error);
  }
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

runButton.addEventListener("click", run);
document.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.code === "Enter") {
    run();
  }
});

function getVersion() {
  fetch("../assets/examples/cel-input.json")
    .then((response) => response.json())
    .then(({ versions }) => {
      document.getElementById("version").innerText = versions["cel-go"];
    })
    .catch((err) => {
      console.error(err);
    });
}

getVersion();
