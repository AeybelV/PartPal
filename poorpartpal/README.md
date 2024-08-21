# PoorPartPal

PoorPartPal is a CLI tool that acts as the front-end for the PartPal library. It
allows users to optimize their Bill of Materials (BOM) for cost, shipping, and
availability by querying multiple electronic component distributors
like Digi-Key, Mouser, etc.

## Features

- **BOM Optimization**: Optimize your BOM for cost, shipping, and availability
  across multiple distributors.
- **Multi-Distributor Support**: Supports distributors like Digi-Key and Mouser,
  with easy integration for more.
- **Interactive Interface**: Configure your BOM within the TUI interface of the tool

## Usage

```sh
poorpartpal optimize my_bom.csv
```

This will look up the prices and availability of components in my_bom.csv and
optimize it based on cost and stock availability.

You can provide more options to strategize for.

- --output, -o: Specify the output file for the optimized BOM.
- --distributors, -d: Specify which distributors to query (e.g., digikey, mouser).
- --strategy, -s: Choose the optimization strategy (e.g., cost, shipping).
- --verbose, -v: Enable verbose output for debugging and detailed logs.

Example:

```sh
poorpartpal optimize my_bom.csv -o optimized_bom.html -d digikey mouser -s cost
```

For a full list of commands and options:

```sh
poorpartpal --help
```

## Configuration

PoorPartPal uses a configuration file to store API keys and preferences for distributors.
By default, this file is located at `~/.poorpartpal/config.yaml`.

Here’s a sample configuration file:

```yaml
distributors:
  digikey:
    api_key: YOUR_DIGIKEY_API_KEY
  mouser:
    api_key: YOUR_MOUSER_API_KEY
```

## Contirbuting

Contributions are welcome! Please feel free to submit a pull request
or open an issue to discuss improvements or bugs.

## License

PoorPartPal is licenced under the GNU GPL.

## Acknowledgements

- Mathew Yu for the idea
