#!/usr/bin/env python3
import json
import os

# Read current card data
with open('data/card_data.json', 'r', encoding='utf-8') as f:
    cards = json.load(f)

# Remove duplicates by creating a dict with card number as key
unique_cards = {}
for card in cards:
    unique_cards[card['number']] = card

# Read backup to get promo cards
if os.path.exists('data/card_data_backup.json'):
    with open('data/card_data_backup.json', 'r', encoding='utf-8') as f:
        backup_cards = json.load(f)
    
    # Add promo cards (with (P) in number) from backup
    for card in backup_cards:
        if '(P)' in card['number'] and card['number'] not in unique_cards:
            unique_cards[card['number']] = card
            print(f"Added promo card: {card['number']} - {card['name']}")

# Convert back to list and sort
cards_list = list(unique_cards.values())

# Sort by extracting number
def get_sort_key(card):
    num = card['number'].replace('(P)', '').replace('-P', '').strip()
    parts = num.split('-')
    if len(parts) >= 2:
        try:
            base_num = int(parts[1])
            # Add small offset for variants
            if '-P' in card['number']:
                base_num += 0.2
            elif '(P)' in card['number']:
                base_num += 0.1
            return base_num
        except:
            return 999
    return 999

cards_list.sort(key=get_sort_key)

# Save cleaned data
with open('data/card_data.json', 'w', encoding='utf-8') as f:
    json.dump(cards_list, f, ensure_ascii=False, indent=4)

print(f"Total unique cards: {len(cards_list)}")
print("Duplicates removed and data cleaned!")