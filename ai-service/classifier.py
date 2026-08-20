"""
Core classification logic.

Given a raw transaction string, figure out:
- merchant name (canonical, cleaned up)
- category
- whether it's likely a subscription
- confidence level + how we matched it (for transparency in the demo)
"""

import re
from rapidfuzz import fuzz, process

from merchants import MERCHANTS, SUBSCRIPTION_KEYWORDS

# Flatten the merchant list into a lookup: pattern -> (name, category)
_PATTERN_TO_MERCHANT = {}
for entry in MERCHANTS:
    for pattern in entry["patterns"]:
        _PATTERN_TO_MERCHANT[pattern] = (entry["name"], entry["category"])

_ALL_PATTERNS = list(_PATTERN_TO_MERCHANT.keys())


def _clean_text(raw: str) -> str:
    """Lowercase and strip noise like trailing reference codes/numbers."""
    text = raw.lower().strip()
    text = re.sub(r"\s+", " ", text)
    return text


def _extract_amount(raw: str):
    """Pull the first dollar/naira-style number out of the string, if present."""
    match = re.search(r"(\d+[.,]\d{2})", raw)
    if match:
        try:
            return float(match.group(1).replace(",", ""))
        except ValueError:
            return None
    return None


def classify_transaction(raw_text: str) -> dict:
    text = _clean_text(raw_text)
    amount = _extract_amount(raw_text)

    # 1. Exact / substring match
    for pattern, (name, category) in _PATTERN_TO_MERCHANT.items():
        if pattern in text:
            return {
                "raw_text": raw_text,
                "merchant": name,
                "category": category,
                "amount": amount,
                "is_subscription": True,
                "confidence": 0.95,
                "match_type": "exact",
            }

    # 2. Fuzzy match fallback (catches typos / unusual formatting)
    best_match = process.extractOne(text, _ALL_PATTERNS, scorer=fuzz.partial_ratio)
    if best_match and best_match[1] >= 80:  # similarity threshold
        pattern = best_match[0]
        name, category = _PATTERN_TO_MERCHANT[pattern]
        return {
            "raw_text": raw_text,
            "merchant": name,
            "category": category,
            "amount": amount,
            "is_subscription": True,
            "confidence": round(best_match[1] / 100, 2),
            "match_type": "fuzzy",
        }

    # 3. Heuristic fallback: unknown merchant, but text hints at subscription
    if any(keyword in text for keyword in SUBSCRIPTION_KEYWORDS):
        return {
            "raw_text": raw_text,
            "merchant": raw_text.strip(),
            "category": "Unknown",
            "amount": amount,
            "is_subscription": True,
            "confidence": 0.4,
            "match_type": "heuristic",
        }

    # 4. No match — not detected as a subscription
    return {
        "raw_text": raw_text,
        "merchant": None,
        "category": None,
        "amount": amount,
        "is_subscription": False,
        "confidence": 0.0,
        "match_type": "none",
    }


def classify_batch(raw_texts: list[str]) -> list[dict]:
    return [classify_transaction(t) for t in raw_texts]