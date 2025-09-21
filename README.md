# Dog Food Calculator

A modern GUI application for calculating dog food portions using polynomial regression. Built with Go and Fyne.

## Features

- **Polynomial Regression Model**: Uses 3rd-degree polynomial regression to predict food portions based on dog weight
- **Interactive GUI**: Clean, user-friendly interface built with Fyne
- **Data Management**: Edit and save feeding data points
- **Real-time Calculation**: Instant portion calculations
- **Persistent Storage**: Data stored in JSON format
- **Structured Logging**: Comprehensive logging with slog
- **Thread-safe**: Proper mutex usage for concurrent access

## Installation

### Prerequisites

- Go 1.21 or later
- Fyne dependencies (see [Fyne documentation](https://developer.fyne.io/started/))

### Build from Source

```bash
# Clone the repository
git clone https://github.com/odinnordico/fedo.git
cd fedo

# Install dependencies
go mod tidy

# Build the application
go build

# Run the application
./fedo
```

## Usage

1. **Launch the Application**: Run the built executable
2. **Enter Dog Weight**: Input your dog's weight in kilograms
3. **Calculate Portion**: Click "Calculate" to get the recommended food portion
4. **Edit Data**: Click "Edit Data" to modify the feeding guide data points
5. **Save Changes**: Save your edits to persist changes

### Data Format

The feeding data is stored in `feeding_data.json`:

```json
[
  {"Weight_kg": 1.0, "Daily_g": 33},
  {"Weight_kg": 2.25, "Daily_g": 50},
  ...
]
```

## Architecture

The application is organized into several modules:

- `models.go` - Data structures and global variables
- `data.go` - Data loading, saving, and model training
- `calculator.go` - Portion calculation logic
- `ui.go` - User interface components
- `main.go` - Application entry point

## Development

### Project Structure

```
fedo/
├── main.go           # Application entry point
├── models.go         # Data structures
├── data.go           # Data management and ML
├── calculator.go     # Calculation logic
├── ui.go             # User interface
├── feeding_data.json # Feeding data
├── go.mod            # Go module file
├── go.sum            # Go dependencies
├── fedo.png          # Application icon
├── README.md         # This file
├── LICENSE           # MIT License
├── CONTRIBUTING.md   # Contribution guidelines
└── CODE_OF_CONDUCT.md # Code of conduct
```

### Key Technologies

- **Go**: Programming language
- **Fyne**: GUI framework
- **Gonum**: Numerical computing library
- **Slog**: Structured logging

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Fyne](https://fyne.io/)
- Numerical computations powered by [Gonum](https://gonum.org/)