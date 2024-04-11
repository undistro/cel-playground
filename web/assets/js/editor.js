/**
 * Copyright 2023 Undistro Authors
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

const EDITOR_DEFAULTS = {
  "cel-input": {
    theme: "ace/theme/clouds",
    mode: "ace/mode/javascript",
  },
  "data-input": {
    theme: "ace/theme/clouds",
    mode: "ace/mode/yaml",
  },
};

class AceEditor {
  constructor(id) {
    this.editor = ace.edit(id);
    this.editor.setTheme(EDITOR_DEFAULTS[id].theme);
    this.editor.setShowPrintMargin(false);
    this.editor.getSession().setMode(EDITOR_DEFAULTS[id].mode);
    this.editor.getSession().setUseWorker(false);
  }

  setValue(value, cursorPosition = 0) {
    this.editor.setValue(value, cursorPosition);
  }

  getValue() {
    return this.editor.getValue();
  }

  setExpressionSyntax(syntax) {
    this.editor.getSession().setMode(`ace/mode/${syntax ?? "javascript"}`);
  }
}

export { AceEditor };
