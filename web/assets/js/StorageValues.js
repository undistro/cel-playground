const KEY = "values";

export class StorageValues {
  constructor() {
    this.values = JSON.parse(localStorage.getItem(KEY));
  }

  setValues(values) {
    localStorage.setItem(KEY, JSON.stringify(values));
  }

  getValues() {
    return this.values;
  }
}
