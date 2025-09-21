# Contributing to Dog Food Calculator

Thank you for your interest in contributing to the Dog Food Calculator project! We welcome contributions from the community.

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker to report bugs or request features
- Provide detailed information including:
  - Steps to reproduce the issue
  - Expected vs. actual behavior
  - Your environment (OS, Go version, etc.)
  - Screenshots if applicable

### Development Process

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/odinnordico/fedo.git`
3. **Create** a feature branch: `git checkout -b feature/your-feature-name`
4. **Make** your changes
5. **Test** your changes: `go test ./...`
6. **Build** and verify: `go build`
7. **Commit** your changes: `git commit -am 'Add some feature'`
8. **Push** to the branch: `git push origin feature/your-feature-name`
9. **Create** a Pull Request

### Code Guidelines

#### Go Code Style

- Follow standard Go formatting: `go fmt`
- Use `gofmt` and `goimports` for consistent formatting
- Follow Go naming conventions
- Write clear, concise comments
- Use meaningful variable and function names

#### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in imperative mood (e.g., "Add", "Fix", "Update")
- Keep the first line under 50 characters
- Add detailed description if needed

#### Testing

- Write tests for new features
- Ensure all tests pass before submitting
- Test edge cases and error conditions

### Project Structure

```
fedo/
├── main.go           # Application entry point
├── models.go         # Data structures and globals
├── data.go           # Data management and ML training
├── calculator.go     # Calculation logic
├── ui.go             # User interface components
├── feeding_data.json # Default feeding data
└── fedo.png          # Application icon
```

### Key Areas for Contribution

- **UI Improvements**: Enhance the user interface
- **Algorithm Enhancements**: Improve the regression model
- **Data Validation**: Add better input validation
- **Performance**: Optimize calculations and memory usage
- **Documentation**: Improve code comments and documentation
- **Testing**: Add comprehensive test coverage
- **Cross-platform**: Improve compatibility across platforms

### Pull Request Process

1. Ensure your code follows the guidelines above
2. Update documentation if needed
3. Add tests for new functionality
4. Ensure CI/CD checks pass
5. Request review from maintainers

### Code Review

All submissions require review. We will:

- Review code for style and functionality
- Test the changes
- Provide constructive feedback
- Merge approved changes

## Getting Help

- Check existing issues and documentation first
- Ask questions in GitHub discussions
- Contact maintainers for guidance

## License

By contributing to this project, you agree that your contributions will be licensed under the same MIT License that covers the project.