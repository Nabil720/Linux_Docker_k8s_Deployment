from dataclasses import dataclass
from typing import Optional

@dataclass
class Employee:
    name: str
    id: str
    position: str
    
    def to_dict(self):
        return {
            "name": self.name,
            "id": self.id,
            "position": self.position
        }
    
    @classmethod
    def from_dict(cls, data):
        return cls(
            name=data.get('name'),
            id=data.get('id'),
            position=data.get('position')
        )