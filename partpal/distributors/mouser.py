import requests
import json
import urllib.parse
from typing import Tuple, Dict, Any

from .distributors_base import Distributor


class MouserDistributor(Distributor):
    # TODO: Replace this with a config manager
    def __init__(self, api_key) -> None:
        super().__init__(name="Mouser")
        self.api_key = api_key
        # TODO: Support in future for different locales and currencies
        self.locale = "US"
        self.currency = "USD"

        self.base_url = "https://api.mouser.com/api/v1/"

    def get_product_information(self, partnumber: str) -> Tuple[bool, Dict[str, Any]]:
        # URL Encode the part number
        url_apikey = urllib.parse.quote(self.api_key, safe="")
        url = f"{self.base_url}/search/partnumber?apikey={url_apikey}"

        # Create the JSON payload
        data = {
            "SearchByPartRequest": {
                "mouserPartNumber": partnumber,
                "partSearchOptions": "string",
            }
        }
        response = requests.post(url, json=data)

        json_response = response.json()

        if response.status_code == 200:
            product_data = {
                "partNumber": json_response["SearchResults"]["Parts"][0][
                    "MouserPartNumber"
                ],
                "mfrPartNumber": json_response["SearchResults"]["Parts"][0][
                    "ManufacturerPartNumber"
                ],
                "manufacturer": json_response["SearchResults"]["Parts"][0][
                    "Manufacturer"
                ],
                "unitPrice": json_response["SearchResults"]["Parts"][0]["PriceBreaks"][
                    0
                ]["Price"],
                "availability": json_response["SearchResults"]["Parts"][0][
                    "AvailabilityInStock"
                ],
                "description": json_response["SearchResults"]["Parts"][0][
                    "Description"
                ],
                "datasheet": json_response["SearchResults"]["Parts"][0]["DataSheetUrl"],
                "productUrl": json_response["SearchResults"]["Parts"][0][
                    "ProductDetailUrl"
                ],
            }
            return True, product_data
        else:
            return False, json_response
