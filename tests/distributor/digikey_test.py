import pytest

import partpal.core.config as pp_config
from partpal.distributors.digikey import DigiKeyDistributor


def test_get_oauth_token_success():
    """Tests whether can get OAuth2 Token"""
    config = pp_config.load_config()
    digikey = DigiKeyDistributor(
        config["distributors"]["digikey"]["client_id"],
        config["distributors"]["digikey"]["client_secret"],
        config["distributors"]["digikey"]["sandbox"],
    )

    # Make sure the Distributor is correct
    assert digikey.name == "DigiKey"

    status = digikey.get_access_token()

    # Check if a token was able to be requested
    assert status == True
    # Make sure the token is set
    assert digikey.token != None


def test_get_product_details_digikeypn():
    """Tests whether can query for product details via Digikey PN"""
    config = pp_config.load_config()
    digikey = DigiKeyDistributor(
        config["distributors"]["digikey"]["client_id"],
        config["distributors"]["digikey"]["client_secret"],
        config["distributors"]["digikey"]["sandbox"],
    )

    # Make sure the Distributor is correct
    assert digikey.name == "DigiKey"

    status = digikey.get_access_token()

    # Check if a token was able to be requested
    assert status == True
    # Make sure the token is set
    assert digikey.token != None

    # Get the product information for a STM32C031G4U6
    status, response = digikey.get_product_information("497-STM32C031G4U6-ND")

    # Successful Response
    assert status == True

    # Check the Product Desciption is correct to make sure the query worked
    assert response["description"] == "IC MCU 32BIT 16KB FLASH 28UFQFPN"

    print(response)


def test_get_product_details_mfrpn():
    """Tests whether can query for product details via Manafacturer PN"""
    config = pp_config.load_config()
    digikey = DigiKeyDistributor(
        config["distributors"]["digikey"]["client_id"],
        config["distributors"]["digikey"]["client_secret"],
        config["distributors"]["digikey"]["sandbox"],
    )

    # Make sure the Distributor is correct
    assert digikey.name == "DigiKey"

    status = digikey.get_access_token()

    # Check if a token was able to be requested
    assert status == True
    # Make sure the token is set
    assert digikey.token != None

    # Get the product information for a STM32C031G4U6
    status, response = digikey.get_product_information("STM32C031G4U6")

    # Successful Response
    assert status == True

    # Check the Product Desciption is correct to make sure the query worked
    assert response["description"] == "IC MCU 32BIT 16KB FLASH 28UFQFPN"

    print(response)
