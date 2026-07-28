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
import {
  COPY_FEEDBACK_MS,
  copyToClipboard,
} from "../../utils/clipboard.js";

const outputResultEl = document.getElementById("editor__output-result");
const holderEl = document.querySelector(".editor__output-holder");

// Everything the evaluator returns is built from what the user typed -- CEL
// errors quote the source expression verbatim, and messageExpression results
// are whatever string the policy produced. It all reaches innerHTML below, and
// share.js encodes the editor state into a link, so it must be escaped.
function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

// sectionLabel names a result section, for the row and for the copy button's
// accessible name. item.name is a validation or variable name straight out of
// the policy.
function sectionLabel(item, name, i) {
  return item.name ? `${name}.${item.name}` : `${name}[${i}]`;
}

// aria-label wins the accessible-name computation over title, so both have to
// move together or the icon swap stays visual-only.
function setCopyLabel(button, label) {
  button.title = label;
  button.setAttribute("aria-label", label);
}

function createCopyButton(text, label) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = "result-accordion-copy";

  const icon = document.createElement("i");
  icon.className = "ph ph-copy";
  button.appendChild(icon);

  // Every section renders one of these, so an undifferentiated "Copy to
  // clipboard" repeated N times is not navigable.
  const idleLabel = label ? `Copy ${label}` : "Copy to clipboard";
  setCopyLabel(button, idleLabel);

  let resetTimer;
  const showFeedback = (copied) => {
    icon.className = copied ? "ph ph-check" : "ph ph-x";
    button.classList.toggle("copied", copied);
    button.classList.toggle("copy-failed", !copied);
    setCopyLabel(button, copied ? "Copied" : "Could not copy");

    window.clearTimeout(resetTimer);
    resetTimer = window.setTimeout(() => {
      icon.className = "ph ph-copy";
      button.classList.remove("copied", "copy-failed");
      setCopyLabel(button, idleLabel);
    }, COPY_FEEDBACK_MS);
  };

  button.onclick = (event) => {
    // The whole row toggles the accordion, so copying must not also collapse
    // the section the user is reading.
    event.stopPropagation();
    copyToClipboard(text).then(showFeedback);
  };

  return button;
}

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
  // The cost and the copy button share a right-hand group, so the label still
  // sits opposite them under the row's space-between.
  const actions = document.createElement("div");
  actions.className = "result-accordion-actions";
  const costSpan = document.createElement("span");
  costSpan.innerHTML = `Cost: ${result?.cost ?? "-"}`;
  actions.appendChild(costSpan);
  const { text, html } = renderResult(result);
  if (text !== "") {
    const label = sectionLabel(result, name, index);
    actions.appendChild(createCopyButton(text, label));
  }
  accordionContent.appendChild(actions);

  const expansibleContent = document.createElement("div");
  expansibleContent.className = "result-accordion-expansible-content";
  expansibleContent.innerHTML = `<span>${html}</span>`;

  listItem.appendChild(accordionContent);
  listItem.appendChild(expansibleContent);

  outputResultEl.appendChild(listItem);
}
// renderResult produces both forms of a section body from one branch structure:
// the markup for the panel, and the plain text the copy button puts on the
// clipboard. They cannot drift apart while they are built together.
//
// text is "" wherever the markup is a placeholder describing the value rather
// than being it, and no copy button is offered for those.
function renderResult(result) {
  if (result.isError) {
    const text = result.error == null ? "" : String(result.error);
    return {
      text,
      html: `<span class="result-error">${escapeHtml(text)}</span>`,
    };
  }
  if ("message" in result) return renderText(result.message);
  if ("value" in result) {
    if (typeof result.value === "object") {
      const text = JSON.stringify(result.value, null, 2);
      return { text, html: `<pre>${escapeHtml(text)}</pre>` };
    }
    return renderText(result.value);
  }

  // omitempty drops the key entirely when the expression yielded nothing.
  if (result.result === undefined) {
    return { text: "", html: `<span class="result-empty">(no value)</span>` };
  }
  return renderText(result.result);
}

// An expression can legitimately produce nothing to print, and an open section
// with an empty body reads as a broken panel rather than as a result. The two
// ways of producing nothing stay separate placeholders: optional.none() and an
// expression returning "" are different outcomes.
function renderText(value) {
  const text = String(value);
  if (text === "") {
    return { text, html: `<span class="result-empty">(empty string)</span>` };
  }
  return { text, html: escapeHtml(text) };
}

function createLabel(item, name, i) {
  const parentContainer = document.createElement("div");
  parentContainer.className = "result-accordion-label";

  const arrowIcon = document.createElement("i");
  arrowIcon.className = "ph ph-caret-right ph-bold result-arrow";

  const span = document.createElement("span");
  span.innerHTML = escapeHtml(sectionLabel(item, name, i));

  parentContainer.appendChild(arrowIcon);

  if (item?.isError) {
    const errorIcon = document.createElement("i");
    errorIcon.className = "ph ph-x-circle ph-fill result-accordion-error-icon";

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
  holderEl.scrollTo({ top: 0, behavior: "smooth" });
}

export function handleRenderAccordions(result) {
  outputResultEl.innerHTML = "";
  outputResultEl.style.display = "flex";

  holderEl.scrollTo({ top: 0, behavior: "smooth" });
  holderEl.style.overflowY = "auto";
  holderEl.style.overflowX = "hidden";

  Object.entries(result).forEach(([key, values]) => {
    renderAccordions(key, values);
  });
}
