#!/usr/bin/env python3
"""
ReSet Integration Test Suite
============================
Tests the subscription detection classifier and optionally the Go backend API.

Usage:
    # Test only the Python classifier (no services needed)
    python test_runner.py --mode classifier

    # Test classifier + Go API (requires backend running on localhost:8080)
    python test_runner.py --mode full --api-url http://localhost:8080

    # Test with custom Python service URL
    python test_runner.py --mode full --python-url http://localhost:8001 --api-url http://localhost:8080
"""

import argparse
import json
import sys
import time
from pathlib import Path

import requests

# Colors for terminal output
class Colors:
    GREEN = "\033[92m"
    RED = "\033[91m"
    YELLOW = "\033[93m"
    BLUE = "\033[94m"
    CYAN = "\033[96m"
    BOLD = "\033[1m"
    END = "\033[0m"

def color(text, c):
    return f"{c}{text}{Colors.END}"

def load_json(filename):
    with open(Path(__file__).parent / filename) as f:
        return json.load(f)

def print_header(title):
    print(f"\n{color('═' * 60, Colors.CYAN)}")
    print(f"{color('  ' + title, Colors.BOLD + Colors.CYAN)}")
    print(f"{color('═' * 60, Colors.CYAN)}")

def print_result(merchant, category, confidence, match_type, amount, is_sub, expected=None):
    status = "✓" if is_sub else "✗"
    status_color = Colors.GREEN if is_sub else Colors.RED

    conf_color = Colors.GREEN if confidence >= 0.8 else (Colors.YELLOW if confidence >= 0.5 else Colors.RED)
    match_badge = {
        "exact": color("EXACT", Colors.GREEN),
        "fuzzy": color("FUZZY", Colors.YELLOW),
        "heuristic": color("HEURISTIC", Colors.YELLOW),
        "none": color("NONE", Colors.RED),
    }.get(match_type, match_type)

    print(f"  {color(status, status_color)} {merchant or 'Unknown':<25} | {category or 'N/A':<18} | "
          f"{color(f'{confidence:.0%}', conf_color):<8} | {match_badge:<12} | "
          f"₦{amount or 'N/A'}")
    if expected and merchant != expected:
        print(f"    {color('⚠ Expected: ' + expected, Colors.YELLOW)}")

