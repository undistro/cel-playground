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

import { hideAccordions } from "../components/accordions/result.js";
import { localStorageModeKey } from "../constants.js";
import { AceEditor } from "../editor.js";
import { setEditorTheme } from "../theme.js";
import { getCurrentMode } from "./localStorage.js";

const examplesList = document.getElementById("examples");
const selectInstance = NiceSelect.bind(examplesList);

const dataEditorInputClassNames = ".editor__input.data__input";

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

  examplesByCategory.forEach((example) => {
    const optGroup = document.createElement("optgroup");
    optGroup.label = example.label;

    example.value.forEach((item) => {
      const itemName = item.name;
      const option = document.createElement("option");
      option.className = "examples__option";
      option.value = itemName;
      option.innerText = itemName;
      optGroup.appendChild(option);
    });

    if (example.label === "Blank" || example.label === "default") {
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
    if (event.target.value === "Blank") {
      const currentMode = getCurrentMode();
      const exprEditor = new AceEditor(currentMode);
      exprEditor.setValue("", -1);
      setEditorTheme(exprEditor);
      document
        .querySelectorAll(dataEditorInputClassNames)
        .forEach((container) => {
          const containerId = container.id;
          const inputEditor = new AceEditor(containerId);
          inputEditor.setValue("");
          setEditorTheme(inputEditor);
        });
      hideAccordions();
      output.value = "";
    } else {
      if (!example) return;
      handleFillExpressionContent(mode, example);
      handleFillTabContent(mode, example);
    }

    setCost("");
    output.value = "";
    hideAccordions();
  });
}

export function setCost(cost) {
  const costElem = document.getElementById("cost");
  costElem.innerText = cost || "-";
}

export function handleFillExpressionContent(mode, example) {
  const exprEditor = new AceEditor(mode.id);
  exprEditor.setValue(example[mode.id], -1);
  exprEditor.editor.setTheme();
  setEditorTheme(exprEditor);
}

export function handleFillTabContent(mode, example) {
  resetTabs();
  mode.tabs.forEach((tab) => {
    const containerId = tab.id;
    const inputEditor = new AceEditor(containerId);
    inputEditor.setValue(example[containerId], -1);
    setEditorTheme(inputEditor);
  });
}

export function renderExpressionContent(mode, examples) {
  const exprInput = document.querySelector(".editor__input.expr__input");
  exprInput.id = mode.id;

  const currentExample = getCurrentExample(mode, examples);
  const exprEditor = new AceEditor(mode.id);
  exprEditor.setSyntax(mode.mode);
  exprEditor.setValue(currentExample?.[mode.id] ?? mode[mode.id], -1);
}

export function renderTabs(mode, examples) {
  const { tabs } = mode;

  const holderElement = document.getElementById("tab");
  holderElement.innerHTML = "";
  const divParent = document.createElement("div");
  divParent.className = "tabs";
  divParent.id = "tabs";
  divParent.setAttribute("data-tab-active", 0);
  divParent.setAttribute("data-tab-length", tabs.length);

  document.querySelectorAll(dataEditorInputClassNames)?.forEach((editor) => {
    editor.remove();
  });

  tabs.forEach((tab, idx) => {
    const currentExample = getCurrentExample(mode, examples);

    if (!currentExample) return;

    const containerId = tab.id;
    const editorContainer = createEditorContainer(containerId);
    const inputEditor = new AceEditor(containerId);
    inputEditor.setSyntax(tab.mode);
    inputEditor.setValue(currentExample[containerId], -1);

    const tabButton = document.createElement("button");
    tabButton.innerHTML = `<span>${tab.name}</span>`;
    tabButton.className = "tabs-button";
    tabButton.title = tab.name;

    tabButton.onclick = () => {
      document
        .querySelectorAll(dataEditorInputClassNames)
        ?.forEach((editor) => {
          editor.style.display = "none";
        });
      editorContainer.style.display = "block";
      const allButtons = divParent?.querySelectorAll(".tabs-button");
      allButtons.forEach(removeActiveClass);
      if (tabButton.classList.contains("active")) removeActiveClass(tabButton);
      else addActiveClass(tabButton);
      divParent.setAttribute("style", `--current-tab: ${idx}`);
      divParent.setAttribute("data-tab-active", idx);
    };
    if (idx === 0) addActiveClass(tabButton);
    divParent.appendChild(tabButton);
  });

  holderElement.appendChild(divParent);
  handleFillTabContent(mode, examples[0]);

  if (tabs.length <= 1) holderElement.style.visibility = "hidden";
  else holderElement.style.visibility = "visible";

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

function resetTabs() {
  document.querySelectorAll(".tabs-button").forEach((tabButton, i) => {
    if (i === 0) {
      tabButton.parentElement.setAttribute("style", `--current-tab: ${i}`);
      tabButton.classList.add("active");
    } else tabButton.classList.remove("active");
  });

  document.querySelectorAll(dataEditorInputClassNames).forEach((editor, i) => {
    if (i === 0) editor.style.display = "block";
    else editor.style.display = "none";
  });
}
