const costTooltipElement = document.getElementById("cost-tooltip");
const costInfoElement = document.getElementById("cost-info");

function openCostInfoTooltip() {
  costTooltipElement.style.display = "initial";
}

costInfoElement.addEventListener("mouseover", openCostInfoTooltip);

costInfoElement.addEventListener("mouseleave", () => {
  costTooltipElement.style.display = "none";
});