def test_classifier(python_url):
    """Test the Python AI classifier service directly."""
    print_header("TESTING PYTHON CLASSIFIER")

    transactions = load_json("test_transactions.json")

    print(f"\n{color('→', Colors.BLUE)} Testing exact matches (should be 95% confidence)...")
    exact_pass = 0
    for tx in transactions["exact_matches"]:
        try:
            resp = requests.post(f"{python_url}/classify-transaction", 
                                json={"text": tx}, timeout=5)
            r = resp.json()
            print_result(r["merchant"], r["category"], r["confidence"], 
                        r["match_type"], r["amount"], r["is_subscription"])
            if r["is_subscription"] and r["confidence"] >= 0.9:
                exact_pass += 1
        except Exception as e:
            print(f"  {color('✗ ERROR', Colors.RED)} {tx[:40]}... -> {e}")

    print(f"\n{color('→', Colors.BLUE)} Testing fuzzy matches (typos, abbreviations)...")
    fuzzy_pass = 0
    for tx in transactions["fuzzy_matches"]:
        try:
            resp = requests.post(f"{python_url}/classify-transaction", 
                                json={"text": tx}, timeout=5)
            r = resp.json()
            print_result(r["merchant"], r["category"], r["confidence"], 
                        r["match_type"], r["amount"], r["is_subscription"])
            if r["is_subscription"] and r["match_type"] == "fuzzy":
                fuzzy_pass += 1
        except Exception as e:
            print(f"  {color('✗ ERROR', Colors.RED)} {tx[:40]}... -> {e}")

    print(f"\n{color('→', Colors.BLUE)} Testing heuristic matches (unknown merchants with keywords)...")
    heuristic_pass = 0
    for tx in transactions["heuristic_matches"]:
        try:
            resp = requests.post(f"{python_url}/classify-transaction", 
                                json={"text": tx}, timeout=5)
            r = resp.json()
            print_result(r["merchant"], r["category"], r["confidence"], 
                        r["match_type"], r["amount"], r["is_subscription"])
            if r["is_subscription"] and r["match_type"] == "heuristic":
                heuristic_pass += 1
        except Exception as e:
            print(f"  {color('✗ ERROR', Colors.RED)} {tx[:40]}... -> {e}")

    print(f"\n{color('→', Colors.BLUE)} Testing non-subscriptions (should reject)...")
    reject_pass = 0
    for tx in transactions["non_subscriptions"]:
        try:
            resp = requests.post(f"{python_url}/classify-transaction", 
                                json={"text": tx}, timeout=5)
            r = resp.json()
            print_result(r["merchant"], r["category"], r["confidence"], 
                        r["match_type"], r["amount"], r["is_subscription"])
            if not r["is_subscription"]:
                reject_pass += 1
        except Exception as e:
            print(f"  {color('✗ ERROR', Colors.RED)} {tx[:40]}... -> {e}")

    # Summary
    total_exact = len(transactions["exact_matches"])
    total_fuzzy = len(transactions["fuzzy_matches"])
    total_heuristic = len(transactions["heuristic_matches"])
    total_reject = len(transactions["non_subscriptions"])

    print(f"\n{color('─' * 60, Colors.CYAN)}")
    print(f"  {color('CLASSIFIER SUMMARY', Colors.BOLD)}")
    print(f"  Exact matches:     {color(f'{exact_pass}/{total_exact}', Colors.GREEN if exact_pass==total_exact else Colors.YELLOW)}")
    print(f"  Fuzzy matches:     {color(f'{fuzzy_pass}/{total_fuzzy}', Colors.GREEN if fuzzy_pass==total_fuzzy else Colors.YELLOW)}")
    print(f"  Heuristic matches: {color(f'{heuristic_pass}/{total_heuristic}', Colors.GREEN if heuristic_pass==total_heuristic else Colors.YELLOW)}")
    print(f"  Rejected (good):   {color(f'{reject_pass}/{total_reject}', Colors.GREEN if reject_pass==total_reject else Colors.RED)}")
    print(f"{color('─' * 60, Colors.CYAN)}")

def test_batch_classifier(python_url):
    """Test batch classification endpoint."""
    print_header("TESTING BATCH CLASSIFICATION")

    transactions = load_json("test_transactions.json")
    batch = transactions["exact_matches"][:5]

    try:
        resp = requests.post(f"{python_url}/classify-batch", 
                            json={"texts": batch}, timeout=10)
        results = resp.json()["results"]
        print(f"{color('✓', Colors.GREEN)} Batch processed {len(results)} transactions in one request")
        for r in results:
            print(f"  • {r['merchant']:<20} ({r['match_type']}, {r['confidence']:.0%})")
    except Exception as e:
        print(f"{color('✗ Batch failed:', Colors.RED)} {e}")

