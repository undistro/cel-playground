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

const celEditor = new AceEditor("cel-input");
const dataEditor = new AceEditor("data-input");

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

window.addEventListener("load", () => {
  let theme = localStorage.getItem("theme");
  if (theme === "dark") {
    toggleMode("dark");
  }
});

const toggleBtn = document.getElementsByClassName("toggle-theme")[0];
toggleBtn.addEventListener("click", function () {
  let currTheme = localStorage.getItem("theme");
  if (currTheme === "dark") toggleMode("light");
  else toggleMode("dark");
});

function toggleMode(theme) {
  let toggleIcon = document.getElementsByClassName("toggle-theme__icon")[0];
  let celLogo = document.getElementsByClassName("cel-logo")[0];
  let copyIcon = document.querySelectorAll(".editor-copy-icon");

  if (theme === "dark") {
    document.body.classList.add("dark");
    toggleIcon.src = "./assets/img/moon.svg";
    celEditor.editor.setTheme("ace/theme/tomorrow_night");
    dataEditor.editor.setTheme("ace/theme/tomorrow_night");
    celLogo.src = "./assets/img/logo-dark.svg";
    copyIcon[0].src = "./assets/img/copy-dark.svg";
    copyIcon[1].src = "./assets/img/copy-dark.svg";
    localStorage.setItem("theme", "dark");
  } else {
    document.body.classList.remove("dark");
    toggleIcon.src = "./assets/img/sun.svg";
    celEditor.editor.setTheme("ace/theme/clouds");
    dataEditor.editor.setTheme("ace/theme/clouds");
    celLogo.src = "./assets/img/logo.svg";
    copyIcon[0].src = "./assets/img/copy.svg";
    copyIcon[1].src = "./assets/img/copy.svg";
    localStorage.setItem("theme", "light");
  }
}

function share() {
  const data = dataEditor.getValue();
  const expression = celEditor.getValue();

  const obj = {
    data: data,
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
    dataEditor.setValue(obj.data, -1);
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
})
celInput.addEventListener("mouseleave", () => {
  celCopyIcon.style.display = "none";
})

celCopyIcon.addEventListener("click", () => {
  let value = celEditor.editor.getValue();
  navigator.clipboard.writeText(value).catch(console.error);
  celCopyHover.style.display = "none"
  celCopyClick.style.display = "flex";
  setTimeout(() => {
    celCopyClick.style.display = "none";
  }, 1000)
})

celCopyIcon.addEventListener("mouseover", () => {
  celCopyHover.style.display = "flex"
})

celCopyIcon.addEventListener("mouseleave", () => {
  celCopyHover.style.display = "none"
})

let dataCopyIcon = document.getElementById("data-copy-icon");
let dataCopyHover = document.getElementById("data-copy-hover");
let dataCopyClick = document.getElementById("data-copy-click");
let dataInput = document.getElementById("data-cont");

dataInput.addEventListener("mouseover", () => {
  dataCopyIcon.style.display = "inline"
})

dataInput.addEventListener("mouseleave", () => {
  dataCopyIcon.style.display = "none"
})

dataCopyIcon.addEventListener("click", () => {
  let value = dataEditor.editor.getValue();
  navigator.clipboard.writeText(value);
  dataCopyHover.style.display = "none";
  dataCopyClick.style.display = "flex";
  setTimeout(() => {
    dataCopyClick.style.display = "none";
  }, 1000);
})

dataCopyIcon.addEventListener("mouseover", () => {
  dataCopyHover.style.display = "flex";
})

dataCopyIcon.addEventListener("mouseleave", () => {
  dataCopyHover.style.display = "none";
})


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

fetch("../assets/data.json")
  .then((response) => response.json())
  .then(({ examples, versions }) => {
    // Dynamically set the CEL Go version
    document.getElementById("version").innerText = versions["cel-go"];

    // Load the examples into the select element
    const examplesList = document.getElementById("examples");

    const groupByCategory = examples.reduce((acc, example) => {
      return {
        ...acc,
        [example.category]: [...(acc[example.category] ?? []), example],
      };
    }, {});

    const examplesByCategory = Object.entries(groupByCategory).map(
      ([key, value]) => ({ label: key, value })
    );

    examplesByCategory.forEach((example) => {
      const optGroup = document.createElement("optgroup");
      optGroup.label = example.label;

      example.value.forEach((item) => {
        const option = document.createElement("option");
        const itemName = item.name;

        option.value = itemName;
        option.innerText = itemName;
        optGroup.appendChild(option);
      });

      if (example.label === "default") {
        if (!urlParams.has("content")) {
          celEditor.setValue(example.value[0].cel, -1);
          dataEditor.setValue(example.value[0].data, -1);
        }
      } else if (example.label === "Blank") {
        return;
      } else {
        examplesList.appendChild(optGroup);
      }
    });

    const blankOption = document.createElement("option");
    blankOption.innerText = "Blank";
    blankOption.value = "Blank";
    examplesList.appendChild(blankOption);

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
