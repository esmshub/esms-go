# ESMS Go!

ESMS Go! is a unified command-line interface for the ESMS fantasy football engine that 
wraps the core match engine and supporting utilities that were traditionally separate executables.

## Guiding Principles

### Cross-Platform Support
- Develop core application logic using platform-agnostic APIs and libraries.
- Establish a build pipeline to generate executables for all targeted operating systems (Windows, macOS, Linux).

### Single-Executable Deployment
- Package all application components and dependencies into a single executable file for each supported OS.
- Implement seamless installation-free usage—download & run with zero configuration or setup.

### Plugin Support (Create-Your-Own Engine)
- Design and implement a modular plugin architecture with clearly defined APIs and lifecycle hooks.
- Provide SDKs, templates, and documentation to enable external developers to create custom plugins.

### Config Consolidation
- Centralize all configuration into a well-structured, consistent format (preferably JSON or YAML).
- Define strict configuration schemas with validation rules to ensure predictability and correctness.
- Maintain backward compatibility or migration paths for legacy configuration formats.

### Maintainability
- Establish and enforce coding standards, patterns, and best practices across all areas.