def test_go_api(api_url, python_url):
    """Test the Go backend API endpoints."""
    print_header("TESTING GO BACKEND API")

    # Health check
    print(f"\n{color('→', Colors.BLUE)} Health check...")
    try:
        resp = requests.get(f"{api_url}/health", timeout=5)
        if resp.status_code == 200:
            print(f"  {color('✓', Colors.GREEN)} API is healthy")
        else:
            print(f"  {color('✗', Colors.RED)} Health returned {resp.status_code}")
    except Exception as e:
        print(f"  {color('✗', Colors.RED)} Cannot reach API: {e}")
        print(f"  {color('Hint:', Colors.YELLOW)} Make sure the Go server is running on {api_url}")
        return

    # Detect endpoint (tracked subscription flow)
    print(f"\n{color('→', Colors.BLUE)} Testing /subscriptions/detect (AI detection flow)...")
    test_cases = [
        ("demo@reset.ng", "NETFLIX.COM 4500.00 - Monthly subscription"),
        ("demo@reset.ng", "SPOTIFY PREMIUM 1500.00 NGN"),
        ("tunde@example.com", "MTN NIGERIA 3000.00 - 10GB data plan"),
    ]

    detected_ids = []
    for email, text in test_cases:
        try:
            resp = requests.post(f"{api_url}/subscriptions/detect", 
                                data={"email": email, "text": text}, timeout=10)
            if resp.status_code == 200:
                # Parse HTML response to extract ID (simplified)
                html = resp.text
                if "Tracked ID" in html:
                    # Extract ID from the HTML
                    import re
                    id_match = re.search(r'Tracked ID</span>\s*<span[^>]*>([^<]+)</span>', html)
                    if id_match:
                        detected_ids.append(id_match.group(1).strip())
                print(f"  {color('✓', Colors.GREEN)} Detected: {text[:40]}...")
            else:
                print(f"  {color('✗', Colors.RED)} Detect failed: {resp.status_code}")
        except Exception as e:
            print(f"  {color('✗', Colors.RED)} Error: {e}")

    # List tracked subscriptions
    print(f"\n{color('→', Colors.BLUE)} Testing /tracked (list subscriptions)...")
    try:
        resp = requests.get(f"{api_url}/tracked?email=demo@reset.ng", timeout=5)
        if resp.status_code == 200:
            print(f"  {color('✓', Colors.GREEN)} Listed tracked subscriptions for demo@reset.ng")
        else:
            print(f"  {color('✗', Colors.RED)} List failed: {resp.status_code}")
    except Exception as e:
        print(f"  {color('✗', Colors.RED)} Error: {e}")

    # Cancel a tracked subscription (if we have IDs)
    if detected_ids:
        print(f"\n{color('→', Colors.BLUE)} Testing /subscriptions/tracked/cancel...")
        try:
            resp = requests.post(f"{api_url}/subscriptions/tracked/cancel",
                                data={"id": detected_ids[0]}, timeout=10)
            if resp.status_code == 200:
                print(f"  {color('✓', Colors.GREEN)} Cancelled subscription {detected_ids[0][:8]}...")
            else:
                print(f"  {color('✗', Colors.RED)} Cancel failed: {resp.status_code}")
        except Exception as e:
            print(f"  {color('✗', Colors.RED)} Error: {e}")

