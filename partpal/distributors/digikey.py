import requests
import json
import urllib.parse
from typing import Tuple, Dict, Any


class DigiKeyDistributor:
    # TODO: Replace this with a config manager
    def __init__(self, client_id, client_secret, sandbox=False) -> None:
        self.client_id = client_id
        self.client_secret = client_secret
        # TODO: Support in future for different locales and currencies
        self.locale = "US"
        self.currency = "USD"

        # Set the base URL based on whether sandbox mode is enabled
        if sandbox:
            self.base_url = "https://sandbox-api.digikey.com"
        else:
            self.base_url = "https://api.digikey.com"

    def get_access_token(self) -> bool:
        """Gets a Digikey oauth2 token, returns true for success"""

        url = f"{self.base_url}/v1/oauth2/token"
        url_data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "redirect_uri": "https://localhost",
            "grant_type": "client_credentials",
        }

        response = requests.post(url, data=url_data)

        # TODO: The token refreshes every so often, add a function to refetch a new one
        if response.status_code == 200:
            self.token = response.json()
            return True
        else:
            self.token = None
            return False

    def get_product_information(self, partnumber: str) -> Tuple[bool, Dict[str, Any]]:
        # URL Encode the part number
        pn = urllib.parse.quote(partnumber, safe="")
        url = f"{self.base_url}/products/v4/search/{pn}/productdetails"

        url_header = {
            "X-DIGIKEY-Locale-Site": self.locale,
            "X-DIGIKEY-Locale-Currency": self.currency,
            "Authorization": f"{self.token['token_type']} {self.token['access_token']}",
            "X-DIGIKEY-Client-Id": self.client_id,
        }

        response = requests.get(url, headers=url_header)

        if response.status_code == 200:
            return True, response.json()
        else:
            return False, response.json()
