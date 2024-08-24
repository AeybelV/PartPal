import json
import os


def load_config(config_path: str = "partpal.json") -> dict:
    """
    Loads the configuration from the specified JSON file.

    Args:
        config_path (str): Path to the configuration file.

    Returns:
        dict: The loaded configuration dictionary.
    """
    # Check if the config file exists
    if not os.path.exists(config_path):
        raise FileNotFoundError(f"Config file not found: {config_path}")

    with open(config_path, "r") as config_file:
        config = json.load(config_file)

    return config
