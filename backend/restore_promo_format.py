#!/usr/bin/env python3
import json
import shutil

# First backup current data
shutil.copy('data/card_data.json', 'data/card_data_before_restore.json')

# Read current data
with open('data/card_data.json', 'r', encoding='utf-8') as f:
    cards = json.load(f)

# Read official data to understand which are promo vs parallel
with open('data/official_cardlist.json', 'r', encoding='utf-8') as f:
    official = json.load(f)

# Create a set of official promo cards (those with -P in official site)
official_promos = {c['card_number'] for c in official if '-P' in c['card_number']}

# Fix the data
fixed_cards = []
seen = set()

for card in cards:
    card_key = (card['number'], card.get('name', ''))
    
    # Skip duplicates
    if card_key in seen:
        continue
    seen.add(card_key)
    
    # If it's a -P card from official site, it's actually a promo
    # Convert to (P) format as the user specified
    if card['number'] in official_promos:
        # This is a promo card, should use (P) format
        base_number = card['number'].replace('-P', '')
        card['number'] = f"{base_number} (P)"
        # Make sure it has regular rarity (not -P rarity)
        if card['local_image_path']:
            # Fix image path if needed
            if '_R-P' in card['local_image_path']:
                card['local_image_path'] = card['local_image_path'].replace('_R-P', '-P_R')
            elif '_SR-P' in card['local_image_path']:
                card['local_image_path'] = card['local_image_path'].replace('_SR-P', '-P_SR')
            elif '_C-P' in card['local_image_path']:
                card['local_image_path'] = card['local_image_path'].replace('_C-P', '-P_C')
    
    # True parallel cards keep -P in number and have -P in rarity
    # These are cards like F-016-P with SR-P rarity
    elif '-P' in card['number'] and card['number'] not in official_promos:
        # This is a true parallel card
        # Should have -P in both number and rarity
        pass
    
    fixed_cards.append(card)

# Fix F-068 name
for card in fixed_cards:
    if card['number'] == 'F-068':
        card['name'] = 'デコーレーション'  # Fix the dash character

# Sort and save
fixed_cards.sort(key=lambda x: (x['number'].replace(' (P)', '-P.1').replace('-P', '-P.2'), x['name']))

with open('data/card_data.json', 'w', encoding='utf-8') as f:
    json.dump(fixed_cards, f, ensure_ascii=False, indent=4)

print(f"Fixed {len(fixed_cards)} cards")
print(f"Promo cards (with (P)): {sum(1 for c in fixed_cards if '(P)' in c['number'])}")
print(f"Parallel cards (with -P): {sum(1 for c in fixed_cards if '-P' in c['number'] and '(P)' not in c['number'])}")