import random
from typing import Dict

class EvasionEngine:
    @staticmethod
    def get_stealth_headers() -> Dict[str, str]:
        return {"User-Agent": "Mozilla/5.0 HLFA-Ultimate", "X-Stealth": "Active"}
    @staticmethod
    def obfuscate_payload(payload: str) -> str:
        return payload.encode('utf-8').hex()
