import {
  renderExamplesInSelectInstance,
  renderTabs,
} from "../utils/render-functions.js";

const toggleModeButton = document.getElementById("toggle-mode");

const playgroundModesModalEl = document.getElementById(
  "playground-modes__modal"
);
const closePlaygroundModesModalButton = document.querySelector(
  ".playground-modes__modal-close-btn"
);
toggleModeButton.addEventListener("click", () => {
  playgroundModesModalEl.style.display = "block";
});

renderModeOptions();

closePlaygroundModesModalButton.addEventListener("click", closeModal);

window.onclick = function (event) {
  if (event.target === playgroundModesModalEl) {
    closeModal();
  }
};

function handleModeClick(event, mode, element) {
  const { value } = event.target;

  document
    .querySelectorAll(".playground-modes__options--option")
    .forEach((option) => option.classList.remove("active"));

  element.classList.add("active");
  renderUIChangesByMode(mode);
  localStorage.setItem(localStorageKey, value);
  setTimeout(() => closeModal(), 1000);
}

function renderModeOptions() {
  const el = document.querySelector(".playground-modes__options");

  getModes()
    .then((modes) => {
      modes.forEach((mode, i) => {
        const divOption = createParentElement(mode);
        const label = createLabelElement(mode);
        const input = createInputElement(mode);
        input.onclick = (e) => handleModeClick(e, mode, divOption);

        const modeSaved = localStorage.getItem(localStorageKey);

        if (!modeSaved && i === 0) {
          divOption.classList.add("active");
          renderUIChangesByMode(modes.find((mode) => mode.id === "cel"));
        }
        if (modeSaved === mode.id) {
          divOption.classList.add("active");
          renderUIChangesByMode(mode);
        }

        divOption.appendChild(label);
        divOption.appendChild(input);
        el.appendChild(divOption);
      });
    })
    .catch((err) => console.log(err));
}

function createParentElement(mode) {
  const divOption = document.createElement("div");
  divOption.className = "playground-modes__options--option";
  divOption.setAttribute("data-mode", mode.id);
  return divOption;
}

function createLabelElement(mode) {
  const label = document.createElement("label");
  label.htmlFor = mode.id;
  label.innerHTML = mode.name;
  return label;
}

function createInputElement(mode) {
  const input = document.createElement("input");
  input.className = "playground-modes__options--option-input";
  input.type = "radio";
  input.name = mode.id;
  input.id = mode.id;
  input.value = mode.id;
  return input;
}

async function getModes() {
  const response = await fetch("../../assets/modes.json");
  const modes = await response.json();
  return modes;
}

function closeModal() {
  playgroundModesModalEl.style.display = "none";
}

function renderUIChangesByMode(mode) {
  const titleEl = document.querySelector(".title.expression__square");
  const toggleModeHolder = document.querySelector(".modes__container-holder");
  const titleInputSquareEl = document.querySelector(".title.input__square");

  titleEl.innerHTML = mode.name;
  toggleModeHolder.innerHTML = mode.name;
  titleInputSquareEl.innerHTML = mode.inputs.length > 1 ? "Inputs: " : "Input";

  renderExamplesInSelectInstance(mode, callbackFns);
  callbackFns();

  function callbackFns() {
    renderTabs(mode);
  }
}
