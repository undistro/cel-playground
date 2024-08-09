# CEL Playground
![GitHub](https://img.shields.io/github/license/undistro/cel-playground)
[![Go Report Card](https://goreportcard.com/badge/github.com/undistro/cel-playground)](https://goreportcard.com/report/github.com/undistro/cel-playground)
[![slack](https://img.shields.io/badge/Slack-Join-4a154b?logo=slack)](https://join.slack.com/t/undistrocommunity/shared_invite/zt-21slyrao4-dTW_XtOB90QVj05txOX6rA)

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

## Community

To engage with our community, you can use the following resources:
- Give us a [star :star:](https://github.com/undistro/cel-playground/stargazers) - If you want us to continue developing and improving CEL Playground
- [Contributing to CEL Playground](https://github.com/undistro/cel-playground/blob/main/CONTRIBUTING.md) - Start here if you're interested in contributing to the project
- [Code of Conduct](https://github.com/undistro/cel-playground/blob/main/CODE_OF_CONDUCT.md) - Learn about the guidelines that govern our community interactions
- [Slack Channel](https://join.slack.com/t/undistrocommunity/shared_invite/zt-21slyrao4-dTW_XtOB90QVj05txOX6rA) - Join us on Slack to get support or discuss the project
- [Community Sessions](https://tinyurl.com/undistro-community-calendar) - Join our monthly community meetings and bi-weekly office hours ([agenda and meeting notes](https://docs.google.com/document/d/13AhGiyIiX58UJMw7CDJi_T8e1_SC7f7p1kE2PcyDwRU/edit#heading=h.7k7sl4hlyyqw))
- [Roadmap](https://github.com/undistro/cel-playground/blob/main/roadmap.md) - Discover what's next for the project
- [Adopters](https://github.com/undistro/cel-playground/blob/main/ADOPTERS.md) - Is your company using CEL Playground? Let us know, and we'll feature you here!

## License

CEL Playground is available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
