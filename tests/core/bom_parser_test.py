import pytest

import pathlib
import partpal.core.bom_parser as bom_parser


def test_readBOM():
    test_file_path = pathlib.Path(__file__).parent.resolve()
    test_csv_path = test_file_path.joinpath("test_bom.csv")
    bom = bom_parser.parse_bom_csv(test_csv_path.absolute())
    print(bom)
