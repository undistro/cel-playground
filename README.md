# CEL Playground
![GitHub](https://img.shields.io/github/license/undistro/cel-playground)
[![Go Report Card](https://goreportcard.com/badge/github.com/undistro/cel-playground)](https://goreportcard.com/report/github.com/undistro/cel-playground)

CEL Playground is an interactive WebAssembly (Wasm) powered environment to explore and experiment with the [Common Expression Language (CEL)](https://github.com/google/cel-spec).
It provides a simple and user-friendly interface to write and quickly evaluate CEL expressions.

## CEL libraries

CEL Playground is built by compiling Go code to WebAssembly and includes the following libraries that are available in the environment:

- CEL [extended string function library](https://pkg.go.dev/github.com/google/cel-go/ext#Strings)
- [Kubernetes list library](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-list-library)
- [Kubernetes regex library](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-regex-library)
- [Kubernetes URL library](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-url-library)

Take a look at [all the environment options](eval/eval.go#L26).

## Development

Build the Wasm binary:
```shell
make build
```

Serve the static files:
```shell
make serve
```

## Contributing

We appreciate your contribution.
Please refer to our [contributing guideline](https://github.com/undistro/cel-playground/blob/main/CONTRIBUTING.md) for further information.
This project adheres to the Contributor Covenant [code of conduct](https://github.com/undistro/cel-playground/blob/main/CODE_OF_CONDUCT.md).

## License

CEL Playground is available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
