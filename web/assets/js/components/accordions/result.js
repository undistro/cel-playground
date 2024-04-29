/**
 * Copyright 2024 Undistro Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { createTooltip } from "../tooltips/index.js";

const outputResultEl = document.getElementById("editor__output-result");
const holderEl = document.querySelector(".editor__output-holder");

function createAccordionItemsByResults(name, result, index) {
  const listItem = document.createElement("li");
  listItem.className = "editor__output-result-accordion";
  listItem.setAttribute("data-open", "false");
  listItem.onclick = (e) => {
    const isAccordionOpen = listItem.getAttribute("data-open") === "true";
    if (isAccordionOpen) listItem.setAttribute("data-open", "false");
    else listItem.setAttribute("data-open", "true");
  };

  const accordionContent = document.createElement("div");
  accordionContent.className = "result-accordion-content";
  accordionContent.appendChild(createLabel(result, name, index));
  const costSpan = document.createElement("span");
  costSpan.innerHTML = `Cost: ${result?.cost ?? "-"}`;
  accordionContent.appendChild(costSpan);

  const expansibleContent = document.createElement("div");
  expansibleContent.className = "result-accordion-expansible-content";
  expansibleContent.innerHTML = `<span>${getResultValue(result)}</span>`;

  listItem.appendChild(accordionContent);
  listItem.appendChild(expansibleContent);

  outputResultEl.appendChild(listItem);
}
function getResultValue(result) {
  if (result.isError) {
    return `<span style="color:#e01e5a">${result.error}</span>`;
  } else if ("value" in result) {
    if (typeof result.value === "object")
      return `<pre>${JSON.stringify(result.value, null, 2)}</pre>`;
    return String(result.value);
  }

  return result.result;
}

function createLabel(item, name, i) {
  const parentContainer = document.createElement("div");
  parentContainer.style =
    "display: flex; align-items: center; gap: 0.5rem; position:relative";

  const arrowIcon = document.createElement("i");
  arrowIcon.className = "ph ph-caret-right ph-bold result-arrow";

  const span = document.createElement("span");
  span.innerHTML = `${item.name ? `${name}.${item.name}` : `${name}[${i}]`}`;

  parentContainer.appendChild(arrowIcon);

  if (item?.isError) {
    const errorIcon = document.createElement("i");
    errorIcon.className = "ph ph-x-circle ph-fill";
    errorIcon.style =
      "color: #e01e5a; z-index:999999; display:flex;align-items:center;justify-content: center;";

    const errorIconWithTooltip = createTooltip({
      contentText: "Validation compilation failed.",
      triggerElement: errorIcon,
      position: {
        left: 50,
        top: -10,
      },
    });
    parentContainer.appendChild(errorIconWithTooltip);
  }

  parentContainer.appendChild(span);

  return parentContainer;
}

function renderAccordions(key, values, index = 0) {
  if (Array.isArray(values))
    values.forEach((value, index) => renderAccordions(key, value, index));
  else createAccordionItemsByResults(key, values, index);
}

export function hideAccordions() {
  outputResultEl.style.display = "none";
}

export function handleRenderAccordions(result) {
  outputResultEl.innerHTML = "";
  outputResultEl.style.display = "flex";
  outputResultEl.scrollTo({ top: 0, behavior: "smooth" });
  holderEl.style.overflowY = "auto";
  holderEl.style.overflowX = "hidden";

  Object.entries(result).forEach(([key, values]) => {
    renderAccordions(key, values);
  });
}
