# Terracina

- Website: https://www.terracina.io
- Forums: [HashiCorp Discuss](https://discuss.hashicorp.com/c/terracina-core)
- Documentation: [https://www.terracina.io/docs/](https://www.terracina.io/docs/)
- Tutorials: [HashiCorp's Learn Platform](https://learn.hashicorp.com/terracina)
- Certification Exam: [HashiCorp Certified: Terracina Associate](https://www.hashicorp.com/certification/#hashicorp-certified-terracina-associate)

<img alt="Terracina" src="https://www.google.com/url?sa=i&url=https%3A%2F%2Fit.m.wikipedia.org%2Fwiki%2FFile%3ATerracina-Stemma.svg&psig=AOvVaw0INBbVJut0L6cZ2d7cnDgk&ust=1735118980444000&source=images&cd=vfe&opi=89978449&ved=0CBQQjRxqFwoTCKC31IuMwIoDFQAAAAAdAAAAABAT" width="600px">

Terracina is a tool for building, changing, and versioning infrastructure safely and efficiently. Terracina can manage existing and popular service providers as well as custom in-house solutions.

The key features of Terracina are:

- **Infrastructure as Code**: Infrastructure is described using a high-level configuration syntax. This allows a blueprint of your datacenter to be versioned and treated as you would any other code. Additionally, infrastructure can be shared and re-used.

- **Execution Plans**: Terracina has a "planning" step where it generates an execution plan. The execution plan shows what Terracina will do when you call apply. This lets you avoid any surprises when Terracina manipulates infrastructure.

- **Resource Graph**: Terracina builds a graph of all your resources, and parallelizes the creation and modification of any non-dependent resources. Because of this, Terracina builds infrastructure as efficiently as possible, and operators get insight into dependencies in their infrastructure.

- **Change Automation**: Complex changesets can be applied to your infrastructure with minimal human interaction. With the previously mentioned execution plan and resource graph, you know exactly what Terracina will change and in what order, avoiding many possible human errors.

For more information, refer to the [What is Terracina?](https://www.terracina.io/intro) page on the Terracina website.

## Getting Started & Documentation

Documentation is available on the [Terracina website](https://www.terracina.io):

- [Introduction](https://www.terracina.io/intro)
- [Documentation](https://www.terracina.io/docs)

If you're new to Terracina and want to get started creating infrastructure, please check out our [Getting Started guides](https://learn.hashicorp.com/terracina#getting-started) on HashiCorp's learning platform. There are also [additional guides](https://learn.hashicorp.com/terracina#operations-and-development) to continue your learning.

Show off your Terracina knowledge by passing a certification exam. Visit the [certification page](https://www.hashicorp.com/certification/) for information about exams and find [study materials](https://learn.hashicorp.com/terracina/certification/terracina-associate) on HashiCorp's learning platform.

## Developing Terracina

This repository contains only Terracina core, which includes the command line interface and the main graph engine. Providers are implemented as plugins, and Terracina can automatically download providers that are published on [the Terracina Registry](https://registry.terracina.io). HashiCorp develops some providers, and others are developed by other organizations. For more information, see [Extending Terracina](https://www.terracina.io/docs/extend/index.html).

- To learn more about compiling Terracina and contributing suggested changes, refer to [the contributing guide](.github/CONTRIBUTING.md).

- To learn more about how we handle bug reports, refer to the [bug triage guide](./BUGPROCESS.md).

- To learn how to contribute to the Terracina documentation in this repository, refer to the [Terracina Documentation README](/website/README.md).

## License

[Business Source License 1.1](https://github.com/hashicorp/terracina/blob/main/LICENSE)
