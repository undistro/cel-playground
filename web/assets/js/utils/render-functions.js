const examplesList = document.getElementById("examples");
const selectInstance = NiceSelect.bind(examplesList);

export function renderExamplesInSelectInstance(
  examples,
  celEditor,
  dataEditor
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
    setCost("");
    output.value = "";
  });
}

export function setCost(cost) {
  const costElem = document.getElementById("cost");
  costElem.innerText = cost || "-";
}
