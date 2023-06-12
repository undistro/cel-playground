const EDITOR_ELEMENTS = ["cel-input", "data-input"];

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

const SAMPLE_DATA = `{
    "object": {
        "replicas": 42,
        "name": "sample",
        "message": "Hello, world!",
        "items": [1, 2, 3, 4],
        "spec": {
            "foo": "bar",
            "key": "value"
        },
        "enabled": true,
        "status": null
    }
}`;

const WASM_URL = "main.wasm.gz";

const EXAMPLES = [
  {
    name: "Not allowed hostPort",
    src: "examples/example.cel",
    data: "examples/example.yaml",
  },
  {
    name: "Not allowed seccomp profile",
    src: "examples/example.cel",
    data: "examples/example.yaml",
  },
  {
    name: "Privileged access to the Windows node",
    src: "examples/example.cel",
    data: "examples/example.yaml",
  },
  {
    name: "Forbidden seccomp profile",
    src: "examples/example.cel",
    data: "examples/example.yaml",
  },
  {
    name: "Automounted service account token",
    src: "examples/example.cel",
    data: "examples/example.yaml",
  },
];

export { EDITOR_ELEMENTS, EDITOR_DEFAULTS, SAMPLE_DATA, WASM_URL, EXAMPLES};
