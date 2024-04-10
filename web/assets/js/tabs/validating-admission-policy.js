const vapTabsEl = document.getElementById("vap__tabs");

const vapTabButtonsEl = document.querySelectorAll(".vap__tabs-button");

vapTabButtonsEl.forEach((button) => {
  button.addEventListener("click", () => {
    vapTabButtonsEl.forEach((button) => {
      button.classList.remove("active");
    });

    button.classList.add("active");
  });
});

function renderVAPTabs() {
  const savedMode = localStorage.getItem("@cel-playground:mode") ?? "CEL";
  if (savedMode === "VAP") {
    vapTabsEl.style.display = "flex";
    vapTabsEl.previousElementSibling.innerHTML = "Inputs: ";
  } else {
    vapTabsEl.style.display = "none";
    vapTabsEl.previousElementSibling.innerHTML = "Input";
  }
}
