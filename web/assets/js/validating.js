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

const celEditor = new AceEditor("vap-input");
const dataEditorOriginal = new AceEditor("data-input-original");
const dataEditorUpdated = new AceEditor("data-input-updated");

function run() {
  const dataOriginal = dataEditorOriginal.getValue();
  const dataUpdated = dataEditorUpdated.getValue();
  const expression = celEditor.getValue();
  const output = document.getElementById("output");

  output.value = "Evaluating...";
  const result = vapEval(expression, dataOriginal, dataUpdated);

  const { output: resultOutput, isError } = result;
  output.value = `${resultOutput}`;
  output.style.color = isError ? "red" : "white";
}


function share() {
  const dataOriginal = dataEditorOriginal.getValue();
  const dataUpdated = dataEditorUpdated.getValue();
  const expression = celEditor.getValue();

  const obj = {
    dataOriginal: dataOriginal,
    dataUpdated: dataUpdated,
    expression: expression,
  };

  const str = JSON.stringify(obj);
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
    celEditor.setValue(obj.expression, -1);
    dataEditorOriginal.setValue(obj.dataOriginal, -1);
    dataEditorUpdated.setValue(obj.dataUpdated, -1);
  } catch (error) {
    console.error(error);
  }
}

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

function showTab(evt, tabName) {
  var index;
  var tabcontent = document.getElementsByClassName("tab-content");
  for (index = 0; index < tabcontent.length; index++) {
    tabcontent[index].style.display = "none";
  }

  var tablinks = document.getElementsByClassName("tab-links");
  for (index = 0; index < tablinks.length; index++) {
    tablinks[index].className = tablinks[index].className.replace(" active", "");
  }

  var elem = document.getElementById(tabName);
  elem.style.flex = "auto";
  elem.style.display = "flex";
  evt.currentTarget.className += " active";
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
        "Press 'Run' to evaluate the ValidatingAdmissionPolicy.";
    })
    .catch((err) => {
      console.error(err);
    });
})();

const runButton = document.getElementById("run");
const shareButton = document.getElementById("share");
const copyButton = document.getElementById("copy");
const originalResourceButton = document.getElementById("original-resource-button");
const updatedResourceButton = document.getElementById("updated-resource-button");

runButton.addEventListener("click", run);
shareButton.addEventListener("click", share);
copyButton.addEventListener("click", copy);
originalResourceButton.addEventListener("click", (event) => {showTab(event, "original-resource-tab")})
updatedResourceButton.addEventListener("click", (event) => {showTab(event, "updated-resource-tab")})
document.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.code === "Enter") {
    run();
  }
});

updatedResourceButton.click()

fetch("../assets/validating_data.json")
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

      if (example.name === "default") {
        if (!urlParams.has("content")) {
          celEditor.setValue(example.vap, -1);
          dataEditorOriginal.setValue(example.dataOriginal, -1);
          dataEditorUpdated.setValue(example.dataUpdated, -1);
        }
      } else {
        examplesList.appendChild(option);
      }
    });

    selectInstance.update();

    examplesList.addEventListener("change", (event) => {
      const example = examples.find(
        (example) => example.name === event.target.value
      );
      celEditor.setValue(example.vap, -1);
      dataEditorOriginal.setValue(example.dataOriginal, -1);
      dataEditorUpdated.setValue(example.dataUpdated, -1);
    });
  })
  .catch((err) => {
    console.error(err);
  });