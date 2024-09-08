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

        if response.status_code == 200:
            return True, response.json()
        else:
            return False, response.json()
