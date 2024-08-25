from abc import ABC, abstractmethod
from typing import Tuple, Dict, Any


class Distributor(ABC):
    def __init__(self, name: str):
        self.name = name

    @abstractmethod
    def get_product_information(self, partnumber: str) -> Tuple[bool, Dict[str, Any]]:
        pass
