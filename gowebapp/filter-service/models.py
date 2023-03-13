from dataclasses import dataclass
from typing import List, Optional


@dataclass
class Filters:
    """Class for keeping track of filters."""
    platforms: Optional[List[str]]
    levels: Optional[List[str]]
    start_time: Optional[str]
    end_time: Optional[str]
    message_contains: Optional[str]
    sourcetopic: str
    sinktopic: str

    # quantity_on_hand: int = 0

    def __init__(
            self, 
            platforms: List[str] = None, 
            levels: List[str] = None, 
            sinktopic: str = "",
            sourcetopic: str = "",
            start_time: str = None,
            end_time: str = None,
            message_contains: str = None
        ):

        self.levels = levels
        self.platforms = platforms
        self.sinktopic = sinktopic
        self.sourcetopic = sourcetopic
        self.start_time = start_time
        self.end_time = end_time
        self.message_contains = message_contains

    @classmethod
    def from_json(cls, json_dict):
        return cls(**json_dict)
