import { AceEditor } from "../editor.js";

const celEditor = new AceEditor("cel-input");
const dataEditor = new AceEditor("data-input");
const examplesList = document.getElementById("examples");
const selectInstance = NiceSelect.bind(examplesList);

export function renderExamplesInSelectInstance(mode, callbackFn) {
  examplesList.innerHTML = `<option data-display="Examples" value="" disabled selected hidden>
      Examples
    </option>`;

  const urlParams = new URLSearchParams(window.location.search);
  const examples = mode.examples;
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
        celEditor.setValue(example.value[0]?.id ?? "", -1);
        dataEditor.setValue(example.value[0]?.data ?? "", -1);
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

    if (example) {
      celEditor.setExpressionSyntax(mode.syntax);
      localStorage.setItem("example-selected", example.id);
      const tabButtonActived = document.querySelector("#tab .vap__tabs-button");
      fetchTabData(mode, example.id, tabButtonActived);
    }
    callbackFn();
    setCost("");
    output.value = "";
  });
}

export function setCost(cost) {
  const costElem = document.getElementById("cost");
  costElem.innerText = cost || "-";
}

export function renderTabs(mode) {
  const { inputs } = mode;
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
      const savedExample = localStorage.getItem("example-selected");
      fetchTabData(mode, savedExample, tabButton);
    };
    if (idx === 0) addActiveClass(tabButton);
    dataEditor.setExpressionSyntax(input.syntax);

    divParent.appendChild(tabButton);
  });

  holderElement.appendChild(divParent);
  if (inputs.length === 1) holderElement.innerHTML = "";

  function removeActiveClass(element) {
    element.classList.remove("active");
  }

  function addActiveClass(element) {
    element.classList.add("active");
  }
}

function fetchTabData(mode, exampleID, tabButton) {
  fetch(`../../assets/examples/${mode.id}/${exampleID}.json`)
    .then((response) => response.json())
    .then(({ code, inputs }) => {
      celEditor.setValue(code, -1);
      if (tabButton) {
        dataEditor.setValue(inputs[tabButton.id], -1);
      }
    });
}
