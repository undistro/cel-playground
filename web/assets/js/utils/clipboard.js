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

// document.execCommand is deprecated, but navigator.clipboard only exists in a
// secure context, which the playground is not when reached over a LAN address
// rather than localhost.
function copyWithExecCommand(text) {
  const textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "");
  textarea.style.position = "fixed";
  textarea.style.top = "-1000px";
  textarea.style.opacity = "0";
  // execCommand copies the *document* selection, so this clobbers whatever the
  // user had selected and moves focus to the textarea. Both are put back below.
  const previousFocus = document.activeElement;
  const selection = window.getSelection();
  const previousRanges = [];
  for (let i = 0; i < selection.rangeCount; i++) {
    previousRanges.push(selection.getRangeAt(i));
  }

  document.body.appendChild(textarea);
  textarea.select();
  textarea.setSelectionRange(0, text.length);

  let copied = false;
  try {
    copied = document.execCommand("copy");
  } catch (error) {
    console.error("failed to copy to the clipboard", error);
  }
  document.body.removeChild(textarea);

  selection.removeAllRanges();
  previousRanges.forEach((range) => selection.addRange(range));
  previousFocus?.focus?.();
  return copied;
}

// How long a copy affordance stays in its "copied" state.
export const COPY_FEEDBACK_MS = 1000;

// Resolves to whether the text reached the clipboard, and never rejects:
// callers decide what to show.
export async function copyToClipboard(text) {
  if (navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch (error) {
      console.error("failed to copy to the clipboard", error);
    }
  }
  return copyWithExecCommand(text);
}
