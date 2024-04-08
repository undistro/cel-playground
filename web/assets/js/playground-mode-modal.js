const MODES = Object.freeze({
  "Common Expression Language (CEL)": "CEL",
  "Validating Admission Policy": "VAP",
  "Web Hooks": "WEB_HOOKS",
  "Authentication Claim Mapping": "AUTH_CMAPPING",
  Authentication: "AUTH",
});

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
  const modeSaved = localStorage.getItem("@cel-playground:mode");
  const dataModeOption = option.getAttribute("data-mode");

  if (!modeSaved) self[0].classList.add("active");

  if (modeSaved === dataModeOption) option.classList.add("active");

  const input = option.querySelector("input[type=radio]");

  input.addEventListener("click", (e) => {
    playgroundModeOptions.forEach((option) => {
      option.classList.remove("active");
    });
    const { value } = e.target;
    if (value === dataModeOption) option.classList.add("active");
    localStorage.setItem("@cel-playground:mode", value);
    setTimeout(() => closeModal(), 1000);
  });
});

function renderModeOptions() {
  const el = document.querySelector(".playground-modes__options");
  const playCelModeKeys = Object.keys(MODES);
  playCelModeKeys.forEach((key) => {
    const mode = MODES[key];

    const divOption = document.createElement("div");
    divOption.className = "playground-modes__options--option";
    divOption.setAttribute("data-mode", mode);

    const label = document.createElement("label");
    label.htmlFor = mode;
    label.innerHTML = key;

    const input = document.createElement("input");
    input.className = "playground-modes__options--option-input";
    input.type = "radio";
    input.name = mode;
    input.id = mode;
    input.value = mode;

    divOption.appendChild(label);
    divOption.appendChild(input);

    el.appendChild(divOption);
  });
}

function closeModal() {
  playgroundModesModalEl.style.display = "none";
}
