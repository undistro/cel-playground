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
import { ExampleService } from "../services/examples.js";

const celEditor = new AceEditor("cel-input");
const dataEditor = new AceEditor("data-input");
const examplesList = document.getElementById("examples");
const selectInstance = NiceSelect.bind(examplesList);

export function renderExamplesInSelectInstance(examples, callbackFn) {
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
    if (i === 0) {
      const [firstValue] = example.value;
      setEditors(firstValue.data, firstValue.inputs[0].data);
      renderTabs(firstValue);
    }

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
      if (!urlParams.has("content"))
        setEditors(example.value[0].data, example.value[0].inputs[0].data);
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

    if (example) {
      setEditors(example.data, example.inputs[0].data);
      callbackFn(example);
    }
    setCost("");
    output.value = "";
  });
}

export function setCost(cost) {
  const costElem = document.getElementById("cost");
  costElem.innerText = cost || "-";
}

export function renderTabs(example) {
  const { inputs } = example;
  const holderElement = document.getElementById("tab");
  holderElement.innerHTML = "";

  const divParent = document.createElement("div");
  divParent.className = "vap__tabs";
  divParent.id = "vap__tabs";

  inputs.forEach((input, idx) => {
    const tabButton = document.createElement("button");
    tabButton.innerHTML = input.name;
    tabButton.className = "vap__tabs-button";
    tabButton.id = input.id;
    tabButton.onclick = () => {
      const allButtons = divParent?.querySelectorAll(".vap__tabs-button");
      allButtons.forEach(removeActiveClass);
      if (tabButton.classList.contains("active")) removeActiveClass(tabButton);
      else addActiveClass(tabButton);
      divParent.setAttribute("style", `--current-tab: ${idx}`);
      setEditors(example.data, input.data);
    };
    if (idx === 0) addActiveClass(tabButton);

    divParent.appendChild(tabButton);
  });

  holderElement.appendChild(divParent);
  if (inputs.length <= 1) holderElement.innerHTML = "";

  function removeActiveClass(element) {
    element.classList.remove("active");
  }

  function addActiveClass(element) {
    element.classList.add("active");
  }
}

function setEditors(expressionEditorValue, inputEditorValue) {
  celEditor.setValue(expressionEditorValue, -1);
  dataEditor.setValue(inputEditorValue, -1);
}
