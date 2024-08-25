import csv
from typing import List, Dict

# Define the type aliases
Component = Dict[str, str]
BOM = List[Component]

# When reading in a BOM, column names could differ but ultimately stand for the same thing
# This map contains a mapping of possible variaitons of column name to a unified definition we use in PartPal
COLUMN_MAP = {
    "part_number": ["P/N", "Part Number", "Part_No", "PartNumber"],
    "quantity": ["Qty", "Quantity", "QTY", "Q"],
    "name": ["Reference", "Component", "Name", "Item"],
    "description": ["Description", "Desc", "Details"],
    "cost": ["Cost", "Price", "Unit Price"],
    "distributor": ["Distributor", "Vendor"],
}

# Fields in a BOM we are interested in
FIELDS_OF_INTEREST = [
    "name",
    "quantity",
    "part_number",
    "description",
    "cost",
    "distributor",
]


def normalize_column_name(header: str, column_map: Dict[str, List[str]]) -> str:
    """Normalize the column names using the COLUMN_MAP.

    Args:
        header: The original header name from the CSV.
        column_map: A dictionary mapping standard names to possible CSV headers.

    Returns:
        The normalized column name, or the original if no match is found.
    """
    for key, possible_names in column_map.items():
        if header in possible_names:
            return key
    return header  # Return the original if no match is found


def parse_bom_csv(file_path: str) -> BOM:
    """Parse the BOM CSV file and return a list of dictionaries representing components,
    including only the fields of interest.

    Args:
        file_path: The path to the BOM CSV file.

    Returns:
        A list of dictionaries, each representing a component in the BOM,
        with only the fields of interest.
    """
    with open(file_path, mode="r", newline="") as csvfile:
        reader = csv.DictReader(csvfile)
        normalized_headers: List[str] = [
            normalize_column_name(header, COLUMN_MAP) for header in reader.fieldnames
        ]

        # Create a list of dictionaries with normalized keys, filtered by FIELDS_OF_INTEREST
        bom: BOM = []
        for row in reader:
            component: Component = {
                normalized_headers[i]: row[header]
                for i, header in enumerate(reader.fieldnames)
                if normalized_headers[i] in FIELDS_OF_INTEREST
            }
            bom.append(component)

    return bom
