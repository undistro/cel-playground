const examplesList = document.getElementById("examples");
const selectInstance = NiceSelect.bind(examplesList);

export function renderExamplesInSelectInstance(
  examples,
  celEditor,
  dataEditor,
  callbackFn
) {
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
      celEditor.setValue(example.id, -1);
      dataEditor.setValue(example.data, -1);
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

export function renderTabs(inputs) {
  const holderElement = document.getElementById("tab");
  holderElement.innerHTML = "";

  const divParent = document.createElement("div");
  divParent.className = "vap__tabs";
  divParent.id = "vap__tabs";

  inputs.forEach((input, idx) => {
    const buttonTab = document.createElement("button");
    buttonTab.innerHTML = input.name;
    buttonTab.className = "vap__tabs-button";
    buttonTab.onclick = () => {
      const allButtons = divParent?.querySelectorAll(".vap__tabs-button");
      allButtons.forEach(removeActiveClass);
      if (buttonTab.classList.contains("active")) removeActiveClass(buttonTab);
      else addActiveClass(buttonTab);
    };
    if (idx === 0) addActiveClass(buttonTab);
    divParent.appendChild(buttonTab);
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