def simulate_real_life_scenario():
    """Print a step-by-step real life scenario."""
    print_header("REAL LIFE SCENARIO SIMULATION")

    scenario = """
  👤 USER: Chioma Nwosu (chioma@example.com)
  📍 LOCATION: Lagos, Nigeria
  🏦 BANK: Access Bank

  STEP 1 — THE PROBLEM
  ─────────────────────
  Chioma checks her bank statement and sees several deductions she doesn't 
  recognize or forgot about:

    • "NETFLIX.COM 4500.00" — she signed up 6 months ago for one movie
    • "SPOTIFY 1500.00" — she uses Apple Music now, not Spotify  
    • "ADOBE 25000.00" — freelancer project ended 3 months ago
    • "MTN NG 3000.00" — auto-renews every month, she has WiFi at home now
    • "DSTV 9000.00" — she travels often, barely watches

  Total leaking: ₦42,000/month = ₦504,000/year

  STEP 2 — USING ReSet DETECT
  ────────────────────────────
  Chioma opens ReSet and pastes her transaction texts one by one:

    Input:  "DEBIT: NGN 4,500.00 - NETFLIX.COM - Monthly subscription"
    → AI detects: Netflix, Streaming, 95% confidence (exact match)
    → Saved to her "Tracked" list

    Input:  "SPOTIFY PREMIUM 1,500.00 NGN - Recurring payment"  
    → AI detects: Spotify, Streaming, 95% confidence (exact match)
    → Saved to tracked list

    Input:  "ADOBE CREATIVE CLOUD 25,000.00" 
    → AI detects: Adobe Creative Cloud, Productivity, 95% confidence
    → Saved to tracked list

    Input:  "MTN NIGERIA 3,000.00 - 10GB monthly data"
    → AI detects: MTN, Telecom, 95% confidence
    → Saved to tracked list

    Input:  "DSTV/MULTICHOICE 9,000.00 - Compact package"
    → AI detects: DStv, Cable TV, 95% confidence
    → Saved to tracked list

  STEP 3 — REVIEWING TRACKED SUBSCRIPTIONS
  ─────────────────────────────────────────
  Chioma visits the "Tracked" page. She sees:

    ┌─────────────────────────────────────────────┐
    │ Netflix        Streaming    ₦4,500   [Cancel]│
    │ Spotify        Streaming    ₦1,500   [Cancel]│
    │ Adobe CC       Productivity ₦25,000  [Cancel]│
    │ MTN            Telecom      ₦3,000   [Cancel]│
    │ DStv           Cable TV     ₦9,000   [Cancel]│
    └─────────────────────────────────────────────┘

    Monthly total: ₦42,000

  STEP 4 — CANCELLING (2 paths)
  ─────────────────────────────

  PATH A — Direct/Paystack subscriptions:
  If the merchant uses Paystack (e.g., local Nigerian services):
    • Chioma clicks [Cancel] on MTN
    • ReSet calls Paystack API to disable the subscription
    • Paystack stops future deductions
    • Status changes to "Canceled" in ReSet
    • Cache is invalidated immediately

  PATH B — Manual/tracked only (Netflix, Spotify, Adobe):
  These are international services not on Paystack:
    • Chioma clicks [Cancel] in ReSet
    • ReSet marks it as "Canceled" locally for tracking
    • Chioma is reminded to cancel manually on the merchant's site
    • ReSet keeps the record so she knows it's "handled"

  STEP 5 — VERIFICATION
  ─────────────────────
  A week later, Chioma checks her new bank statement:
    ✓ No more Netflix charges
    ✓ No more Spotify charges  
    ✓ No more Adobe charges
    ✓ No more MTN auto-renewal
    ✓ No more DStv charges

  She saved ₦42,000/month. Over a year, that's ₦504,000 kept in her pocket.

  STEP 6 — ONGOING MONITORING (Future feature)
  ────────────────────────────────────────────
  ReSet could periodically scan her email for bank alerts, or she could
  paste new transactions monthly. Any NEW subscription gets flagged:

    Input: "UNKNOWN MERCHANT 5000.00 - Monthly premium plan"
    → AI detects: heuristic match (unknown merchant + "monthly" + "plan")
    → Confidence: 40% — flagged for review
    → Chioma decides: "Oh, that's my new gym. Keep it."
"""
    print(scenario)

def main():
    parser = argparse.ArgumentParser(description="ReSet Test Runner")
    parser.add_argument("--mode", choices=["classifier", "full", "scenario"], 
                       default="classifier", help="Test mode")
    parser.add_argument("--python-url", default="http://localhost:8001",
                       help="Python AI service URL")
    parser.add_argument("--api-url", default="http://localhost:8080",
                       help="Go backend API URL")
    args = parser.parse_args()

    print(color("""
    ╔══════════════════════════════════════════════════════════════╗
    ║                                                              ║
    ║           ReSet — Subscription Detection Test Suite          ║
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
    """, Colors.BOLD + Colors.CYAN))

    if args.mode in ("classifier", "full"):
        # Test Python classifier
        try:
            requests.get(f"{args.python_url}/health", timeout=3)
            test_classifier(args.python_url)
            test_batch_classifier(args.python_url)
        except requests.ConnectionError:
            print(f"\n{color('✗ Python service not running at ' + args.python_url, Colors.RED)}")
            print(f"{color('  Start it with: uvicorn main:app --reload --port 8001', Colors.YELLOW)}")
            if args.mode == "classifier":
                sys.exit(1)

        if args.mode == "full":
            test_go_api(args.api_url, args.python_url)

    if args.mode in ("full", "scenario"):
        simulate_real_life_scenario()

    print(f"\n{color('✓ Test run complete!', Colors.GREEN + Colors.BOLD)}\n")

if __name__ == "__main__":
    main()