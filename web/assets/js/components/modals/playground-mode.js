import {
  renderExamplesInSelectInstance,
  renderTabs,
} from "../../utils/render-functions.js";
import { AceEditor } from "../../editor.js";
import { ModesService } from "../../services/modes.js";

const celEditor = new AceEditor("cel-input");
const dataEditor = new AceEditor("data-input");

const playgroundModesModalEl = document.getElementById(
  "playground-modes__modal"
);

var modal = modal || {
  set: Array.from(document.querySelectorAll("[data-modal]")),
  openTriggers: Array.from(document.querySelectorAll("[data-modal-trigger]")),
  closeTriggers: Array.from(document.querySelectorAll("[data-modal-close]")),
  focusable:
    "a[href], area[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), button:not([disabled]), object, embed, *[tabindex], *[contenteditable], label, div.playground-modes__options",
  focused: "",
};

modal.init = function () {
  modal.set.forEach((modal) => {
    modal.setAttribute("aria-hidden", "true");
  });
  modal.openTriggers.forEach((trigger) => {
    trigger.addEventListener("click", function (e) {
      e.preventDefault();
      let name = this.dataset.modalTrigger;
      modal.el = modal.set.find(function (value) {
        return value.dataset.modal === name;
      });
      modal.show();
    });
  });
  modal.closeTriggers.forEach((trigger) => {
    trigger.addEventListener("click", function (e) {
      e.preventDefault();
      modal.hide();
    });
  });
};

modal.show = function () {
  document.body.classList.add("has-modal");
  document.querySelector(".main-content").setAttribute("aria-hidden", true);
  modal.focused = document.activeElement;
  modal.el.setAttribute("aria-hidden", "false");
  modal.el.classList.add("modal-show");
  modal.focusableChildren = Array.from(
    modal.el.querySelectorAll(modal.focusable) ?? []
  );
  modal.focusableChildren[0].focus();
  modal.el.onkeydown = function (e) {
    modal.trap(e);
  };
};

modal.hide = function () {
  document.body.classList.remove("has-modal");
  document.querySelector(".main-content").setAttribute("aria-hidden", false);
  modal.el.setAttribute("aria-hidden", "true");
  modal.el.classList.remove("modal-show");
  modal.focused.focus();
};

window.onclick = (e) => {
  if (e.target === playgroundModesModalEl) modal.hide();
};

modal.trap = function (e) {
  if (e.which == 27) {
    modal.hide();
  }
  if (e.which == 9) {
    let currentFocus = document.activeElement;
    let totalOfFocusable = modal.focusableChildren.length;
    let focusedIndex = modal.focusableChildren.indexOf(currentFocus);
    if (e.shiftKey) {
      if (focusedIndex === 0) {
        e.preventDefault();
        modal.focusableChildren[totalOfFocusable - 1].focus();
      }
    } else {
      if (focusedIndex == totalOfFocusable - 1) {
        e.preventDefault();
        modal.focusableChildren[0].focus();
      }
    }
  }
};

renderModeOptions();
modal.init();

function handleModeClick(event, mode, element) {
  const { value } = event.target;

  document
    .querySelectorAll(".playground-modes__options--option")
    .forEach((option) => option.classList.remove("active"));

  localStorage.removeItem("example-selected");

  element.classList.add("active");
  renderUIChangesByMode(mode);
  localStorage.setItem(localStorageKey, value);
  celEditor.setValue("", -1);
  dataEditor.setValue("", -1);
  setTimeout(() => modal.hide(), 1000);
}

function renderModeOptions() {
  const el = document.querySelector(".playground-modes__options");

  ModesService.getModes()
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
