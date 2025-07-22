#!/usr/bin/env python3
import json

# Read the backup file with promo cards
with open('data/card_data_backup.json', 'r', encoding='utf-8') as f:
    promo_cards = json.load(f)

# Read the current card data
with open('data/card_data.json', 'r', encoding='utf-8') as f:
    current_cards = json.load(f)

# Create a set of current card numbers to avoid duplicates
current_numbers = {card['number'] for card in current_cards}

# Add promo cards that aren't already in the current data
added_count = 0
for card in promo_cards:
    if '(P)' in card['number'] and card['number'] not in current_numbers:
        current_cards.append(card)
        added_count += 1
        print(f"Added: {card['number']} - {card['name']}")

# Sort by card number
def get_sort_key(card):
    # Extract the base number for sorting
    num = card['number'].replace('(P)', '').replace('-P', '').strip()
    parts = num.split('-')
    if len(parts) >= 2:
        try:
            return int(parts[1])
        except:
            return 999
    return 999

current_cards.sort(key=get_sort_key)

# Save the updated card data
with open('data/card_data.json', 'w', encoding='utf-8') as f:
    json.dump(current_cards, f, ensure_ascii=False, indent=4)

print(f"\nTotal promo cards added: {added_count}")
print(f"Total cards now: {len(current_cards)}")