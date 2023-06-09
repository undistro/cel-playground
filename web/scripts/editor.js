import { EDITOR_DEFAULTS } from "./constants.js";

class AceEditor {
  constructor(id) {
    this.editor = ace.edit(id);
    this.editor.setTheme(EDITOR_DEFAULTS[id].theme);
    this.editor.getSession().setMode(EDITOR_DEFAULTS[id].mode);
    this.editor.getSession().setUseWorker(false);
  }

  setValue(value, cursorPosition = 0) {
    this.editor.setValue(value, cursorPosition);
  }

  getValue() {
    return this.editor.getValue();
  }
}

export { AceEditor };
