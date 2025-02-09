# Artestian

![Artestian Logo](https://github.com/user-attachments/assets/a2a45048-b3f8-4d8c-9127-c02b8babf7cd)

Artestian is your AI-powered assistant for reducing the boilerplate and setup frustrations of writing tests. While it won’t replace your careful thought on designing tests, it helps you focus on what's important: ensuring your code works as expected.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Quick Start](#quick-start)
  - [Configuration Details](#configuration-details)
  - [Context Files](#context-files)
- [CLI Flags and Environment Variables](#cli-flags-and-environment-variables)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Minimize Boilerplate:** Automatically generate test files and setup code.
- **AI-Powered Assistance:** Uses example templates and contextual information to help generate tests.
- **Flexible Configuration:** Customize the behavior through a simple JSON configuration file.
- **Multiple Test Runners:** Supports both Jest (for TypeScript) and Go testing frameworks.

---

## Installation

### Prerequisites

- **Go**: Version 1.23 or later.
- **API Keys**: 
  - OpenAI API key (if using the OpenAI provider)
  - Anthropic API key (if using the Anthropic provider)

Set your environment variables:

```bash
export OPENAI_API_KEY=your_openai_key_here
export ANTHROPIC_API_KEY=your_anthropic_key_here
```

### Clone the Repository

```bash
git clone https://github.com/username/artestian.git
cd artestian
```

---

## Usage

Artestian uses a JSON configuration file to manage test examples and settings.

### Quick Start

1. **Create a Configuration File**

   Create a file (e.g., `config.json`) with the following structure:

   ```json
   {
     "version": "1.0",
     "examples": [
       {
         "name": "Basic Unit Test",
         "type": "unit",
         "file_path": "./examples/unit-test.example.ts",
         "description": "Basic unit test for a simple function"
       }
     ],
     "settings": {
       "default_test_directory": "__tests__",
       "language": "typescript",
       "test_runner": "jest"
     },
     "context": {
       "files": [
         {
           "path": "./types/interfaces.ts",
           "description": "Core type definitions and interfaces",
           "type": "types"
         }
       ]
     }
   }
   ```

2. **Add Example Test Files**

   Create the test example files referenced in your configuration. For instance, create `./examples/unit-test.example.ts`:

   ```typescript
   // examples/unit-test.example.ts
   import { describe, it, expect } from '@jest/globals';
   import { MyService } from './my-service';

   describe('MyService', () => {
     it('should perform the expected operation', () => {
       const service = new MyService();
       const result = service.doSomething();
       expect(result).toBe(true);
     });
   });
   ```

3. **Run Artestian**

   Execute the tool by pointing it to your configuration file:

   ```bash
   artestian -config ./config.json
   ```

### Configuration Details

Artestian’s behavior is driven by a JSON configuration file. The configuration contains the following sections:

#### Required Fields

- **version**: The configuration version (e.g., `"1.0"`).
- **examples**: List of test examples:
  - `name`: Unique identifier for the test.
  - `type`: Type of test (choose from `"unit"`, `"integration"`, `"worker"`, or `"prompt"`).
  - `file_path`: Relative path to the example test file.
  - `description`: A brief explanation of the test.
- **settings**: Global settings:
  - `default_test_directory`: Directory where tests will be generated.
  - `language`: Programming language (either `"typescript"` or `"go"`).
  - `test_runner`: Test runner to use (either `"jest"` or `"go test"`).

#### Optional Fields

- **context**: Additional files to provide richer context for test generation.
  - `files`: An array of files used as supporting context:
    - `path`: File location.
    - `description`: Description of its content.
    - `type`: Type of context (e.g., `"types"`, `"utils"`, `"constants"`).

### Context Files

Context files can help the AI better understand your codebase. For example:

```json
{
  "context": {
    "files": [
      {
        "path": "./types/interfaces.ts",
        "description": "Core type definitions and interfaces",
        "type": "types"
      },
      {
        "path": "./utils/test-helpers.ts",
        "description": "Common test utilities and helper functions",
        "type": "utils"
      },
      {
        "path": "./constants/test-data.ts",
        "description": "Test data constants and fixtures",
        "type": "constants"
      }
    ]
  }
}
```

---

## CLI Flags and Environment Variables

### Command Line Flags

- `-config`: Path to your configuration JSON file.
- `-ai`: Selects the AI provider (`"openai"` or `"anthropic"`). Default is `"anthropic"`.
- `-log-level`: Sets the log verbosity (`"debug"`, `"info"`, `"warn"`, `"error"`). Default is `"info"`.

### Environment Variables

- `OPENAI_API_KEY`: Your OpenAI API key (if using OpenAI).
- `ANTHROPIC_API_KEY`: Your Anthropic API key (if using Anthropic).

---

## Project Structure

An overview of the project architecture:

```
.
├── cmd/
│   └── artestian/        # CLI entry point
├── pkg/                  # Core modules and logic
└── types/                # Shared types and interfaces
```

---

## Contributing

Contributions are warmly welcomed! If you'd like to contribute:

1. Fork the repository.
2. Create a branch for your feature or bug fix.
3. Submit a pull request with a clear description of your changes.

Please adhere to our coding standards and include tests where applicable.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

*Happy Testing!*