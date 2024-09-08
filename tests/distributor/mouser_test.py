import pytest

import partpal.core.config as pp_config
from partpal.distributors.mouser import MouserDistributor


def test_init():
    """Tests whether can initialize"""
    config = pp_config.load_config()
    digikey = MouserDistributor(
        config["distributors"]["mouser"]["api_key"],
    )

    # Make sure the Distributor is correct
    assert digikey.name == "Mouser"


def test_get_product_details_mouserpn():
    """Tests whether can query for product details via Mouser PN"""

    config = pp_config.load_config()
    digikey = MouserDistributor(
        config["distributors"]["mouser"]["api_key"],
    )

    # Make sure the Distributor is correct
    assert digikey.name == "Mouser"
    # Get the product information for a STM32C031G4U6
    status, response = digikey.get_product_information("511-STM32C031G4U6")

    # Successful Response
    assert status == True

    # Check the Product Desciption is correct to make sure the query worked
    assert (
        response["Product"]["Description"]["ProductDescription"]
        == "IC MCU 32BIT 16KB FLASH 28UFQFPN"
    )

    print(response)


def test_get_product_details_mfrpn():
    """Tests whether can query for product details via Manafacturer PN"""
    config = pp_config.load_config()
    digikey = MouserDistributor(
        config["distributors"]["mouser"]["api_key"],
    )

    # Make sure the Distributor is correct
    assert digikey.name == "Mouser"
    # Get the product information for a STM32C031G4U6
    status, response = digikey.get_product_information("STM32C031G4U6")

    # Successful Response
    assert status == True

    # Check the Product Desciption is correct to make sure the query worked
    assert (
        response["Product"]["Description"]["ProductDescription"]
        == "IC MCU 32BIT 16KB FLASH 28UFQFPN"
    )

    print(response)
