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

const playgroundModeOptions = document.querySelectorAll(
  ".playground-modes__options--option"
);
playgroundModeOptions.forEach((option, _, self) => {
  const modeSaved = localStorage.getItem(localStorageKey);
  const dataModeOption = option.getAttribute("data-mode");
  const currentMode = MODES[dataModeOption];

  if (!modeSaved) {
    self[0].classList.add("active");
    renderUIChanges(MODES["CEL"]);
  }

  if (modeSaved === dataModeOption) {
    option.classList.add("active");
    renderUIChanges(currentMode);
  }

  renderVAPTabs();

  const input = option.querySelector("input[type=radio]");

  input.addEventListener("click", (e) => {
    playgroundModeOptions.forEach((option) => {
      option.classList.remove("active");
    });
    const { value } = e.target;
    if (value === dataModeOption) {
      option.classList.add("active");
      renderUIChanges(currentMode);
    }
    localStorage.setItem(localStorageKey, value);
    renderVAPTabs();
    setTimeout(() => closeModal(), 1000);
  });
});

function renderModeOptions() {
  const el = document.querySelector(".playground-modes__options");
  const playCelModeKeys = Object.keys(MODES);
  playCelModeKeys.forEach((key) => {
    const { value: modeValue, title } = MODES[key];

    const divOption = document.createElement("div");
    divOption.className = "playground-modes__options--option";
    divOption.setAttribute("data-mode", modeValue);

    const label = document.createElement("label");
    label.htmlFor = modeValue;
    label.innerHTML = title;

    const input = document.createElement("input");
    input.className = "playground-modes__options--option-input";
    input.type = "radio";
    input.name = modeValue;
    input.id = modeValue;
    input.value = modeValue;

    divOption.appendChild(label);
    divOption.appendChild(input);

    el.appendChild(divOption);
  });
}

function closeModal() {
  playgroundModesModalEl.style.display = "none";
}

function renderUIChanges(changes) {
  const titleEl = document.querySelector(".title");
  const toggleModeHolder = document.querySelector(".modes__container-holder");

  titleEl.innerHTML = changes.title;
  toggleModeHolder.innerHTML = changes.title;
}