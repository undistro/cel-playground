const localStorageKey = "@cel-playground:mode";

const MODES = Object.freeze({
  CEL: { value: "CEL", title: "CEL Expression", html: "" },
  VAP: {
    value: "VAP",
    title: "Validating Admission Policy",
    html: [{ selector: "", string: "" }],
  },
  WEB_HOOKS: { value: "WEB_HOOKS", title: "Web Hooks", html: "" },
  AUTH_CMAPPING: {
    value: "AUTH_CMAPPING",
    title: "Authentication Claim Mapping",
    html: "",
  },
  AUTH: { value: "AUTH", title: "Authentication", html: "" },
});
