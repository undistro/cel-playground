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
