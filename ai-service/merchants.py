"""
Known merchant dictionary for subscription detection.

Each entry maps a set of raw-text patterns (keywords that tend to appear
in bank/transaction statements) to a canonical merchant name and category.

Matching strategy (handled in classifier.py):
1. Exact/substring match against these patterns (case-insensitive)
2. Fuzzy match fallback for near-misses (typos, extra codes, etc.)
"""

MERCHANTS = [
    # Streaming
    {"patterns": ["netflix"], "name": "Netflix", "category": "Streaming"},
    {"patterns": ["spotify"], "name": "Spotify", "category": "Streaming"},
    {"patterns": ["amazon prime", "amzn prime", "prime video"], "name": "Amazon Prime", "category": "Streaming"},
    {"patterns": ["disney+", "disney plus"], "name": "Disney+", "category": "Streaming"},
    {"patterns": ["youtube premium", "youtube*premium", "google *youtube"], "name": "YouTube Premium", "category": "Streaming"},
    {"patterns": ["showmax"], "name": "Showmax", "category": "Streaming"},
    {"patterns": ["hbo max", "hbomax"], "name": "HBO Max", "category": "Streaming"},
    {"patterns": ["apple music"], "name": "Apple Music", "category": "Streaming"},
    {"patterns": ["deezer"], "name": "Deezer", "category": "Streaming"},

    # Cloud / SaaS / Productivity
    {"patterns": ["icloud", "apple.com/bill"], "name": "iCloud", "category": "Cloud Storage"},
    {"patterns": ["google one", "google storage"], "name": "Google One", "category": "Cloud Storage"},
    {"patterns": ["dropbox"], "name": "Dropbox", "category": "Cloud Storage"},
    {"patterns": ["microsoft 365", "office 365", "msft *365"], "name": "Microsoft 365", "category": "Productivity"},
    {"patterns": ["notion"], "name": "Notion", "category": "Productivity"},
    {"patterns": ["canva"], "name": "Canva", "category": "Productivity"},
    {"patterns": ["adobe", "adobe creative"], "name": "Adobe Creative Cloud", "category": "Productivity"},
    {"patterns": ["github"], "name": "GitHub", "category": "Developer Tools"},
    {"patterns": ["chatgpt", "openai"], "name": "ChatGPT Plus", "category": "AI Tools"},
    {"patterns": ["claude.ai", "anthropic"], "name": "Claude", "category": "AI Tools"},

    # Fitness / Health
    {"patterns": ["planet fitness"], "name": "Planet Fitness", "category": "Fitness"},
    {"patterns": ["gymfit", "fitness club", "gym membership"], "name": "Gym Membership", "category": "Fitness"},
    {"patterns": ["strava"], "name": "Strava", "category": "Fitness"},

    # Telecom / Utilities
    {"patterns": ["mtn", "mtn nigeria"], "name": "MTN", "category": "Telecom"},
    {"patterns": ["airtel"], "name": "Airtel", "category": "Telecom"},
    {"patterns": ["glo"], "name": "Glo", "category": "Telecom"},
    {"patterns": ["9mobile", "etisalat"], "name": "9mobile", "category": "Telecom"},
    {"patterns": ["dstv", "multichoice"], "name": "DStv", "category": "Cable TV"},
    {"patterns": ["gotv"], "name": "GOtv", "category": "Cable TV"},
    {"patterns": ["startimes"], "name": "StarTimes", "category": "Cable TV"},

    # Internet / ISPs (Nigeria)
    {"patterns": ["spectranet"], "name": "Spectranet", "category": "Internet"},
    {"patterns": ["smile communications", "smile 4g", "smilecoms"], "name": "Smile", "category": "Internet"},
    {"patterns": ["ipnx"], "name": "ipNX", "category": "Internet"},
    {"patterns": ["swift 4g", "swift networks"], "name": "Swift 4G", "category": "Internet"},

    # Streaming / entertainment (Nigeria / Africa)
    {"patterns": ["irokotv", "iroko tv"], "name": "IrokoTV", "category": "Streaming"},
    {"patterns": ["boomplay"], "name": "Boomplay", "category": "Streaming"},
    {"patterns": ["audiomack"], "name": "Audiomack", "category": "Streaming"},

    # Shopping / delivery memberships
    {"patterns": ["jumia prime"], "name": "Jumia Prime", "category": "Shopping"},
    {"patterns": ["chowdeck plus", "chowdeck premium"], "name": "Chowdeck+", "category": "Food Delivery"},

    # Fintech savings / investment plans (recur like subscriptions on statements)
    {"patterns": ["piggyvest", "piggy vest"], "name": "PiggyVest", "category": "Savings/Investment"},
    {"patterns": ["cowrywise"], "name": "Cowrywise", "category": "Savings/Investment"},
    {"patterns": ["risevest", "rise vest"], "name": "RiseVest", "category": "Savings/Investment"},

    # Finance / news
    {"patterns": ["nytimes", "new york times"], "name": "NYTimes", "category": "News"},
    {"patterns": ["medium.com", "medium membership"], "name": "Medium", "category": "News"},
    {"patterns": ["premium times"], "name": "Premium Times", "category": "News"},
    {"patterns": ["punch newspapers", "punchng"], "name": "Punch", "category": "News"},
]

# Keywords that hint "this looks like a subscription" even for unknown merchants
SUBSCRIPTION_KEYWORDS = [
    "subscription", "premium", "membership", "monthly", "recurring",
    "renewal", "plan", "pro plan", "annual plan",
]