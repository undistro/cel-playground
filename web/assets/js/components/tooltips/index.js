const tooltipContainers = document.querySelectorAll(".tooltip__container");

tooltipContainers.forEach((container) => {
  const tooltipTrigger = container.querySelector("#tooltip__trigger");
  const tooltipContent = container.querySelector("#tooltip__content");
  tooltipTrigger?.addEventListener("mouseover", () => {
    tooltipContent.style.display = "initial";
  });

  tooltipTrigger?.addEventListener("mouseleave", () => {
    tooltipContent.style.display = "none";
  });
});
