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

const examplesList = document.getElementById("examples");
const selectInstance = NiceSelect.bind(examplesList);

export function renderExamplesInSelectInstance(mode, examples) {
  examplesList.innerHTML = `<option data-display="Examples" value="" disabled selected hidden>
      Examples
    </option>`;

  const urlParams = new URLSearchParams(window.location.search);
  const groupByCategory = examples.reduce((acc, example) => {
    return {
      ...acc,
      [example.category]: [...(acc[example.category] ?? []), example],
    };
  }, {});

  const examplesByCategory = Object.entries(groupByCategory).map(
    ([key, value]) => ({ label: key, value })
  );

  examplesByCategory.forEach((example, i) => {
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
      }
      // setEditors(example.value[0].data, example.value[0].inputs[0].data);
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
    if (event.target.value === "Blank") return;
    if (!example) return;

    handleFillExpressionContent(mode, example);
    handleFillTabContent(mode, example);
    setCost("");
    output.value = "";
  });
}

export function setCost(cost) {
  const costElem = document.getElementById("cost");
  costElem.innerText = cost || "-";
}

export function handleFillExpressionContent(mode, example) {
  const exprEditor = new AceEditor(mode.id, mode.mode);
  exprEditor.setValue(example[mode.id], -1);
}

export function handleFillTabContent(mode, example) {
  document.querySelectorAll(".editor__input.data__input")[0].style.display =
    "block";
  mode.tabs.forEach((tab) => {
    const containerId = tab.id;
    const inputEditor = new AceEditor(containerId, tab.mode);
    inputEditor.setValue(example[containerId], -1);
  });
}

export function renderExpressionContent(mode, examples) {
  const exprInput = document.querySelector(".editor__input.expr__input");
  exprInput.id = mode.id;

  const currentExample = getCurrentExample(mode, examples);

  const exprEditor = new AceEditor(mode.id, mode.mode);
  exprEditor.setValue(currentExample?.[mode.id] ?? mode[mode.id], -1);
}

export function renderTabs(mode, examples) {
  const { tabs } = mode;

  const holderElement = document.getElementById("tab");
  holderElement.innerHTML = "";
  const divParent = document.createElement("div");
  divParent.className = "tabs";
  divParent.id = "tabs";

  document.querySelectorAll(".editor__input.data__input")?.forEach((editor) => {
    editor.remove();
  });

  tabs.forEach((tab, idx) => {
    const currentExample = getCurrentExample(mode, examples);

    if (!currentExample) return;

    const containerId = tab.id;
    const editorContainer = createEditorContainer(containerId);
    const inputEditor = new AceEditor(containerId, tab.mode);
    inputEditor.setValue(currentExample[containerId], -1);

    const tabButton = document.createElement("button");
    tabButton.innerHTML = tab.name;
    tabButton.className = "tabs-button";

    tabButton.onclick = () => {
      document
        .querySelectorAll(".editor__input.data__input")
        ?.forEach((editor) => {
          editor.style.display = "none";
        });
      editorContainer.style.display = "block";
      const allButtons = divParent?.querySelectorAll(".tabs-button");
      allButtons.forEach(removeActiveClass);
      if (tabButton.classList.contains("active")) removeActiveClass(tabButton);
      else addActiveClass(tabButton);
      divParent.setAttribute("style", `--current-tab: ${idx}`);
    };
    if (idx === 0) addActiveClass(tabButton);
    divParent.appendChild(tabButton);
  });

  holderElement.appendChild(divParent);
  handleFillTabContent(mode, examples[0]);

  if (tabs.length <= 1) holderElement.style.visibility = "hidden";

  function removeActiveClass(element) {
    element.classList.remove("active");
  }

  function addActiveClass(element) {
    element.classList.add("active");
  }
}

function createEditorContainer(containerId) {
  const holderElement = document.getElementById("data-cont");

  const parent = document.createElement("div");
  parent.className = "editor__input data__input";
  parent.id = containerId;
  parent.style.display = "none";
  holderElement.appendChild(parent);
  return parent;
}

function getCurrentExample(mode, examples) {
  const currentExample = examples.find((example) =>
    Object.keys(example).includes(mode.id)
  );

  return currentExample;
}
